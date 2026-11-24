package valueobjects

import (
	"errors"
	"strings"
)

const (
	MAX_EMAIL_ADDRESS_LENGTH            = 255
	EMAIL_ADDRESS_LENGTH_ERROR_MESSAGE  = "メールアドレスの文字数が多すぎます"
	EMAIL_ADDRESS_EMPTY_ERROR_MESSAGE   = "メールアドレスが空です"
	INVALID_EMAIL_ADDRESS_ERROR_MESSAGE = "メールアドレスが無効です"
)

type EmailAddress struct {
	emailAddress string
}

func NewEmailAddress(emailAddress string) *EmailAddress {
	return &EmailAddress{
		emailAddress: emailAddress,
	}
}

func (e *EmailAddress) String() string {
	if e == nil {
		return ""
	}
	return e.emailAddress
}

func (e *EmailAddress) validateEmailAddress() error {
	if e.emailAddress == "" {
		return errors.New(EMAIL_ADDRESS_EMPTY_ERROR_MESSAGE)
	}

	if len(e.emailAddress) > MAX_EMAIL_ADDRESS_LENGTH {
		return errors.New(EMAIL_ADDRESS_LENGTH_ERROR_MESSAGE)
	}

	if !strings.Contains(e.emailAddress, "@") {
		return errors.New(INVALID_EMAIL_ADDRESS_ERROR_MESSAGE)
	}

	return nil
}
