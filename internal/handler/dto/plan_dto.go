package dto

import (
	"errors"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
	"github.com/google/uuid"
)

type PlansZipcodeRequest struct {
	Zipcode string `json:"zipcode"`
}

type PlanIdRequest struct {
	Id string `json:"id"`
}

type PlansResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type PlanResponseData struct {
	Id             string  `json:"id"`
	Name           string  `json:"name"`
	Provider       string  `json:"provider"`
	Tier           string  `json:"tier"`
	MonthlyPremium float64 `json:"monthly_premium"`
	Deductible     float64 `json:"deductible"`
	OutOfPocketMax float64 `json:"out_of_pocket_max"`
}

func (req *PlansZipcodeRequest) Validate() error {
	if req.Zipcode == "" {
		return errors.New("zipcode is required")
	}

	if len(req.Zipcode) < 0 || len(req.Zipcode) > 5 {
		return errors.New("zipcode is invalid")
	}

	if !shared.CheckNumericOnly(req.Zipcode) {
		return errors.New("zipcode must be number only")
	}

	return nil
}

func (req *PlanIdRequest) Validate() error {
	if req.Id == "" {
		return errors.New("id is required")
	}

	err := uuid.Validate(req.Id)
	if err != nil {
		return errors.New("invalid id format")
	}

	return nil
}
