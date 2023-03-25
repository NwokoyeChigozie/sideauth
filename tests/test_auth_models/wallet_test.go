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

func TestCreateWallet(t *testing.T) {
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
		Name         string
		RequestBody  models.CreateWalletRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK create wallet",
			RequestBody: models.CreateWalletRequest{
				AccountID: us.AccountID,
				Currency:  "NGN",
				Available: 200,
			},
			ExpectedCode: http.StatusCreated,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "OK without available",
			RequestBody: models.CreateWalletRequest{
				AccountID: us.AccountID,
				Currency:  "NGN",
			},
			ExpectedCode: http.StatusCreated,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no account id",
			RequestBody: models.CreateWalletRequest{
				Currency:  "NGN",
				Available: 200,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no currency",
			RequestBody: models.CreateWalletRequest{
				AccountID: us.AccountID,
				Available: 200,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name:         "no input",
			RequestBody:  models.CreateWalletRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
	}

	auth_model := auth_model.Controller{Db: db, Validator: validatorRef}

	authTypeUrl := r.Group(fmt.Sprintf("%v/auth", "v2"), middleware.Authorize(db, middleware.AppType))
	{
		authTypeUrl.POST("/create_wallet", auth_model.CreateWallet)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/create_wallet"}

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
func TestUpdateWalletBalance(t *testing.T) {
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

	wallet := models.WalletBalance{
		AccountID: int(us.AccountID),
		Currency:  "NGN",
		Available: 200,
	}
	err = wallet.CreateWalletBalance(db.Auth)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name         string
		RequestBody  models.UpdateWalletRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK update wallet",
			RequestBody: models.UpdateWalletRequest{
				ID:        wallet.ID,
				Available: 300,
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no account id",
			RequestBody: models.UpdateWalletRequest{
				Available: 200,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no available",
			RequestBody: models.UpdateWalletRequest{
				ID: wallet.ID,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name:         "no input",
			RequestBody:  models.UpdateWalletRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
	}

	auth_model := auth_model.Controller{Db: db, Validator: validatorRef}

	authTypeUrl := r.Group(fmt.Sprintf("%v/auth", "v2"), middleware.Authorize(db, middleware.AppType))
	{
		authTypeUrl.PATCH("/update_wallet_balance", auth_model.UpdateWalletBalance)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/update_wallet_balance"}

			req, err := http.NewRequest(http.MethodPatch, URI.String(), &b)
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
func TestGetWalletByAccountIDAndCurrency(t *testing.T) {
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
	ra := gin.Default()
	tst.SignupUser(t, ra, auth, userSignUpData)
	_, accountID := tst.GetLoginTokenAndAccountID(t, ra, auth, loginData)
	us := models.User{AccountID: uint(accountID)}
	_, err := us.GetUserByAccountID(db.Auth)
	if err != nil {
		log.Panic(err.Error())
	}
	r := gin.Default()

	wallet := models.WalletBalance{
		AccountID: int(us.AccountID),
		Currency:  "NGN",
		Available: 200,
	}
	err = wallet.CreateWalletBalance(db.Auth)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name         string
		RequestBody  interface{}
		ExpectedCode int
		Headers      map[string]string
		Message      string
		AccountID    int
		currency     string
	}{
		{
			Name:         "OK get wallet",
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
			AccountID: int(us.AccountID),
			currency:  "NGN",
		},
		{
			Name:         "wrong input",
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
			AccountID: 220,
			currency:  "NOT",
		},
	}

	auth_model := auth_model.Controller{Db: db, Validator: validatorRef}

	authTypeUrl := r.Group(fmt.Sprintf("%v/auth", "v2"), middleware.Authorize(db, middleware.AppType))
	{
		authTypeUrl.GET("/get_wallet/:account_id/:currency", auth_model.GetWalletByAccountIDAndCurrency)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: fmt.Sprintf("/v2/get_wallet/%v/%v", strconv.Itoa(test.AccountID), test.currency)}

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

func TestCreateWalletHistory(t *testing.T) {
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
		Name         string
		RequestBody  models.CreateWalletHistoryRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK create wallet history",
			RequestBody: models.CreateWalletHistoryRequest{
				AccountID:        int(us.AccountID),
				Reference:        utility.RandomString(20),
				Amount:           200,
				Currency:         "NGN",
				Type:             "credit",
				AvailableBalance: 250,
			},
			ExpectedCode: http.StatusCreated,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no account id",
			RequestBody: models.CreateWalletHistoryRequest{
				Reference:        utility.RandomString(20),
				Amount:           200,
				Currency:         "NGN",
				Type:             "credit",
				AvailableBalance: 250,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no reference",
			RequestBody: models.CreateWalletHistoryRequest{
				AccountID:        int(us.AccountID),
				Amount:           200,
				Currency:         "NGN",
				Type:             "credit",
				AvailableBalance: 250,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no amount",
			RequestBody: models.CreateWalletHistoryRequest{
				AccountID:        int(us.AccountID),
				Reference:        utility.RandomString(20),
				Currency:         "NGN",
				Type:             "credit",
				AvailableBalance: 250,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no currency",
			RequestBody: models.CreateWalletHistoryRequest{
				AccountID:        int(us.AccountID),
				Reference:        utility.RandomString(20),
				Amount:           200,
				Type:             "credit",
				AvailableBalance: 250,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no type",
			RequestBody: models.CreateWalletHistoryRequest{
				AccountID:        int(us.AccountID),
				Reference:        utility.RandomString(20),
				Amount:           200,
				Currency:         "NGN",
				AvailableBalance: 250,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "wrong type",
			RequestBody: models.CreateWalletHistoryRequest{
				AccountID:        int(us.AccountID),
				Reference:        utility.RandomString(20),
				Amount:           200,
				Currency:         "NGN",
				Type:             "wrong",
				AvailableBalance: 250,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no available balancce",
			RequestBody: models.CreateWalletHistoryRequest{
				AccountID: int(us.AccountID),
				Reference: utility.RandomString(20),
				Amount:    200,
				Currency:  "NGN",
				Type:      "credit",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name:         "no input",
			RequestBody:  models.CreateWalletHistoryRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
	}

	auth_model := auth_model.Controller{Db: db, Validator: validatorRef}

	authTypeUrl := r.Group(fmt.Sprintf("%v/auth", "v2"), middleware.Authorize(db, middleware.AppType))
	{
		authTypeUrl.POST("/create_wallet_history", auth_model.CreateWalletHistory)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/create_wallet_history"}

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
func TestCreateWalletTransaction(t *testing.T) {
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

	falseValue := false
	tests := []struct {
		Name         string
		RequestBody  models.CreateWalletTransactionRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK create wallet transaction with first approval",
			RequestBody: models.CreateWalletTransactionRequest{
				SenderAccountID:   int(us.AccountID),
				ReceiverAccountID: int(us.AccountID),
				SenderAmount:      200,
				ReceiverAmount:    3000,
				SenderCurrency:    "USD",
				ReceiverCurrency:  "NGN",
				Approved:          "pending",
				FirstApproval:     false,
			},
			ExpectedCode: http.StatusCreated,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "OK create wallet transaction with first and second approval",
			RequestBody: models.CreateWalletTransactionRequest{
				SenderAccountID:   int(us.AccountID),
				ReceiverAccountID: int(us.AccountID),
				SenderAmount:      200,
				ReceiverAmount:    3000,
				SenderCurrency:    "USD",
				ReceiverCurrency:  "NGN",
				Approved:          "pending",
				FirstApproval:     false,
				SecondApproval:    &falseValue,
			},
			ExpectedCode: http.StatusCreated,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no sender account id",
			RequestBody: models.CreateWalletTransactionRequest{
				ReceiverAccountID: int(us.AccountID),
				SenderAmount:      200,
				ReceiverAmount:    3000,
				SenderCurrency:    "USD",
				ReceiverCurrency:  "NGN",
				Approved:          "pending",
				FirstApproval:     false,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no receiver account id",
			RequestBody: models.CreateWalletTransactionRequest{
				SenderAccountID:  int(us.AccountID),
				SenderAmount:     200,
				ReceiverAmount:   3000,
				SenderCurrency:   "USD",
				ReceiverCurrency: "NGN",
				Approved:         "pending",
				FirstApproval:    false,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no sender amount",
			RequestBody: models.CreateWalletTransactionRequest{
				SenderAccountID:   int(us.AccountID),
				ReceiverAccountID: int(us.AccountID),
				ReceiverAmount:    3000,
				SenderCurrency:    "USD",
				ReceiverCurrency:  "NGN",
				Approved:          "pending",
				FirstApproval:     false,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no receiver amount",
			RequestBody: models.CreateWalletTransactionRequest{
				SenderAccountID:   int(us.AccountID),
				ReceiverAccountID: int(us.AccountID),
				SenderAmount:      200,
				SenderCurrency:    "USD",
				ReceiverCurrency:  "NGN",
				Approved:          "pending",
				FirstApproval:     false,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no sender currency",
			RequestBody: models.CreateWalletTransactionRequest{
				SenderAccountID:   int(us.AccountID),
				ReceiverAccountID: int(us.AccountID),
				SenderAmount:      200,
				ReceiverAmount:    3000,
				ReceiverCurrency:  "NGN",
				Approved:          "pending",
				FirstApproval:     false,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no receiver currency",
			RequestBody: models.CreateWalletTransactionRequest{
				SenderAccountID:   int(us.AccountID),
				ReceiverAccountID: int(us.AccountID),
				SenderAmount:      200,
				ReceiverAmount:    3000,
				SenderCurrency:    "USD",
				Approved:          "pending",
				FirstApproval:     false,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "no approved",
			RequestBody: models.CreateWalletTransactionRequest{
				SenderAccountID:   int(us.AccountID),
				ReceiverAccountID: int(us.AccountID),
				SenderAmount:      200,
				ReceiverAmount:    3000,
				SenderCurrency:    "USD",
				ReceiverCurrency:  "NGN",
				FirstApproval:     false,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "wrong approved value",
			RequestBody: models.CreateWalletTransactionRequest{
				SenderAccountID:   int(us.AccountID),
				ReceiverAccountID: int(us.AccountID),
				SenderAmount:      200,
				ReceiverAmount:    3000,
				SenderCurrency:    "USD",
				Approved:          "wrong",
				FirstApproval:     false,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name:         "no input",
			RequestBody:  models.CreateWalletTransactionRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
	}

	auth_model := auth_model.Controller{Db: db, Validator: validatorRef}

	authTypeUrl := r.Group(fmt.Sprintf("%v/auth", "v2"), middleware.Authorize(db, middleware.AppType))
	{
		authTypeUrl.POST("/create_wallet_transaction", auth_model.CreateWalletTransaction)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/create_wallet_transaction"}

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
