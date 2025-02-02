package controller

import (
	"apollo-image-processor/internal/models"
	"apollo-image-processor/internal/repository"
	"context"
)

type ImageController interface {
	UploadImages(context.Context, []models.UploadedFile) error
}

type imageController struct {
	imageRepo repository.ImageRepository
}

func NewImageController(repo repository.ImageRepository) ImageController {
	return &imageController{
		imageRepo: repo,
	}
}

func (ic imageController) UploadImages(ctx context.Context, files []models.UploadedFile) error {

	batchID, imageIDs, err := ic.imageRepo.UploadImages(ctx, files)
	if err != nil {
		return err
	}

	err = ic.imageRepo.PublishMessage(batchID, imageIDs)
	if err != nil {
		return err
	}

	return nil
}
