package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	usermodel "time_management/domain/user/model"
	userports "time_management/domain/user/ports"
)

func TestUserHandlerReturnsBadRequestWhenRequestValidationFails(t *testing.T) {
	handler := NewUserHandler(stubUserUsecase{})
	request := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"user","password":"password"}`))
	response := httptest.NewRecorder()

	handler.RegisterUser(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestTimeTrackingHandlerReturnsBadRequestWhenRequestValidationFails(t *testing.T) {
	handler := NewTimeTrackingHandler(nil, nil, nil, nil, nil, nil, nil)
	request := httptest.NewRequest(http.MethodPost, "/work-logs", strings.NewReader(`{"user_id":"00000000-0000-4000-8000-000000000101"}`))
	response := httptest.NewRecorder()

	handler.CreateWorkLog(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

type stubUserUsecase struct{}

func (s stubUserUsecase) RegisterUser(input userports.RegisterUserInput) (usermodel.User, error) {
	return usermodel.User{}, nil
}
