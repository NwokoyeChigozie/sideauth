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
	"github.com/vesicash/auth-ms/pkg/middleware"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	tst "github.com/vesicash/auth-ms/tests"
	"github.com/vesicash/auth-ms/utility"
)

func TestUpgradeTier(t *testing.T) {
	logger := tst.Setup()
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
	token, _ := tst.GetLoginTokenAndAccountID(t, r, auth, loginData)

	type requestBody struct {
		Tier int `json:"tier"`
	}

	tests := []struct {
		Name         string
		RequestBody  requestBody
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK tier 1",
			RequestBody: requestBody{
				Tier: 1,
			},
			ExpectedCode: http.StatusOK,
			Message:      "Upgraded",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "OK tier 2",
			RequestBody: requestBody{
				Tier: 2,
			},
			ExpectedCode: http.StatusOK,
			Message:      "Upgraded",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "no tier specified",
			RequestBody:  requestBody{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "incorrect tier",
			RequestBody: requestBody{
				Tier: 6,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	authTypeUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, middleware.AuthType))
	{
		authTypeUrl.POST("/user/upgrade_tier", auth.UpgradeUserTier)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/user/upgrade_tier"}

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

func TestGetUserRestrictions(t *testing.T) {
	logger := tst.Setup()
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
	token, _ := tst.GetLoginTokenAndAccountID(t, r, auth, loginData)

	tests := []struct {
		Name         string
		RequestBody  interface{}
		ExpectedCode int
		Tier         int
		Headers      map[string]string
		Message      string
	}{
		{
			Name:         "OK get restrictions tier 0",
			RequestBody:  nil,
			ExpectedCode: http.StatusOK,
			Message:      "success",
			Tier:         0,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "OK get restrictions tier 1",
			RequestBody:  nil,
			ExpectedCode: http.StatusOK,
			Message:      "success",
			Tier:         1,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "OK get restrictions tier 2",
			RequestBody:  nil,
			ExpectedCode: http.StatusOK,
			Message:      "success",
			Tier:         2,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	authTypeUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, middleware.AuthType))
	{
		authTypeUrl.GET("/user/restrictions", auth.GetUserRestrictions)
		authTypeUrl.POST("/user/upgrade_tier", auth.UpgradeUserTier)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if test.Tier != 0 {
				type requestBody struct {
					Tier int `json:"tier"`
				}
				upgradeReq := requestBody{Tier: test.Tier}

				requestResetURI := url.URL{Path: "/v2/user/upgrade_tier"}
				var f bytes.Buffer
				json.NewEncoder(&f).Encode(upgradeReq)
				fReq, err := http.NewRequest(http.MethodPost, requestResetURI.String(), &f)
				if err != nil {
					t.Fatal(err)
				}
				fReq.Header.Set("Content-Type", "application/json")
				fReq.Header.Set("Authorization", "Bearer "+token)

				frr := httptest.NewRecorder()
				r.ServeHTTP(frr, fReq)
			}

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/user/restrictions"}

			req, err := http.NewRequest(http.MethodGet, URI.String(), &b)
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

func TestUpgradeAccount(t *testing.T) {
	logger := tst.Setup()
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
	token, _ := tst.GetLoginTokenAndAccountID(t, r, auth, loginData)

	type requestBody struct {
		BusinessType string `json:"business_type" validate:"required,oneof=ecommerce social_commerce marketplace"`
		BusinessName string `json:"business_name" validate:"required"`
		WebhookUri   string `json:"webhook_uri"`
	}

	tests := []struct {
		Name         string
		RequestBody  requestBody
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK upgrade account ecommerce",
			RequestBody: requestBody{
				BusinessType: "ecommerce",
				BusinessName: "yy",
				WebhookUri:   "http://link_to_webhook_uri",
			},
			ExpectedCode: http.StatusOK,
			Message:      "Upgraded",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "OK upgrade account social_commerce",
			RequestBody: requestBody{
				BusinessType: "social_commerce",
				BusinessName: "yy",
				WebhookUri:   "http://link_to_webhook_uri",
			},
			ExpectedCode: http.StatusOK,
			Message:      "Upgraded",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "OK upgrade account marketplace",
			RequestBody: requestBody{
				BusinessType: "marketplace",
				BusinessName: "yy",
				WebhookUri:   "http://link_to_webhook_uri",
			},
			ExpectedCode: http.StatusOK,
			Message:      "Upgraded",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "incorrect business type",
			RequestBody: requestBody{
				BusinessType: "business",
				BusinessName: "yy",
				WebhookUri:   "http://link_to_webhook_uri",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "no business type",
			RequestBody: requestBody{
				BusinessName: "yy",
				WebhookUri:   "http://link_to_webhook_uri",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "no business name",
			RequestBody: requestBody{
				BusinessType: "marketplace",
				WebhookUri:   "http://link_to_webhook_uri",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "no request body",
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	authTypeUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, middleware.AuthType))
	{
		authTypeUrl.POST("/user/upgrade/account", auth.UpgradeAccount)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/user/upgrade/account"}

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
