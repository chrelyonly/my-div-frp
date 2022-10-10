package plugin

import (
	"io"
	"net"
	"net/http"

	"github.com/gorilla/mux"

	frpNet "github.com/fatedier/frp/pkg/util/net"
)

const PluginStaticFile = "static_file"

func init() {
	Register(PluginStaticFile, NewStaticFilePlugin)
}

type StaticFilePlugin struct {
	localPath   string
	stripPrefix string
	httpUser    string
	httpPasswd  string

	l *Listener
	s *http.Server
}

func NewStaticFilePlugin(params map[string]string) (Plugin, error) {
	localPath := params["plugin_local_path"]
	stripPrefix := params["plugin_strip_prefix"]
	httpUser := params["plugin_http_user"]
	httpPasswd := params["plugin_http_passwd"]

	listener := NewProxyListener()

	sp := &StaticFilePlugin{
		localPath:   localPath,
		stripPrefix: stripPrefix,
		httpUser:    httpUser,
		httpPasswd:  httpPasswd,

		l: listener,
	}
	var prefix string
	if stripPrefix != "" {
		prefix = "/" + stripPrefix + "/"
	} else {
		prefix = "/"
	}

	router := mux.NewRouter()
	router.Use(frpNet.NewHTTPAuthMiddleware(httpUser, httpPasswd).Middleware)
	router.PathPrefix(prefix).Handler(frpNet.MakeHTTPGzipHandler(http.StripPrefix(prefix, http.FileServer(http.Dir(localPath))))).Methods("GET")
	sp.s = &http.Server{
		Handler: router,
	}
	go func() {
		_ = sp.s.Serve(listener)
	}()
	return sp, nil
}

func (sp *StaticFilePlugin) Handle(conn io.ReadWriteCloser, realConn net.Conn, extraBufToLocal []byte) {
	wrapConn := frpNet.WrapReadWriteCloserToConn(conn, realConn)
	_ = sp.l.PutConn(wrapConn)
}

func (sp *StaticFilePlugin) Name() string {
	return PluginStaticFile
}

func (sp *StaticFilePlugin) Close() error {
	sp.s.Close()
	sp.l.Close()
	return nil
}
