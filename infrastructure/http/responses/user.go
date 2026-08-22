package responses

import usermodel "time_management/domain/user/model"

type UserResponse struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	RoleID   string `json:"role_id"`
	RoleName string `json:"role_name"`
}

func User(v usermodel.User) UserResponse {
	return UserResponse{UserID: v.UserID, Name: v.Name, Email: v.Email, RoleID: v.Role.RoleID, RoleName: v.Role.RoleName}
}
