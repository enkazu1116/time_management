package handlers

import (
	"net/http"

	userports "time_management/domain/user/ports"
	userusecases "time_management/domain/user/usecases"
	"time_management/infrastructure/http/requests"
	"time_management/infrastructure/http/responses"
	"time_management/infrastructure/util/decoder"
	"time_management/infrastructure/util/id"
	"time_management/infrastructure/util/response"
)

type UserHandler struct {
	userUsecase userusecases.UserUsecase
}

func NewUserHandler(userUsecase userusecases.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request requests.RegisterUserRequest
	if err := decoder.DecodeJSON(r, &request); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := request.Validate(); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if request.UserID == "" {
		request.UserID = id.NewID()
	}
	if request.RoleName == "" {
		request.RoleName = "guest"
	}
	if request.RoleID == "" {
		request.RoleID = defaultRoleID(request.RoleName)
	}

	user, err := h.userUsecase.RegisterUser(userports.RegisterUserInput{
		UserID:       request.UserID,
		Password:     request.Password,
		Name:         request.Name,
		Email:        request.Email,
		RegularStart: request.RegularStart,
		RegularEnd:   request.RegularEnd,
		RoleID:       request.RoleID,
		RoleName:     request.RoleName,
	})
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	response.WriteJSON(w, http.StatusCreated, responses.User(user))
}

func defaultRoleID(roleName string) string {
	switch roleName {
	case "admin":
		return "00000000-0000-4000-8000-000000000001"
	default:
		return "00000000-0000-4000-8000-000000000002"
	}
}
