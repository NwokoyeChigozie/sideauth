package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type BusinessProfile struct {
	ID                                uint    `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID                         int     `gorm:"column:account_id; type:int; not null" json:"account_id"`
	BusinessName                      string  `gorm:"column:business_name; type:varchar(250)" json:"business_name"`
	BusinessType                      string  `gorm:"column:business_type; type:varchar(250)" json:"business_type"`
	LogoUri                           string  `gorm:"column:logo_uri; type:varchar(250)" json:"logo_uri"`
	Website                           string  `gorm:"column:website; type:varchar(250)" json:"website"`
	Country                           string  `gorm:"column:country; type:varchar(250)" json:"country"`
	BusinessAddress                   string  `gorm:"column:business_address; type:varchar(250)" json:"business_address"`
	PaymentGateway                    string  `gorm:"column:payment_gateway; type:varchar(250)" json:"payment_gateway"`
	EscrowChargeOld                   float32 `gorm:"column:escrow_charge_old; type:decimal(20,2)" json:"escrow_charge_old"`
	DisbursementGateway               string  `gorm:"column:disbursement_gateway; type:varchar(255)" json:"disbursement_gateway"`
	AutoTransactionStatusSettings     bool    `gorm:"column:auto_transaction_status_settings; type:bool; default:false; not null" json:"auto_transaction_status_settings"`
	DisbursementSettings              string  `gorm:"column:disbursement_settings; type:varchar(255); not null; default:'instant'; comment: instant or accumulate" json:"disbursement_settings"`
	State                             string  `gorm:"column:state; type:varchar(255)" json:"state"`
	City                              string  `gorm:"column:city; type:varchar(255)" json:"city"`
	Webhook_uri                       string  `gorm:"column:webhook_uri; type:varchar(255)" json:"webhook_uri"`
	Currency                          string  `gorm:"column:currency; type:varchar(255); not null; default:'USD'" json:"currency"`
	IsRegistered                      bool    `gorm:"column:is_registered; type:bool" json:"is_registered"`
	DefaultDeliveryPeriod             string  `gorm:"column:default_delivery_period; type:varchar(255)" json:"default_delivery_period"`
	BusinessIgnoredNotifications      string  `gorm:"column:business_ignored_notifications; type:text; comment:This holds a JSON of ignored notifications for a business" json:"business_ignored_notifications"`
	BusinessCancellationFee           string  `gorm:"column:business_cancellation_fee; type:varchar(255); default:'0'" json:"business_cancellation_fee"`
	BusinessProcessingFee             string  `gorm:"column:business_processing_fee; type:varchar(255); default:'0'" json:"business_processing_fee"`
	AutoAggregateTransactionsSettings bool    `gorm:"column:auto_aggregate_transactions_settings; type:bool; not null; default:false" json:"auto_aggregate_transactions_settings"`
	DefaultChargeBearer               string  `gorm:"column:default_charge_bearer; type:varchar(255)" json:"default_charge_bearer"`
	IsVerificationWaved               bool    `gorm:"column:is_verification_waved; type:bool; default:false" json:"is_verification_waved"`

	BusinessGivenNotifications string  `gorm:"column:business_given_notifications; type:text; comment: List of notifications a business will receive" json:"business_given_notifications"`
	Units                      float32 `gorm:"column:units; type:decimal(8,2); default:1" json:"units"`

	RedirectUrl                   string    `gorm:"column:redirect_url; type:varchar(255)" json:"redirect_url"`
	BusinessDisabledNotifications bool      `gorm:"column:business_disabled_notifications; type:bool; default: false" json:"business_disabled_notifications"`
	IsBankTransferFeeWaved        bool      `gorm:"column:is_bank_transfer_fee_waved; type:bool; not null; default false" json:"is_bank_transfer_fee_waved"`
	EscrowCharge                  jsonmap   `gorm:"column:escrow_charge; type:json; not null; default: '{\"type\":\"percentage\",\"value\":\"0.05\"}'" json:"escrow_charge"`
	Bio                           string    `gorm:"column:bio; type:text" json:"bio"`
	BusinessEmail                 string    `gorm:"column:business_email; type:varchar(255)" json:"business_email"`
	DeletedAt                     time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt                     time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt                     time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

type GetBusinessProfileModel struct {
	ID        uint `json:"id"`
	AccountID uint `json:"account_id" pgvalidate:"exists=auth$users$account_id"`
}

func (b *BusinessProfile) CreateBusinessProfile(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &b)
	if err != nil {
		return fmt.Errorf("business Profile failed: %v", err.Error())
	}
	return nil
}

func (b *BusinessProfile) GetByAccountID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &b, "account_id = ? ", b.AccountID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (b *BusinessProfile) GetByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &b, "id = ? ", b.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
