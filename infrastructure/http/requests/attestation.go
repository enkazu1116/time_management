package requests

import (
	"errors"

	sharedmodel "time_management/domain/shared/model"
	"time_management/infrastructure/util/messages"
	"time_management/infrastructure/util/parser"
)

type DailyAttestationRequest struct {
	UserID string `json:"user_id"`
	Date   string `json:"date"`
}

func (request DailyAttestationRequest) Validate() error {
	if request.UserID == "" {
		return errors.New(messages.RequestUserIDRequired)
	}
	if err := sharedmodel.ValidateUUIDString(request.UserID); err != nil {
		return err
	}
	_, err := parser.ParseDate(request.Date)
	return err
}
