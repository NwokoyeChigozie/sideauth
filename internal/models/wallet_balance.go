package models

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type WalletBalance struct {
	ID        uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID int       `gorm:"column:account_id; type:int; not null" json:"account_id"`
	Available float64   `gorm:"column:available; type:decimal(20,2); not null" json:"available"`
	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	Currency  string    `gorm:"column:currency; type:varchar(255)" json:"currency"`
}

type CreateWalletRequest struct {
	AccountID uint    `json:"account_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
	Currency  string  `json:"currency" validate:"required"`
	Available float64 `json:"available"`
}
type UpdateWalletRequest struct {
	ID        uint    `json:"id" validate:"required" pgvalidate:"exists=auth$wallet_balances$id"`
	Available float64 `json:"available"`
}
type GetWalletsRequest struct {
	Currencies []string `json:"currencies" validate:"required"`
}

func (w *WalletBalance) GetWalletBalanceByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &w, "id = ? ", w.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (w *WalletBalance) GetWalletBalanceByAccountIDAndCurrency(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &w, "account_id = ? and LOWER(currency)=?", w.AccountID, strings.ToLower(w.Currency))
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (w *WalletBalance) GetWalletBalancesByAccountIDAndCurrencies(db *gorm.DB, currencies []string) ([]WalletBalance, error) {
	lowerCurrencies := []string{}
	for _, v := range currencies {
		lowerCurrencies = append(lowerCurrencies, strings.ToLower(v))
	}
	fmt.Println("lowerCurrencies", lowerCurrencies)
	wallets := []WalletBalance{}
	err := postgresql.SelectAllFromDb(db, "asc", &wallets, "account_id = ? and LOWER(currency) IN (?) ", w.AccountID, lowerCurrencies)
	if err != nil {
		return wallets, err
	}
	fmt.Println("wallets", wallets)
	return wallets, nil
}

func (w *WalletBalance) CreateWalletBalance(db *gorm.DB) error {
	w.Currency = strings.ToUpper(w.Currency)
	err := postgresql.CreateOneRecord(db, &w)
	if err != nil {
		return fmt.Errorf("wallet creation failed: %v", err.Error())
	}
	return nil
}

func (w *WalletBalance) Update(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &w)
	return err
}
