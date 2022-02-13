package main

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	cognitoPoolID string
	userData      UserDataInterface
	imageData     ImageDataInterface
	dataStorage   DataStorageInterface
)

func init() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	cognitoPoolID = os.Getenv("POOL_ID")
	userData = NewCognitoData(sess)
	imageData = NewS3Data(sess)
	dataStorage = NewDynamoDBStorage(sess, os.Getenv("DYNAMODB_IMAGE_TABLE"))
}

func main() {
	lambda.Start(handler)
}

func handler(events events.S3Event) {
	for _, record := range events.Records {
		ImageID := idFromKey(record.S3.Object.Key)

		s3Metadata, err := validateHash(record.S3.Bucket.Name, record.S3.Object.Key)
		if err != nil {
			log.Printf("WARNING: Object '%s' can't be processed because: %s", ImageID, err.Error())
			continue
		}

		err = createRecord(record.S3.Object.Key, s3Metadata)
		if err != nil {
			log.Printf("WARNING: Record for '%s' object can't be created because: %s", ImageID, err.Error())
			continue
		}

		log.Printf("INFO: Record for '%s' object created", ImageID)
	}
}

func idFromKey(key string) string {
	return strings.Split(key, "/")[1]
}

func validateHash(bucket string, key string) (map[string]string, error) {
	s3Meta, err := imageData.GetMeta(bucket, key)
	if err != nil {
		return map[string]string{}, err
	}

	userHash, err := userData.GetHash(s3Meta["User-Sub"])
	if err != nil {
		return map[string]string{}, err
	}

	if userHash != s3Meta["User-Hash"] {
		return map[string]string{}, errors.New("user hash differ of s3 header hash")
	}

	return s3Meta, nil
}

func createRecord(key string, s3Meta map[string]string) error {
	record := ImageRecord{
		ImageID:     idFromKey(key),
		OwnerSub:    s3Meta["User-Sub"],
		OwnerGPS:    gpsFromHeaders(s3Meta),
		OwnerReport: false,
		PageMeta:    metaFromHeaders(s3Meta),
		CreatedAt:   s3Meta["CreateAt"],
		UpdateAt:    s3Meta["CreateAt"],
		OwnerQRCode: s3Meta["Qr-Code"],
	}

	return dataStorage.Save(record)
}

func gpsFromHeaders(s3Meta map[string]string) GPSRecord {
	lat, _ := strconv.ParseFloat(s3Meta["Latitude"], 64)
	lon, _ := strconv.ParseFloat(s3Meta["Longitude"], 64)
	acu, _ := strconv.ParseFloat(s3Meta["Accuracy"], 64)

	return GPSRecord{
		Latitude:  lat,
		Longitude: lon,
		Accuracy:  acu,
	}
}

func metaFromHeaders(s3Meta map[string]string) PageMeta {
	qr := s3Meta["Qr-Code"]
	if len(qr) != 21 {
		return PageMeta{}
	}

	PageType, _ := strconv.ParseUint(qr[0:2], 10, 8)
	LocationStateCode, _ := strconv.ParseUint(qr[2:4], 10, 8)
	LocationMunicipalityCode, _ := strconv.ParseUint(qr[4:7], 10, 8)
	LocationZoneCode, _ := strconv.ParseUint(qr[7:9], 10, 8)
	LocationPlace, _ := strconv.ParseUint(qr[9:11], 10, 8)
	LocationTable, _ := strconv.ParseUint(qr[11:14], 10, 8)
	PageNumber, _ := strconv.ParseUint(qr[14:16], 10, 8)

	return PageMeta{
		LocationStateCode:        uint8(LocationStateCode),
		LocationMunicipalityCode: uint8(LocationMunicipalityCode),
		LocationZoneCode:         uint8(LocationZoneCode),
		LocationPlace:            uint16(LocationPlace),
		LocationTable:            uint16(LocationTable),
		PageNumber:               uint8(PageNumber),
		PageType:                 uint8(PageType),
		PageQR:                   qr,
	}
}
