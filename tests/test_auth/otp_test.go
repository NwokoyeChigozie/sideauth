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

func TestSendOTP(t *testing.T) {
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
	token, accountID := tst.GetLoginTokenAndAccountID(t, r, auth, loginData)
	accessToken := tst.GetAccessToken(accountID, db.Auth)

	tests := []struct {
		Name         string
		RequestBody  models.SendOtpTokenReq
		ExpectedCode int
		Path         string
		Headers      map[string]string
		Message      string
	}{
		{
			Name:         "OK otp with authentication token",
			RequestBody:  models.SendOtpTokenReq{},
			ExpectedCode: http.StatusOK,
			Message:      "OTP Generated",
			Path:         "/v2/auth/send_otp",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name:         "otp without authentication token",
			RequestBody:  models.SendOtpTokenReq{},
			ExpectedCode: http.StatusUnauthorized,
			Path:         "/v2/auth/send_otp",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, {
			Name:         "OK otp with access token",
			RequestBody:  models.SendOtpTokenReq{AccountID: accountID},
			ExpectedCode: http.StatusOK,
			Message:      "OTP Generated",
			Path:         "/v2/auth/api/send_otp",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": accessToken.PrivateKey,
				"v-public-key":  accessToken.PublicKey,
			},
		}, {
			Name:         "OK otp without access token",
			RequestBody:  models.SendOtpTokenReq{AccountID: accountID},
			ExpectedCode: http.StatusOK,
			Message:      "OTP Generated",
			Path:         "/v2/auth/otp/send_otp",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, {
			Name:         "no account_id for otp without access token",
			RequestBody:  models.SendOtpTokenReq{},
			ExpectedCode: http.StatusBadRequest,
			Path:         "/v2/auth/otp/send_otp",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, {
			Name:         "otp without access token",
			RequestBody:  models.SendOtpTokenReq{AccountID: accountID},
			ExpectedCode: http.StatusUnauthorized,
			Path:         "/v2/auth/api/send_otp",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, {
			Name:         "otp without account id",
			RequestBody:  models.SendOtpTokenReq{},
			ExpectedCode: http.StatusBadRequest,
			Path:         "/v2/auth/api/send_otp",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": accessToken.PrivateKey,
				"v-public-key":  accessToken.PublicKey,
			},
		},
	}

	authUrl := r.Group(fmt.Sprintf("%v/auth", "v2"))
	{
		authUrl.POST("/otp/send_otp", auth.SendOTPAPI)

	}

	authTypeUrl := r.Group(fmt.Sprintf("%v/auth", "v2"), middleware.Authorize(db, middleware.AuthType))
	{
		authTypeUrl.POST("/send_otp", auth.SendOTP)

	}
	authApiUrl := r.Group(fmt.Sprintf("%v/auth/api", "v2"), middleware.Authorize(db, middleware.ApiType))
	{
		authApiUrl.POST("/send_otp", auth.SendOTPAPI)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: test.Path}

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
