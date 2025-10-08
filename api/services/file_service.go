package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/grundigdev/club/models"
	"github.com/grundigdev/club/requests"
	"gorm.io/gorm"
)

type FileService struct {
	DB *gorm.DB
}

func NewFileService(db *gorm.DB) *FileService {
	return &FileService{DB: db}
}

func (c FileService) CreateFile(data *requests.CreateFileRequest) (*models.FileModel, error) {
	fileCreated := &models.FileModel{
		ID:              data.ID,
		FilePathVideo:   data.FilePathVideo,
		FilePathContext: data.FilePathContext,
		Status:          "CREATED",
	}

	result := c.DB.FirstOrCreate(fileCreated)
	if result.Error != nil {
		return nil, result.Error
	}
	return fileCreated, nil
}

func (c FileService) GetFile(uuid uuid.UUID) (*models.FileModel, error) {
	var file models.FileModel

	result := c.DB.Where("ID = ?", uuid).First(&file)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("file with that uuid not found")
		}
		return nil, result.Error
	}

	return &file, nil
}

func (c FileService) UpdateFile(data *requests.CreateFileRequest) (*models.FileModel, error) {
	var existingFile models.FileModel
	if err := c.DB.First(&existingFile, data.ID).Error; err != nil {
		return nil, err
	}

	existingFile.FilePathVideo = data.FilePathVideo
	existingFile.FilePathContext = data.FilePathContext
	existingFile.Status = data.Status

	if err := c.DB.Save(&existingFile).Error; err != nil {
		return nil, err
	}

	return &existingFile, nil
}
