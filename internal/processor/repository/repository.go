package procrepository

import (
	"apollo-image-processor/internal/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProcessorRepository interface {
	GetImage(string) ([]byte, error)
	InsertImage(string, string, []byte) error
	UpdateBatchStatus(string, models.BatchStatus) error
	UpdateImageStatus(string, models.ImageStatus) error
}

type processorRepository struct {
	db *sql.DB
}

func NewProcessorRepository(db *sql.DB) ProcessorRepository {
	return &processorRepository{
		db: db,
	}
}

// TODO update the status of the batch and image when getting the image instead of calling a separate db call
func (p *processorRepository) GetImage(imageID string) ([]byte, error) {

	imageUUID, err := uuid.Parse(imageID)
	if err != nil {
		return nil, fmt.Errorf("error parsing image uuid: %w", err)
	}

	query := fmt.Sprintf("SELECT image FROM images WHERE image_id = '%s'", imageUUID)

	var srcimage []byte
	err = p.db.QueryRow(query).Scan(&srcimage)
	if err != nil {
		return nil, fmt.Errorf("error fetching image %s from db: %w", imageID, err)
	}

	return srcimage, nil
}

func (p *processorRepository) InsertImage(imageID string, batchID string, procImage []byte) error {

	imageUUID, err := uuid.Parse(imageID)
	if err != nil {
		return fmt.Errorf("error parsing image uuid: %w", err)
	}

	status := models.ImageStatusCompleted
	query := fmt.Sprintf("UPDATE images SET (image_proc_bw, status, processed_at) = ($1, $2, $3) WHERE image_id = '%s'", imageUUID)
	_, err = p.db.Exec(query, procImage, status, time.Now())
	if err != nil {
		return fmt.Errorf("error inserting updated image %s: %w", imageID, err)
	}

	return nil
}

func (p *processorRepository) UpdateBatchStatus(batchID string, status models.BatchStatus) error {

	batchUUID, err := uuid.Parse(batchID)
	if err != nil {
		return fmt.Errorf("error parsing batch uuid: %w", err)
	}

	var currentStatus models.BatchStatus
	query1 := fmt.Sprintf("SELECT status FROM batches WHERE batch_id = '%s'", batchUUID)
	err = p.db.QueryRow(query1).Scan(&currentStatus)
	if err != nil {
		return fmt.Errorf("error fetching status for batch %s: %w", batchUUID, err)
	}

	if currentStatus == models.BatchStatusFailed {
		return fmt.Errorf("error updating status for batch %s: status is failed", batchUUID)
	}

	query2 := fmt.Sprintf("UPDATE batches SET status = '%s' WHERE batch_id = '%s'", status, batchUUID)
	_, err = p.db.Exec(query2)
	if err != nil {
		return fmt.Errorf("error updating status for batch %s: %w", batchUUID, err)
	}

	return nil
}

func (p *processorRepository) UpdateImageStatus(imageID string, status models.ImageStatus) error {

	imageUUID, err := uuid.Parse(imageID)
	if err != nil {
		return fmt.Errorf("error parsing image uuid: %w", err)
	}

	query := fmt.Sprintf("UPDATE images SET status = '%s' WHERE image_id = '%s'", status, imageUUID)
	_, err = p.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error updating status for image %s: %w", imageUUID, err)
	}

	return nil
}
