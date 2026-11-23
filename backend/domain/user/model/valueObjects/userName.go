package valueobjects

import (
	"errors"
	"strings"
)

const (
	MAX_FIRST_NAME_LENGTH     = 20
	MAX_LAST_NAME_LENGTH      = 20
	NAME_LENGTH_ERROR_MESSAGE = "名前の文字数が多すぎます"
	NAME_EMPTY_ERROR_MESSAGE  = "名前が空です"
)

type UserName struct {
	firstName string
	lastName  string
}

func NewUserName(firstName string, lastName string) *UserName {
	return &UserName{
		firstName: firstName,
		lastName:  lastName,
	}
}

func (u *UserName) GetFullName() string {
	if u == nil {
		return ""
	}
	if err := u.validateUserName(); err != nil {
		return err.Error()
	}

	return u.firstName + " " + u.lastName
}

// NewUserNameFromFullName reconstructs UserName from a single DB column value.
func NewUserNameFromFullName(fullName string) *UserName {
	firstName, lastName, _ := strings.Cut(fullName, " ")
	return NewUserName(firstName, lastName)
}

func (u *UserName) validateUserName() error {

	if u.firstName == "" || u.lastName == "" {
		return errors.New(NAME_EMPTY_ERROR_MESSAGE)
	}

	if len(u.firstName) > 20 || len(u.lastName) > 20 {
		return errors.New(NAME_LENGTH_ERROR_MESSAGE)
	}

	return nil
}
