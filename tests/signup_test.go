package tests

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
	"github.com/vesicash/auth-ms/internal/config"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/internal/models/migrations"
	"github.com/vesicash/auth-ms/pkg/controller/auth"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
)

func Setup() {
	config := config.Setup("../config")
	db := postgresql.ConnectToDatabases(config.TestDatabases)
	if config.TestDatabases.Migrate {
		migrations.RunAllMigrations(db)
	}
}

func parseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	res := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&res)
	fmt.Println(res)
	return res
}

func assertStatusCode(t *testing.T, got, expected int) {
	if got != expected {
		t.Errorf("handler returned wrong status code: got status %d expected status %d", got, expected)
	}
}

func assertResponseMessage(t *testing.T, got, expected string) {
	if got != expected {
		t.Errorf("handler returned wrong message: got message: %q expected: %q", got, expected)
	}
}

func TestSignup(t *testing.T) {
	Setup()
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

	auth := auth.Controller{Db: db, Validator: validatorRef}

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

			data := parseResponse(rr)

			code := int(data["code"].(float64))
			assertStatusCode(t, code, test.ExpectedCode)
			fmt.Println("code", code)

			if test.Message != "" {
				message := data["message"]
				fmt.Println(message)
				if message != nil {
					assertResponseMessage(t, message.(string), test.Message)
				} else {
					assertResponseMessage(t, "", test.Message)
				}

			}

		})

	}

}
