package myConfig

// 是否网络模式
var IsNet bool

// frp配置信息
var ServerAddr string
var ServerPort int
var Token string
var Comment string
var FrpType string
var LocalIp string
var LocalPort int
var RemotePort int
var CustomDomains string

// R 转换结构
type FrpConfig struct {
	ServerAddr string `json:"server_addr"`
	ServerPort int    `json:"server_port"`
	Token      string `json:"token"`
	Comment    string `json:"comment"`
	FrpType    string `json:"frp_type"`
	LocalIp    string `json:"local_ip"`
	LocalPort  int    `json:"local_port"`
	RemotePort string `json:"remote_port"`
}
