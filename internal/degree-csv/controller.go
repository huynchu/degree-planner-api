package degree

import (
	"encoding/json"
	"net/http"
	"strings"
)

type DegreeCsvController struct {
	storage *DegreeCsvStorage
}

func NewDegreeCsvController(storage *DegreeCsvStorage) *DegreeCsvController {
	return &DegreeCsvController{
		storage: storage,
	}
}

func (s *DegreeCsvController) UploadDegreeCsv(w http.ResponseWriter, r *http.Request) {
	// Checks that request body is in form-data format
	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		http.Error(w, "request body is not in form-data format", http.StatusBadRequest)
		return
	}

	// Parse form data
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20+512)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get fileheaders from form data: UploadedFiles
	fileHeader, exists := r.MultipartForm.File["file"]
	if !exists || len(fileHeader) != 1 {
		http.Error(w, "Request body must contain one file", http.StatusBadRequest)
		return
	}

	// Get file from fileheader
	file, err := fileHeader[0].Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Upload file to s3
	uploadResult, err := s.storage.Upload(fileHeader[0].Filename, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode json response
	response := struct {
		FileUrl  string      `json:"file_url"`
		Metadata interface{} `json:"location"`
	}{
		Metadata: uploadResult.ResultMetadata,
		FileUrl:  "https://" + s.storage.bucket + ".s3.amazonaws.com/" + fileHeader[0].Filename,
	}

	// Respond with json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
