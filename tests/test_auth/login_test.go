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

func TestLogin(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		loginPath      = "/v2/login"
		loginURI       = url.URL{Path: loginPath}
		muuid, _       = uuid.NewV4()
		duuid, _       = uuid.NewV4()
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

	tests := []struct {
		Name         string
		RequestBody  models.LoginUserRequestModel
		ExpectedCode int
		Message      string
	}{
		{
			Name: "OK username login successful",
			RequestBody: models.LoginUserRequestModel{
				Username: userSignUpData.Username,
				Password: userSignUpData.Password,
			},
			ExpectedCode: http.StatusOK,
			Message:      "login successful",
		}, {
			Name: "OK email login successful",
			RequestBody: models.LoginUserRequestModel{
				EmailAddress: userSignUpData.EmailAddress,
				Password:     userSignUpData.Password,
			},
			ExpectedCode: http.StatusOK,
			Message:      "login successful",
		}, {
			Name: "OK phone login successful",
			RequestBody: models.LoginUserRequestModel{
				PhoneNumber: userSignUpData.PhoneNumber,
				Password:    userSignUpData.Password,
			},
			ExpectedCode: http.StatusOK,
			Message:      "login successful",
		}, {
			Name:         "password not provided",
			RequestBody:  models.LoginUserRequestModel{},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "username or phone or email not provided",
			RequestBody: models.LoginUserRequestModel{
				Password: userSignUpData.Password,
			},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "username does not exist",
			RequestBody: models.LoginUserRequestModel{
				Username: duuid.String(),
				Password: userSignUpData.Password,
			},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "email does not exist",
			RequestBody: models.LoginUserRequestModel{
				EmailAddress: duuid.String(),
				Password:     userSignUpData.Password,
			},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "phone does not exist",
			RequestBody: models.LoginUserRequestModel{
				PhoneNumber: duuid.String(),
				Password:    userSignUpData.Password,
			},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "incorrect password",
			RequestBody: models.LoginUserRequestModel{
				Username: userSignUpData.Username,
				Password: "incorrect",
			},
			ExpectedCode: http.StatusBadRequest,
		},
	}

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	r.POST(loginPath, auth.Login)

	tst.SignupUser(t, r, auth, userSignUpData)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, loginURI.String(), &b)
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

func TestLoginPhone(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		loginPath      = "/v2/login-phone"
		loginURI       = url.URL{Path: loginPath}
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

	type requestStruct struct {
		PhoneNumber string `json:"phone_number"`
	}

	tests := []struct {
		Name         string
		RequestBody  requestStruct
		ExpectedCode int
		Message      string
	}{
		{
			Name: "OK phone login successful",
			RequestBody: requestStruct{
				PhoneNumber: userSignUpData.PhoneNumber,
			},
			ExpectedCode: http.StatusOK,
			Message:      "login successful",
		}, {
			Name:         "phone not provided",
			RequestBody:  requestStruct{},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "phone number does not exist",
			RequestBody: requestStruct{
				PhoneNumber: fmt.Sprintf("+233%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
			},
			ExpectedCode: http.StatusBadRequest,
		},
	}

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	r.POST(loginPath, auth.PhoneOtpLogin)

	tst.SignupUser(t, r, auth, userSignUpData)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, loginURI.String(), &b)
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

func TestLogout(t *testing.T) {
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
		Headers      map[string]string
		Message      string
	}{
		{
			Name:         "OK logout",
			RequestBody:  nil,
			ExpectedCode: http.StatusOK,
			Message:      "logout successful",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "OK wrong token",
			RequestBody:  nil,
			ExpectedCode: http.StatusUnauthorized,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	authTypeUrl := r.Group(fmt.Sprintf("%v/auth", "v2"), middleware.Authorize(db, middleware.AuthType))
	{
		authTypeUrl.POST("/logout", auth.Logout)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/logout"}

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

func TestValidateToken(t *testing.T) {
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
		Headers      map[string]string
		Message      string
	}{
		{
			Name:         "OK validate token",
			RequestBody:  nil,
			ExpectedCode: http.StatusOK,
			Message:      "token valid",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	authTypeUrl := r.Group(fmt.Sprintf("%v/auth", "v2"), middleware.Authorize(db, middleware.AuthType))
	{
		authTypeUrl.POST("/validate-token", auth.ValidateToken)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/validate-token"}

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

func TestGetAccessToken(t *testing.T) {
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
		Headers      map[string]string
		Message      string
	}{
		{
			Name:         "OK get access token",
			RequestBody:  nil,
			ExpectedCode: http.StatusOK,
			Message:      "success",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	authTypeUrl := r.Group(fmt.Sprintf("%v/auth", "v2"), middleware.Authorize(db, middleware.AuthType))
	{
		authTypeUrl.POST("/user/security/get_access_token", auth.GetAccessToken)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/user/security/get_access_token"}

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
