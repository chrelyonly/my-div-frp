package consts

var (
	// proxy status
	Idle    = "idle"
	Working = "working"
	Closed  = "closed"
	Online  = "online"
	Offline = "offline"

	// proxy type
	TCPProxy    = "tcp"
	UDPProxy    = "udp"
	TCPMuxProxy = "tcpmux"
	HTTPProxy   = "http"
	HTTPSProxy  = "https"
	STCPProxy   = "stcp"
	XTCPProxy   = "xtcp"
	SUDPProxy   = "sudp"

	// authentication method
	TokenAuthMethod = "token"
	OidcAuthMethod  = "oidc"

	// TCP multiplexer
	HTTPConnectTCPMultiplexer = "httpconnect"
)
