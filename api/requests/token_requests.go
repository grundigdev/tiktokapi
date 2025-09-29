package requests

import "time"

type CreateTokenRequest struct {
	AccessToken  string    `json:"access_token" validate:"required"`
	RefreshToken string    `json:"refresh_token" validate:"required"`
	ExpiresAt    time.Time `json:"expires_at" validate:"required"`
}

type CheckTokenRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
}
