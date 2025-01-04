package panel

type UpdateSiteConfigRequest struct {
	GoogleCaptchaSecretKey string `json:"google_captcha_secret_key" validate:"omitempty"`
	GoogleCaptchaSiteKey   string `json:"google_captcha_site_key" validate:"omitempty"`
}

type GetPanelConfigResponse struct {
	Init                 bool   `json:"init"`
	GoogleCaptchaSiteKey string `json:"google_captcha_site_key"`
}

type GetFullPanelConfigResponse struct {
	GoogleCaptchaSecretKey string `json:"google_captcha_secret_key"`
	GoogleCaptchaSiteKey   string `json:"google_captcha_site_key"`
}

type CreateSiteConfigRequest struct {
	GoogleCaptchaSecretKey string `json:"google_captcha_secret_key" validate:"omitempty"`
	GoogleCaptchaSiteKey   string `json:"google_captcha_site_key" validate:"omitempty"`
}
