package models

import (
	"fmt"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type UserTracking struct {
	ID         uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID  int       `gorm:"column:account_id; type:int; not null" json:"account_id"`
	IpAddress  string    `gorm:"column:ip_address; type:varchar(250)" json:"ip_address"`
	DeviceType string    `gorm:"column:device_type; type:varchar(250)" json:"device_type"`
	Location   string    `gorm:"column:location; type:varchar(250)" json:"location"`
	Browser    string    `gorm:"column:browser; type:varchar(250)" json:"browser"`
	LoginTime  time.Time `gorm:"column:login_time; autoCreateTime" json:"login_time"`
	CreatedAt  time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (t *UserTracking) CreateUserTracking(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &t)
	if err != nil {
		return fmt.Errorf("user tracking creation failed: %v", err.Error())
	}
	return nil
}

func (t *UserTracking) GetAllByAccountID(db *gorm.DB) ([]UserTracking, error) {
	tracking := []UserTracking{}
	err := postgresql.SelectAllFromDb(db, "asc", &tracking, "account_id = ? ", t.AccountID)
	if err != nil {
		return tracking, err
	}
	return tracking, nil
}
