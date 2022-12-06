package models

import (
	"net/http"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type EscrowCharge struct {
	ID             uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	BusinessID     int       `gorm:"column:business_id; type:int; not null" json:"business_id"`
	BusinessCharge string    `gorm:"column:business_charge; type:varchar(250); not null" json:"business_charge"`
	VesicashCharge string    `gorm:"column:vesicash_charge; type:varchar(250); not null" json:"vesicash_charge"`
	IsTermsAgreed  bool      `gorm:"column:is_terms_agreed; type:bool;default:false" json:"is_terms_agreed"`
	CreatedAt      time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (e *EscrowCharge) GetByBusinessID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &e, "business_id = ? ", e.BusinessID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
