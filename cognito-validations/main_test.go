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
		c.Equal(fmt.Sprintf("username(Cédula):'%s' inválido. Debe ser el número de la cédula, sólo numeros", badUser), err.Error())
	}

	response, err = handler(eventRequest("1234567890", "+573211234567", "user@example.com"))
	c.Nil(err)
	c.NotNil(response)
}

func Test_phoneNumber(t *testing.T) {
	c := require.New(t)

	var response events.CognitoEventUserPoolsPreSignup
	var err error

	response, err = handler(eventRequest("1234567890", "+573211234567", "user@example.com"))
	c.Nil(err)
	c.NotNil(response)

	badPhoneNumbers := []string{"faketext", "1234567890", "-1234567890", "15.34", "13 45"}

	for _, badNumber := range badPhoneNumbers {
		response, err = handler(eventRequest("1234567890", badNumber, "user@example.com"))
		c.NotNil(err)
		c.NotNil(response)
		c.Equal(fmt.Sprintf("Phone Number(Celular):'%s' inválido. Formato: '+############' Símbolo de suma (+) y sólo números incluyendo el código de pais. ej: +573211234567", badNumber), err.Error())
	}

	response, err = handler(eventRequest("1234567890", "+573211234567", "user@example.com"))
	c.Nil(err)
	c.NotNil(response)
}

func Test_multiple(t *testing.T) {
	c := require.New(t)

	var response events.CognitoEventUserPoolsPreSignup
	var err error

	response, err = handler(eventRequest("1234567890", "+573211234567", "user@example.com"))
	c.Nil(err)
	c.NotNil(response)

	response, err = handler(eventRequest("badUser", "badNumber", "user@example.com"))
	c.NotNil(err)
	c.NotNil(response)
	c.Equal("username(Cédula):'badUser' inválido. Debe ser el número de la cédula, sólo numeros\nPhone Number(Celular):'badNumber' inválido. Formato: '+############' Símbolo de suma (+) y sólo números incluyendo el código de pais. ej: +573211234567", err.Error())

	response, err = handler(eventRequest("1234567890", "+573211234567", "user@example.com"))
	c.Nil(err)
	c.NotNil(response)
}
