package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	sedKey string
	bytes  = []byte{193, 188, 236, 45, 124, 242, 218, 180, 143, 130, 226, 135, 210, 120, 10, 70}
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
	toHash := event.Request.UserAttributes["username"] + " " + event.Request.UserAttributes["name"] + " " + event.Request.UserAttributes["email"] + " " + event.Request.UserAttributes["phone_number"]

	cipherUsername, err := Encrypt(toHash, sedKey)
	if err != nil {
		return event, err
	}

	event.Request.UserAttributes["custom:hash"] = cipherUsername

	return event, nil
}

func init() {
	sedKey = os.Getenv("SEED_KEY")
}

func main() {
	lambda.Start(handler)
}
