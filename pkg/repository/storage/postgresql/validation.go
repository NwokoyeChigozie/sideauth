package postgresql

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/vesicash/auth-ms/utility"
)

type Validation interface {
	LogModelData(*utility.Logger)
}

type ValidationError struct {
	Field string
	Error string
}

var (
	regexpEmail      = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	ErrEmptyField    = errors.New("Field cannot be empty")
	ErrInvalidEmail  = errors.New("Email is invalid")
	ErrInvalidPass   = errors.New("Field length should be greater than 8")
	ErrNil           = errors.New("Nil")
	ValidationNeeded = "Input validation failed on some fields"
)

func ValidateRequest(V interface{}) error {

	var err []ValidationError
	if reflect.ValueOf(V).Kind() == reflect.Struct {
		t := reflect.TypeOf(V)
		v := reflect.ValueOf(V)

		for i := 0; i < t.NumField(); i++ {
			FieldT := t.Field(i)
			FieldV := v.Field(i)
			// reflect.ValueOf(V).Field(i).Type()

			validateFields := FieldT.Tag.Get("pgvalidate")
			splitFields := strings.Split(validateFields, ",")
			if validateFields == "_" || validateFields == "" {
				continue
			}

			for j := 0; j < len(splitFields); j++ {
				splitFieldsStr := strings.ToLower(splitFields[j])
				if strings.Contains(splitFieldsStr, "notexists") {
					value, status := ValidateNext(FieldV)
					if status {
						firstSplit := strings.Split(splitFieldsStr, "=")
						if len(firstSplit) == 2 {
							secondSplit := strings.Split(firstSplit[1], "$")
							if len(secondSplit) == 3 {
								db := ReturnDatabase(secondSplit[0])
								tableName := secondSplit[1]
								columnName := secondSplit[2]
								if CheckExistsInTable(db, tableName, fmt.Sprintf("%v = ?", columnName), value) {
									err = append(err, ValidationError{
										Field: FieldT.Name,
										Error: fmt.Sprintf("%v exists in %v table", columnName, tableName),
									})
								}

							}

						}
					}
				} else if strings.Contains(splitFieldsStr, "exists") {
					value, status := ValidateNext(FieldV)
					if status {
						firstSplit := strings.Split(splitFieldsStr, "=")
						if len(firstSplit) == 2 {
							secondSplit := strings.Split(firstSplit[1], "$")
							if len(secondSplit) == 3 {
								db := ReturnDatabase(secondSplit[0])
								tableName := secondSplit[1]
								columnName := secondSplit[2]
								if !CheckExistsInTable(db, tableName, fmt.Sprintf("%v = ?", columnName), value) {
									err = append(err, ValidationError{
										Field: FieldT.Name,
										Error: fmt.Sprintf("%v does not exist in %v table", columnName, tableName),
									})
								}

							}

						}
					}
				} else if strings.Contains(splitFieldsStr, "email") {
					if FieldV.String() != "" {
						if !regexpEmail.Match([]byte(FieldV.String())) {
							err = append(err, ValidationError{
								Field: FieldT.Name,
								Error: ErrInvalidEmail.Error(),
							})
						}
					}
				}
			}
		}
	}

	errString := ""
	if len(err) < 1 {
		return nil
	} else {
		for _, v := range err {
			errString += v.Field + ": " + v.Error + " ;"
		}
	}
	return fmt.Errorf(errString)
}

func ValidateNext(value reflect.Value) (interface{}, bool) {
	if value.Type().Kind() == reflect.Int {
		return value.Int(), value.Int() != 0
	} else if value.Type().Kind() == reflect.Int8 {
		return value.Int(), value.Int() != 0
	} else if value.Type().Kind() == reflect.Int16 {
		return value.Int(), value.Int() != 0
	} else if value.Type().Kind() == reflect.Int32 {
		return value.Int(), value.Int() != 0
	} else if value.Type().Kind() == reflect.Int64 {
		return value.Int(), value.Int() != 0
	} else if value.Type().Kind() == reflect.Uint {
		return value.Uint(), value.Uint() != 0
	} else if value.Type().Kind() == reflect.Uint8 {
		return value.Int(), value.Uint() != 0
	} else if value.Type().Kind() == reflect.Uint16 {
		return value.Int(), value.Uint() != 0
	} else if value.Type().Kind() == reflect.Uint32 {
		return value.Int(), value.Uint() != 0
	} else if value.Type().Kind() == reflect.Uint64 {
		return value.Int(), value.Uint() != 0
	} else if value.Type().Kind() == reflect.Uintptr {
		return value.Int(), value.Uint() != 0
	} else if value.Type().Kind() == reflect.Float32 {
		return value.Float(), value.Float() != 0
	} else if value.Type().Kind() == reflect.Float64 {
		return value.Float(), value.Float() != 0
	} else if value.Type().Kind() == reflect.Bool {
		return value.Bool(), true
	} else if value.Type().Kind() == reflect.String {
		return value.String(), value.String() != ""
	} else {
		return value.String(), value.String() != ""
	}
}
