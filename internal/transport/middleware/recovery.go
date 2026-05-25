package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
)

func RecoveryMiddleware(goEnv string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if err := recover(); err != nil {
						stackTrace := debug.Stack()

						if goEnv == "development" {
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

						shared.SendJSONError(
							w,
							shared.ErrorResponse{Success: false, Message: "Internal Server Error."},
							http.StatusInternalServerError,
						)
					}
				}()

				next.ServeHTTP(w, r)
			},
		)
	}
}
