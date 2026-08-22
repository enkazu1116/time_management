package model

import "testing"

func TestUserValidateRequiresDomainInvariants(t *testing.T) {
	user := NewUser("00000000-0000-4000-8000-000000000101", "password", "user", "user@example.com", "08:00", "22:00", NewRole("00000000-0000-4000-8000-000000000201", "guest"))

	if err := user.Validate(); err != nil {
		t.Fatalf("expected valid user: %v", err)
	}

	user.Role = Role{}
	if err := user.Validate(); err != ErrRoleRequired {
		t.Fatalf("expected role error, got %v", err)
	}
}

func TestUserValidateRequiresValidEmailDomain(t *testing.T) {
	role := NewRole("00000000-0000-4000-8000-000000000201", "guest")
	tests := []struct {
		name    string
		email   string
		wantErr error
	}{
		{name: "domain without dot", email: "user@example", wantErr: ErrInvalidEmailDomain},
		{name: "empty domain label", email: "user@example..com", wantErr: ErrInvalidEmail},
		{name: "domain starts with hyphen", email: "user@-example.com", wantErr: ErrInvalidEmailDomain},
		{name: "domain ends with hyphen", email: "user@example-.com", wantErr: ErrInvalidEmailDomain},
		{name: "numeric top level domain", email: "user@example.123", wantErr: ErrInvalidEmailDomain},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewUser("00000000-0000-4000-8000-000000000101", "password", "user", tt.email, "08:00", "22:00", role)

			if err := user.Validate(); err != tt.wantErr {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}
