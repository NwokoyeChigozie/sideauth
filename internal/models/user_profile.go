package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type UserProfile struct {
	ID         uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID  int       `gorm:"column:account_id; type:int; unique; not null" json:"account_id"`
	Address    string    `gorm:"column:address; type:varchar(250)" json:"address"`
	State      string    `gorm:"column:state; type:varchar(250)" json:"state"`
	City       string    `gorm:"column:city; type:varchar(250)" json:"city"`
	Country    string    `gorm:"column:country; type:varchar(250)" json:"country"`
	Dob        string    `gorm:"column:dob; type:varchar(250)" json:"dob"`
	Currency   string    `gorm:"column:currency; type:varchar(250);default:'USD'; not null" json:"currency"`
	IpAddress  string    `gorm:"column:ip_address; type:varchar(250)" json:"ip_address"`
	Sex        string    `gorm:"column:sex; type:varchar(250)" json:"sex"`
	Profession string    `gorm:"column:profession; type:varchar(250)" json:"profession"`
	Age        uint      `gorm:"column:age; type:int" json:"age"`
	Bio        string    `gorm:"column:bio; type:text" json:"bio"`
	DeletedAt  time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt  time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

type GetUserProfileModel struct {
	ID        uint `json:"id"`
	AccountID uint `json:"account_id" pgvalidate:"exists=auth$users$account_id"`
}

func (u *UserProfile) CreateUserProfile(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &u)
	if err != nil {
		return fmt.Errorf("user Profile creation failed: %v", err.Error())
	}
	return nil
}

func (u *UserProfile) GetByAccountID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &u, "account_id = ? ", u.AccountID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
func (u *UserProfile) GetByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &u, "id = ? ", u.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
