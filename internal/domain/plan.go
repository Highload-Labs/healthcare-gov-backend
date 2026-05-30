package domain

import "time"

type Plan struct {
	ID             string    `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Provider       string    `json:"provider" db:"provider"`
	Tier           string    `json:"tier" db:"tier"`
	MonthlyPremium float64   `json:"monthly_premium" db:"monthly_premium"`
	Deductible     float64   `json:"deductible" db:"deductible"`
	OutOfPocket    float64   `json:"out_of_pocket" db:"out_of_pocket"`
	State          string    `json:"state" db:"state"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
