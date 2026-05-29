package domain

import "time"

type Coverage struct {
	Id           string    `db:"id" json:"id"`
	State        string    `db:"state" json:"state"`
	ZipcodeStart string    `db:"zipcode_start" json:"zipcode_start"`
	ZipcodeEnd   string    `db:"zipcode_end" json:"zipcode_end"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
