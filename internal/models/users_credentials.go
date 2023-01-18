package models

import (
	"fmt"
	"net/http"
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

type GetUserCredentialModel struct {
	ID                 uint   `json:"id"`
	AccountID          uint   `json:"account_id" pgvalidate:"exists=auth$users$account_id"`
	IdentificationType string `json:"identification_type"`
}

type CreateUserCredentialModel struct {
	AccountID          uint   `json:"account_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
	Bvn                string `json:"bvn"`
	IdentificationType string `json:"identification_type"`
	IdentificationData string `json:"identification_data"`
}
type UpdateUserCredentialModel struct {
	ID                 uint   `json:"id" validate:"required"`
	AccountID          uint   `json:"account_id" pgvalidate:"exists=auth$users$account_id"`
	IdentificationType string `json:"identification_type"`
	Bvn                string `json:"bvn"`
	IdentificationData string `json:"identification_data"`
}

func (u *UsersCredential) GetUserCredentialByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &u, "id = ? ", u.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (u *UsersCredential) GetUserCredentialByAccountIdAndType(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &u, "account_id = ? and identification_type = ?", u.AccountID, u.IdentificationType)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
func (u *UsersCredential) GetUserCredentialByAccountID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &u, "account_id = ?", u.AccountID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (u *UsersCredential) CreateUsersCredential(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &u)
	if err != nil {
		return fmt.Errorf("user credential failed: %v", err.Error())
	}
	return nil
}

func (u *UsersCredential) Update(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &u)
	return err
}
