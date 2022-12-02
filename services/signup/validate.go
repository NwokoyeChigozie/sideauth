package signup

import (
	"fmt"
	"strings"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
)

func ValidateSignupRequest(req models.CreateUserRequestModel, dbs postgresql.Databases) (models.CreateUserRequestModel, error) {
	if req.PhoneNumber == "" && req.EmailAddress == "" {
		return req, fmt.Errorf("Please provide at-least one input for either e-mail address or phone number.")
	}

	user := models.User{}

	if req.EmailAddress != "" {
		req.EmailAddress = strings.ToLower(req.EmailAddress)
		if !utility.EmailValid(req.EmailAddress) {
			return req, fmt.Errorf("email address is invalid")
		}
		_, SErr := postgresql.SelectOneFromDb(dbs.Auth, &user, "email_address = ?", req.EmailAddress)
		if SErr == nil {
			return req, fmt.Errorf("email address is already in use")
		}
	}

	if req.PhoneNumber != "" {
		phone, status := utility.PhoneValid(req.PhoneNumber)
		if !status {
			return req, fmt.Errorf("phone number is invalid")
		}
		req.PhoneNumber = phone
		_, SErr := postgresql.SelectOneFromDb(dbs.Auth, &user, "phone_number = ?", req.PhoneNumber)
		if SErr == nil {
			return req, fmt.Errorf("phone number is already in use")
		}
	}

	return req, nil
}
