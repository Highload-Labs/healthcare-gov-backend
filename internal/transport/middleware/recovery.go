package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	cfg := config.GetConfig()

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					stackTrace := debug.Stack()
					if cfg.GoEnv == "development" {
						slog.Error("panic recovered", "error", err, "method", r.Method, "path", r.URL.Path)
						fmt.Printf("\n--- STACK TRACE ---\n%s\n-------------------\n", stackTrace)
					} else {
						slog.Error(
							"panic recovered",
							"error", err,
							"method", r.Method,
							"path", r.URL.Path,
							"stack", string(stackTrace),
						)
					}

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		},
	)
}
