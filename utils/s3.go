package utils

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awsCreds "github.com/aws/aws-sdk-go/aws/credentials"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

	file, err := os.Open(entry.FullPath)
	if err != nil {
		return err
	}
	bucket := viper.GetString("s3.bucket")
	key := viper.GetString("s3.key_prefix") + entry.RelativePath
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

func createS3Downloader() (*s3manager.Downloader, error) {
	session, err := createAWSSession()
	if err != nil {
		return nil, err
	}
	return s3manager.NewDownloader(session), nil
}

func downloadFileFromS3(entry Entry, uploader *s3manager.Downloader) error {
	ValidateConfigs("s3.bucket")

	Verbose.Println(entry.FullPath)
	file, err := os.OpenFile(entry.FullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	bucket := viper.GetString("s3.bucket")
	key := viper.GetString("s3.key_prefix") + entry.RelativePath
	Verbose.Printf("Downloading from S3 bucket %s key %s to %s", bucket, key, entry.FullPath)
	_, err = uploader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}

// getKeys finds matching keys and calls a function in response to each page
// of data. The callback function takes in an array of Entries which matches the number keys in
// the ListObjectsOutput and a boolean indicating if it is the last page. Return true or false
// to determine if we should stop processing more pages
func getKeys(fn func(entries []Entry, page *s3.ListObjectsOutput, lastPage bool) bool) error {
	session, err := createAWSSession()
	if err != nil {
		return err
	}

	client := s3.New(session)

	basePath := viper.GetString("file.base")
	prefix := viper.GetString("s3.key_prefix")
	client.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(viper.GetString("s3.bucket")),
		Prefix: aws.String(prefix),
	}, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		var entries []Entry
		for _, item := range page.Contents {
			relativePath := strings.Replace(*item.Key, prefix, "", 1)
			entries = append(entries, Entry{
				RelativePath: relativePath,
				BasePath:     basePath,
				FullPath:     path.Join(basePath, relativePath),
			})
		}
		return fn(entries, page, lastPage)
	})
	return nil
}

// ListExistingKeys prints keys that already exists out to stdout
func ListExistingKeys() error {
	return getKeys(func(entries []Entry, page *s3.ListObjectsOutput, lastPage bool) bool {
		for index, entry := range entries {
			fmt.Println(entry.RelativePath)
			Verbose.Printf("S3 Key: %s", *page.Contents[index].Key)
		}
		return true
	})
}

// DownloadEntry downloads a single file from S3
func DownloadEntry(entry Entry) error {
	downloader, err := createS3Downloader()
	if err != nil {
		return err
	}

	return downloadFileFromS3(entry, downloader)
}

// DownloadAll downloads all files from the remote
func DownloadAll() error {
	var innerError error
	err := getKeys(func(entries []Entry, page *s3.ListObjectsOutput, lastPage bool) bool {
		for _, entry := range entries {
			fmt.Println("Downloading", entry.RelativePath)
			if err := DownloadEntry(entry); err != nil {
				innerError = err
				return false
			}
		}
		return true
	})
	if err != nil {
		return err
	}
	if innerError != nil {
		return innerError
	}
	return nil
}
