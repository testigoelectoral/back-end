package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	var errs []string

	err := validateUsername(event.UserName)
	if err != nil {
		errs = append(errs, err.Error())
	}

	err = validatePhoneNumber(event.Request.UserAttributes["phone_number"])
	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return event, errors.New(strings.Join(errs, "\n"))
	}

	return event, nil
}

func validateUsername(username string) error {
	errorMessage := fmt.Sprintf("username(Cédula):'%s' inválido. Debe ser el número de la cédula, sólo numeros", username)

	intValue, err := strconv.Atoi(username)
	if err != nil {
		return errors.New(errorMessage)
	}

	compareString := strconv.Itoa(int(math.Abs(float64(intValue))))
	if compareString != username {
		return errors.New(errorMessage)
	}

	return nil
}

func validatePhoneNumber(PhoneNumber string) error {
	errorMessage := fmt.Sprintf("Phone Number(Celular):'%s' inválido. Formato: '+############' Símbolo de suma (+) y sólo números incluyendo el código de pais. ej: +573211234567", PhoneNumber)

	intValue, err := strconv.Atoi(PhoneNumber)
	if err != nil {
		return errors.New(errorMessage)
	}

	compareString := "+" + strconv.Itoa(int(math.Abs(float64(intValue))))
	if compareString != PhoneNumber {
		return errors.New(errorMessage)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
