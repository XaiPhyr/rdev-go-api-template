package aws_test

import (
	"context"
	"testing"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestAWS(t *testing.T) {
	t.Run("test aws successful upload", func(t *testing.T) {
		mockSvc := &aws.MockAWSService{
			PutObjectFunc: func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, nil
			},
		}

		svc, err := aws.NewAWSService("ap-southeast-1", "test-bucket", aws.WithS3Client(mockSvc))

		url, err := svc.UploadToS3(context.Background(), "key-aws", []byte(`test-upload`))

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		expectedURL := "https://test-bucket.s3.ap-southeast-1.amazonaws.com/key-aws"
		if url != expectedURL {
			t.Errorf("expected URL %q, got %q", expectedURL, url)
		}
	})
}
