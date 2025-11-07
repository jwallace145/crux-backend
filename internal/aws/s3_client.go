package aws

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

var (
	S3Client *s3.Client
	logger   *zap.Logger
)

// InitS3Client initializes the S3 client with AWS credentials
func InitS3Client(ctx context.Context, log *zap.Logger) error {
	logger = log

	logger.Info("Initializing S3 client")

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Error("Failed to load AWS configuration",
			zap.Error(err),
		)
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	// Create S3 client
	S3Client = s3.NewFromConfig(cfg)

	logger.Info("S3 client initialized successfully")

	return nil
}

// UploadFile uploads a file to S3 and returns the S3 URI
func UploadFile(ctx context.Context, bucket, key string, body io.Reader, contentType string) (string, error) {
	if S3Client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}

	logger.Info("Uploading file to S3",
		zap.String("bucket", bucket),
		zap.String("key", key),
		zap.String("content_type", contentType),
	)

	// Upload file to S3
	_, err := S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		logger.Error("Failed to upload file to S3",
			zap.Error(err),
			zap.String("bucket", bucket),
			zap.String("key", key),
		)
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Return S3 URI
	s3URI := fmt.Sprintf("s3://%s/%s", bucket, key)

	logger.Info("File uploaded successfully to S3",
		zap.String("bucket", bucket),
		zap.String("key", key),
		zap.String("s3_uri", s3URI),
	)

	return s3URI, nil
}

// GeneratePresignedURL generates a presigned URL for accessing an S3 object
func GeneratePresignedURL(ctx context.Context, bucket, key string, expirationMinutes int) (string, error) {
	if S3Client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}

	logger.Debug("Generating presigned URL",
		zap.String("bucket", bucket),
		zap.String("key", key),
		zap.Int("expiration_minutes", expirationMinutes),
	)

	// Create presign client
	presignClient := s3.NewPresignClient(S3Client)

	// Generate presigned URL
	presignedURL, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(time.Duration(expirationMinutes)*time.Minute))
	if err != nil {
		logger.Error("Failed to generate presigned URL",
			zap.Error(err),
			zap.String("bucket", bucket),
			zap.String("key", key),
		)
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	logger.Debug("Presigned URL generated successfully",
		zap.String("bucket", bucket),
		zap.String("key", key),
		zap.String("url", presignedURL.URL),
	)

	return presignedURL.URL, nil
}

// DeleteFile deletes a file from S3
func DeleteFile(ctx context.Context, bucket, key string) error {
	if S3Client == nil {
		return fmt.Errorf("S3 client not initialized")
	}

	logger.Info("Deleting file from S3",
		zap.String("bucket", bucket),
		zap.String("key", key),
	)

	_, err := S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		logger.Error("Failed to delete file from S3",
			zap.Error(err),
			zap.String("bucket", bucket),
			zap.String("key", key),
		)
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	logger.Info("File deleted successfully from S3",
		zap.String("bucket", bucket),
		zap.String("key", key),
	)

	return nil
}
