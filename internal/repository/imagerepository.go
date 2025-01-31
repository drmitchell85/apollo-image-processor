package repository

import "database/sql"

type ImageRepository interface {
	UploadImages() error
}

type imageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) ImageRepository {
	return &imageRepository{
		db: db,
	}
}

func (ir *imageRepository) UploadImages() error {
	return nil
}
