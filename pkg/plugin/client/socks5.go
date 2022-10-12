package plugin

import (
	"io"
	"log"
	"net"

	gosocks5 "github.com/armon/go-socks5"

	frpNet "github.com/fatedier/frp/pkg/util/net"
)

const PluginSocks5 = "socks5"

func init() {
	Register(PluginSocks5, NewSocks5Plugin)
}

type Socks5Plugin struct {
	Server *gosocks5.Server
}

func NewSocks5Plugin(params map[string]string) (p Plugin, err error) {
	user := params["plugin_user"]
	passwd := params["plugin_passwd"]

	cfg := &gosocks5.Config{
		Logger: log.New(io.Discard, "", log.LstdFlags),
	}
	if user != "" || passwd != "" {
		cfg.Credentials = gosocks5.StaticCredentials(map[string]string{user: passwd})
	}
	sp := &Socks5Plugin{}
	sp.Server, err = gosocks5.New(cfg)
	p = sp
	return
}

func (sp *Socks5Plugin) Handle(conn io.ReadWriteCloser, realConn net.Conn, extraBufToLocal []byte) {
	defer conn.Close()
	wrapConn := frpNet.WrapReadWriteCloserToConn(conn, realConn)
	_ = sp.Server.ServeConn(wrapConn)
}

func (sp *Socks5Plugin) Name() string {
	return PluginSocks5
}

func (sp *Socks5Plugin) Close() error {
	return nil
}
