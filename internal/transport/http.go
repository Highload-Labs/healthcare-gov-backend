package transport

import (
	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/transport/middleware"
	"log"
	"net/http"
	"time"
)

type HTTP struct {
	mux *http.ServeMux
	srv *http.Server

	handler *handler.Handler
}

func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func NewHTTP() *HTTP {
	mux := http.NewServeMux()

	wrappedMux := chain(
		mux,
		middleware.RecoveryMiddleware,
		middleware.RequestIDMiddleware,
		middleware.LoggingMiddleware,
	)

	h := handler.NewHandler(mux)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           wrappedMux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	return &HTTP{
		mux:     mux,
		srv:     srv,
		handler: h,
	}
}

func (h *HTTP) SetupAndServe() {
	h.handler.InitializeRoutes()

	log.Println("server running on :8080")
	if err := h.srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
