package auth

import (
	"net/http"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func CreateBankDetailService(req models.CreateBankRequest, db postgresql.Databases) (models.BankDetail, int, error) {
	bankDetail := models.BankDetail{
		AccountID:           req.AccountID,
		BankID:              req.BankID,
		AccountName:         req.AccountName,
		AccountNo:           req.AccountNo,
		MobileMoneyOperator: req.MobileMoneyOperator,
		SwiftCode:           req.SwiftCode,
		SortCode:            req.SortCode,
		BankAddress:         req.BankAddress,
		BankName:            req.BankName,
		MobileMoneyNumber:   req.MobileMoneyNumber,
		Country:             req.Country,
		Currency:            req.Currency,
	}
	code, err := bankDetail.CreateBankDetail(db.Auth)
	if err != nil {
		return bankDetail, code, err
	}
	return bankDetail, http.StatusOK, nil

}
