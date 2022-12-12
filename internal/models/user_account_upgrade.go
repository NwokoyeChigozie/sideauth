package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type UserAccountUpgrade struct {
	ID           uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID    int       `gorm:"column:account_id; type:int; not null" json:"account_id"`
	BusinessType string    `gorm:"column:business_type; type:varchar(250); not null" json:"business_type"`
	CreatedAt    time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (u *UserAccountUpgrade) GetByAccountID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &u, "account_id = ? ", u.AccountID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (u *UserAccountUpgrade) CreateUserAccountUpgrade(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &u)
	if err != nil {
		return fmt.Errorf("user account upgrade creation failed: %v", err.Error())
	}
	return nil
}
