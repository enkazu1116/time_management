package parser

import (
	"errors"
	"time"

	"time_management/infrastructure/util/messages"
)

func ParseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, errors.New(messages.DateRequired)
	}
	return time.Parse("2006-01-02", value)
}

func ParseYearMonth(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, errors.New(messages.YearMonthRequired)
	}
	if parsed, err := time.Parse("2006-01", value); err == nil {
		return parsed, nil
	}
	return time.Parse("2006-01-02", value)
}

func ParseDateTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, errors.New(messages.DateTimeRequired)
	}
	return time.Parse(time.RFC3339, value)
}
