package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/stretchr/testify/require"
)

type cognitoifaceFake struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

func (c *cognitoifaceFake) ListUsers(input *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
	if aws.StringValue(input.Filter) == "sub = \"notfound\"" {
		return &cognitoidentityprovider.ListUsersOutput{Users: []*cognitoidentityprovider.UserType{}}, nil
	}

	if aws.StringValue(input.Filter) == "sub = \"error\"" {
		return nil, errors.New("COGNITO ERROR")
	}

	output := &cognitoidentityprovider.ListUsersOutput{
		Users: []*cognitoidentityprovider.UserType{{
			Attributes: []*cognitoidentityprovider.AttributeType{{
				Name:  aws.String("custom:hash"),
				Value: aws.String("hash"),
			}},
		}},
	}

	return output, nil
}

func init() {
}

func Test_GetHash(t *testing.T) {
	c := require.New(t)

	cognitodataTest := &CognitoData{service: &cognitoifaceFake{}}

	result, err := cognitodataTest.GetHash("sub")
	c.Nil(err)
	c.NotNil(result)
	c.Equal("hash", result)

	result, err = cognitodataTest.GetHash("notfound")
	c.NotNil(err)
	c.Empty(result)
	c.Equal("User not found or repeated sub", err.Error())

	result, err = cognitodataTest.GetHash("error")
	c.NotNil(err)
	c.Empty(result)
	c.Equal("COGNITO ERROR", err.Error())
}
