package panel

type UpdateSiteConfigRequest struct {
	GoogleCaptchaSecretKey string `json:"google_captcha_secret_key" validate:"omitempty"`
	GoogleCaptchaSiteKey   string `json:"google_captcha_site_key" validate:"omitempty"`
}
