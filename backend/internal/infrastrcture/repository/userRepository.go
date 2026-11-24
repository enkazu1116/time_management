package repository

import (
	"context"
	"uuid"

	User "time_management/backend/internal/domain/user/model"
	valueobjects "time_management/backend/internal/domain/user/model/valueObjects"
	"time_management/backend/internal/infrastrcture/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) *UserRepository {
	return &UserRepository{q: q}
}

func (r *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*User.User, error) {
	user, err := r.q.GetUserById(ctx, pgtypeConvertUUID(id))
	if err != nil {
		return nil, err
	}
	return toDomainUser(user), nil
}

func (r *UserRepository) GetUserByUserName(ctx context.Context, userName string) ([]*User.User, error) {
	rows, err := r.q.GetUserByUserName(ctx, userName)
	if err != nil {
		return nil, err
	}

	users := make([]*User.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, toDomainUser(row))
	}
	return users, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*User.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return toDomainUser(user), nil
}

func (r *UserRepository) GetUserByUserNameAndEmail(ctx context.Context, userName string, email string) (*User.User, error) {
	user, err := r.q.GetUserByUserNameAndEmail(ctx, sqlc.GetUserByUserNameAndEmailParams{
		UserName: userName,
		Email:    email,
	})
	if err != nil {
		return nil, err
	}
	return toDomainUser(user), nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *User.User) error {
	_, err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		ID:        pgtypeConvertUUID(user.UserID),
		UserName:  user.UserName.GetFullName(),
		Password:  user.Password.String(),
		Email:     user.Email.String(),
		CreatedAt: pgtype.Timestamp{Time: user.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: user.UpdatedAt, Valid: true},
		DeletedAt: pgtype.Timestamp{Time: user.DeletedAt, Valid: !user.DeletedAt.IsZero()},
	})
	return err
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *User.User) error {
	_, err := r.q.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:        pgtypeConvertUUID(user.UserID),
		UserName:  user.UserName.GetFullName(),
		Password:  user.Password.String(),
		Email:     user.Email.String(),
		UpdatedAt: pgtype.Timestamp{Time: user.UpdatedAt, Valid: true},
	})
	return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.q.DeleteUser(ctx, pgtypeConvertUUID(id))
	return err
}

func toDomainUser(user sqlc.User) *User.User {
	return User.NewUser(
		uuid.UUID(user.ID.Bytes),
		valueobjects.NewUserNameFromFullName(user.UserName),
		valueobjects.NewPassword(user.Password),
		valueobjects.NewEmailAddress(user.Email),
		user.CreatedAt.Time,
		user.UpdatedAt.Time,
		user.DeletedAt.Time,
	)
}
