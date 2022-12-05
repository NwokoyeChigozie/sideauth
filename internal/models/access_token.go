package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type AccessToken struct {
	ID            uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID     int       `gorm:"column:account_id; type:int; not null" json:"account_id"`
	PublicKey     string    `gorm:"column:public_key; type:varchar(250); not null" json:"public_key"`
	PrivateKey    string    `gorm:"column:private_key; type:varchar(250); not null" json:"private_key"`
	IsLive        bool      `gorm:"column:is_live; type:bool; default:false; not null" json:"is_live"`
	IsTermsAgreed bool      `gorm:"column:is_terms_agreed; type:bool;default:false" json:"is_terms_agreed"`
	CreatedAt     time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (a *AccessToken) GetAccessTokens(db *gorm.DB) error {
	err := postgresql.SelectFirstFromDb(db, &a)
	if err != nil {
		return fmt.Errorf("token selection failed: %v", err.Error())
	}
	return nil
}

func (a *AccessToken) LiveTokensWithPublicOrPrivateKey(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &a, "(public_key = ? or private_key = ?) and is_live = ?", a.PublicKey, a.PrivateKey, a.IsLive)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
