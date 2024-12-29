package user

type LoginRequest struct {
	Username   string `json:"username" validate:"required,min=2,max=16"`
	Password   string `json:"password" validate:"required,min=2,max=16"`
	RememberMe bool   `json:"remember_me" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=2,max=16"`
	NewPassword string `json:"new_password" validate:"required,min=2,max=16"`
}
