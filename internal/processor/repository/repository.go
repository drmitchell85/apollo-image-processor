package procrepository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func GetImage(imageID string, db *sql.DB) ([]byte, error) {

	imageUUID, err := uuid.Parse(imageID)
	if err != nil {
		return nil, fmt.Errorf("error parsing image uuid: %w", err)
	}

	query := fmt.Sprintf("SELECT image FROM images WHERE image_id = '%s'", imageUUID)

	var srcimage []byte
	err = db.QueryRow(query).Scan(&srcimage)
	if err != nil {
		return nil, fmt.Errorf("error fetching image %s from db: %w", imageID, err)
	}

	return srcimage, nil
}

func InsertImage(imageID string, procImage []byte, db *sql.DB) error {

	imageUUID, err := uuid.Parse(imageID)
	if err != nil {
		return fmt.Errorf("error parsing image uuid: %w", err)
	}

	// (image_proc_bw) VALUES ($1) WHERE image_id
	query := fmt.Sprintf("UPDATE images SET image_proc_bw = $1 WHERE image_id = '%s'", imageUUID)
	_, err = db.Exec(query, procImage)
	if err != nil {
		return fmt.Errorf("error inserting updated image %s: %w", imageID, err)
	}

	return nil
}
