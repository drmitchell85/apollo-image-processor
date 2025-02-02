package handler

import (
	"apollo-image-processor/internal/controller"
	"apollo-image-processor/internal/models"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type ImageHandler struct {
	imageController controller.ImageController
}

func NewImageHandler(imageController controller.ImageController) *ImageHandler {
	return &ImageHandler{
		imageController: imageController,
	}
}

/*
Receive and Respond to the Multipart form Request
https://andrew-mccall.com/blog/2024/06/golang-send-multipart-form-data-to-api-endpoint/

- validate the request method
- parse the form data
- process the requests file
- construct a payload to respond with
*/
func (ih ImageHandler) UploadImages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := r.ParseMultipartForm(32 << 10) // 32mb
	if err != nil {
		respondFailure(w, http.StatusInternalServerError, fmt.Errorf("error parsing form: %w", err))
		return
	}

	// var newFile uploadedFile
	var files []models.UploadedFile

	for fKey, fHeaders := range r.MultipartForm.File {

		for fNum, headers := range fHeaders {

			var newFile models.UploadedFile

			// return the associated file
			file, err := headers.Open()
			if err != nil {
				respondFailure(w, http.StatusInternalServerError, fmt.Errorf("error opening form: %w", err))
				return
			}
			defer file.Close()

			// detect content type
			buff := make([]byte, 512)
			file.Read(buff)
			file.Seek(0, 0) // reset to start

			contentType := http.DetectContentType(buff)
			newFile.ContentType = contentType

			// get file size
			var sizeBuff bytes.Buffer
			fileSize, err := sizeBuff.ReadFrom(file)
			if err != nil {
				respondFailure(w, http.StatusInternalServerError, fmt.Errorf("error getting file size: %w", err))
				return
			}

			file.Seek(0, 0) // reset to start
			newFile.Size = fileSize
			newFile.Filename = headers.Filename

			newFile.FileContent, err = io.ReadAll(file)
			if err != nil {
				respondFailure(w, http.StatusInternalServerError, fmt.Errorf("error reading content: %w", err))
				return
			}

			newFile.Key = fKey
			newFile.ImageNum = fNum
			files = append(files, newFile)

		}

	}

	err = ih.imageController.UploadImages(ctx, files)
	if err != nil {
		respondFailure(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusOK, files)
}
