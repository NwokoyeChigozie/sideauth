package models

import (
	"fmt"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type PasswordResetToken struct {
	ID        uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID int       `gorm:"column:account_id; type:int; not null" json:"account_id"`
	Token     int       `gorm:"column:token; type:varchar(250); not null" json:"token"`
	ExpiresAt string    `gorm:"column:expires_at; type:varchar(250); not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (p *PasswordResetToken) CreatePasswordResetToken(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &p)
	if err != nil {
		return fmt.Errorf("password reset token creation failed: %v", err.Error())
	}
	return nil
}
