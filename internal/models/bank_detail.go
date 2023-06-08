package models

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type BankDetail struct {
	ID                  uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID           int       `gorm:"column:account_id; type:int; not null" json:"account_id"`
	BankID              int       `gorm:"column:bank_id; type:int" json:"bank_id"`
	AccountName         string    `gorm:"column:account_name; type:varchar(250); not null" json:"account_name"`
	AccountNo           string    `gorm:"column:account_no; type:varchar(250); not null" json:"account_no"`
	MobileMoneyOperator string    `gorm:"column:mobile_money_operator; type:varchar(250)" json:"mobile_money_operator"`
	SwiftCode           string    `gorm:"column:swift_code; type:varchar(250)" json:"swift_code"`
	SortCode            string    `gorm:"column:sort_code; type:varchar(250)" json:"sort_code"`
	BankAddress         string    `gorm:"column:bank_address; type:varchar(250)" json:"bank_address"`
	BankName            string    `gorm:"column:bank_name; type:varchar(250)" json:"bank_name"`
	MobileMoneyNumber   string    `gorm:"column:mobile_money_number; type:varchar(250)" json:"mobile_money_number"`
	Country             string    `gorm:"column:country; type:varchar(250); not null; default:'NG'" json:"country"`
	Currency            string    `gorm:"column:currency; type:varchar(250); not null; default:'NGN'" json:"currency"`
	CreatedAt           time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

type GetBankDetailModel struct {
	ID                    uint   `json:"id" pgvalidate:"exists=auth$bank_details$id"`
	AccountID             uint   `json:"account_id" pgvalidate:"exists=auth$users$account_id"`
	Country               string `json:"country"`
	Currency              string `json:"currency"`
	IsMobileMoneyOperator bool   `json:"is_mobile_money_operator"`
}

func (b *BankDetail) GetByAccountID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &b, "account_id = ? ", b.AccountID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (b *BankDetail) GetBankDetailByQuery(db *gorm.DB, isMobileMoneyOperator bool) (int, error) {
	query := ""
	if b.AccountID != 0 {
		if query != "" {
			query += " and "
		}
		query += fmt.Sprintf(" account_id = %v ", b.AccountID)
	}

	if b.Country != "" {
		if query != "" {
			query += " and "
		}
		query += fmt.Sprintf(" LOWER(country) = '%v' ", strings.ToLower(b.Country))
	}

	if b.Currency != "" {
		if query != "" {
			query += " and "
		}
		query += fmt.Sprintf(" LOWER(currency) = '%v' ", strings.ToLower(b.Currency))
	}

	if isMobileMoneyOperator {
		if query != "" {
			query += " and "
		}
		query += " (mobile_money_operator != '' and mobile_money_operator is not null)"
	}

	err, nilErr := postgresql.SelectOneFromDb(db, &b, query)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (b *BankDetail) GetByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &b, "id = ? ", b.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (b *BankDetail) GetAllByAccountID(db *gorm.DB) ([]BankDetail, error) {
	details := []BankDetail{}
	err := postgresql.SelectAllFromDb(db, "asc", &details, "account_id = ? ", b.AccountID)
	if err != nil {
		return details, err
	}
	return details, nil
}

type CreateBankRequest struct {
	AccountID           int    `json:"account_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
	BankID              int    `json:"bank_id" validate:"required"`
	AccountName         string `json:"account_name" validate:"required"`
	AccountNo           string `json:"account_no" validate:"required"`
	MobileMoneyOperator string `json:"mobile_money_operator"`
	SwiftCode           string `json:"swift_code"`
	SortCode            string `json:"sort_code"`
	BankAddress         string `json:"bank_address"`
	BankName            string `json:"bank_name"`
	MobileMoneyNumber   string `json:"mobile_money_number"`
	Country             string `json:"country"`
	Currency            string `json:"currency"`
}

func (b *BankDetail) CreateBankDetail(db *gorm.DB) (int, error) {
	if b.AccountID == 0 {
		return http.StatusBadRequest, fmt.Errorf("account id not provided")
	}
	if b.Country != "" && b.Currency != "" {
		country := Country{CountryCode: b.Country, CurrencyCode: b.Currency}
		code, err := country.FindWithCurrencyAndCode(db)
		if err != nil {
			return code, err
		}
		b.Country, b.Currency = country.CountryCode, country.CurrencyCode
	} else {
		b.Country, b.Currency = "NG", "NGN"
	}
	err := postgresql.CreateOneRecord(db, &b)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("user creation failed: %v", err.Error())
	}
	return http.StatusOK, nil
}
