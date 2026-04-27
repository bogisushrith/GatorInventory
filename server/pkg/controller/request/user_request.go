package request

import "ims-intro/pkg/service/dto"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role"`
}

func (request *SignUpRequest) ToDtoModel() dto.UserCreate {
	return dto.UserCreate{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
		Role:     request.Role,
	}
}
