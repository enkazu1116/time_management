package model

import "testing"

func TestUUIDValidateRequiresValidUUID(t *testing.T) {
	if err := NewUUID("00000000-0000-4000-8000-000000000001").Validate(); err != nil {
		t.Fatalf("expected valid uuid: %v", err)
	}

	if err := NewUUID("user-1").Validate(); err != ErrInvalidUUID {
		t.Fatalf("expected invalid uuid error, got %v", err)
	}
}

func TestNewDeterministicUUIDReturnsValidStableUUID(t *testing.T) {
	first := NewDeterministicUUID("weekly-log-task")
	second := NewDeterministicUUID("weekly-log-task")

	if first != second {
		t.Fatalf("expected stable uuid, got %s and %s", first, second)
	}
	if err := first.Validate(); err != nil {
		t.Fatalf("expected generated uuid to be valid: %v", err)
	}
}
