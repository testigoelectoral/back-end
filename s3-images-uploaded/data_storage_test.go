package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/require"
)

type dynamodbifaceFake struct {
	dynamodbiface.DynamoDBAPI
}

func (d *dynamodbifaceFake) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	if aws.StringValue(input.Item["ImageID"].S) == "ERROR" {
		return nil, errors.New("DYNAMODB ERROR")
	}

	return &dynamodb.PutItemOutput{}, nil
}

func init() {}

func Test_Save(t *testing.T) {
	c := require.New(t)

	dynamodbstorageTest := &DynamoDBStorage{service: &dynamodbifaceFake{}, table: "table-name"}

	err := dynamodbstorageTest.Save(ImageRecord{})
	c.Nil(err)

	err = dynamodbstorageTest.Save(ImageRecord{ImageID: "ERROR"})
	c.NotNil(err)
	c.Equal("DYNAMODB ERROR", err.Error())
}
