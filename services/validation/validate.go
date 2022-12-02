package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
)

func Validate(req interface{}) (interface{}, interface{}, error) {
	validatorRef := validator.New()
	reqObj := req

	err := validatorRef.Struct(&reqObj)
	if err != nil {
		return reqObj, utility.ValidationResponse(err, validatorRef), err
	}

	err = postgresql.ValidateRequest(reqObj)
	if err != nil {
		return reqObj, err.Error(), err
	}
	return reqObj, nil, nil
}
