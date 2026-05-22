package transport

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/transport/middleware"
)

type HTTP struct {
	srv *http.Server
}

func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func setupPprof() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := &http.Server{
		Addr:              "127.0.0.1:6060",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	ln, err := net.Listen("tcp4", srv.Addr)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		_ = srv.Serve(ln)
	}()
}

func NewHTTP(authRegisterSvc service.AuthRegisterService, authLoginSvc service.AuthLoginService) *HTTP {
	mux := http.NewServeMux()

	cfg := config.GetConfig()

	if cfg.GoEnv == "development" {
		setupPprof()
	}

	h := handler.NewHandler(mux, cfg, authRegisterSvc, authLoginSvc)
	h.InitializeRoutes()

	wrappedMux := chain(
		mux,
		middleware.RecoveryMiddleware,
		middleware.RequestIDMiddleware,
		middleware.LoggingMiddleware,
	)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", config.GetConfig().ServerPort),
		Handler:           wrappedMux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	return &HTTP{
		srv: srv,
	}
}

func (h *HTTP) Serve() {
	log.Println("server running on :8080")

	if err := h.srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
