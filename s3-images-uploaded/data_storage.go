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

type PageMeta struct {
	LocationStateCode        uint8  // 34 States
	LocationMunicipalityCode uint8  // Max 125 Municipalities per state (Based on 2018 elections that was Antioquia)
	LocationZoneCode         uint8  // Max 36 Zones on same Municipality (Based on 2018 elections that was Cali)
	LocationPlace            uint16 // max 642 Places on same Municipality (Based on 2018 elections that was Bogota)
	LocationTable            uint16 // max 394 Tables on same Place (Based on 2018 elections that was Corferias)
	PageNumer                uint8  // Max 13 pages
	PageType                 uint8  // Senado(71)/Camara(72)/... etc
	PageQR                   string //
}

type ImageRecord struct {
	ImageID     string
	OwnerSub    string
	OwnerGPS    GPSRecord
	OwnerReport bool
	PageMeta    PageMeta
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
