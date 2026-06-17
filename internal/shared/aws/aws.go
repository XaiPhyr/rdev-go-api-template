package aws

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3API interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type AWSService interface {
	UploadToS3(ctx context.Context, key string, data []byte) (string, error)
}

type service struct {
	s3Client S3API
	bucket   string
	region   string
}

func WithS3Client(client S3API) func(*service) {
	return func(s *service) {
		s.s3Client = client
	}
}

func NewAWSService(region, bucket string, opts ...func(*service)) (*service, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	svc := &service{
		s3Client: s3.NewFromConfig(cfg),
		bucket:   bucket,
		region:   region,
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc, nil
}

func (s *service) UploadToS3(ctx context.Context, key string, data []byte) (string, error) {
	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(http.DetectContentType(data)),
	})

	if err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)

	return publicURL, nil
}
