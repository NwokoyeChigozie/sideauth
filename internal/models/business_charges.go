package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type BusinessCharge struct {
	ID                  uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	BusinessId          int       `gorm:"column:business_id; type:int; not null; comment:'same as account_id'" json:"account_id"`
	Country             string    `gorm:"column:country; type:varchar(250)" json:"country"`
	Currency            string    `gorm:"column:currency; type:varchar(250)" json:"currency"`
	BusinessCharge      string    `gorm:"column:business_charge; type:varchar(250); not null; default:'0'" json:"business_charge"`
	VesicashCharge      string    `gorm:"column:vesicash_charge; type:varchar(250); not null; default:'0'" json:"vesicash_charge"`
	ProcessingFee       string    `gorm:"column:processing_fee; type:varchar(250); not null; default:'0'" json:"processing_fee"`
	CancellationFee     string    `gorm:"column:cancellation_fee; type:varchar(250); default:'0'" json:"cancellation_fee"`
	DisbursementCharge  string    `gorm:"column:disbursement_charge; type:varchar(250); default:'0'" json:"disbursement_charge"`
	PaymentGateway      string    `gorm:"column:payment_gateway; type:varchar(250)" json:"payment_gateway"`
	DisbursementGateway string    `gorm:"column:disbursement_gateway; type:varchar(250)" json:"disbursement_gateway"`
	ChargeMin           jsonmap   `gorm:"column:charge_min; type:varchar(250)" json:"charge_min"`
	ChargeMid           jsonmap   `gorm:"column:charge_mid; type:varchar(250)" json:"charge_mid"`
	ChargeMax           jsonmap   `gorm:"column:charge_max; type:varchar(250)" json:"charge_max"`
	ProcessingFeeMode   string    `gorm:"column:processing_fee_mode; type:varchar(250); default:'fixed'" json:"processing_fee_mode"`
	DeletedAt           time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt           time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (b *BusinessCharge) CreateBusinessCharge(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &b)
	if err != nil {
		return fmt.Errorf("business Charge failed: %v", err.Error())
	}
	return nil
}

func (b *BusinessCharge) GetByAccountID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &b, "business_id = ? ", b.BusinessId)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (b *BusinessCharge) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &b)
	return err
}

func (b *BusinessCharge) August(db *gorm.DB) (int, error) {
	code, err := b.GetByAccountID(db)
	if err != nil {
		return code, err
	}

	b.ChargeMin = jsonmap{
		"amount": 20000,
		"charge": 250,
	}
	b.ChargeMid = jsonmap{
		"amount": 20000,
		"charge": 500,
	}
	b.ChargeMax = jsonmap{
		"amount": 20000,
		"charge": 500,
	}
	err = b.UpdateAllFields(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
