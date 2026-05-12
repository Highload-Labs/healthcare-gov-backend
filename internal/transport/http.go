package transport

import (
	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler"
	"log"
	"net/http"
	"time"
)

type HTTP struct {
	mux *http.ServeMux
	srv *http.Server

	handler *handler.Handler
}

func NewHTTP() *HTTP {
	mux := http.NewServeMux()

	h := handler.NewHandler(mux)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
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
