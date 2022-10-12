package client

import (
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gorilla/mux"

	"github.com/fatedier/frp/assets"
	frpNet "github.com/fatedier/frp/pkg/util/net"
)

var (
	httpServerReadTimeout  = 60 * time.Second
	httpServerWriteTimeout = 60 * time.Second
)

func (svr *Service) RunAdminServer(address string) (err error) {
	// url router
	router := mux.NewRouter()

	router.HandleFunc("/healthz", svr.healthz)

	// debug
	if svr.cfg.PprofEnable {
		router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		router.HandleFunc("/debug/pprof/trace", pprof.Trace)
		router.PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)
	}

	subRouter := router.NewRoute().Subrouter()
	user, passwd := svr.cfg.AdminUser, svr.cfg.AdminPwd
	subRouter.Use(frpNet.NewHTTPAuthMiddleware(user, passwd).Middleware)

	// api, see admin_api.go
	subRouter.HandleFunc("/api/reload", svr.apiReload).Methods("GET")
	subRouter.HandleFunc("/api/status", svr.apiStatus).Methods("GET")
	subRouter.HandleFunc("/api/config", svr.apiGetConfig).Methods("GET")
	subRouter.HandleFunc("/api/config", svr.apiPutConfig).Methods("PUT")

	// view
	subRouter.Handle("/favicon.ico", http.FileServer(assets.FileSystem)).Methods("GET")
	subRouter.PathPrefix("/static/").Handler(frpNet.MakeHTTPGzipHandler(http.StripPrefix("/static/", http.FileServer(assets.FileSystem)))).Methods("GET")
	subRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/", http.StatusMovedPermanently)
	})

	server := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  httpServerReadTimeout,
		WriteTimeout: httpServerWriteTimeout,
	}
	if address == "" {
		address = ":http"
	}
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	go func() {
		_ = server.Serve(ln)
	}()
	return
}
