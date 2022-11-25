package entity

// FrpConfig R 转换结构
type FrpConfig struct {
	ServerAddr string    `json:"server_addr"`
	ServerPort string    `json:"server_port"`
	Token      string    `json:"token"`
	ListMap    []ListMap `json:"listMap"`
}

type ListMap struct {
	LocalIp       string `json:"local_ip"`
	LocalPort     string `json:"local_port"`
	RemotePort    string `json:"remote_port"`
	CustomDomains string `json:"custom_domains"`
	Comment       string `json:"comment"`
	Type          string `json:"type"`
}
