package auth

import (
	"net/http"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
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

	return token, http.StatusOK, nil

}
