package requests

import "github.com/google/uuid"

type CreateFileRequest struct {
	ID              uuid.UUID `json:"id" validate:"required"`
	FilePathVideo   string    `json:"filepath_video" validate:"required"`
	FilePathContext string    `json:"filepath_context" validate:"required"`
	Status          string    `json:"status"`
}

type GetFileRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}
