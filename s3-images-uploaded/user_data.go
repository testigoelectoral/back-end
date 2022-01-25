package main

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type UserDataInterface interface {
	GetHash(string) (string, error)
}

type CognitoData struct {
	service cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

func NewCognitoData(sess *session.Session) *CognitoData {
	return &CognitoData{
		service: cognitoidentityprovider.New(sess),
	}
}

func (u *CognitoData) GetHash(userSub string) (string, error) { //nolint:revive
	attr := []string{"custom:hash"}
	filter := &cognitoidentityprovider.ListUsersInput{
		UserPoolId:      aws.String(cognitoPoolID),
		Filter:          aws.String(fmt.Sprintf("sub = \"%s\"", userSub)),
		AttributesToGet: aws.StringSlice(attr),
	}

	result, err := u.service.ListUsers(filter)
	if err != nil {
		return "", err
	}

	if len(result.Users) != 1 {
		return "", errors.New("User not found or repeated sub")
	}

	return string(*result.Users[0].Attributes[0].Value), nil
}
