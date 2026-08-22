package model

import (
	"testing"
	"time"
)

func TestWorkLogValidateRequiresDomainInvariants(t *testing.T) {
	workLog := NewWorkLog("00000000-0000-4000-8000-000000000501", "00000000-0000-4000-8000-000000000101", NewWorkType("00000000-0000-4000-8000-000000000301", "study"), day(), at(8), at(9), 0)

	if err := workLog.Validate(); err != nil {
		t.Fatalf("expected valid work log: %v", err)
	}

	workLog.EndedAt = at(8)
	if err := workLog.Validate(); err != ErrStartedAtAfterEndedAt {
		t.Fatalf("expected time order error, got %v", err)
	}
}

func TestBreakLogValidateRequiresWorkLogID(t *testing.T) {
	breakLog := NewBreakLog("00000000-0000-4000-8000-000000000601", "", "00000000-0000-4000-8000-000000000101", at(9), at(10), 0, nil, nil)

	if err := breakLog.Validate(); err != ErrWorkLogIDRequired {
		t.Fatalf("expected work log id error, got %v", err)
	}
}

func at(hour int) time.Time {
	return time.Date(2026, 8, 20, hour, 0, 0, 0, time.UTC)
}

func day() time.Time {
	return time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC)
}
