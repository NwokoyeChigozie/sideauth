package models

import (
	"fmt"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type UsersCredential struct {
	ID                 uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID          int       `gorm:"column:account_id; type:int; not null" json:"account_id"`
	Bvn                string    `gorm:"column:bvn; type:varchar(250)" json:"bvn"`
	IdentificationType string    `gorm:"column:identification_type; type:varchar(250)" json:"identification_type"`
	IdentificationData string    `gorm:"column:identification_data; type:varchar(250)" json:"identification_data"`
	DeletedAt          time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt          time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (u *UsersCredential) CreateUsersCredential(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &u)
	if err != nil {
		return fmt.Errorf("user credential failed: %v", err.Error())
	}
	return nil
}
