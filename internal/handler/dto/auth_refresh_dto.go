package dto

type AuthRefreshResponseData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AuthRefreshResponse struct {
	Success bool                     `json:"success"`
	Data    *AuthRefreshResponseData `json:"data,omitempty"`
}
