package auth_model

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func GetBusinessChargeService(req models.GetBusinessChargeModel, db postgresql.Databases) (*models.BusinessCharge, int, error) {
	var (
		businessCharge = models.BusinessCharge{
			ID:         req.ID,
			BusinessId: int(req.BusinessID),
			Currency:   req.Currency,
			Country:    req.Country,
		}
	)

	if req.ID != 0 {
		code, err := businessCharge.GetByID(db.Auth)
		if err != nil {
			return &models.BusinessCharge{}, code, err
		}
	} else if req.BusinessID != 0 {
		if req.Country == "" && req.Currency == "" {
			return &models.BusinessCharge{}, http.StatusBadRequest, fmt.Errorf("specify either country or currency")
		}
		code, err := businessCharge.GetByBusinessIDAndOthers(db.Auth)
		if err != nil {
			return &models.BusinessCharge{}, code, err
		}

	} else {
		return &models.BusinessCharge{}, http.StatusBadRequest, fmt.Errorf("error occured please check your input")
	}

	return &businessCharge, http.StatusOK, nil
}

func InitBusinessChargeService(req models.InitBusinessChargeModel, db postgresql.Databases) (*models.BusinessCharge, int, error) {
	var (
		businessCharge = models.BusinessCharge{
			BusinessId: int(req.BusinessID),
			Currency:   req.Currency,
		}
	)

	country := models.Country{CurrencyCode: req.Currency}
	code, err := country.FindWithCurrency(db.Auth)
	if err != nil {
		return &models.BusinessCharge{}, code, err
	}

	businessCharge.Country = country.CountryCode
	businessCharge.BusinessCharge = "2.5"
	businessCharge.VesicashCharge = "1"
	businessCharge.ProcessingFee = "0"
	businessCharge.DisbursementCharge = "0"

	if strings.ToUpper(country.CountryCode) == "NG" {
		businessCharge.PaymentGateway = "rave"
		businessCharge.DisbursementCharge = "rave"
		businessCharge.BusinessCharge = "2.5"
	} else {
		businessCharge.PaymentGateway = "rave"
		businessCharge.DisbursementCharge = "rave_momo"
		businessCharge.BusinessCharge = "5"
	}

	err = businessCharge.CreateBusinessCharge(db.Auth)
	if err != nil {
		return &models.BusinessCharge{}, http.StatusInternalServerError, err
	}

	return &businessCharge, http.StatusOK, nil
}
