package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var tagName = "validate"
var mailRegx = regexp.MustCompile(`\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`)

type ValidationErr struct {
	Field string
	Msg   string
}

func (v ValidationErr) Error() string {
	return fmt.Sprintf("%s %s", v.Field, v.Msg)
}

type Validator interface {
	Validate(any) (bool, error)
}
type defaultValidator struct{}

func (v defaultValidator) Validate(val any) (bool, error) {
	return true, nil
}

type requiredValidator struct{}

func (v requiredValidator) Validate(val any) (bool, error) {
	data := val.(string)

	if data == "" {
		return false, errors.New("field is required.")
	}
	return true, nil
}

type emailValidator struct{}

func (v emailValidator) Validate(val any) (bool, error) {
	data := val.(string)

	if !mailRegx.Match([]byte(data)) {
		return false, errors.New("field is not a valid email.")
	}
	return true, nil
}

func getValidatorFromTag(tag string) []Validator {
	args := strings.Split(tag, ",")
	validators := []Validator{}

	for _, item := range args {
		switch item {
		case "required":
			validators = append(validators, requiredValidator{})
		case "email":
			validators = append(validators, emailValidator{})
		}
	}
	if len(validators) == 0 {
		return append(validators, defaultValidator{})
	}
	return validators
}

func ValidateStruct(s any) []error {
	errs := []error{}

	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}

		validators := getValidatorFromTag(tag)
		for _, item := range validators {
			isValid, err := item.Validate(v.Field(i).Interface())
			if !isValid && err != nil {
				jsonName := v.Type().Field(i).Tag.Get("json")
				// fieldName := v.Type().Field(i).Name

				err := ValidationErr{
					Field: strings.Split(jsonName, ",")[0],
					Msg:   err.Error(),
				}
				errs = append(errs, err)
			}
		}

	}
	return errs
}