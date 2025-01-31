package controller

import "apollo-image-processor/internal/repository"

type ImageController interface {
	UploadImages() error
}

type imageController struct {
	imageRepo repository.ImageRepository
}

func NewImageController(repo repository.ImageRepository) ImageController {
	return &imageController{
		imageRepo: repo,
	}
}

func (ic imageController) UploadImages() error {
	return nil
}
