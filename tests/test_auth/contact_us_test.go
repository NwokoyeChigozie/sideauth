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
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/controller/auth"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	tst "github.com/vesicash/auth-ms/tests"
)

func TestContactUs(t *testing.T) {
	tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()

	auth := auth.Controller{Db: db, Validator: validatorRef}
	r := gin.Default()

	tests := []struct {
		Name         string
		RequestBody  models.ContactUsCreateModel
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK contact us",
			RequestBody: models.ContactUsCreateModel{
				FirstName:    "first name",
				LastName:     "last name",
				Email:        "email@email.com",
				WebsiteUrl:   "https://st.com",
				BusinessType: "ecommerce",
				Country:      "nigeria",
				Message:      "message",
			},
			ExpectedCode: http.StatusOK,
			Message:      "Contact-Us submitted",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "no first name",
			RequestBody: models.ContactUsCreateModel{
				LastName:     "last name",
				Email:        "email@email.com",
				WebsiteUrl:   "https://st.com",
				BusinessType: "ecommerce",
				Country:      "nigeria",
				Message:      "message",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "no last name",
			RequestBody: models.ContactUsCreateModel{
				FirstName:    "first name",
				Email:        "email@email.com",
				WebsiteUrl:   "https://st.com",
				BusinessType: "ecommerce",
				Country:      "nigeria",
				Message:      "message",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "no email",
			RequestBody: models.ContactUsCreateModel{
				FirstName:    "first name",
				LastName:     "last name",
				WebsiteUrl:   "https://st.com",
				BusinessType: "ecommerce",
				Country:      "nigeria",
				Message:      "message",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "no website url",
			RequestBody: models.ContactUsCreateModel{
				FirstName:    "first name",
				LastName:     "last name",
				Email:        "email@email.com",
				BusinessType: "ecommerce",
				Country:      "nigeria",
				Message:      "message",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "invalid website url",
			RequestBody: models.ContactUsCreateModel{
				FirstName:    "first name",
				LastName:     "last name",
				Email:        "email@email.com",
				WebsiteUrl:   "st.com",
				BusinessType: "ecommerce",
				Country:      "nigeria",
				Message:      "message",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "no business type",
			RequestBody: models.ContactUsCreateModel{
				FirstName:  "first name",
				LastName:   "last name",
				Email:      "email@email.com",
				WebsiteUrl: "https://st.com",
				Country:    "nigeria",
				Message:    "message",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "invalid business type",
			RequestBody: models.ContactUsCreateModel{
				FirstName:    "first name",
				LastName:     "last name",
				Email:        "email@email.com",
				WebsiteUrl:   "https://st.com",
				BusinessType: "ee",
				Country:      "nigeria",
				Message:      "message",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "no country",
			RequestBody: models.ContactUsCreateModel{
				FirstName:    "first name",
				LastName:     "last name",
				Email:        "email@email.com",
				WebsiteUrl:   "https://st.com",
				BusinessType: "ecommerce",
				Message:      "message",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "no message",
			RequestBody: models.ContactUsCreateModel{
				FirstName:    "first name",
				LastName:     "last name",
				Email:        "email@email.com",
				WebsiteUrl:   "https://st.com",
				BusinessType: "ecommerce",
				Country:      "nigeria",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	authTypeUrl := r.Group(fmt.Sprintf("%v/auth", "v2"))
	{
		authTypeUrl.POST("/contact-us", auth.ContactUs)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/auth/contact-us"}

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
