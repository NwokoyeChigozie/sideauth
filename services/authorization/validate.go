package authorization

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/vesicash/auth-ms/internal/config"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/middleware"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func ValidateAuthorizationService(req models.ValidateAuthorizationReq, db postgresql.Databases) (interface{}, string, bool, int, error) {
	switch req.Type {
	case string(middleware.ApiType):
		msg, status := validateApiType(db, req.VPrivateKey, req.VPublicKey)
		return nil, msg, status, http.StatusOK, nil
	case string(middleware.AppType):
		msg, status := validateAppType(db, req.VApp)
		return nil, msg, status, http.StatusOK, nil
	case string(middleware.AuthType):
		data, msg, status := validateAuthType(db, req.AuthorizationToken)
		return data, msg, status, http.StatusOK, nil
	case string(middleware.BusinessAdmin):
		msg, status := validateBusinessAdminType(db, req.VPrivateKey, req.VPublicKey)
		return nil, msg, status, http.StatusOK, nil
	case string(middleware.Business):
		msg, status := validateBusinessType(db, req.VPrivateKey, req.VPublicKey)
		return nil, msg, status, http.StatusOK, nil
	default:
		return nil, "not implemented", false, http.StatusBadRequest, fmt.Errorf("not implemented")
	}
}

func validateAuthType(db postgresql.Databases, bearerToken string) (interface{}, string, bool) {
	var invalidToken = "Your request was made with invalid credentials."
	if bearerToken == "" {
		return nil, invalidToken, false
	}

	token, err := middleware.TokenValid(bearerToken)
	if err != nil {
		return nil, invalidToken, false
	}

	claims := token.Claims.(jwt.MapClaims)
	activeUserType, ok := claims["type"].(string) //convert the interface to string
	if !ok {
		return nil, invalidToken, false
	}

	activeUserAccountID, ok := claims["account_id"].(float64) //convert the interface to float
	if !ok {
		return nil, invalidToken, false
	}

	authoriseStatus, ok := claims["authorised"].(bool) //check if token is authorised for middleware
	if !ok && !authoriseStatus {
		return nil, invalidToken, false
	}

	myIdentity := models.UserIdentity{
		AccountID: int(activeUserAccountID),
		Type:      activeUserType,
	}

	user := models.User{AccountID: uint(myIdentity.AccountID)}
	code, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		if code == http.StatusInternalServerError {
			return nil, err.Error(), false
		}
		return nil, "user does not exist", false
	}

	if user.LoginAccessToken != bearerToken {
		return nil, invalidToken, false
	}

	if user.LoginAccessTokenExpiresIn == "" {
		return nil, invalidToken, false
	}

	parseInt, err := strconv.Atoi(user.LoginAccessTokenExpiresIn)
	if err != nil {
		return nil, invalidToken, false
	}

	unixTimeUTC := time.Unix(int64(parseInt), 0)
	if time.Now().After(unixTimeUTC) {
		return nil, "expired token", false
	}

	return user, "authorized", true
}

func validateBusinessType(db postgresql.Databases, privateKey, publicKey string) (string, bool) {
	_, msg, status := checkAccessTokens(db, privateKey, publicKey)
	return msg, status
}

func validateAppType(db postgresql.Databases, appKey string) (string, bool) {
	config := config.GetConfig().App
	if appKey == "" {
		return "missing app key", false
	}

	if appKey != config.Key {
		return "invalid app key", false
	}

	return "authorized", true
}

func validateBusinessAdminType(db postgresql.Databases, privateKey, publicKey string) (string, bool) {
	token, msg, status := checkAccessTokens(db, privateKey, publicKey)
	if !status {
		return msg, status
	}

	user := models.User{AccountID: uint(token.AccountID)}
	_, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "access denied", false
		} else {
			return "server error", false
		}
	}
	if user.AccountType != "admin" {
		return "access denied", false
	}
	return "authorized", true
}

func validateApiType(db postgresql.Databases, privateKey, publicKey string) (string, bool) {
	_, msg, status := checkAccessTokens(db, privateKey, publicKey)
	return msg, status
}

func checkAccessTokens(db postgresql.Databases, privateKey, publicKey string) (models.AccessToken, string, bool) {

	if privateKey == "" && publicKey == "" {
		return models.AccessToken{}, "missing api keys", false
	}

	if privateKey == "" || publicKey == "" {
		return models.AccessToken{}, "either public or private key is missing", false
	}

	token := models.AccessToken{PublicKey: publicKey, PrivateKey: privateKey, IsLive: true}
	_, err := token.LiveTokensWithPublicOrPrivateKey(db.Auth)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return token, "invalid keys", false
		} else {
			return token, "server error", false
		}
	}
	return token, "authorized", true
}
