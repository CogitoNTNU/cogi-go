package util

import (
	"bytes"
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

func GetS3Client(e *env.EnvConfig, logger *logrus.Entry, ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(e.AWSRegion))
	if err != nil {
		logger.Errorf("unable to load AWS SDK config, %v", err)
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

func UploadFile(e *env.EnvConfig, ctx context.Context, client *s3.Client, key string, body []byte) error {
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(e.AWSBucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})
	if err != nil {
		return err
	}
	return nil
}

func UploadLargeFile(e *env.EnvConfig, ctx context.Context, client *s3.Client, key string, body []byte) error {
	largeBuffer := bytes.NewReader(body)

	var partMiBs int64 = 10
	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(e.AWSBucketName),
		Key:    aws.String(key),
		Body:   largeBuffer,
	})
	if err != nil {
		return err
	}
	return nil
}

func DeleteFile(e *env.EnvConfig, ctx context.Context, client *s3.Client, key string) error {
	_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(e.AWSBucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	return nil
}

func GetPublicURL(e *env.EnvConfig, key string) string {
	escapedKey := strings.ReplaceAll(key, " ", "%20")
	return "https://" + e.AWSBucketName + ".s3." + e.AWSRegion + ".amazonaws.com/" + escapedKey
}