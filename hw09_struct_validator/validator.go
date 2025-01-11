package hw09structvalidator

import (
	"bytes"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var ErrInputNotStruct = errors.New("input not struct")

type Elem interface {
	int | string
}

type ValidationError struct {
	Field string

	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString("errors:")
	for i := 0; i < len(v); i++ {
		buffer.WriteString("\n")
		buffer.WriteString(v[i].Field + ":" + v[i].Err.Error())
	}

	return buffer.String()
}

func Validate(v interface{}) error {
	refValue := reflect.ValueOf(v)

	if refValue.Type().Kind() != reflect.Struct {
		return ErrInputNotStruct
	}

	validationErrors, err := validateRef(refValue)
	if err != nil {
		return err
	}

	if len(validationErrors) != 0 {
		return errors.New(validationErrors.Error())
	}

	return nil
}

func validateRef(refValue reflect.Value) (ValidationErrors, error) {
	validationErrors := make([]ValidationError, 0)

	refType := refValue.Type()

	for i := 0; i < refType.NumField(); i++ {
		currField := refType.Field(i)

		validateTag := currField.Tag.Get("validate")

		if validateTag == "" {
			continue
		}

		curFieldValue := refValue.Field(i)

		var newValidationErrors ValidationErrors

		var err error

		switch kind := curFieldValue.Type().Kind(); kind {
		case reflect.Struct:

			newValidationErrors, err = validateRef(curFieldValue)

		case reflect.Slice:

			newValidationErrors, err = validateSlice(getValidateTag(currField), curFieldValue, currField.Name)

		case reflect.String:

			newValidationErrors, err = validateSingleElem(getValidateTag(currField), curFieldValue.String(), currField.Name)

		case reflect.Int:

			newValidationErrors, err = validateSingleElem(getValidateTag(currField), curFieldValue.Int(), currField.Name)

		default:

			continue
		}

		if err != nil {
			return nil, err
		}

		if newValidationErrors != nil {
			validationErrors = append(validationErrors, newValidationErrors...)
		}
	}

	return validationErrors, nil
}

func getValidateTag(currField reflect.StructField) []string {
	validateTag := currField.Tag.Get("validate")

	if validateTag == "" {
		return nil
	}

	return strings.Split(validateTag, "|")
}

func validateSlice(validators []string, curFieldValue reflect.Value, sliceName string) (ValidationErrors, error) {
	elemKind := curFieldValue.Index(0).Kind()

	validationErrors := make([]ValidationError, 0)

	for i := 0; i < curFieldValue.Len(); i++ {
		curElemValue := curFieldValue.Index(i)

		var newValidationErrors ValidationErrors

		var err error

		if elemKind == reflect.String {
			//nolint:lll
			newValidationErrors, err = validateSingleElem(validators, curElemValue.String(), sliceName+"["+curElemValue.String()+"]")
		} else if elemKind == reflect.Int {
			//nolint:lll
			newValidationErrors, err = validateSingleElem(validators, curElemValue.String(), sliceName+"["+strconv.FormatInt(curElemValue.Int(), 10)+"]")
		}

		if err != nil {
			return nil, err
		}

		if newValidationErrors != nil {
			validationErrors = append(validationErrors, newValidationErrors...)
		}
	}

	return validationErrors, nil
}

func validateSingleElem(validators []string, elem any, elemName string) (ValidationErrors, error) {
	validationErrors := make([]ValidationError, 0)

	for i := 0; i < len(validators); i++ {
		validatorSeps := strings.Split(validators[i], ":")

		if len(validatorSeps) != 2 {
			return nil, errors.New("validator " + validators[i] + " is not in right format")
		}

		var notValid, err error

		switch v := elem.(type) {
		case int64:

			notValid, err = singleIntValidator(validatorSeps[0], validatorSeps[1], v)

		case string:

			notValid, err = singleStringValidator(validatorSeps[0], validatorSeps[1], v)
		}

		if err != nil {
			return nil, err
		}

		if notValid != nil {
			validationErrors = append(validationErrors, ValidationError{elemName, notValid})
		}
	}

	return validationErrors, nil
}

func singleStringValidator(validatorKey string, validatorValue string, s string) (notValid error, err error) {
	switch valType := validatorKey; valType {
	case "len":

		return validateStringLen(validatorValue, s)

	case "regexp":

		return validateStringRegExp(validatorValue, s)

	case "in":

		return validateStringIn(validatorValue, s), nil

	default:

		return nil, nil
	}
}

func validateStringLen(validatorValue string, s string) (notValid error, err error) {
	stringLen, err := strconv.Atoi(validatorValue)
	if err != nil {
		return nil, errors.New("String length " + validatorValue + "not int value")
	}

	if len(s) != stringLen {
		return errors.New(s + " len not " + validatorValue), nil
	}

	return nil, nil
}

func validateStringRegExp(regExp string, s string) (notValid error, err error) {
	matched, err := regexp.MatchString(regExp, s)
	if err != nil {
		return nil, errors.New("Error while matching for regular expression " + regExp)
	}

	if !matched {
		return errors.New(s + " not matches " + regExp), nil
	}

	return nil, nil
}

func validateStringIn(validatorValue string, s string) error {
	inValues := strings.Split(validatorValue, ",")

	var matched bool

	for i := 0; i < len(inValues); i++ {
		if s == inValues[i] {
			matched = true

			break
		}
	}

	if !matched {
		return errors.New(s + " not in " + validatorValue)
	}

	return nil
}

func singleIntValidator(validatorKey string, validatorValue string, i int64) (notValid error, err error) {
	switch valType := validatorKey; valType {
	case "min":

		return validateIntMin(validatorValue, i)

	case "max":

		return validateIntMax(validatorValue, i)

	case "in":

		return validateIntIn(validatorValue, i)

	default:

		return nil, nil
	}
}

func validateIntMin(validatorValue string, i int64) (notValid error, err error) {
	minValue, err := strconv.Atoi(validatorValue)
	if err != nil {
		return nil, errors.New("Min " + validatorValue + "not int value")
	}

	if i < int64(minValue) {
		return errors.New(strconv.FormatInt(i, 10) + " lesser than " + validatorValue), nil
	}

	return nil, nil
}

func validateIntMax(validatorValue string, i int64) (notValid error, err error) {
	maxValue, err := strconv.Atoi(validatorValue)
	if err != nil {
		return nil, errors.New("Max " + validatorValue + "not int value")
	}

	if i > int64(maxValue) {
		return errors.New(strconv.FormatInt(i, 10) + " greater than " + validatorValue), nil
	}

	return nil, nil
}

func validateIntIn(validatorValue string, i int64) (notValid error, err error) {
	inValues := strings.Split(validatorValue, ",")

	var matched bool

	for j := 0; j < len(inValues); j++ {
		intValue, err := strconv.Atoi(inValues[j])
		if err != nil {
			return nil, errors.New("Min " + validatorValue + "not int value")
		}

		if i == int64(intValue) {
			matched = true

			break
		}
	}

	if !matched {
		return errors.New(strconv.FormatInt(i, 10) + " not in " + validatorValue), nil
	}

	return nil, nil
}
