package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Authorize struct {
	ID           uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID    int       `gorm:"column:account_id; type:int" json:"account_id"`
	Authorized   bool      `gorm:"column:authorized; type:bool" json:"authorized"`
	Token        string    `gorm:"column:token; type:varchar(250); not null" json:"token"`
	IpAddress    string    `gorm:"column:ip_address; type:varchar(250); not null" json:"ip_address"`
	Browser      string    `gorm:"column:browser; type:varchar(250); not null" json:"browser"`
	Os           string    `gorm:"column:os; type:varchar(250)" json:"os"`
	Location     string    `gorm:"column:location; type:varchar(250); not null" json:"location"`
	Attempt      int       `gorm:"column:attempt; type:int; default: 0" json:"attempt"`
	AuthorizedAt time.Time `gorm:"column:authorized_at" json:"authorized_at"`
	CreatedAt    time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	DeletedAt    time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

type GetAuthorizeModel struct {
	ID         uint   `json:"id" pgvalidate:"exists=auth$authorizes$id"`
	AccountID  uint   `json:"account_id" pgvalidate:"exists=auth$users$account_id"`
	Authorized bool   `json:"authorized"`
	IpAddress  string `json:"ip_address"`
	Browser    string `json:"browser"`
}

type CreateAuthorizeModel struct {
	AccountID  uint   `json:"account_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
	Authorized bool   `json:"authorized"`
	Token      string `json:"token"`
	IpAddress  string `json:"ip_address"`
	Browser    string `json:"browser"`
	Os         string `json:"os"`
	Location   string `json:"location"`
}

type UpdateAuthorizeModel struct {
	ID         uint   `json:"id" validate:"required" pgvalidate:"exists=auth$authorizes$id"`
	AccountID  uint   `json:"account_id" pgvalidate:"exists=auth$users$account_id"`
	Authorized bool   `json:"authorized"`
	Token      string `json:"token"`
	IpAddress  string `json:"ip_address"`
	Browser    string `json:"browser"`
	Os         string `json:"os"`
	Location   string `json:"location"`
	Attempt    int    `json:"attempt"`
}

func (a *Authorize) GetAuthorizeByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &a, "id = ? ", a.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (a *Authorize) GetAuthorizeByAccountIDAndOneOrAll(db *gorm.DB) (int, error) {
	query := `
		account_id = ? 
		and authorized = ?
	`

	if a.IpAddress != "" {
		query += ` and ip_address = '` + a.IpAddress + `'`
	}
	if a.Browser != "" {
		query += ` and browser = '` + a.Browser + `'`
	}

	err, nilErr := postgresql.SelectOneFromDb(db, &a, query, a.AccountID, a.Authorized)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (a *Authorize) GetAuthorizeByAccountIDAndOneOrAll2(db *gorm.DB) (int, error) {
	query := `
		account_id = ? 
	`

	if a.IpAddress != "" {
		query += ` and ip_address = '` + a.IpAddress + `'`
	}
	if a.Browser != "" {
		query += ` and browser = '` + a.Browser + `'`
	}

	err, nilErr := postgresql.SelectOneFromDb(db, &a, query, a.AccountID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (a *Authorize) CreateAuthorize(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &a)
	if err != nil {
		return fmt.Errorf("create authorize failed: %v", err.Error())
	}
	return nil
}

func (a *Authorize) Update(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &a)
	return err
}
