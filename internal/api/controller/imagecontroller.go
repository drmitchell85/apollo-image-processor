package controller

import (
	"apollo-image-processor/internal/api/messenger"
	"apollo-image-processor/internal/api/repository"
	"apollo-image-processor/internal/models"
	"context"
	"fmt"
)

type ImageController interface {
	UploadImages(context.Context, []models.UploadedFile) error
}

type imageController struct {
	imageRepo      repository.ImageRepository
	messengerQueue messenger.MessengerQueue
}

func NewImageController(repo repository.ImageRepository, messengerQueue messenger.MessengerQueue) ImageController {
	return &imageController{
		imageRepo:      repo,
		messengerQueue: messengerQueue,
	}
}

func (ic imageController) UploadImages(ctx context.Context, files []models.UploadedFile) error {

	batch_id, imageIDs, err := ic.imageRepo.UploadImages(ctx, files)
	if err != nil {
		return err
	}

	err = ic.messengerQueue.PublishMessage(batch_id, imageIDs)
	if err != nil {
		return fmt.Errorf("error sending batch %s to messenger: %w", batch_id, err)
	}

	return nil
}
