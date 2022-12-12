package auth

import (
	"fmt"
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

func GetBusinessCustomersBankDetailsService(db postgresql.Databases, accountID int) (interface{}, int, error) {
	var (
		resp        = []map[string]interface{}{}
		bankDetails = []models.BankDetail{}
	)
	user := models.User{AccountID: uint(accountID)}
	code, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		return resp, code, err
	}

	users, err := user.SelectByBusinessID(db.Auth)
	if err != nil {
		return resp, http.StatusInternalServerError, err
	}

	for _, u := range users {
		cBank := models.BankDetail{AccountID: int(u.AccountID)}
		banks, err := cBank.GetAllByAccountID(db.Auth)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			bankDetails = append(bankDetails, banks...)
		}
	}

	for _, b := range bankDetails {
		record := map[string]interface{}{
			"account_id":   b.AccountID,
			"account_name": b.AccountName,
			"bank_id":      b.BankID,
			"account_no":   b.AccountNo,
		}
		resp = append(resp, record)
	}

	return resp, http.StatusOK, nil
}
