package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func (ReferralPromo) TableName() string {
	return "referral_promo"
}

type ReferralPromo struct {
	ID           uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	ReferralCode string    `gorm:"column:referral_code; type:varchar(250); not null" json:"referral_code"`
	PromoCode    string    `gorm:"column:promo_code; type:varchar(250); not null" json:"promo_code"`
	CreatedAt    time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (r *ReferralPromo) GetReferralPromoByCode(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &r, "referral_code = ? ", r.ReferralCode)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (r *ReferralPromo) ActivatePromoCode(db *gorm.DB, accountID int) (int, error) {
	businessCharge := BusinessCharge{BusinessId: accountID}
	switch r.ReferralCode {
	case "promo_august":
		return businessCharge.August(db)
	default:
		return http.StatusBadRequest, fmt.Errorf("promo code not implemented")
	}
}
