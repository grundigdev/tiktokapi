package services

import (
	"github.com/grundigdev/club/models"
	"github.com/grundigdev/club/requests"
	"gorm.io/gorm"
)

type UploadService struct {
	DB *gorm.DB
}

func NewUploadService(db *gorm.DB) *UploadService {
	return &UploadService{DB: db}
}

func (c UploadService) CreateUpload(data *requests.CreateUploadRequest) (*models.UploadModel, error) {
	uploadCreated := &models.UploadModel{
		Title:          data.Title,
		PrivacyLevel:   data.PrivacyLevel,
		FilePath:       data.FilePath,
		FileSize:       data.FileSize,
		CoverTimestamp: data.CoverTimestamp,
	}
	result := c.DB.FirstOrCreate(uploadCreated, models.UploadModel{FilePath: data.FilePath})
	if result.Error != nil {
		return nil, result.Error
	}
	return uploadCreated, nil
}

func (c UploadService) GetAllUploads() ([]models.UploadModel, error) {
	var uploads []models.UploadModel
	result := c.DB.Find(&uploads)

	if result.Error != nil {
		return nil, result.Error
	}

	return uploads, nil
}
