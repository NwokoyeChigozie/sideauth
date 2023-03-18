package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func (WalletHistory) TableName() string {
	return "wallet_histories"
}

type WalletHistory struct {
	ID               uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID        string    `gorm:"column:account_id; type:varchar(255); not null" json:"account_id"`
	Reference        string    `gorm:"column:reference; type:varchar(255); not null" json:"reference"`
	Amount           float64   `gorm:"column:amount; type:decimal(20,2); not null" json:"amount"`
	Currency         string    `gorm:"column:currency; type:varchar(255); not null; comment: NGN,ESCROW_NGN,USD,ESCROW_USD,GBP,ESCROW_GBP" json:"currency"`
	Type             string    `gorm:"column:type; type:varchar(255); not null; comment: credit,debit" json:"type"`
	AvailableBalance float64   `gorm:"column:available_balance; type:decimal(20,2); not null" json:"available_balance"`
	CreatedAt        time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	DeletedAt        time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

type CreateWalletHistoryRequest struct {
	AccountID        int     `json:"account_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
	Reference        string  `json:"reference" validate:"required"`
	Amount           float64 `json:"amount" validate:"required"`
	Currency         string  `json:"currency" validate:"required"`
	Type             string  `json:"type" validate:"required,oneof=credit debit"`
	AvailableBalance float64 `json:"available_balance" validate:"required"`
}

func (w *WalletHistory) CreateWalletHistory(db *gorm.DB) error {
	w.Currency = strings.ToUpper(w.Currency)
	err := postgresql.CreateOneRecord(db, &w)
	if err != nil {
		return fmt.Errorf("wallet history creation failed: %v", err.Error())
	}
	return nil
}
