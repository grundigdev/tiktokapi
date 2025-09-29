package models

import "time"

type TokenModel struct {
	BaseModel
	AccessToken  string    `json:"access_token"` // Foreign key to AbsenceCategoryModel
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (TokenModel) TableName() string {
	return "tokens"
}
