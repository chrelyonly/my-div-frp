package plugin

import (
	"fmt"
	"io"
	"net"

	frpIo "github.com/fatedier/golib/io"
)

const PluginUnixDomainSocket = "unix_domain_socket"

func init() {
	Register(PluginUnixDomainSocket, NewUnixDomainSocketPlugin)
}

type UnixDomainSocketPlugin struct {
	UnixAddr *net.UnixAddr
}

func NewUnixDomainSocketPlugin(params map[string]string) (p Plugin, err error) {
	unixPath, ok := params["plugin_unix_path"]
	if !ok {
		err = fmt.Errorf("plugin_unix_path not found")
		return
	}

	unixAddr, errRet := net.ResolveUnixAddr("unix", unixPath)
	if errRet != nil {
		err = errRet
		return
	}

	p = &UnixDomainSocketPlugin{
		UnixAddr: unixAddr,
	}
	return
}

func (uds *UnixDomainSocketPlugin) Handle(conn io.ReadWriteCloser, realConn net.Conn, extraBufToLocal []byte) {
	localConn, err := net.DialUnix("unix", nil, uds.UnixAddr)
	if err != nil {
		return
	}
	if len(extraBufToLocal) > 0 {
		if _, err := localConn.Write(extraBufToLocal); err != nil {
			return
		}
	}

	frpIo.Join(localConn, conn)
}

func (uds *UnixDomainSocketPlugin) Name() string {
	return PluginUnixDomainSocket
}

func (uds *UnixDomainSocketPlugin) Close() error {
	return nil
}
