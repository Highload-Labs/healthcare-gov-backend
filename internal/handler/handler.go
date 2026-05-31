package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler/dto"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/transport/middleware"
)

type Handler struct {
	mux    *http.ServeMux
	config *config.Config

	authService       service.AuthService
	coverageService   service.CoverageService
	planService       service.PlanService
	enrollmentService service.EnrollmentService

	authorizationMiddleware *middleware.AuthorizationMiddleware
}

func NewHandler(
	mux *http.ServeMux,
	cfg *config.Config,
	authService service.AuthService,
	coverageService service.CoverageService,
	planService service.PlanService,
	enrollmentService service.EnrollmentService,
	authorizationMiddleware *middleware.AuthorizationMiddleware,
) *Handler {
	return &Handler{
		mux:                     mux,
		config:                  cfg,
		authService:             authService,
		coverageService:         coverageService,
		planService:             planService,
		enrollmentService:       enrollmentService,
		authorizationMiddleware: authorizationMiddleware,
	}
}

func (h *Handler) InitializeRoutes() {
	h.mux.HandleFunc("GET /healthz", h.HealthzGetHandler)
	h.mux.HandleFunc("POST /auth/register", h.AuthRegisterPostHandler)
	h.mux.HandleFunc("POST /auth/login", h.AuthLoginPostHandler)
	h.mux.HandleFunc("GET /auth/refresh", h.AuthRefreshHandler)

	h.mux.HandleFunc("GET /coverage/{zipcode}", h.CoverageGetByZipcodeHandler)

	h.mux.HandleFunc(
		"GET /plans",
		h.authorizationMiddleware.Authorization(http.HandlerFunc(h.PlansGetByZipcode)).ServeHTTP,
	)
	h.mux.HandleFunc(
		"GET /plans/{id}",
		h.authorizationMiddleware.Authorization(http.HandlerFunc(h.PlanGetById)).ServeHTTP,
	)
	h.mux.HandleFunc(
		"POST /enrollments",
		h.authorizationMiddleware.Authorization(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req dto.EnrollPlanRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				shared.SendJSONError(
					w,
					shared.ErrorResponse{
						Success: false,
						Message: "Missing required fields",
					},
					http.StatusBadRequest,
				)
				return
			}

			ctx := r.Context()
			claims := ctx.Value("claims").(*service.Claims)

			enrollData, err := h.enrollmentService.EnrollPlan(ctx, service.EnrollPlanInput{
				UserID: claims.Subject,
				PlanID: req.PlanID,
			})

			if err != nil {
				shared.SendJSONError(w, shared.ErrorResponse{Message: "Internal Server Error."}, http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(&dto.EnrollPlanResponse{
				Success: true,
				Data: &dto.EnrollPlanData{
					ID:            enrollData.ID,
					PlanID:        enrollData.PlanID,
					EffectiveDate: enrollData.EffectiveDate,
					EndDate:       enrollData.EndDate,
				},
			})
		})).ServeHTTP,
	)
}
