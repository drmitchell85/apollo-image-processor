package procrepository

import (
	"apollo-image-processor/internal/models"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProcessorRepository interface {
	GetImage(string, string) ([]byte, error)
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

func (p *processorRepository) GetImage(imageID string, batchID string) ([]byte, error) {

	imageUUID, err := uuid.Parse(imageID)
	if err != nil {
		return nil, fmt.Errorf("error parsing image uuid: %w", err)
	}

	batchUUID, err := uuid.Parse(batchID)
	if err != nil {
		return nil, fmt.Errorf("error parsing batch uuid: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tx, err := p.db.BeginTx(ctx, nil)

	imageStatus := models.ImageStatusProcessing
	query1 := fmt.Sprintf("UPDATE images SET status = '%s' WHERE image_id = '%s' RETURNING image", imageStatus, imageUUID)

	var srcimage []byte
	err = tx.QueryRow(query1).Scan(&srcimage)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error fetching image %s from db: %w", imageID, err)
	}

	batchStatus := models.BatchStatusProcessing
	query2 := fmt.Sprintf("UPDATE batches SET status = '%s' WHERE batch_id = '%s'", batchStatus, batchUUID)
	_, err = tx.Exec(query2)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error updating status for batch %s: %w", batchUUID, err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return srcimage, nil
}

func (p *processorRepository) InsertImage(imageID string, batchID string, procImage []byte) error {

	imageUUID, err := uuid.Parse(imageID)
	if err != nil {
		return fmt.Errorf("error parsing image uuid: %w", err)
	}

	batchUUID, err := uuid.Parse(batchID)
	if err != nil {
		return fmt.Errorf("error parsing batch uuid: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tx, err := p.db.BeginTx(ctx, nil)

	currentTime := time.Now()
	imageStatus := models.ImageStatusCompleted
	query := "UPDATE images SET image_proc_bw = $1, status = $2, processed_at = $3 WHERE image_id = $4"
	_, err = tx.Exec(query, procImage, imageStatus, currentTime, imageUUID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting updated image %s: %w", imageID, err)
	}

	batchStatus := models.BatchStatusCompleted
	batchQuery := `
		UPDATE batches
		SET
			processed_images = processed_images + 1,
			status = CASE
				WHEN processed_images + 1 = total_images THEN $1::batch_status
				ELSE status
			END,
			completed_at = CASE
				WHEN processed_images + 1 = total_images THEN $3
				ELSE completed_at
			END
		WHERE batch_id = $2
	`
	_, err = tx.Exec(batchQuery, batchStatus, batchUUID, currentTime)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating batch status %s: %w", batchID, err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error committing transaction: %w", err)
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
