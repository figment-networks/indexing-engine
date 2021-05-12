package datalake

import (
	"bytes"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type s3Storage struct {
	region     string
	bucket     string
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

var _ Storage = (*s3Storage)(nil)

// NewS3Storage creates an Amazon S3 storage
func NewS3Storage(region string, bucket string) Storage {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	return &s3Storage{
		region:     region,
		bucket:     bucket,
		client:     s3.New(sess),
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
	}
}

func (ss *s3Storage) Store(data []byte, path ...string) error {
	_, err := ss.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(ss.bucket),
		Key:    aws.String(ss.storageKey(path)),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return err
	}

	return nil
}

func (ss *s3Storage) IsStored(path ...string) (bool, error) {
	_, err := ss.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(ss.bucket),
		Key:    aws.String(ss.storageKey(path)),
	})

	if err != nil {
		reqErr, ok := err.(awserr.RequestFailure)
		if ok && reqErr.StatusCode() == 404 {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (ss *s3Storage) Retrieve(path ...string) ([]byte, error) {
	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := ss.downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(ss.bucket),
		Key:    aws.String(ss.storageKey(path)),
	})

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (ss *s3Storage) storageKey(path []string) string {
	return strings.Join(path, "/")
}
