package repository

import (
	"apollo-image-processor/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type ImageRepository interface {
	UploadImages(context.Context, []models.UploadedFile) (string, []string, error)
}

type imageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) ImageRepository {
	return &imageRepository{
		db: db,
	}
}

func (ir *imageRepository) UploadImages(ctx context.Context, files []models.UploadedFile) (string, []string, error) {

	var batch_id string
	var imageIDs []string

	tx, err := ir.db.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return batch_id, imageIDs, fmt.Errorf("error starting transaction: %w", err)
	}

	q1 := `INSERT INTO batches (status, total_images) VALUES ($1, $2) RETURNING batch_id`
	err = tx.QueryRow(q1, models.BatchStatusCreated, len(files)).Scan(&batch_id)
	if err != nil {
		tx.Rollback()
		return batch_id, imageIDs, fmt.Errorf("error inserting into batches table: %w", err)
	}

	numFiles := len(files)
	paramPlaceholders := make([]string, numFiles)
	paramData := make([]interface{}, 0, numFiles*4)
	for i := 0; i < numFiles; i++ {
		// Calculate starting parameter index for this row
		start := i*4 + 1

		// Create the placeholder group for this row
		paramPlaceholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d)",
			start, start+1, start+2, start+3)

		// Add the actual parameters for this row
		paramData = append(paramData,
			batch_id,
			models.ImageStatusPending,
			files[i].Filename,
			files[i].FileContent)
	}

	q2 := fmt.Sprintf(
		"INSERT INTO images (batch_id, status, image_name, image) VALUES %s RETURNING image_id",
		strings.Join(paramPlaceholders, ", "))

	// _, err = tx.ExecContext(ctx, q2, paramData...)
	rows, err := tx.Query(q2, paramData...)
	if err != nil {
		tx.Rollback()
		return batch_id, imageIDs, fmt.Errorf("error inserting images: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var imageID string
		err = rows.Scan(&imageID)
		if err != nil {
			tx.Rollback()
			return batch_id, imageIDs, fmt.Errorf("error getting imageID: %w", err)
		}

		imageIDs = append(imageIDs, imageID)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return batch_id, imageIDs, fmt.Errorf("error committing transaction: %w", err)
	}

	return batch_id, imageIDs, nil
}
