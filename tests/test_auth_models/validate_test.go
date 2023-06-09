package test_auth_models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

func TestValidateOnDb(t *testing.T) {
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

	tests := []struct {
		Name              string
		RequestBody       models.ValidateOnDBReq
		ExpectedCode      int
		Headers           map[string]string
		Message           string
		CheckResponseData bool
		Response          bool
	}{
		{
			Name: "OK exists validate on db with value",
			RequestBody: models.ValidateOnDBReq{
				Table: "users",
				Type:  "exists",
				Query: "email_address = ?",
				Value: userSignUpData.EmailAddress,
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
			CheckResponseData: true,
			Response:          true,
		}, {
			Name: "OK exists validate on db without value",
			RequestBody: models.ValidateOnDBReq{
				Table: "users",
				Type:  "exists",
				Query: "account_id = " + strconv.Itoa(accountID),
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
			CheckResponseData: true,
			Response:          true,
		}, {
			Name: "OK notexists validate on db with value",
			RequestBody: models.ValidateOnDBReq{
				Table: "users",
				Type:  "notexists",
				Query: "email_address = ?",
				Value: userSignUpData.EmailAddress + "not",
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
			CheckResponseData: true,
			Response:          true,
		}, {
			Name: "OK exists validate on db without value",
			RequestBody: models.ValidateOnDBReq{
				Table: "users",
				Type:  "notexists",
				Query: "account_id = " + strconv.Itoa(accountID),
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
			CheckResponseData: true,
			Response:          false,
		}, {
			Name: "table omitted",
			RequestBody: models.ValidateOnDBReq{
				Type:  "notexists",
				Query: "account_id = " + strconv.Itoa(accountID),
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "type omitted",
			RequestBody: models.ValidateOnDBReq{
				Table: "users",
				Query: "account_id = " + strconv.Itoa(accountID),
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "query omitted",
			RequestBody: models.ValidateOnDBReq{
				Table: "users",
				Type:  "notexists",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name:         "no input",
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
	}

	auth_model := auth_model.Controller{Db: db, Validator: validatorRef}

	authTypeUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, middleware.AppType))
	{
		authTypeUrl.POST("/validate_on_db", auth_model.ValidateOnDB)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/validate_on_db"}

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

			if test.CheckResponseData {
				resData := data["data"].(bool)
				tst.AssertBool(t, resData, test.Response)
			}

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

func TestValidateAuthorization(t *testing.T) {
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
	token, accountID := tst.GetLoginTokenAndAccountID(t, r, auth, loginData)
	us := models.User{AccountID: uint(accountID)}
	_, err := us.GetUserByAccountID(db.Auth)
	if err != nil {
		log.Panic(err.Error())
	}
	us.AccountType = "admin"
	err = us.Update(db.Auth)
	if err != nil {
		log.Panic(err.Error())
	}

	accessToken := models.AccessToken{AccountID: int(us.AccountID)}
	code, err := accessToken.GetByAccountID(db.Auth)
	if err != nil {
		if code == http.StatusInternalServerError {
			log.Panic(err.Error())
		}
		err = accessToken.CreateAccessToken(db.Auth)
		if err != nil {
			log.Panic(err.Error())
		}
	}

	tests := []struct {
		Name              string
		RequestBody       models.ValidateAuthorizationReq
		ExpectedCode      int
		Headers           map[string]string
		Message           string
		CheckResponseData bool
		Response          interface{}
	}{
		{
			Name: "OK auth type",
			RequestBody: models.ValidateAuthorizationReq{
				Type:               "auth",
				AuthorizationToken: token,
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
			CheckResponseData: true,
			Response:          us,
		}, {
			Name: "OK app type",
			RequestBody: models.ValidateAuthorizationReq{
				Type:               "app",
				AuthorizationToken: app.Key,
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "OK api type",
			RequestBody: models.ValidateAuthorizationReq{
				Type:        "api",
				VPrivateKey: accessToken.PrivateKey,
				VPublicKey:  accessToken.PublicKey,
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "OK business type",
			RequestBody: models.ValidateAuthorizationReq{
				Type:        "business",
				VPrivateKey: accessToken.PrivateKey,
				VPublicKey:  accessToken.PublicKey,
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "OK business_admin type",
			RequestBody: models.ValidateAuthorizationReq{
				Type:        "business_admin",
				VPrivateKey: accessToken.PrivateKey,
				VPublicKey:  accessToken.PublicKey,
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name:         "type not provided",
			RequestBody:  models.ValidateAuthorizationReq{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "auth without token",
			RequestBody: models.ValidateAuthorizationReq{
				Type: "auth",
			},
			ExpectedCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "app without app key",
			RequestBody: models.ValidateAuthorizationReq{
				Type: "app",
			},
			ExpectedCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "api without api keys",
			RequestBody: models.ValidateAuthorizationReq{
				Type: "api",
			},
			ExpectedCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "business without api keys",
			RequestBody: models.ValidateAuthorizationReq{
				Type: "business",
			},
			ExpectedCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "business_admin without api keys",
			RequestBody: models.ValidateAuthorizationReq{
				Type: "business_admin",
			},
			ExpectedCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
	}

	auth_model := auth_model.Controller{Db: db, Validator: validatorRef}

	authTypeUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, middleware.AppType))
	{
		authTypeUrl.POST("/validate_authorization", auth_model.ValidateAuthorization)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/validate_authorization"}

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

			if test.CheckResponseData && test.RequestBody.Type == "auth" {
				resData1 := data["data"].(map[string]interface{})
				resData := resData1["data"].(map[string]interface{})
				reqData := test.Response.(models.User)
				accountID := uint(resData["account_id"].(float64))
				if accountID != reqData.AccountID {
					t.Errorf("handler returned wrong status code: got status %v expected status %v", resData, reqData)
				}
			}

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
