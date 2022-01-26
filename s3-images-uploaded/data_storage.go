package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type GPSRecord struct {
	Latitude  float64
	Longitude float64
	Accuracy  float64
}

type ImageRecord struct {
	ImageID     string
	OwnerSub    string
	OwnerGPS    GPSRecord
	OwnerReport bool
	CreatedAt   string
	UpdateAt    string
	OwnerQRCode string
}

type DataStorageInterface interface {
	Save(ImageRecord) error
}

type DynamoDBStorage struct {
	service dynamodbiface.DynamoDBAPI
	table   string
}

func NewDynamoDBStorage(sess *session.Session, tablename string) *DynamoDBStorage {
	return &DynamoDBStorage{
		service: dynamodb.New(sess),
		table:   tablename,
	}
}

func (d *DynamoDBStorage) Save(recod ImageRecord) error {
	item, err := dynamodbattribute.MarshalMap(recod)
	if err != nil {
		return err
	}

	itemImput := &dynamodb.PutItemInput{
		TableName: aws.String(d.table),
		Item:      item,
	}

	_, err = d.service.PutItem(itemImput)
	if err != nil {
		return err
	}

	return nil
}
