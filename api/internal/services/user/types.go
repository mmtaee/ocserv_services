package user

type CreateAdminUserRequest struct {
	Username string `json:"username" validate:"required,min=2,max=16" example:"john_doe" `
	Password string `json:"password" validate:"required,min=2,max=16" example:"doe123456"`
}

type CreateAdminUserResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Username   string `json:"username" validate:"required,min=2,max=16"`
	Password   string `json:"password" validate:"required,min=2,max=16"`
	RememberMe bool   `json:"remember_me"`
}

type LoginResponse struct {
	Token string `json:"token" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=2,max=16"`
	NewPassword string `json:"new_password" validate:"required,min=2,max=16"`
}
