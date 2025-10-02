package requests

type CreateUploadRequest struct {
	Title          string `json:"title" validate:"required"`
	PrivacyLevel   string `json:"privacy_level" validate:"required"`
	FilePath       string `json:"file_path" validate:"required"`
	FileSize       int64  `json:"file_size" validate:"required"`
	CoverTimestamp int    `json:"cover_timestamp" validate:"required"`
}
