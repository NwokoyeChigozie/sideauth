package auth

import (
	"fmt"
	"github.com/vesicash/auth-ms/internal/config"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
	"net/http"
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

func UpdateUserMorSettings(db postgresql.Databases, req models.EnableMORReq) (models.User, int, error) {

	var (
		userDetails = models.User{AccountID: uint(req.AccountID)}
	)

	code, err := userDetails.GetUserByAccountID(db.Auth)
	if err != nil {
		return models.User{}, code, err
	}

	userDetails.IsMorEnabled = req.Status
	err = userDetails.Update(db.Auth)
	if err != nil {
		return userDetails, http.StatusInternalServerError, nil
	}

	return userDetails, http.StatusOK, nil

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
