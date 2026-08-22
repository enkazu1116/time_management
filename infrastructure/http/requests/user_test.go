package requests

import "testing"

func TestRegisterUserRequestValidateRequiresHTTPFields(t *testing.T) {
	request := RegisterUserRequest{
		Name:     "user",
		Password: "password",
		Email:    "user@example.com",
	}

	if err := request.Validate(); err != nil {
		t.Fatalf("expected valid request: %v", err)
	}

	request.Email = ""
	if err := request.Validate(); err != ErrEmailRequired {
		t.Fatalf("expected email error, got %v", err)
	}
}
