package dto

import "time"

type EnrollPlanRequest struct {
	PlanID string `json:"plan_id"`
}

type EnrollPlanResponse struct {
	Success bool            `json:"success"`
	Data    *EnrollPlanData `json:"data"`
}

type EnrollPlanData struct {
	ID            string    `json:"id"`
	PlanID        string    `json:"plan_id"`
	EffectiveDate time.Time `json:"effective_date"`
	EndDate       time.Time `json:"end_date"`
}
