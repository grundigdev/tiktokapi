package models

import (
	"time"

	"github.com/google/uuid"
)

type FileModel struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	FilePathVideo   string    `json:"filepath_video"`
	FilePathContext string    `json:"filepath_context"`
	Status          string    `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (FileModel) TableName() string {
	return "file"
}
