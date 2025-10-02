package models

type UploadModel struct {
	BaseModel
	Title          string `json:"title"`
	PrivacyLevel   string `json:"privacy_level"`
	FilePath       string `json:"file_path"`
	FileSize       int64  `json:"file_size"`
	CoverTimestamp int    `json:"cover_timestamp"`
}

func (UploadModel) TableName() string {
	return "uploads"
}
