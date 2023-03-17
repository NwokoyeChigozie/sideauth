package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type WalletTransaction struct {
	ID                uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	SenderAccountID   string    `gorm:"column:sender_account_id; type:varchar(255); not null" json:"sender_account_id"`
	ReceiverAccountID string    `gorm:"column:receiver_account_id; type:varchar(255); not null" json:"receiver_account_id"`
	SenderAmount      float64   `gorm:"column:sender_amount; type:decimal(20,2); not null" json:"sender_amount"`
	ReceiverAmount    float64   `gorm:"column:receiver_amount; type:decimal(20,2); not null" json:"receiver_amount"`
	SenderCurrency    string    `gorm:"column:sender_currency; type:varchar(255); not null" json:"sender_currency"`
	ReceiverCurrency  string    `gorm:"column:receiver_currency; type:varchar(255); not null" json:"receiver_currency"`
	Approved          string    `gorm:"column:approved; type:varchar(255); not null; default: pending; comment: yes,no,pending" json:"approved"`
	CreatedAt         time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	DeletedAt         time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	FirstApproval     bool      `gorm:"column:first_approval; type:bool; default:false;not null" json:"first_approval"`
	SecondApproval    bool      `gorm:"column:second_approval; type:bool" json:"second_approval"`
}

type CreateWalletTransactionRequest struct {
	SenderAccountID   int     `json:"sender_account_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
	ReceiverAccountID int     `json:"receiver_account_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
	SenderAmount      float64 `json:"sender_amount" validate:"required"`
	ReceiverAmount    float64 `json:"receiver_amount" validate:"required"`
	SenderCurrency    string  `json:"sender_currency" validate:"required"`
	ReceiverCurrency  string  `json:"receiver_currency" validate:"required"`
	Approved          string  `json:"approved" validate:"required,oneof=yes no pending"`
	FirstApproval     bool    `json:"first_approval"`
	SecondApproval    *bool   `json:"second_approval"`
}

func (w *WalletTransaction) CreateWalletTransaction(db *gorm.DB) error {
	w.SenderCurrency = strings.ToUpper(w.SenderCurrency)
	w.ReceiverCurrency = strings.ToUpper(w.ReceiverCurrency)
	err := postgresql.CreateOneRecord(db, &w)
	if err != nil {
		return fmt.Errorf("wallet transaction creation failed: %v", err.Error())
	}
	return nil
}
