package test_auth_models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/vesicash/auth-ms/internal/config"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/controller/auth"
	"github.com/vesicash/auth-ms/pkg/controller/auth_model"
	"github.com/vesicash/auth-ms/pkg/middleware"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	tst "github.com/vesicash/auth-ms/tests"
	"github.com/vesicash/auth-ms/utility"
)

func TestGetBank(t *testing.T) {
	logger := tst.Setup()
	app := config.GetConfig().App
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _       = uuid.NewV4()
		userSignUpData = models.CreateUserRequestModel{
			EmailAddress: fmt.Sprintf("testuser%v@qa.team", muuid.String()),
			PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
			AccountType:  "individual",
			Firstname:    "test",
			Lastname:     "user",
			Password:     "password",
			Country:      "nigeria",
			Username:     fmt.Sprintf("test_username%v", muuid.String()),
		}
		loginData = models.LoginUserRequestModel{
			Username:     userSignUpData.Username,
			EmailAddress: userSignUpData.EmailAddress,
			PhoneNumber:  userSignUpData.PhoneNumber,
			Password:     userSignUpData.Password,
		}
	)

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	tst.SignupUser(t, r, auth, userSignUpData)
	_, accountID := tst.GetLoginTokenAndAccountID(t, r, auth, loginData)
	us := models.User{AccountID: uint(accountID)}
	_, err := us.GetUserByAccountID(db.Auth)
	if err != nil {
		log.Panic(err.Error())
	}

	bank := models.Bank{
		Name:    "Test Bank",
		Code:    "899",
		Country: "Nigeria",
	}
	code, err := bank.GetBankByQuery(db.Auth)
	if err != nil {
		if code == http.StatusInternalServerError {
			log.Panic(err.Error())
		}
		_, err = bank.CreateBank(db.Auth)
		if err != nil {
			log.Panic(err.Error())
		}
	}

	tests := []struct {
		Name         string
		RequestBody  models.GetBankRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK get bank by id",
			RequestBody: models.GetBankRequest{
				ID: bank.ID,
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "OK get bank",
			RequestBody: models.GetBankRequest{
				ID:      bank.ID,
				Name:    bank.Name,
				Code:    bank.Code,
				Country: bank.Country,
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
	}

	auth_model := auth_model.Controller{Db: db, Validator: validatorRef}

	authTypeUrl := r.Group(fmt.Sprintf("%v/auth", "v2"), middleware.Authorize(db, middleware.AppType))
	{
		authTypeUrl.POST("/get_bank", auth_model.GetBank)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/auth/get_bank"}

			req, err := http.NewRequest(http.MethodPost, URI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}

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
