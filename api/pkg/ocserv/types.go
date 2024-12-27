package ocserv

type OcctlUser struct {
	Username    string `json:"Username"`
	Hostname    string `json:"Hostname"`
	Device      string `json:"Device"`
	RemoteIP    string `json:"Remote IP"`
	UserAgent   string `json:"User-Agent"`
	Since       string `json:"_Connected at"`
	ConnectedAt string `json:"Connected at"`
	AverageRX   string `json:"Average RX"`
	AverageTX   string `json:"Average TX"`
}

type IPBan struct {
	IP    string `json:"IP"`
	Since string `json:"Since"`
	Score string `json:"Score"`
}

type IRoute struct {
	ID       string `json:"ID"`
	Username string `json:"Username"`
	Vhost    string `json:"vhost"`
	Device   string `json:"Device"`
	IP       string `json:"IP"`
	IRoute   string `json:"iRoutes"`
}

type Sync struct {
	Username string `json:"username"`
	Group    string `json:"group"`
}
