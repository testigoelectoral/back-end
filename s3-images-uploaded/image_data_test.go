package main

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/require"
)

type s3ifaceFake struct {
	s3iface.S3API
}

func (s *s3ifaceFake) HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	if aws.StringValue(input.Key) == "error" {
		return nil, errors.New("S3 ERROR")
	}

	output := &s3.HeadObjectOutput{
		LastModified: aws.Time(time.Now()),
		Metadata: map[string]*string{
			"User-Hash": aws.String("hash"),
		},
	}

	return output, nil
}

func init() {}

func Test_GetMeta(t *testing.T) {
	c := require.New(t)

	s3dataTest := &S3Data{service: &s3ifaceFake{}}

	result, err := s3dataTest.GetMeta("bucket", "key")
	c.Nil(err)
	c.NotNil(result)
	c.Equal("hash", result["User-Hash"])

	result, err = s3dataTest.GetMeta("bucket", "error")
	c.NotNil(result)
	c.NotNil(err)
	c.Equal("S3 ERROR", err.Error())
}
