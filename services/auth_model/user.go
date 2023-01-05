package auth_model

import (
	"fmt"
	"net/http"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func GetUserService(req models.GetUserModel, db postgresql.Databases) (*models.User, int, error) {
	user := models.User{}
	if req.ID != 0 {
		user.ID = req.ID
		code, err := user.GetUserByID(db.Auth)
		if err != nil {
			return nil, code, err
		}
		return &user, http.StatusOK, nil
	} else if req.AccountID != 0 {
		user.AccountID = req.AccountID
		code, err := user.GetUserByAccountID(db.Auth)
		if err != nil {
			return nil, code, err
		}
		return &user, http.StatusOK, nil
	} else if req.EmailAddress != "" {
		user.EmailAddress = req.EmailAddress
		code, err := user.GetUserByUsernameEmailOrPhone(db.Auth)
		if err != nil {
			return nil, code, err
		}
		return &user, http.StatusOK, nil
	} else if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
		code, err := user.GetUserByUsernameEmailOrPhone(db.Auth)
		if err != nil {
			return nil, code, err
		}
		return &user, http.StatusOK, nil
	} else if req.Username != "" {
		user.Username = req.Username
		code, err := user.GetUserByUsernameEmailOrPhone(db.Auth)
		if err != nil {
			return nil, code, err
		}
		return &user, http.StatusOK, nil
	} else {
		return nil, http.StatusBadRequest, fmt.Errorf("no request values provided")
	}
}
