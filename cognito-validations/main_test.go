package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
)

func eventRequest(userName string, phoneNumber string, email string) events.CognitoEventUserPoolsPreSignup {
	event := events.CognitoEventUserPoolsPreSignup{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			UserName: userName,
		},
		Request: events.CognitoEventUserPoolsPreSignupRequest{
			UserAttributes: map[string]string{
				"phone_number": phoneNumber,       // Max 15 chars
				"name":         "Nombre Apellido", // Max 128 chars
				"email":        email,             // Max 128 chars
			},
			ClientMetadata: map[string]string{},
		},
		Response: events.CognitoEventUserPoolsPreSignupResponse{},
	}

	return event
}
func Test_userName(t *testing.T) {
	c := require.New(t)

	var response events.CognitoEventUserPoolsPreSignup
	var err error

	response, err = handler(eventRequest("1234567890", "+573211234567", "user@example.com"))
	c.Nil(err)
	c.NotNil(response)

	badUsernames := []string{"faketext", "+1234567890", "-1234567890", "15.34", "13 45"}

	for _, badUser := range badUsernames {
		response, err = handler(eventRequest(badUser, "+573211234567", "user@example.com"))
		c.NotNil(err)
		c.NotNil(response)
		c.Equal(fmt.Sprintf("username(Cédula):%s invlaido. Debe ser el número de la cédula, y debe contener sólo numeros", badUser), err.Error())
	}

	response, err = handler(eventRequest("1234567890", "+573211234567", "user@example.com"))
	c.Nil(err)
	c.NotNil(response)
}
