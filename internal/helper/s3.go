package helper

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/Dzikuri/openidea-segokuning/internal/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	MaxFileSize = 1024 * 1024 * 2
	MinFileSize = 1024 * 10
	JPG         = ".jpg"
	JPEG        = ".jpeg"
)

func UploadFileToS3(ctx context.Context, filename string, file *multipart.FileHeader) (string, error) {
	awsRegion := os.Getenv("S3_REGION")
	awsAccessKeyId := os.Getenv("S3_ID")
	awsSecretAccessKey := os.Getenv("S3_SECRET_KEY")

	bucketName := os.Getenv("S3_BUCKET_NAME")

	// NOTE Check File Size
	if file.Size > MaxFileSize || file.Size < MinFileSize {
		return "", model.ErrFileSizeNotValid
	}
	// NOTE Check File Extension
	if file.Filename != filepath.Base(file.Filename) {
		return "", model.ErrExtensionNotValid
	}
	// NOTE Check File Extension
	ext := filepath.Ext(filepath.Base(file.Filename))
	if ext != JPG && ext != JPEG {
		return "", model.ErrExtensionNotValid
	}

	contentType := file.Header.Get("Content-Type")

	var timeout time.Duration
	conf := &aws.Config{
		Endpoint:    nil,
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyId, awsSecretAccessKey, ""),
	}
	s3Session := session.Must(session.NewSession(conf))

	ctx = context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}

	if cancelFn != nil {
		defer cancelFn()
	}

	stream, err := file.Open()
	if condition := err != nil; condition {
		return "", err

	}

	uploader := s3manager.NewUploader(s3Session)

	var uploadInput s3manager.UploadInput

	uploadInput = s3manager.UploadInput{
		Bucket: aws.String(bucketName), // bucket's name
		// NOTE Change location to folder filenames
		Key:         aws.String(filename),    // files destination location
		Body:        stream,                  // content of the file
		ContentType: aws.String(contentType), // content type
		ACL:         aws.String("public-read"),
	}

	doUpload := &uploadInput

	output, err := uploader.UploadWithContext(ctx, doUpload)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return output.Location, nil

}
