package login

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

func TestLogin(t *testing.T) {
	tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		loginPath      = "/v2/auth/login"
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

	auth := auth.Controller{Db: db, Validator: validatorRef}
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
