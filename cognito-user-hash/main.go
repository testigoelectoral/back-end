package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

var (
	sedKey         string
	bytes          = []byte{193, 188, 236, 45, 124, 242, 218, 180, 143, 130, 226, 135, 210, 120, 10, 70}
	cognitoService userUpdater
)

func Encrypt(text string, secretSeed string) (string, error) {
	block, err := aes.NewCipher([]byte(secretSeed))
	if err != nil {
		log.Printf("Err: %v", err.Error())
		return "", err
	}

	plainText := []byte(text)
	cipherText := make([]byte, len(plainText))

	cfb := cipher.NewCFBEncrypter(block, bytes)
	cfb.XORKeyStream(cipherText, plainText)
	cipherPlain := base64.StdEncoding.EncodeToString(cipherText)

	return cipherPlain, nil
}

func handler(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	toHash := event.Request.UserAttributes["custom:document"] + " " + event.Request.UserAttributes["name"] + " " + event.Request.UserAttributes["email"] + " " + event.Request.UserAttributes["phone_number"]

	cipherUsername, err := Encrypt(toHash, sedKey)
	if err != nil {
		return event, err
	}

	userAttributes := map[string]string{
		"custom:hash": cipherUsername,
	}

	err = cognitoService.UpdateUser(event.UserName, event.UserPoolID, userAttributes)
	if err != nil {
		return event, err
	}

	event.Request.UserAttributes["custom:hash"] = cipherUsername

	return event, nil
}

type userUpdater interface {
	UpdateUser(string, string, map[string]string) error
}

type cognitoUpdater struct {
	service cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

func (u *cognitoUpdater) UpdateUser(userName string, poolId string, attributes map[string]string) error { //nolint:revive
	changedAttributes := []*cognitoidentityprovider.AttributeType{}
	for name, value := range attributes {
		changedAttributes = append(changedAttributes, &cognitoidentityprovider.AttributeType{
			Name:  aws.String(name),
			Value: aws.String(value),
		})
	}

	parameters := &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserPoolId:     &poolId,
		UserAttributes: changedAttributes,
		Username:       &userName,
	}

	_, err := u.service.AdminUpdateUserAttributes(parameters)

	return err
}

func init() {
	sedKey = os.Getenv("SEED_KEY")

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	cognitoService = &cognitoUpdater{
		service: cognitoidentityprovider.New(sess),
	}
}

func main() {
	lambda.Start(handler)
}
