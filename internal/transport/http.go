package transport

import (
	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/transport/middleware"
	"log"
	"net/http"
	"time"
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

func NewHTTP(authRegisterSvc service.AuthRegisterService) *HTTP {
	mux := http.NewServeMux()

	h := handler.NewHandler(mux, authRegisterSvc)
	h.InitializeRoutes()

	wrappedMux := chain(
		mux,
		middleware.RecoveryMiddleware,
		middleware.RequestIDMiddleware,
		middleware.LoggingMiddleware,
	)

	srv := &http.Server{
		Addr:              ":8080",
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
