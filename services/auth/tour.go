package auth

import (
	"net/http"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func UpdateTourStatusService(db postgresql.Databases, status bool, accountID int) (*bool, int, error) {
	user := models.User{AccountID: uint(accountID)}
	user.GetUserByAccountID(db.Auth)
	code, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		return nil, code, err
	}

	user.HasSeenTour = status
	err = user.Update(db.Auth)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &user.HasSeenTour, http.StatusOK, nil
}
