package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type ImageDataInterface interface {
	GetMeta(string, string) (map[string]string, error)
}

type S3Data struct {
	service s3iface.S3API
}

func NewS3Data(sess *session.Session) *S3Data {
	return &S3Data{
		service: s3.New(sess),
	}
}

func (s *S3Data) GetMeta(bucket string, key string) (map[string]string, error) {
	s3ObjectFilter := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	s3HeadData, err := s.service.HeadObject(s3ObjectFilter)
	if err != nil {
		return map[string]string{}, err
	}

	metaReturn := make(map[string]string)

	for k, v := range s3HeadData.Metadata {
		metaReturn[k] = aws.StringValue(v)
	}

	metaReturn["CreateAt"] = s3HeadData.LastModified.Format(time.RFC3339)

	return metaReturn, nil
}
