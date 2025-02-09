package procrepository

import (
	"database/sql"
	"fmt"
)

func GetImage(imageID string, db *sql.DB) ([]byte, error) {

	q1 := fmt.Sprintf("SELECT image FROM images WHERE image_id = '%s'", imageID)

	var srcimage []byte
	err := db.QueryRow(q1).Scan(&srcimage)
	if err != nil {
		return nil, fmt.Errorf("error fetching image %s from db: %w", imageID, err)
	}

	return srcimage, nil
}
