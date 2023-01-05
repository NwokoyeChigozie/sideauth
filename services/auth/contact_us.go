package auth

import (
	"net/http"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func ContactUsService(db postgresql.Databases, req models.ContactUsCreateModel) (int, error) {
	contactUs := models.ContactUs{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		WebsiteUrl:   req.WebsiteUrl,
		BusinessType: req.BusinessType,
		Country:      req.Country,
		Message:      req.Message,
	}
	err := contactUs.CreateContactUs(db.Auth)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
