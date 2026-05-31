package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler/dto"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
	"github.com/google/uuid"
)

func (h *Handler) EnrollmentPostHandler(w http.ResponseWriter, r *http.Request) {
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

	err = uuid.Validate(req.PlanID)
	if err != nil {
		shared.SendJSONError(
			w,
			shared.ErrorResponse{
				Success: false,
				Message: "Invalid enrollment request.",
			},
			http.StatusBadRequest,
		)
		return
	}

	ctx := r.Context()
	claims, ok := ctx.Value("claims").(*service.Claims)
	if !ok {
		shared.SendJSONError(
			w, shared.ErrorResponse{
				Success: false,
				Message: "Invalid or expired access token.",
			}, http.StatusUnauthorized,
		)
		return
	}

	enrollData, err := h.enrollmentService.EnrollPlan(
		ctx, service.EnrollPlanInput{
			UserID: claims.Subject,
			PlanID: req.PlanID,
		},
	)

	if err != nil {
		if errors.Is(err, service.InvalidUserOrPlan) {
			shared.SendJSONError(
				w,
				shared.ErrorResponse{
					Success: false,
					Message: "Plan not found.",
				},
				http.StatusNotFound,
			)
			return
		}

		if errors.Is(err, service.ErrUserHasActivePlan) {
			shared.SendJSONError(
				w,
				shared.ErrorResponse{
					Success: false,
					Message: "User already has an active enrollment.",
				},
				http.StatusConflict,
			)
			return
		}

		shared.SendJSONError(
			w,
			shared.ErrorResponse{Message: "Internal Server Error."},
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(
		&dto.EnrollPlanResponse{
			Success: true,
			Data: &dto.EnrollPlanData{
				ID:            enrollData.ID,
				PlanID:        enrollData.PlanID,
				EffectiveDate: enrollData.EffectiveDate,
				EndDate:       enrollData.EndDate,
			},
		},
	)
}
