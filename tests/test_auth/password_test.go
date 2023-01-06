package test_auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

func TestRequestPasswordReset(t *testing.T) {
	tst.Setup()
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
	)

	type requestBody struct {
		EmailAddress string `json:"email_address"`
		PhoneNumber  string `json:"phone_number"`
	}

	auth := auth.Controller{Db: db, Validator: validatorRef}
	r := gin.Default()
	tst.SignupUser(t, r, auth, userSignUpData)

	tests := []struct {
		Name         string
		RequestBody  requestBody
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK with email and phone",
			RequestBody: requestBody{
				PhoneNumber:  userSignUpData.PhoneNumber,
				EmailAddress: userSignUpData.EmailAddress,
			},
			ExpectedCode: http.StatusOK,
			Message:      "Request Sent",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "OK with phone",
			RequestBody: requestBody{
				PhoneNumber: userSignUpData.PhoneNumber,
			},
			ExpectedCode: http.StatusOK,
			Message:      "Request Sent",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "OK with email",
			RequestBody: requestBody{
				EmailAddress: userSignUpData.EmailAddress,
			},
			ExpectedCode: http.StatusOK,
			Message:      "Request Sent",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name:         "no email or phone",
			RequestBody:  requestBody{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	authUrl := r.Group(fmt.Sprintf("%v/auth", "v2"))
	{
		authUrl.POST("/reset-password", auth.RequestPasswordReset)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/auth/reset-password"}

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

func TestUpdatePasswordWithToken(t *testing.T) {
	tst.Setup()
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
		resetReq = struct {
			EmailAddress string `json:"email_address"`
			PhoneNumber  string `json:"phone_number"`
		}{
			PhoneNumber:  userSignUpData.PhoneNumber,
			EmailAddress: userSignUpData.EmailAddress,
		}
	)

	type requestBody struct {
		AccountID int    `json:"account_id"`
		Token     int    `json:"token"`
		Password  string `json:"password"`
	}

	auth := auth.Controller{Db: db, Validator: validatorRef}
	r := gin.Default()
	tst.SignupUser(t, r, auth, userSignUpData)
	_, accountID := tst.GetLoginTokenAndAccountID(t, r, auth, loginData)

	tests := []struct {
		Name         string
		RequestBody  requestBody
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK password update",
			RequestBody: requestBody{
				AccountID: accountID,
				Password:  "new_password",
			},
			ExpectedCode: http.StatusOK,
			Message:      "Password Updated",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "incorrect account id",
			RequestBody: requestBody{
				AccountID: utility.GetRandomNumbersInRange(700, 9099),
				Password:  "new_password",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "incorrect no password",
			RequestBody: requestBody{
				AccountID: accountID,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "incorrect token",
			RequestBody: requestBody{
				AccountID: accountID,
				Password:  "new_password",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name:         "no request body",
			RequestBody:  requestBody{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	authUrl := r.Group(fmt.Sprintf("%v/auth", "v2"))
	{
		authUrl.POST("/reset-password", auth.RequestPasswordReset)
		authUrl.POST("/reset-password/change-password", auth.UpdatePasswordWithToken)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			token := 0
			if test.Name != "incorrect token" && test.Name != "no request body" {
				requestResetURI := url.URL{Path: "/v2/auth/reset-password"}
				var b bytes.Buffer
				json.NewEncoder(&b).Encode(resetReq)
				req, err := http.NewRequest(http.MethodPost, requestResetURI.String(), &b)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")

				rr := httptest.NewRecorder()
				r.ServeHTTP(rr, req)
				prToken := models.PasswordResetToken{AccountID: accountID}
				prToken.GetLatestByAccountID(db.Auth)
				tokenInt, _ := strconv.Atoi(prToken.Token)
				token = tokenInt

			}
			requestB := test.RequestBody
			requestB.Token = token
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(requestB)
			URI := url.URL{Path: "/v2/auth/reset-password/change-password"}

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

func TestUpdatePassword(t *testing.T) {
	tst.Setup()
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

	type requestBody struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required"`
	}

	auth := auth.Controller{Db: db, Validator: validatorRef}
	r := gin.Default()
	tst.SignupUser(t, r, auth, userSignUpData)
	token, _ := tst.GetLoginTokenAndAccountID(t, r, auth, loginData)

	tests := []struct {
		Name         string
		RequestBody  requestBody
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK password update",
			RequestBody: requestBody{
				OldPassword: userSignUpData.Password,
				NewPassword: "new_password",
			},
			ExpectedCode: http.StatusOK,
			Message:      "Password Updated",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "incorrect password",
			RequestBody: requestBody{
				OldPassword: "old_password",
				NewPassword: "new_password",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "no input provided",
			RequestBody:  requestBody{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "new password not provided",
			RequestBody: requestBody{
				OldPassword: "old_password",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "old password not provided",
			RequestBody: requestBody{
				NewPassword: "new_password",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	authUrl := r.Group(fmt.Sprintf("%v/auth", "v2"), middleware.Authorize(db, middleware.AuthType))
	{
		authUrl.POST("/user/security/update_password", auth.UpdatePassword)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/auth/user/security/update_password"}

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
