package services

import (
	"errors"

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

func (c TokenService) GetLastToken() (*models.TokenModel, error) {
	var token models.TokenModel

	// Order by created_at descending and get the first record
	result := c.DB.Order("created_at DESC").First(&token)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found")
		}
		return nil, result.Error
	}

	return &token, nil
}
