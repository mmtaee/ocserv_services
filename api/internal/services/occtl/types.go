package occtl

type DisconnectRequest struct {
	Username string `json:"username"`
}

type UnBanIPRequest struct {
	IP string `json:"ip"`
}

type ShowStatusResponse struct {
	Status string `json:"status"`
}
