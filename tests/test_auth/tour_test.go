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

func TestTourStatusUpdate(t *testing.T) {
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
		Status bool `json:"status"`
	}

	tests := []struct {
		Name         string
		RequestBody  requestBody
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK update tour status",
			RequestBody: requestBody{
				Status: true,
			},
			ExpectedCode: http.StatusOK,
			Message:      "tour status updated",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "OK update tour status 2",
			RequestBody:  requestBody{},
			ExpectedCode: http.StatusOK,
			Message:      "tour status updated",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	authTypeUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, middleware.AuthType))
	{
		authTypeUrl.POST("/user/update_tour_status", auth.UpdateTourStatus)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/user/update_tour_status"}

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
