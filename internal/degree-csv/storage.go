package degree

import (
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/huynchu/degree-planner-api/internal/storage"
)

type DegreeCsvStorage struct {
	bucket  string
	storage *storage.S3FileStorage
}

func NewDegreeCsvStorage(bucket string, storage *storage.S3FileStorage) *DegreeCsvStorage {
	return &DegreeCsvStorage{
		bucket:  bucket,
		storage: storage,
	}
}

func (s *DegreeCsvStorage) Upload(key string, file multipart.File) (*s3.PutObjectOutput, error) {
	return s.storage.Upload(s.bucket, key, file)
}

func (s *DegreeCsvStorage) GetFileDownloadLink(key string) (string, error) {
	return s.storage.GetFileDownloadLink(s.bucket, key)
}
