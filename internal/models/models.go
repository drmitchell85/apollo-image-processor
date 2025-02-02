package models

import (
	"time"

	"github.com/google/uuid"
)

type UploadedFile struct {
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	Filename    string `json:"filename"`
	// FileContent string `json:"file_content"`
	FileContent []byte `json:"file_content"`
	Key         string `json:"key"`
	ImageNum    int    `json:"image_num"`
}

type BatchMessage struct {
	Batchid string
	Imageid string
}

type BatchStatus string

const (
	BatchStatusCreated    BatchStatus = "created"
	BatchStatusProcessing BatchStatus = "processing"
	BatchStatusCompleted  BatchStatus = "completed"
	BatchStatusFailed     BatchStatus = "failed"
)

type Batch struct {
	BatchID         uuid.UUID   `json:"batch_id" db:"batch_id"`
	Status          BatchStatus `json:"status" db:"status"`
	CreatedAt       time.Time   `json:"created_at" db:"created_at"`
	CompletedAt     *time.Time  `json:"completed_at,omitempty" db:"completed_at"`
	TotalImages     int32       `json:"total_images" db:"total_images"`
	ProcessedImages int32       `json:"processed_images" db:"processed_images"`
}

type ImageStatus string

const (
	ImageStatusPending    ImageStatus = "pending"
	ImageStatusProcessing ImageStatus = "processing"
	ImageStatusCompleted  ImageStatus = "completed"
	ImageStatusFailed     ImageStatus = "failed"
)

type Image struct {
	ImageID     uuid.UUID   `json:"image_id" db:"image_id"`
	BatchID     uuid.UUID   `json:"batch_id" db:"batch_id"`
	Status      ImageStatus `json:"status" db:"status"`
	Error       *string     `json:"error,omitempty" db:"error"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	ProcessedAt *time.Time  `json:"processed_at,omitempty" db:"processed_at"`
}
