package model

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"regexp"

	"time_management/domain/shared/messages"
)

var (
	ErrUUIDRequired = errors.New(messages.UUIDRequired)
	ErrInvalidUUID  = errors.New(messages.InvalidUUID)
)

var uuidPattern = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

type UUID string

func NewUUID(value string) UUID {
	return UUID(value)
}

func (u UUID) String() string {
	return string(u)
}

func (u UUID) Validate() error {
	if u == "" {
		return ErrUUIDRequired
	}
	if !uuidPattern.MatchString(string(u)) {
		return ErrInvalidUUID
	}
	return nil
}

func ValidateUUIDString(value string) error {
	return NewUUID(value).Validate()
}

func NewDeterministicUUID(name string) UUID {
	sum := sha1.Sum([]byte(name))
	sum[6] = (sum[6] & 0x0f) | 0x50
	sum[8] = (sum[8] & 0x3f) | 0x80
	return UUID(fmt.Sprintf("%x-%x-%x-%x-%x", sum[0:4], sum[4:6], sum[6:8], sum[8:10], sum[10:16]))
}
