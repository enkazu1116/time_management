package ports

import "testing"

func TestRegisterUserInputValidateRequiresUsecaseFields(t *testing.T) {
	input := RegisterUserInput{
		UserID:   "00000000-0000-4000-8000-000000000101",
		Name:     "user",
		RoleID:   "00000000-0000-4000-8000-000000000201",
		RoleName: "guest",
	}

	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid input: %v", err)
	}

	input.UserID = ""
	if err := input.Validate(); err != ErrUserIDRequired {
		t.Fatalf("expected user id error, got %v", err)
	}
}
