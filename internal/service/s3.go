package service

import (
	"bytes"
	"context"
	"strings"

	"github.com/CogitoNTNU/cogi-go/internal/util/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

type S3 struct {
	logger *logrus.Entry
	env    *env.EnvConfig
	client *s3.Client
}

func NewS3Service(logger *logrus.Entry, cfg *env.EnvConfig, client *s3.Client) *S3 {
	return &S3{
		logger: logger,
		env:    cfg,
		client: client,
	}
}

func GetS3Client(e *env.EnvConfig, logger *logrus.Entry, ctx context.Context) (*s3.Client, error) {
	region, err := e.Read("AWS_REGION")
	if err != nil {
		logger.Errorf("unable to read AWS_REGION, %v", err)
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		logger.Errorf("unable to load AWS SDK config, %v", err)
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

func (s *S3) UploadFile(ctx context.Context, key string, body []byte) error {
	bucketName, err := s.env.Read("AWS_BUCKET_NAME")
	if err != nil {
		s.logger.Errorf("failed to read AWS_BUCKET_NAME, %v", err)
		return err
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})
	if err != nil {
		s.logger.Errorf("failed to upload file, %v", err)
		return err
	}

	return nil
}

func (s *S3) UploadLargeFile(ctx context.Context, key string, body []byte) error {
	bucketName, err := s.env.Read("AWS_BUCKET_NAME")
	if err != nil {
		s.logger.Errorf("failed to read AWS_BUCKET_NAME, %v", err)
		return err
	}

	largeBuffer := bytes.NewReader(body)

	var partMiBs int64 = 10
	uploader := manager.NewUploader(s.client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})

	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   largeBuffer,
	})
	if err != nil {
		s.logger.Errorf("failed to upload large file, %v", err)
		return err
	}

	return nil
}

func (s *S3) DeleteFile(ctx context.Context, key string) error {
	bucketName, err := s.env.Read("AWS_BUCKET_NAME")
	if err != nil {
		s.logger.Errorf("failed to read AWS_BUCKET_NAME, %v", err)
		return err
	}

	_, err = s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		s.logger.Errorf("failed to delete file, %v", err)
		return err
	}

	return nil
}

func (s *S3) GetPublicURL(key string) string {
	bucketName, err := s.env.Read("AWS_BUCKET_NAME")
	if err != nil {
		s.logger.Errorf("failed to read AWS_BUCKET_NAME, %v", err)
		return ""
	}

	region, err := s.env.Read("AWS_REGION")
	if err != nil {
		s.logger.Errorf("failed to read AWS_REGION, %v", err)
		return ""
	}

	escapedKey := strings.ReplaceAll(key, " ", "%20")
	return "https://" + bucketName + ".s3." + region + ".amazonaws.com/" + escapedKey
}