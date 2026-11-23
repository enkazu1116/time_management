package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"uuid"

	User "time_management/backend/domain/user/model"
	valueobjects "time_management/backend/domain/user/model/valueObjects"
	userusecases "time_management/backend/domain/user/usecases"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type UserHandler struct {
	usecase *userusecases.UserUsecase
}

func NewUserHandler(usecase *userusecases.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) Register(r *mux.Router) {
	r.HandleFunc("/users", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", h.GetByID).Methods(http.MethodGet)
}

type createUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

type userResponse struct {
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	now := time.Now().UTC()
	user := User.NewUser(
		uuid.New(),
		valueobjects.NewUserName(req.FirstName, req.LastName),
		valueobjects.NewPassword(req.Password),
		valueobjects.NewEmailAddress(req.Email),
		now,
		now,
		time.Time{},
	)

	if err := h.usecase.CreateUser(r.Context(), user); err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	WriteJSON(w, http.StatusCreated, toUserResponse(user))
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user, err := h.usecase.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			WriteJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
			return
		}
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	WriteJSON(w, http.StatusOK, toUserResponse(user))
}

func toUserResponse(user *User.User) userResponse {
	return userResponse{
		UserID:    user.UserID.String(),
		UserName:  user.UserName.GetFullName(),
		Email:     user.Email.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
