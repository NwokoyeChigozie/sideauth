package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func (OtpVerification) TableName() string {
	return "otp_verification"
}

type OtpVerification struct {
	ID        uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID int       `gorm:"column:account_id; type:int; not null; comment: account id of the user" json:"account_id"`
	OtpToken  string    `gorm:"column:otp_token; type:varchar(250); not null" json:"otp_token"`
	ExpiresAt time.Time `gorm:"column:expires_at;" json:"expires_at"`
	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

type SendOtpTokenReq struct {
	AccountID int `json:"account_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
}

func (o *OtpVerification) GetLatestByAccountID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectLatestFromDb(db, &o, "account_id = ? ", o.AccountID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (u *OtpVerification) Create(db *gorm.DB) error {
	expiry := time.Now().Add(30 * time.Minute)
	expiry.Format("2006-01-02 15:04:05")
	u.ExpiresAt = expiry
	err := postgresql.CreateOneRecord(db, &u)
	if err != nil {
		return fmt.Errorf("otp creation failed: %v", err.Error())
	}
	return nil
}
func (u *OtpVerification) Delete(db *gorm.DB) error {
	err := postgresql.DeleteRecordFromDb(db, &u)
	if err != nil {
		return fmt.Errorf("otp delete failed: %v", err.Error())
	}
	return nil
}
