package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Key interface {
	Generate() string
}

type S3FileStorage struct {
	client *s3.Client
}

func NewS3FileStorage(s3Client *s3.Client) *S3FileStorage {
	return &S3FileStorage{client: s3Client}
}

func (s *S3FileStorage) Upload(bucket string, key string, body io.Reader) (*s3.PutObjectOutput, error) {
	uploaded, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   body,
	})
	if err != nil {
		return nil, err
	}
	return uploaded, nil
}
