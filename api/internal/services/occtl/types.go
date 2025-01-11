package occtl

type UnBanIPRequest struct {
	IP string `json:"ip"`
}

type ShowStatusResponse struct {
	Status string `json:"status"`
}
