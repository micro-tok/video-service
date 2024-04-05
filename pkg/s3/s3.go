package s3

import (
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/micro-tok/video-service/pkg/config"
)

type AWSService struct {
	bucketName string
	s3Config   *aws.Config
}

func NewAWSService(cfg *config.Config) *AWSService {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.BucketKey, cfg.BucketSecret, ""),
		Endpoint:         aws.String("s3://micro-tok"),
		Region:           aws.String("eu-north-1"),
		S3ForcePathStyle: aws.Bool(false),
	}

	return &AWSService{
		bucketName: cfg.BucketName,
		s3Config:   s3Config,
	}
}

func (s AWSService) UploadFile(file multipart.File, filepath string) (string, error) {

	newSession, err := session.NewSession(s.s3Config)
	if err != nil {
		return "", err
	}

	s3Client := s3.New(newSession)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(filepath),
		Body:   file,
		ACL:    aws.String("public-read"),
	})

	if err != nil {
		return "", err
	}

	return filepath, nil
}
