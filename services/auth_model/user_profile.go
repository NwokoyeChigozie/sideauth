package auth_model

import (
	"fmt"
	"net/http"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func GetUserProfileService(req models.GetUserProfileModel, db postgresql.Databases) (*models.UserProfile, int, error) {
	userProfile := models.UserProfile{ID: req.ID, AccountID: int(req.AccountID)}

	if req.AccountID == 0 && req.ID == 0 {
		return &models.UserProfile{}, http.StatusBadRequest, fmt.Errorf("either id or account_id is required")
	}

	if req.ID != 0 {
		code, err := userProfile.GetByID(db.Auth)
		if err != nil {
			return &models.UserProfile{}, code, err
		}
	} else if req.AccountID != 0 {
		code, err := userProfile.GetByAccountID(db.Auth)
		if err != nil {
			return &models.UserProfile{}, code, err
		}
	} else {
		return &models.UserProfile{}, http.StatusBadRequest, fmt.Errorf("error occured please check your input")
	}

	return &userProfile, http.StatusOK, nil
}
