package auth_model

import (
	"fmt"
	"net/http"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func GetCountryService(req models.GetCountryModel, db postgresql.Databases) (*models.Country, int, error) {
	country := models.Country{ID: req.ID, Name: req.Name, CountryCode: req.CountryCode, CurrencyCode: req.CurrencyCode}

	if req.Name == "" && req.ID == 0 && req.CountryCode == "" && req.CurrencyCode == "" {
		return &models.Country{}, http.StatusBadRequest, fmt.Errorf("no input value, try id, name, country_code, currency_code")
	}

	if req.ID != 0 {
		code, err := country.FindCountryByID(db.Auth)
		if err != nil {
			return &models.Country{}, code, err
		}
	} else if req.CurrencyCode != "" && req.CountryCode != "" {
		code, err := country.FindWithCurrencyAndCode(db.Auth)
		if err != nil {
			return &models.Country{}, code, err
		}
	} else if req.CurrencyCode != "" && req.CountryCode == "" {
		code, err := country.FindWithCurrency(db.Auth)
		if err != nil {
			return &models.Country{}, code, err
		}
	} else if req.Name != "" {
		code, err := country.FindWithNameOrCode(db.Auth)
		if err != nil {
			return &models.Country{}, code, err
		}
	} else {
		if req.CurrencyCode != "" && req.CountryCode == "" {
			return &models.Country{}, http.StatusBadRequest, fmt.Errorf("country_code is needed")
		}
		return &models.Country{}, http.StatusBadRequest, fmt.Errorf("error occured please check your input")
	}

	return &country, http.StatusOK, nil
}
