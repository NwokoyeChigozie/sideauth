package models

import (
	"errors"
	"fmt"
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

type GetCountryModel struct {
	ID           uint   `json:"id" pgvalidate:"exists=auth$countries$id"`
	Name         string `json:"name"`
	CountryCode  string `json:"country_code"`
	CurrencyCode string `json:"currency_code"`
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

func (c *Country) FindWithCurrencyAndCode(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &c, "LOWER(currency_code) = ? and LOWER(country_code) = ?", strings.ToLower(c.CurrencyCode), strings.ToLower(c.CountryCode))
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (c *Country) FindCountryByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &c, "id = ?", c.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (c *Country) GetFirstCountry(db *gorm.DB) error {
	err := postgresql.SelectFirstFromDb(db, &c)
	if err != nil {
		return fmt.Errorf("country selection failed: %v", err.Error())
	}
	return nil
}

func (c *Country) CreateCountry(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &c)
	if err != nil {
		return fmt.Errorf("country creation failed: %v", err.Error())
	}
	return nil
}

func AddCountriesIfNotExist(db *gorm.DB) error {
	countries := []Country{
		{
			Name:         "nigeria",
			CountryCode:  "NG",
			CurrencyCode: "NGN",
		}, {
			Name:         "united states of america",
			CountryCode:  "USA",
			CurrencyCode: "USD",
		},
	}

	for _, v := range countries {
		_, err := v.FindWithNameOrCode(db)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			v.CreateCountry(db)
		}
	}
	return nil
}
