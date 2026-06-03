package domain

type RefreshSession struct {
	UserID    string `json:"user_id" db:"user_id"`
	TokenHash string `json:"token_hash" db:"token_hash"`
}
