package handler

import (
	"apollo-image-processor/internal/controller"
	"bytes"
	"encoding/json"
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
	body := r.Body
	fmt.Printf("body: %+v", body)

	err := r.ParseMultipartForm(32 << 10) // 32mb
	if err != nil {
		respondFailure(w, http.StatusInternalServerError, fmt.Errorf("error parsing form: %w", err))
		return
	}

	name := r.FormValue("name")

	type uploadedFile struct {
		Size        int64  `json:"size"`
		ContentType string `json:"content_type"`
		Filename    string `json:"filename"`
		FileContent string `json:"file_content"`
	}

	// var newFile uploadedFile
	var files []uploadedFile

	for _, fHeaders := range r.MultipartForm.File {

		for _, headers := range fHeaders {

			var newFile uploadedFile

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
			file.Seek(0, 0)

			contentType := http.DetectContentType(buff)
			newFile.ContentType = contentType

			// get file size
			var sizeBuff bytes.Buffer
			fileSize, err := sizeBuff.ReadFrom(file)
			if err != nil {
				respondFailure(w, http.StatusInternalServerError, fmt.Errorf("error getting file size: %w", err))
				return
			}

			file.Seek(0, 0)
			newFile.Size = fileSize
			newFile.Filename = headers.Filename
			contentBuf := bytes.NewBuffer(nil)

			if _, err := io.Copy(contentBuf, file); err != nil {
				respondFailure(w, http.StatusInternalServerError, fmt.Errorf("error with content buffer: %w", err))
				return
			}

			newFile.FileContent = contentBuf.String()

			files = append(files, newFile)

		}

	}
	data := make(map[string]interface{})

	data["form_field_value"] = name
	data["status"] = 200
	data["file_stats"] = files

	fmt.Println("name: ", name)

	if err = json.NewEncoder(w).Encode(data); err != nil {
		respondFailure(w, http.StatusInternalServerError, fmt.Errorf("error encoding: %w", err))
		return
	}

	respondSuccess(w, http.StatusOK, data)
}
