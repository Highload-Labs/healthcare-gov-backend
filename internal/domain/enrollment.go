package domain

import "time"

type Enrollment struct {
	ID            string    `json:"id" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	PlanID        string    `json:"plan_id" db:"plan_id"`
	EffectiveDate time.Time `json:"effective_date" db:"effective_date"`
	EndDate       time.Time `json:"end_date" db:"end_date"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
