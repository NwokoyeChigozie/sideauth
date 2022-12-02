package models

import (
	"net/http"
	"strings"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Country struct {
	ID           uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	Name         string    `gorm:"column:name; type:varchar(250); not null" json:"name"`
	CountryCode  string    `gorm:"column:country_code; type:varchar(250); not null" json:"country_code"`
	CurrencyCode string    `gorm:"column:currency_code; type:varchar(250); not null" json:"currency_code"`
	CreatedAt    time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (c *Country) FindWithNameOrCode(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &c, "name = ? or country_code = ?", c.Name, strings.ToUpper(c.Name))
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
