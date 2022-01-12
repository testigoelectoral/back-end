package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
)

type mockCognitoUpdater struct{}

func (u *mockCognitoUpdater) UpdateUser(string, string, map[string]string) error { //nolint:revive
	return nil
}

func init() {
	sedKey = "2ed5ca7088920f2b11a1f5e2f1ad9d90"
	cognitoService = &mockCognitoUpdater{}
}

func eventRequest() events.CognitoEventUserPoolsPostConfirmation {
	event := events.CognitoEventUserPoolsPostConfirmation{
		Request: events.CognitoEventUserPoolsPostConfirmationRequest{
			UserAttributes: map[string]string{
				"username":     "10234567890",                                                                                                                      // Max 11 chars
				"phone_number": "+12345678901234",                                                                                                                  // Max 15 chars
				"name":         "Adolph Blaine Charles David Earl Frederick Gerald Hubert Irvin John Kenneth Lloyd Martin Nero Oliver Paul Quincy Randolph Sherma", // Max 128 chars
				"email":        "contact-admins-hello-webmaster-services@longest-email-address-known.is-such-a-long-sub-domain-it-could-go-on-forever.pacraig.com", // Max 128 chars
			},
			ClientMetadata: map[string]string{},
		},
		Response: events.CognitoEventUserPoolsPostConfirmationResponse{},
	}

	return event
}
func Test_uploadHandler(t *testing.T) {
	c := require.New(t)

	response, err := handler(eventRequest())
	c.Nil(err)
	c.NotNil(response)

	userData := response.Request.UserAttributes
	original := eventRequest()

	c.Equal(original.Request.UserAttributes["username"], userData["username"])
	c.Equal(original.Request.UserAttributes["phone_number"], userData["phone_number"])
	c.Equal(original.Request.UserAttributes["name"], userData["name"])
	c.Equal(original.Request.UserAttributes["email"], userData["email"])
	c.Equal("aH8XI2fWZoZ/DWIza9VRa+pIdv7w6yQisfelqeS0POoUerpCZWle2JfAtKcdnmR/4jp8fdBpPVKLG6r0oGXEFPj1fzjclQK8+MBeDrX4rP/Z2Fhsh6uUgR/sp1qqy/Bwn6M+LGxiv0tuhwBEuS4nQQmwy6s1GChp1uQiGYKYcyI8HaqG/4WQDMP0SLIBgz+AgvVi88i+aSQ6A05w5rbSArfhRVDrJDBuGestBNu2spgZdXYPbK4352lCRwXh2nqfhV90fxvd2DRGqsRoiTK02T+j9HEXyweyGR8MKeloK4eAmCLXOjKJa7k5356iJb0nx6XWNMGKNAga6QNosrupjBXaNt5JCVHgpkdakUo5Dl9XbCCAyFp2zR6qcJfj", userData["custom:hash"])
}
