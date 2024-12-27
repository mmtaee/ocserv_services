package initialize

type User struct {
	Username string `json:"username" validate:"required,min=2,max=16" example:"john_doe" `
	Password string `json:"password" validate:"required,min=2,max=16" example:"doe123456"`
}
