package transport

import (
	"fmt"
	"log"
	"net/http"
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

func NewHTTP(authRegisterSvc service.AuthRegisterService, authLoginSvc service.AuthLoginService) *HTTP {
	mux := http.NewServeMux()

	cfg := config.GetConfig()

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
