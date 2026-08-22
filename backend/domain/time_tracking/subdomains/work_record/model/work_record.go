package model

import (
	"errors"
	"time"

	"time_management/domain/shared/messages"
	sharedmodel "time_management/domain/shared/model"
)

var (
	ErrWorkTypeRequired       = errors.New(messages.WorkTypeRequired)
	ErrWorkLogIDRequired      = errors.New(messages.WorkLogIDRequired)
	ErrBreakLogIDRequired     = errors.New(messages.BreakLogIDRequired)
	ErrUserIDRequired         = errors.New(messages.UserIDRequired)
	ErrWorkDateRequired       = errors.New(messages.WorkDateRequired)
	ErrStartedAtRequired      = errors.New(messages.StartedAtRequired)
	ErrEndedAtRequired        = errors.New(messages.EndedAtRequired)
	ErrStartedAtAfterEndedAt  = errors.New(messages.StartedAtAfterEndedAt)
	ErrNegativeWorkDuration   = errors.New(messages.NegativeWorkDuration)
	ErrNegativeBreakDuration  = errors.New(messages.NegativeBreakDuration)
	ErrNegativeResumeDuration = errors.New(messages.NegativeResumeDuration)
)

type WorkType struct {
	WorkTypeID   string
	WorkTypeName string
}

type WorkLog struct {
	WorkLogID    string
	UserID       string
	WorkType     WorkType
	WorkDate     time.Time
	StartedAt    time.Time
	EndedAt      time.Time
	WorkDuration time.Duration
}

type BreakLog struct {
	BreakLogID     string
	WorkLogID      string
	UserID         string
	StartedAt      time.Time
	EndedAt        time.Time
	BreakDuration  time.Duration
	ResumedAt      *time.Time
	ResumeDuration *time.Duration
}

func NewWorkType(workTypeID string, workTypeName string) WorkType {
	return WorkType{WorkTypeID: workTypeID, WorkTypeName: workTypeName}
}

func NewWorkLog(workLogID string, userID string, workType WorkType, workDate time.Time, startedAt time.Time, endedAt time.Time, workDuration time.Duration) WorkLog {
	if workDuration == 0 {
		workDuration = endedAt.Sub(startedAt)
	}
	return WorkLog{WorkLogID: workLogID, UserID: userID, WorkType: workType, WorkDate: workDate, StartedAt: startedAt, EndedAt: endedAt, WorkDuration: workDuration}
}

func NewBreakLog(breakLogID string, workLogID string, userID string, startedAt time.Time, endedAt time.Time, breakDuration time.Duration, resumedAt *time.Time, resumeDuration *time.Duration) BreakLog {
	if breakDuration == 0 {
		breakDuration = endedAt.Sub(startedAt)
	}
	return BreakLog{BreakLogID: breakLogID, WorkLogID: workLogID, UserID: userID, StartedAt: startedAt, EndedAt: endedAt, BreakDuration: breakDuration, ResumedAt: resumedAt, ResumeDuration: resumeDuration}
}

func (w WorkLog) Validate() error {
	if w.WorkLogID == "" {
		return ErrWorkLogIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(w.WorkLogID); err != nil {
		return err
	}
	if w.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(w.UserID); err != nil {
		return err
	}
	if w.WorkType.WorkTypeID == "" || w.WorkType.WorkTypeName == "" {
		return ErrWorkTypeRequired
	}
	if err := sharedmodel.ValidateUUIDString(w.WorkType.WorkTypeID); err != nil {
		return err
	}
	if w.WorkDate.IsZero() {
		return ErrWorkDateRequired
	}
	if w.StartedAt.IsZero() {
		return ErrStartedAtRequired
	}
	if w.EndedAt.IsZero() {
		return ErrEndedAtRequired
	}
	if !w.StartedAt.Before(w.EndedAt) {
		return ErrStartedAtAfterEndedAt
	}
	if w.WorkDuration < 0 {
		return ErrNegativeWorkDuration
	}
	return nil
}

func (b BreakLog) Validate() error {
	if b.BreakLogID == "" {
		return ErrBreakLogIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(b.BreakLogID); err != nil {
		return err
	}
	if b.WorkLogID == "" {
		return ErrWorkLogIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(b.WorkLogID); err != nil {
		return err
	}
	if b.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(b.UserID); err != nil {
		return err
	}
	if b.StartedAt.IsZero() {
		return ErrStartedAtRequired
	}
	if b.EndedAt.IsZero() {
		return ErrEndedAtRequired
	}
	if !b.StartedAt.Before(b.EndedAt) {
		return ErrStartedAtAfterEndedAt
	}
	if b.BreakDuration < 0 {
		return ErrNegativeBreakDuration
	}
	if b.ResumeDuration != nil && *b.ResumeDuration < 0 {
		return ErrNegativeResumeDuration
	}
	return nil
}
