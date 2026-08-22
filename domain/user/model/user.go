package model

import (
	"errors"
	"net/mail"
	"strings"
	"unicode"

	"time_management/domain/shared/messages"
	sharedmodel "time_management/domain/shared/model"
)

var (
	ErrUserIDRequired     = errors.New(messages.UserIDRequired)
	ErrUserNameRequired   = errors.New(messages.UserNameRequired)
	ErrUserNameTooLong    = errors.New(messages.UserNameTooLong)
	ErrEmailRequired      = errors.New(messages.EmailRequired)
	ErrInvalidEmail       = errors.New(messages.InvalidEmail)
	ErrInvalidEmailDomain = errors.New(messages.InvalidEmailDomain)
	ErrRoleRequired       = errors.New(messages.RoleRequired)
)

type User struct {
	UserID       string
	Password     string
	Name         string
	Email        string
	RegularStart string
	RegularEnd   string
	Role         Role
}

type UserParams struct {
	UserID       string
	Password     string
	Name         string
	Email        string
	RegularStart string
	RegularEnd   string
	Role         Role
}

func NewUser(userID string, password string, name string, email string, regularStart string, regularEnd string, role Role) User {
	return NewUserFromParams(UserParams{
		UserID:       userID,
		Password:     password,
		Name:         name,
		Email:        email,
		RegularStart: regularStart,
		RegularEnd:   regularEnd,
		Role:         role,
	})
}

func NewUserFromParams(params UserParams) User {
	return User{
		UserID:       params.UserID,
		Password:     params.Password,
		Name:         params.Name,
		Email:        params.Email,
		RegularStart: params.RegularStart,
		RegularEnd:   params.RegularEnd,
		Role:         params.Role,
	}
}

func (u User) Validate() error {
	if u.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(u.UserID); err != nil {
		return err
	}
	if err := validateUserName(u.Name); err != nil {
		return err
	}
	if err := validateEmail(u.Email); err != nil {
		return err
	}
	if err := u.Role.Validate(); err != nil {
		return err
	}
	return nil
}

func validateUserName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrUserNameRequired
	}
	if len([]rune(name)) > 100 {
		return ErrUserNameTooLong
	}
	return nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return ErrEmailRequired
	}
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return ErrInvalidEmail
	}
	emailParts := strings.Split(address.Address, "@")
	if len(emailParts) != 2 || !isValidEmailDomain(emailParts[1]) {
		return ErrInvalidEmailDomain
	}
	return nil
}

func isValidEmailDomain(domain string) bool {
	domain = strings.TrimSpace(strings.ToLower(domain))
	if domain == "" || len(domain) > 253 {
		return false
	}
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false
	}

	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return false
	}
	for _, label := range labels {
		if !isValidEmailDomainLabel(label) {
			return false
		}
	}
	return hasLetter(labels[len(labels)-1])
}

func isValidEmailDomainLabel(label string) bool {
	if label == "" || len(label) > 63 {
		return false
	}
	if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
		return false
	}
	for _, r := range label {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			continue
		}
		return false
	}
	return true
}

func hasLetter(value string) bool {
	for _, r := range value {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}
