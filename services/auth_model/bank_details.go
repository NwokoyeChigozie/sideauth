package auth_model

import (
	"fmt"
	"net/http"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func GetBankDetailService(req models.GetBankDetailModel, db postgresql.Databases) (*models.BankDetail, int, error) {
	bankDetail := models.BankDetail{ID: req.ID, AccountID: int(req.AccountID), Currency: req.Currency, Country: req.Country}

	if req.AccountID == 0 && req.ID == 0 {
		return &models.BankDetail{}, http.StatusBadRequest, fmt.Errorf("either id or account_id is required")
	}

	if req.ID != 0 {
		code, err := bankDetail.GetByID(db.Auth)
		if err != nil {
			return &models.BankDetail{}, code, err
		}
	} else {
		code, err := bankDetail.GetBankDetailByQuery(db.Auth)
		if err != nil {
			return &models.BankDetail{}, code, err
		}
	}

	return &bankDetail, http.StatusOK, nil
}
