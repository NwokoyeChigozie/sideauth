package test_auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/controller/auth"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	tst "github.com/vesicash/auth-ms/tests"
	"github.com/vesicash/auth-ms/utility"
)

func TestSignup(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	// getConfig := config.GetConfig()
	validatorRef := validator.New()
	db := postgresql.Connection()
	requestURI := url.URL{Path: "/v2/auth/signup"}
	iuuid, _ := uuid.NewV4()
	buuid, _ := uuid.NewV4()
	ouuid, _ := uuid.NewV4()
	puuid, _ := uuid.NewV4()

	tests := []struct {
		Name         string
		RequestBody  models.CreateUserRequestModel
		ExpectedCode int
		Message      string
	}{
		{
			Name: "OK individual",
			RequestBody: models.CreateUserRequestModel{
				EmailAddress: fmt.Sprintf("testuser%v@qa.team", iuuid.String()),
				PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
				AccountType:  "individual",
				Firstname:    "test",
				Lastname:     "user",
				Password:     "password",
				Country:      "nigeria",
				Username:     fmt.Sprintf("test_username%v", iuuid.String()),
			},
			ExpectedCode: http.StatusCreated,
			Message:      "signup successful",
		}, {
			Name: "OK business",
			RequestBody: models.CreateUserRequestModel{
				EmailAddress: fmt.Sprintf("testuser%v@qa.team", buuid.String()),
				PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
				AccountType:  "business",
				Firstname:    "test",
				Lastname:     "user",
				Password:     "password",
				Country:      "nigeria",
				Username:     fmt.Sprintf("test_username%v", buuid.String()),
			},
			ExpectedCode: http.StatusCreated,
			Message:      "signup successful",
		}, {
			Name: "OK others",
			RequestBody: models.CreateUserRequestModel{
				EmailAddress: fmt.Sprintf("testuser%v@qa.team", ouuid.String()),
				PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
				AccountType:  "others",
				Firstname:    "test",
				Lastname:     "user",
				Password:     "password",
				Country:      "nigeria",
				Username:     fmt.Sprintf("test_username%v", ouuid.String()),
			},
			ExpectedCode: http.StatusCreated,
			Message:      "signup successful",
		}, {
			Name: "details already exist",
			RequestBody: models.CreateUserRequestModel{
				EmailAddress: fmt.Sprintf("testuser%v@qa.team", iuuid.String()),
				PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
				AccountType:  "individual",
				Firstname:    "test",
				Lastname:     "user",
				Password:     "password",
				Country:      "nigeria",
				Username:     fmt.Sprintf("test_username%v", iuuid.String()),
			},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "invalid email",
			RequestBody: models.CreateUserRequestModel{
				EmailAddress: fmt.Sprintf("testus%v", iuuid.String()),
				PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
				AccountType:  "individual",
				Firstname:    "test",
				Lastname:     "user",
				Password:     "password",
				Country:      "nigeria",
			},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "invalid phone number",
			RequestBody: models.CreateUserRequestModel{
				EmailAddress: fmt.Sprintf("testuser%v@qa.team", puuid.String()),
				PhoneNumber:  "0009",
				AccountType:  "individual",
				Firstname:    "test",
				Lastname:     "user",
				Password:     "password",
				Country:      "nigeria",
				Username:     fmt.Sprintf("test_username%v", iuuid.String()),
			},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "invalid account type",
			RequestBody: models.CreateUserRequestModel{
				EmailAddress: fmt.Sprintf("testuser%v@qa.team", puuid.String()),
				PhoneNumber:  "0009",
				AccountType:  "check",
				Firstname:    "test",
				Lastname:     "user",
				Password:     "password",
				Country:      "nigeria",
				Username:     fmt.Sprintf("test_user_name%v", iuuid.String()),
			},
			ExpectedCode: http.StatusBadRequest,
		},
	}

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		r.POST("/v2/auth/signup", auth.Signup)

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tst.ParseResponse(rr)

			code := int(data["code"].(float64))
			tst.AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					tst.AssertResponseMessage(t, message.(string), test.Message)
				} else {
					tst.AssertResponseMessage(t, "", test.Message)
				}

			}

		})

	}

}

