package auth_model

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func GetAuthorizeService(req models.GetAuthorizeModel, db postgresql.Databases) (*models.Authorize, int, error) {
	authorize := models.Authorize{
		ID:         req.ID,
		AccountID:  int(req.AccountID),
		Authorized: req.Authorized,
		IpAddress:  req.IpAddress,
		Browser:    req.Browser,
	}

	if req.AccountID == 0 && req.ID == 0 {
		return &models.Authorize{}, http.StatusBadRequest, fmt.Errorf("enter either account_id or id")
	}

	if req.ID != 0 {
		code, err := authorize.GetAuthorizeByID(db.Auth)
		if err != nil {
			return &models.Authorize{}, code, err
		}
	} else if req.AccountID != 0 {
		code, err := authorize.GetAuthorizeByAccountIDAndOneOrAll(db.Auth)
		if err != nil {
			return &models.Authorize{}, code, err
		}
	} else {
		return &models.Authorize{}, http.StatusBadRequest, fmt.Errorf("error occured please check your input")
	}

	return &authorize, http.StatusOK, nil
}

func CreateAuthorizeService(req models.CreateAuthorizeModel, db postgresql.Databases) (*models.Authorize, int, error) {
	authorize := models.Authorize{
		AccountID:  int(req.AccountID),
		Authorized: req.Authorized,
		Token:      req.Token,
		IpAddress:  req.IpAddress,
		Browser:    req.Browser,
		Os:         req.Os,
		Location:   req.Location,
		Attempt:    1,
	}

	code, err := authorize.GetAuthorizeByAccountIDAndOneOrAll2(db.Auth)
	if err == nil {
		authorize.AccountID = int(req.AccountID)
		if !authorize.Authorized {
			authorize.Authorized = req.Authorized
		}

		authorize.Token = req.Token
		authorize.IpAddress = req.IpAddress
		authorize.Browser = req.Browser
		authorize.Os = req.Os
		authorize.Location = req.Location
		authorize.Attempt = authorize.Attempt + 1

		if req.Authorized {
			authorize.AuthorizedAt = time.Now()
		}

		err := authorize.Update(db.Auth)
		if err != nil {
			return &models.Authorize{}, http.StatusInternalServerError, err
		}
		return &authorize, http.StatusOK, err
	}

	if code == http.StatusInternalServerError {
		return &models.Authorize{}, code, err
	}

	if req.Authorized {
		authorize.AuthorizedAt = time.Now()
	}
	err = authorize.CreateAuthorize(db.Auth)
	if err != nil {
		return &models.Authorize{}, http.StatusInternalServerError, err
	}

	return &authorize, http.StatusOK, nil
}

func UpdateAuthorizeService(req models.UpdateAuthorizeModel, db postgresql.Databases) (*models.Authorize, int, error) {
	authorize := models.Authorize{
		ID: req.ID,
	}

	code, err := authorize.GetAuthorizeByID(db.Auth)
	if err != nil {
		return &models.Authorize{}, code, err
	}

	if req.AccountID != 0 {
		authorize.AccountID = int(req.AccountID)
	}

	if req.Token != "" {
		authorize.Token = req.Token
	}

	if req.IpAddress != "" {
		authorize.IpAddress = req.IpAddress
	}

	if req.Browser != "" {
		authorize.Browser = req.Browser
	}

	if req.Os != "" {
		authorize.Os = req.Os
	}

	if req.Location != "" {
		authorize.Location = req.Location
	}

	if req.Attempt != 0 {
		authorize.Attempt = req.Attempt
	}

	if !authorize.Authorized {
		authorize.Authorized = req.Authorized
	}

	if req.Authorized {
		authorize.AuthorizedAt = time.Now()
	}
	err = authorize.Update(db.Auth)
	if err != nil {
		return &models.Authorize{}, http.StatusInternalServerError, err
	}

	return &authorize, http.StatusOK, nil
}
