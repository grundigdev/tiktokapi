package services

import (
	"errors"
	"time"

	"github.com/grundigdev/club/models"
	"github.com/grundigdev/club/requests"
	"gorm.io/gorm"
)

type TokenService struct {
	DB *gorm.DB
}

func NewTokenService(db *gorm.DB) *TokenService {
	return &TokenService{DB: db}
}

func (c TokenService) CreateToken(data *requests.CreateTokenRequest) (*models.TokenModel, error) {
	tokenCreated := &models.TokenModel{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiresAt:    data.ExpiresAt,
	}
	result := c.DB.FirstOrCreate(tokenCreated, models.TokenModel{AccessToken: data.AccessToken})
	if result.Error != nil {
		return nil, result.Error
	}
	return tokenCreated, nil
}

func (c TokenService) GetToken(accessToken string) (*models.TokenModel, bool, error) {
	var token models.TokenModel

	// Filter by access token, order by created_at descending and get the first record
	result := c.DB.Where("access_token = ?", accessToken).First(&token)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, false, errors.New("token not found")
		}
		return nil, false, result.Error
	}

	// Check if token has expired
	if token.ExpiresAt.Before(time.Now()) {
		return nil, false, errors.New("token has expired")
	}

	return &token, true, nil
}