func TestBulkSignup(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	// getConfig := config.GetConfig()
	validatorRef := validator.New()
	db := postgresql.Connection()
	requestURI := url.URL{Path: "/v2/auth/signup/bulk"}
	iuuid, _ := uuid.NewV4()
	buuid, _ := uuid.NewV4()
	ouuid, _ := uuid.NewV4()
	// puuid, _ := uuid.NewV4()

	tests := []struct {
		Name         string
		RequestBody  models.BulkCreateUserRequestModel
		ExpectedCode int
		Message      string
	}{
		{
			Name: "OK individual",
			RequestBody: models.BulkCreateUserRequestModel{
				Bulk: []models.CreateUserRequestModel{
					{
						EmailAddress: fmt.Sprintf("testuser%v@qa.team", iuuid.String()),
						PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
						AccountType:  "individual",
						Firstname:    "test",
						Lastname:     "user",
						Password:     "password",
						Country:      "nigeria",
						Username:     fmt.Sprintf("test_username%v", iuuid.String()),
					},
					{
						EmailAddress: fmt.Sprintf("testuser1%v@qa.team", iuuid.String()),
						PhoneNumber:  fmt.Sprintf("+2349%v", utility.GetRandomNumbersInRange(700000000, 909999999)),
						AccountType:  "individual",
						Firstname:    "test",
						Lastname:     "user",
						Password:     "password",
						Country:      "nigeria",
						Username:     fmt.Sprintf("test_username1%v", iuuid.String()),
					},
				},
			},
			ExpectedCode: http.StatusCreated,
			Message:      "users signup successful",
		},
		{
			Name: "OK business",
			RequestBody: models.BulkCreateUserRequestModel{
				Bulk: []models.CreateUserRequestModel{
					{
						EmailAddress: fmt.Sprintf("testuser%v@qa.team", buuid.String()),
						PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
						AccountType:  "business",
						Firstname:    "test",
						Lastname:     "user",
						Password:     "password",
						Country:      "nigeria",
						Username:     fmt.Sprintf("test_username%v", buuid.String()),
					},
					{
						EmailAddress: fmt.Sprintf("testuser1%v@qa.team", buuid.String()),
						PhoneNumber:  fmt.Sprintf("+2349%v", utility.GetRandomNumbersInRange(700000000, 909999999)),
						AccountType:  "business",
						Firstname:    "test",
						Lastname:     "user",
						Password:     "password",
						Country:      "nigeria",
						Username:     fmt.Sprintf("test_username1%v", buuid.String()),
					},
				},
			},
			ExpectedCode: http.StatusCreated,
			Message:      "users signup successful",
		},
		{
			Name: "OK others",
			RequestBody: models.BulkCreateUserRequestModel{
				Bulk: []models.CreateUserRequestModel{
					{
						EmailAddress: fmt.Sprintf("testuser%v@qa.team", ouuid.String()),
						PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
						AccountType:  "others",
						Firstname:    "test",
						Lastname:     "user",
						Password:     "password",
						Country:      "nigeria",
						Username:     fmt.Sprintf("test_username%v", ouuid.String()),
					},
					{
						EmailAddress: fmt.Sprintf("testuser1%v@qa.team", ouuid.String()),
						PhoneNumber:  fmt.Sprintf("+2349%v", utility.GetRandomNumbersInRange(700000000, 909999999)),
						AccountType:  "others",
						Firstname:    "test",
						Lastname:     "user",
						Password:     "password",
						Country:      "nigeria",
						Username:     fmt.Sprintf("test_username1%v", ouuid.String()),
					},
				},
			},
			ExpectedCode: http.StatusCreated,
			Message:      "users signup successful",
		},
	}

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		r.POST("/v2/auth/signup/bulk", auth.BulkSignup)

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tst.ParseResponse(rr)

			code := int(data["code"].(float64))
			tst.AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					tst.AssertResponseMessage(t, message.(string), test.Message)
				} else {
					tst.AssertResponseMessage(t, "", test.Message)
				}

			}

		})

	}

}
