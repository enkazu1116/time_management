package valueobjects

import "errors"

const (
	MIN_PASSWORD_LENGTH              = 8
	PASSWORD_EMPTY_ERROR_MESSAGE     = "パスワードが入力されていません"
	PASSWORD_TOO_SHORT_ERROR_MESSAGE = "パスワードが短すぎます"
)

type Password struct {
	password string
}

func NewPassword(password string) *Password {
	return &Password{
		password: password,
	}
}

func (p *Password) validatePassword() error {
	if p.password == "" {
		return errors.New(PASSWORD_EMPTY_ERROR_MESSAGE)
	}

	if len(p.password) < MIN_PASSWORD_LENGTH {
		return errors.New(PASSWORD_TOO_SHORT_ERROR_MESSAGE)
	}

	return nil
}
