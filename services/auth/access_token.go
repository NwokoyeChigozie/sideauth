package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/internal/config"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
	"net/http"
	"strings"
)

func IssueAccessTokenService(db postgresql.Databases, accountID int) (models.AccessToken, int, error) {
	var (
		token = models.AccessToken{AccountID: accountID}
	)

	user := models.User{AccountID: uint(accountID)}
	code, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		return token, code, err
	}

	code, err = token.GetByAccountID(db.Auth)
	if err != nil {
		if code == http.StatusInternalServerError {
			return token, code, err
		}

		err = token.CreateAccessToken(db.Auth)
		if err != nil {
			return token, http.StatusInternalServerError, err
		}
	}

	app := config.GetConfig().App
	token.IsLive = true
	token.PrivateKey = "v_" + app.Name + "_" + utility.RandomString(50)
	token.PublicKey = "v_" + app.Name + "_" + utility.RandomString(50)

	err = token.Update(db.Auth)
	if err != nil {
		return token, http.StatusInternalServerError, err
	}

	return token, http.StatusOK, nil

}

func UpdateUserMorSettings(db postgresql.Databases, req models.EnableMORReq, accountId int) (models.User, int, error) {

	var (
		userDetails = models.User{AccountID: uint(accountId)}
	)

	code, err := userDetails.GetUserByAccountID(db.Auth)
	if err != nil {
		return models.User{}, code, err
	}

	if req.Status != nil {
		userDetails.IsMorEnabled = *req.Status
	}

	err = userDetails.Update(db.Auth)
	if err != nil {
		return userDetails, http.StatusInternalServerError, nil
	}

	return userDetails, http.StatusOK, nil

}

func GetUserService(db postgresql.Databases, searchParam string, isMorEnabledParam string) (interface{}, int, error) {

	var (
		resp          = []map[string]interface{}{}
		isMorEnabled  *bool
		trueD, falseD = true, false
	)

	if strings.EqualFold(isMorEnabledParam, "true") {
		isMorEnabled = &trueD
	} else if strings.EqualFold(isMorEnabledParam, "false") {
		isMorEnabled = &falseD
	}

	user := models.User{}
	users, err := user.GetUsers(db.Auth, searchParam, isMorEnabled)

	if err != nil {
		return resp, http.StatusInternalServerError, err
	}

	return users, http.StatusOK, nil

}

func ListSelectedCountriesService(db postgresql.Databases) (interface{}, int, error) {

	var (
		resp = []map[string]interface{}{}
	)

	country := models.Country{}
	countries, err := country.GetSelectedCountries(db.Auth)

	if err != nil {
		return resp, http.StatusInternalServerError, err
	}

	return countries, http.StatusOK, nil

}

func RevokeAccessTokenService(db postgresql.Databases, accountID int) (models.AccessToken, int, error) {
	var (
		token = models.AccessToken{AccountID: accountID}
	)

	user := models.User{AccountID: uint(accountID)}
	code, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		return token, code, err
	}

	code, err = token.GetByAccountID(db.Auth)
	if err != nil {
		if code == http.StatusInternalServerError {
			return token, code, err
		}

		return token, code, fmt.Errorf("No token found for this user")
	}

	err = token.RevokeAccessToken(db.Auth)
	if err != nil {
		return token, http.StatusInternalServerError, err
	}

	return token, http.StatusOK, nil

}

func GetUserWalletBalanceService(db postgresql.Databases, accountID int) (interface{}, int, error) {
	var (
		walletBalance = models.WalletBalance{AccountID: accountID}
	)

	user := models.User{AccountID: uint(accountID)}
	code, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		return nil, code, err
	}

	userProfile := models.UserProfile{AccountID: accountID}
	code, err = userProfile.GetByAccountID(db.Auth)
	if err != nil {
		return nil, code, err
	}

	var userCountry *string
	var countryName *string
	currency := userProfile.Currency
	if userProfile.Country != "" {
		userCountry = &userProfile.Country
	}

	if userCountry != nil {
		country := models.Country{CountryCode: *userCountry, Name: *userCountry}
		code, err := country.FindWithNameOrCode(db.Auth)
		if err != nil {
			if code == http.StatusInternalServerError {
				return nil, code, err
			}

		} else {
			countryName = &country.Name
		}

	}

	defaultCurrency := "USD"
	var userWallet = models.WalletBalance{}
	currencies := []string{defaultCurrency, strings.ToUpper(currency)}
	for _, userCurrency := range currencies {
		userCurrencyWallet := models.WalletBalance{AccountID: accountID, Currency: strings.ToUpper(userCurrency)}
		code, err = userCurrencyWallet.GetWalletBalanceByAccountIDAndCurrency(db.Auth)
		if err != nil {
			if code == http.StatusInternalServerError {
				return nil, code, err
			}
			userCurrencyWallet = models.WalletBalance{
				AccountID: accountID,
				Available: 0,
				Currency:  strings.ToUpper(userCurrency),
			}
			err = userCurrencyWallet.CreateWalletBalance(db.Auth)
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}
		} else {
			if strings.EqualFold(userCurrency, userProfile.Currency) {
				userWallet = userCurrencyWallet
			}
		}
	}

	walletBalances, err := walletBalance.GetUserWalletBalances(db.Auth)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return gin.H{
		"balance":  userWallet.Available,
		"currency": userWallet.Currency,
		"country":  countryName,
		"wallets":  walletBalances,
	}, http.StatusOK, nil
	//['balance' => (float) $wallet->available ?? 0, 'currency' => $currency, 'country'=> $countryName, 'wallets' => $wallets]
}
