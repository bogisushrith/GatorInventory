package response

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func NewErrorResponse(errorMessage string) *ErrorResponse {
	return &ErrorResponse{errorMessage}
}

type LoginResponse struct {
	Role string `json:"role"`
}

func NewLoginResponse(role string) *LoginResponse {
	return &LoginResponse{role}
}

type UserSummaryResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
