package ocserv

type OcctlUser struct {
	Username        string   `json:"Username"`
	Hostname        string   `json:"Hostname"`
	Device          string   `json:"Device"`
	RemoteIP        string   `json:"Remote IP"`
	UserAgent       string   `json:"User-Agent"`
	Since           string   `json:"_Connected at"`
	ConnectedAt     string   `json:"Connected at"`
	AverageRX       string   `json:"Average RX"`
	AverageTX       string   `json:"Average TX"`
	GroupName       string   `json:"Groupname"`
	State           string   `json:"State"`
	Vhost           string   `json:"vhost"`
	MTU             string   `json:"MTU"`
	Location        string   `json:"Location"`
	LocalDeviceIP   string   `json:"Local Device IP"`
	IPv4            string   `json:"IPv4"`
	PTPIPv4         string   `json:"P-t-P IPv4"`
	RX              string   `json:"RX"`
	TX              string   `json:"TX"`
	RXConvert       string   `json:"_RX"`
	TXConvert       string   `json:"_TX"`
	DPD             string   `json:"DPD"`
	KeepAlive       string   `json:"KeepAlive"`
	RawConnectedAt  int64    `json:"raw_connected_at"`
	FullSession     string   `json:"Full session"`
	Session         string   `json:"Session"`
	TLSCipherSuite  string   `json:"TLS ciphersuite"`
	DNS             []string `json:"DNS"`
	NBNS            []string `json:"NBNS"`
	SplitDNSDomains []string `json:"Split-DNS-Domains"`
	IRoutes         []string `json:"iRoutes"`
}

type IPBan struct {
	IP    string `json:"IP"`
	Since string `json:"Since"`
	Until string `json:"_Since"`
	Score int    `json:"Score"`
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

type OcGroupConfig struct {
	RxDataPerSec         *string   `json:"rx-data-per-sec"`
	TxDataPerSec         *string   `json:"tx-data-per-sec"`
	MaxSameClients       *int      `json:"max-same-clients"`
	IPv4Network          *string   `json:"ipv4-network"`
	DNS                  *[]string `json:"dns"`
	NoUDP                *bool     `json:"no-udp"`
	KeepAlive            *int      `json:"keepalive"`
	DPD                  *int      `json:"dpd"`
	MobileDPD            *int      `json:"mobile-dpd"`
	TunnelAllDNS         *bool     `json:"tunnel-all-dns"`
	RestrictUserToRoutes *bool     `json:"restrict-user-to-routes"`
	StatsReportTime      *int      `json:"stats-report-time"`
	MTU                  *int      `json:"mtu"`
	IdleTimeout          *int      `json:"idle-timeout"`
	MobileIdleTimeout    *int      `json:"mobile-idle-timeout"`
	SessionTimeout       *int      `json:"session-timeout"`
}

type OcGroupConfigInfo struct {
	Name   string         `json:"name"`
	Path   string         `json:"-"`
	Config *OcGroupConfig `json:"config"`
}
