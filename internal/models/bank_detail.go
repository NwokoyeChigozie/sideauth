package models

import (
	"net/http"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type BankDetail struct {
	ID                    uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID             int       `gorm:"column:account_id; type:int; not null" json:"account_id"`
	BankID                int       `gorm:"column:bank_id; type:int; not null" json:"bank_id"`
	AccountName           string    `gorm:"column:account_name; type:varchar(250); not null" json:"account_name"`
	AccountNo             string    `gorm:"column:account_no; type:varchar(250); not null" json:"account_no"`
	Mobile_money_operator string    `gorm:"column:mobile_money_operator; type:varchar(250)" json:"mobile_money_operator"`
	SwiftCode             string    `gorm:"column:swift_code; type:varchar(250)" json:"swift_code"`
	SortCode              string    `gorm:"column:sort_code; type:varchar(250)" json:"sort_code"`
	BankAddress           string    `gorm:"column:bank_address; type:varchar(250)" json:"bank_address"`
	BankName              string    `gorm:"column:bank_name; type:varchar(250)" json:"bank_name"`
	MobileMoneyNumber     string    `gorm:"column:mobile_money_number; type:varchar(250)" json:"mobile_money_number"`
	Country               string    `gorm:"column:country; type:varchar(250); not null; default:'NG'" json:"country"`
	Currency              string    `gorm:"column:currency; type:varchar(250); not null; default:'NGN'" json:"currency"`
	CreatedAt             time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
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

func (b *BankDetail) GetAllByAccountID(db *gorm.DB) ([]BankDetail, error) {
	details := []BankDetail{}
	err := postgresql.SelectAllFromDb(db, "asc", &details, "account_id = ? ", b.AccountID)
	if err != nil {
		return details, err
	}
	return details, nil
}
