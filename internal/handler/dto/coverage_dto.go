package dto

import (
	"errors"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
)

type CoverageRequest struct {
	Zipcode string `json:"zipcode"`
}

type CoverageResponse struct {
	Success bool                  `json:"success"`
	Data    *CoverageResponseData `json:"data"`
}

type CoverageResponseData struct {
	Zipcode   string `json:"zipcode"`
	State     string `json:"state"`
	Supported bool   `json:"supported"`
}

func (req *CoverageRequest) Validate() error {
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
