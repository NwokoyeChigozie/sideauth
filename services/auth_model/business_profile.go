package auth_model

import (
	"fmt"
	"net/http"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func GetBusinessProfileService(req models.GetBusinessProfileModel, db postgresql.Databases) (*models.BusinessProfile, int, error) {
	businessProfile := models.BusinessProfile{ID: req.ID, AccountID: int(req.AccountID), FlutterwaveMerchantID: req.FlutterwaveMerchantID}

	if req.AccountID == 0 && req.ID == 0 {
		return &models.BusinessProfile{}, http.StatusBadRequest, fmt.Errorf("either id or account_id is required")
	}

	if req.ID != 0 {
		code, err := businessProfile.GetByID(db.Auth)
		if err != nil {
			return &models.BusinessProfile{}, code, err
		}
	} else if req.AccountID != 0 {
		code, err := businessProfile.GetByAccountID(db.Auth)
		if err != nil {
			return &models.BusinessProfile{}, code, err
		}
	} else if req.FlutterwaveMerchantID != "" {
		code, err := businessProfile.GetByFlutterwaveMerchantID(db.Auth)
		if err != nil {
			return &models.BusinessProfile{}, code, err
		}
	} else {
		return &models.BusinessProfile{}, http.StatusBadRequest, fmt.Errorf("error occured please check your input")
	}

	return &businessProfile, http.StatusOK, nil
}
