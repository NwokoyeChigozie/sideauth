package models

import (
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type BannedAccount struct {
	ID        uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID int       `gorm:"column:account_id; type:int; not null" json:"account_id"`
	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (b *BannedAccount) CheckByAccountID(db *gorm.DB) (bool, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &b, "account_id = ? ", b.AccountID)
	if nilErr != nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}
