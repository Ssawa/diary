package utils

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	awsCreds "github.com/aws/aws-sdk-go/aws/credentials"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"
)

// createAWSSession initializes a new aws session using the configured data
func createAWSSession() (*awsSession.Session, error) {
	if err := ValidateConfigs("s3.region", "s3.id", "s3.secret"); err != nil {
		return nil, err
	}

	return awsSession.NewSession(&aws.Config{
		Region: aws.String(viper.GetString("s3.region")),
		Credentials: awsCreds.NewStaticCredentials(
			viper.GetString("s3.id"),
			viper.GetString("s3.secret"),
			"",
		),
	})
}

func createS3Uploader() (*s3manager.Uploader, error) {
	session, err := createAWSSession()
	if err != nil {
		return nil, err
	}
	return s3manager.NewUploader(session), nil
}

func uploadFileToS3(entry Entry, uploader *s3manager.Uploader) error {
	ValidateConfigs("s3.bucket")

	file, err := os.Open(entry.fullPath)
	if err != nil {
		return err
	}
	bucket := viper.GetString("s3.bucket")
	key := viper.GetString("s3.key_prefix") + entry.relativePath
	Verbose.Printf("Uploading to S3 to bucket %s with key %s", bucket, key)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	return err
}

// BackupFile persists a file to S3
func BackupFile(entry Entry) error {
	Verbose.Println("Attempting to upload to s3")
	uploader, err := createS3Uploader()
	if err != nil {
		return err
	}

	return uploadFileToS3(entry, uploader)
}
