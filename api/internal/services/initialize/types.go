package initialize

type CreateAdminUserRequest struct {
	Username string `json:"username" validate:"required,min=2,max=16" example:"john_doe" `
	Password string `json:"password" validate:"required,min=2,max=16" example:"doe123456"`
}

type CreateSiteConfigRequest struct {
	GoogleCaptchaSecretKey string `json:"google_captcha_secret_key" validate:"omitempty"`
	GoogleCaptchaSiteKey   string `json:"google_captcha_site_key" validate:"omitempty"`
}
