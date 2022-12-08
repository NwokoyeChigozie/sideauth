package models

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
	"gorm.io/gorm"
)

type UserIdentity struct {
	AccountID int    `json:"account_id"`
	Type      string `json:"type"`
}

var (
	MyIdentity *UserIdentity
)

type User struct {
	ID                        uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID                 uint      `gorm:"column:account_id; type:int; not null" json:"account_id"`
	AccountType               string    `gorm:"column:account_type; type:varchar(250)" json:"account_type"`
	Firstname                 string    `gorm:"column:firstname; type:varchar(250)" json:"firstname"`
	Lastname                  string    `gorm:"column:lastname; type:varchar(250)" json:"lastname"`
	EmailAddress              string    `gorm:"column:email_address; type:varchar(250)" json:"email_address"`
	PhoneNumber               string    `gorm:"column:phone_number; type:varchar(250)" json:"phone_number"`
	Username                  string    `gorm:"column:username; type:varchar(250)" json:"username"`
	Password                  string    `gorm:"column:password; type:varchar(250)" json:"-"`
	TierType                  int       `gorm:"column:tier_type; type:int" json:"tier_type"`
	DeletedAt                 time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt                 time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt                 time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	LoginAccessToken          string    `gorm:"column:login_access_token; type:text" json:"login_access_token"`
	LoginAccessTokenExpiresIn string    `gorm:"column:login_access_token_expires_in; type:varchar(250)" json:"login_access_token_expires_in"`
	BusinessId                int       `gorm:"column:business_id; type:int" json:"business_id"`
	Middlename                string    `gorm:"column:middlename; type:varchar(250)" json:"middlename"`
	AuthorizationRequired     bool      `gorm:"column:authorization_required; type:bool; default:false;not null" json:"authorization_required"`
	Meta                      string    `gorm:"column:meta; type:text" json:"meta"`
	ThePeerReference          string    `gorm:"column:the_peer_reference; type:varchar(250)" json:"the_peer_reference"`
}

type CreateUserRequestModel struct {
	BusinessID int `json:"business_id" pgvalidate:"exists=auth$users$account_id"`
	// BusinessID      int    `json:"business_id" `
	EmailAddress    string `json:"email_address" validate:"" pgvalidate:"notexists=auth$users$email_address, email"`
	PhoneNumber     string `json:"phone_number" pgvalidate:"notexists=auth$users$phone_number"`
	AccountType     string `json:"account_type" validate:"oneof=business individual others"`
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Username        string `json:"username" pgvalidate:"notexists=auth$users$username"`
	ReferralCode    string `json:"referral_code" pgvalidate:"exists=auth$users$username"`
	Password        string `json:"password"`
	Country         string `json:"country"`
	WebhookURI      string `json:"webhook_uri"`
	BusinessName    string `json:"business_name"`
	BusinessType    string `json:"business_type"`
	BusinessAddress string `json:"business_address"`
}

type LoginUserRequestModel struct {
	Username     string `json:"username"`
	EmailAddress string `json:"email_address"`
	Password     string `json:"password" validate:"required"`
	PhoneNumber  string `json:"phone_number"`
}

type BulkCreateUserRequestModel struct {
	Bulk []CreateUserRequestModel `json:"bulk" validate:"required"`
}

func (u *User) CreateUser(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &u)
	if err != nil {
		return fmt.Errorf("user creation failed: %v", err.Error())
	}
	return nil
}

func (u *User) GetUserByAccountID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &u, "account_id = ? ", u.AccountID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (u *User) Update(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &u)
	return err
}

func (u *User) GetUserByUsernameEmailOrPhone(db *gorm.DB) (int, error) {
	var (
		err, nilErr error
	)
	if u.Username != "" {
		err, nilErr = postgresql.SelectOneFromDb(db, &u, "LOWER(username) = ? ", strings.ToLower(u.Username))
		if nilErr != nil {
			nilErr = fmt.Errorf("username not found")
		}
	} else if u.EmailAddress != "" {
		err, nilErr = postgresql.SelectOneFromDb(db, &u, "LOWER(email_address) = ?", strings.ToLower(u.EmailAddress))
		if nilErr != nil {
			nilErr = fmt.Errorf("email address not found")
		}
	} else if u.PhoneNumber != "" {
		phone, _ := utility.PhoneValid(u.PhoneNumber)
		err, nilErr = postgresql.SelectOneFromDb(db, &u, "phone_number = ? or phone_number = ? ", u.PhoneNumber, phone)
		if nilErr != nil {
			nilErr = fmt.Errorf("phone number not found")
		}
	} else {
		err = fmt.Errorf("no values for GetUserByUsernameEmailOrPhone")
	}

	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil

}
