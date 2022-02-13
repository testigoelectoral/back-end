package main

import (
	"bytes"
	"errors"
	"log"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"
)

type UserDataFake struct{}

func (u *UserDataFake) GetHash(userSub string) (string, error) { //nolint:revive
	if userSub == "cognitoError" {
		return "", errors.New("BAD COGNITO")
	}

	options := map[string]string{
		"sub1":    "hash",
		"sub2":    "hash2",
		"sub3":    "hash",
		"errors3": "hash",
	}

	return options[userSub], nil
}

type ImageDataFake struct{}

func (i *ImageDataFake) GetMeta(bucket string, key string) (map[string]string, error) {
	if key == "uploaded/errorS3" {
		return map[string]string{}, errors.New("BAD S3")
	}

	options := map[string]string{
		"uploaded/name1":   "sub1",
		"uploaded/name2":   "sub2",
		"uploaded/error":   "sub3",
		"uploaded/cognito": "cognitoError",
	}

	result := map[string]string{
		"User-Hash": "hash",
		"User-Sub":  options[key],
		"Qr-Code":   "711600102070110113201",
	}

	return result, nil
}

type DataStorageFake struct{}

func (d *DataStorageFake) Save(i ImageRecord) error {
	if i.ImageID == "error" {
		return errors.New("SAVE ERROR")
	}

	return nil
}

func init() {
	userData = &UserDataFake{}
	imageData = &ImageDataFake{}
	dataStorage = &DataStorageFake{}
}

func Test_Handler(t *testing.T) {
	c := require.New(t)

	buf := new(bytes.Buffer)

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(buf)

	handler(eventRequest("uploaded/name1"))
	c.NotNil(buf.String())
	c.Equal("INFO: Record for 'name1' object created\n", buf.String())

	buf.Reset()
	handler(eventRequest("uploaded/name2"))
	c.NotNil(buf.String())
	c.Equal("WARNING: Object 'name2' can't be processed because: user hash differ of s3 header hash\n", buf.String())

	buf.Reset()
	handler(eventRequest("uploaded/error"))
	c.NotNil(buf.String())
	c.Equal("WARNING: Record for 'error' object can't be created because: SAVE ERROR\n", buf.String())

	buf.Reset()
	handler(eventRequest("uploaded/cognito"))
	c.NotNil(buf.String())
	c.Equal("WARNING: Object 'cognito' can't be processed because: BAD COGNITO\n", buf.String())

	buf.Reset()
	handler(eventRequest("uploaded/errorS3"))
	c.NotNil(buf.String())
	c.Equal("WARNING: Object 'errorS3' can't be processed because: BAD S3\n", buf.String())
}

func eventRequest(key string) events.S3Event {
	return events.S3Event{
		Records: []events.S3EventRecord{{S3: events.S3Entity{
			Bucket: events.S3Bucket{Name: *aws.String("BuckeName")},
			Object: events.S3Object{Key: *aws.String(key)},
		}}},
	}
}

func Test_metaFromHeaders(t *testing.T) {
	c := require.New(t)

	result := metaFromHeaders(map[string]string{"Qr-Code": "711600102070110113201"})
	c.NotNil(result)
	c.IsType(PageMeta{}, result)
	c.Equal(uint8(71), result.PageType)
	c.Equal(uint8(16), result.LocationStateCode)
	c.Equal(uint8(1), result.LocationMunicipalityCode)
	c.Equal(uint8(2), result.LocationZoneCode)
	c.Equal(uint16(7), result.LocationPlace)
	c.Equal(uint16(11), result.LocationTable)
	c.Equal(uint8(1), result.PageNumber)
}
