package models

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Bank struct {
	ID        uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	Name      string    `gorm:"column:name; type:varchar(250); not null" json:"name"`
	Code      string    `gorm:"column:code; type:varchar(250); not null" json:"code"`
	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	Country   string    `gorm:"column:country; type:varchar(250); default: Nigeria" json:"country"`
}

type GetBankRequest struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Country string `json:"country"`
}

func (b *Bank) GetBankByQuery(db *gorm.DB) (int, error) {
	query := ""
	if b.ID != 0 {
		query += fmt.Sprintf(" id = %v ", b.ID)
	}

	if b.Name != "" {
		if query != "" {
			query += " and "
		}
		query += fmt.Sprintf(" LOWER(name) = '%v' ", strings.ToLower(b.Name))
	}

	if b.Code != "" {
		if query != "" {
			query += " and "
		}
		query += fmt.Sprintf(" LOWER(code) = '%v' ", strings.ToLower(b.Code))
	}

	if b.Country != "" {
		if query != "" {
			query += " and "
		}
		query += fmt.Sprintf(" LOWER(country) = '%v' ", strings.ToLower(b.Country))
	}

	err, nilErr := postgresql.SelectOneFromDb(db, &b, query)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (b *Bank) CreateBank(db *gorm.DB) (int, error) {
	b.Country = strings.ToUpper(b.Country)
	err := postgresql.CreateOneRecord(db, &b)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("bank creation failed: %v", err.Error())
	}
	return http.StatusOK, nil
}
