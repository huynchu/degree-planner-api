package storage

import (
	"mime/multipart"
)

type FileStorage interface {
	upload(fileHeader *multipart.FileHeader) (string, error)
}

type MockFileStorage struct {
	files map[string][]byte
}

func NewMockFileStorage() *MockFileStorage {
	return &MockFileStorage{files: make(map[string][]byte)}
}

func (m *MockFileStorage) upload(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileBytes := make([]byte, fileHeader.Size)
	_, err = file.Read(fileBytes)
	if err != nil {
		return "", err
	}

	m.files[fileHeader.Filename] = fileBytes

	return fileHeader.Filename, nil
}
