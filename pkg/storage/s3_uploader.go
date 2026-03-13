package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type Uploader interface {
	UploadImage(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (string, error)
}

type S3Config struct {
	Region        string
	Bucket        string
	AccessKey     string
	SecretKey     string
	PublicBaseURL string
}

type S3Uploader struct {
	client        *s3.Client
	bucket        string
	region        string
	publicBaseURL string
}

func NewS3Uploader(cfg S3Config) (*S3Uploader, error) {
	if cfg.Region == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("s3 region and bucket are required")
	}

	loadOptions := []func(*awsConfig.LoadOptions) error{
		awsConfig.WithRegion(cfg.Region),
	}

	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		loadOptions = append(loadOptions, awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		))
	}

	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background(), loadOptions...)
	if err != nil {
		return nil, err
	}

	return &S3Uploader{
		client:        s3.NewFromConfig(awsCfg),
		bucket:        cfg.Bucket,
		region:        cfg.Region,
		publicBaseURL: strings.TrimSuffix(cfg.PublicBaseURL, "/"),
	}, nil
}

func (s *S3Uploader) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (string, error) {
	if fileHeader == nil {
		return "", fmt.Errorf("file is required")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExt := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowedExt[ext] {
		return "", fmt.Errorf("unsupported image format")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	if folder == "" {
		folder = "uploads"
	}
	folder = strings.Trim(folder, "/")
	key := fmt.Sprintf("%s/%s-%d%s", folder, uuid.NewString(), time.Now().Unix(), ext)
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = inferContentTypeByExt(ext)
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	if s.publicBaseURL != "" {
		return fmt.Sprintf("%s/%s", s.publicBaseURL, key), nil
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key), nil
}

func inferContentTypeByExt(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
