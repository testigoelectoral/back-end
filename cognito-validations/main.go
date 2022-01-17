package main

import (
	"errors"
	"math"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	err := validateUsername(event.UserName)
	if err != nil {
		return event, err
	}

	return event, nil
}

func validateUsername(username string) error {
	intValue, err := strconv.Atoi(username)
	if err != nil {
		return errors.New("username es la cédula, y debe contener sólo numeros")
	}

	compareString := strconv.Itoa(int(math.Abs(float64(intValue))))
	if compareString != username {
		return errors.New("username es la cédula, y debe contener sólo numeros")
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
