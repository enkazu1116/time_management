package requests

import "testing"

func TestCreateWorkLogRequestValidateRequiresHTTPFields(t *testing.T) {
	request := CreateWorkLogRequest{
		UserID:       "00000000-0000-4000-8000-000000000101",
		WorkTypeID:   "00000000-0000-4000-8000-000000000301",
		WorkTypeName: "study",
		WorkDate:     "2026-08-20",
		StartedAt:    "2026-08-20T08:00:00Z",
		EndedAt:      "2026-08-20T09:00:00Z",
	}

	if err := request.Validate(); err != nil {
		t.Fatalf("expected valid request: %v", err)
	}

	request.StartedAt = "invalid"
	if err := request.Validate(); err == nil {
		t.Fatal("expected invalid datetime error")
	}
}

func TestCreateDailyTaskRequestValidateRejectsOutOfRangeProgressRate(t *testing.T) {
	request := CreateDailyTaskRequest{
		WeeklyTaskID: "00000000-0000-4000-8000-000000000701",
		UserID:       "00000000-0000-4000-8000-000000000101",
		WorkDate:     "2026-08-20",
		ProgressRate: 1.1,
	}

	if err := request.Validate(); err != ErrProgressRateOutOfRange {
		t.Fatalf("expected progress rate error, got %v", err)
	}
}

func TestCreateMonthlySummaryRequestValidateRejectsUnsupportedSource(t *testing.T) {
	request := CreateMonthlySummaryRequest{
		UserID:    "00000000-0000-4000-8000-000000000101",
		YearMonth: "2026-08",
		Source:    "daily_logs",
	}

	if err := request.Validate(); err != ErrUnsupportedMonthlySource {
		t.Fatalf("expected source error, got %v", err)
	}
}
