package auth_model

import (
	"fmt"
	"net/http"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func GetUserCredentialsService(req models.GetUserCredentialModel, db postgresql.Databases) (*models.UsersCredential, int, error) {
	userCredential := models.UsersCredential{ID: req.ID, AccountID: int(req.AccountID), IdentificationType: req.IdentificationType}

	if (req.AccountID != 0 && req.IdentificationType == "") || (req.AccountID == 0 && req.IdentificationType != "") {
		return &models.UsersCredential{}, http.StatusBadRequest, fmt.Errorf("missing either account_id or identification_type")
	}

	if req.ID != 0 {
		code, err := userCredential.GetUserCredentialByID(db.Auth)
		if err != nil {
			return &models.UsersCredential{}, code, err
		}
	} else if req.AccountID != 0 && req.IdentificationType != "" {
		code, err := userCredential.GetUserCredentialByAccountIdAndType(db.Auth)
		if err != nil {
			return &models.UsersCredential{}, code, err
		}
	} else {
		return &models.UsersCredential{}, http.StatusBadRequest, fmt.Errorf("error occured please check your input")
	}

	return &userCredential, http.StatusOK, nil
}
func CreateUserCredentialsService(req models.CreateUserCredentialModel, db postgresql.Databases) (*models.UsersCredential, int, error) {
	userCredential := models.UsersCredential{
		AccountID:          int(req.AccountID),
		Bvn:                req.Bvn,
		IdentificationType: req.IdentificationType,
		IdentificationData: req.IdentificationData,
	}

	_, err := userCredential.GetUserCredentialByAccountIdAndType(db.Auth)
	if err == nil {
		userCredential.Bvn = req.Bvn
		userCredential.IdentificationData = req.IdentificationType
		err := userCredential.Update(db.Auth)
		if err != nil {
			return &models.UsersCredential{}, http.StatusInternalServerError, err
		}
		return &userCredential, http.StatusOK, err
	}

	err = userCredential.CreateUsersCredential(db.Auth)
	if err != nil {
		return &models.UsersCredential{}, http.StatusInternalServerError, err
	}

	return &userCredential, http.StatusOK, nil
}
func UpdateUserCredentialsService(req models.UpdateUserCredentialModel, db postgresql.Databases) (*models.UsersCredential, int, error) {
	userCredential := models.UsersCredential{
		ID: req.ID,
	}
	code, err := userCredential.GetUserCredentialByID(db.Auth)
	if err != nil {
		return &models.UsersCredential{}, code, err
	}

	if req.AccountID != 0 {
		userCredential.AccountID = int(req.AccountID)
	}

	if req.IdentificationType != "" {
		userCredential.IdentificationType = req.IdentificationType
	}

	if req.Bvn != "" {
		userCredential.Bvn = req.Bvn
	}

	if req.IdentificationData != "" {
		userCredential.IdentificationData = req.IdentificationData
	}

	err = userCredential.Update(db.Auth)
	if err != nil {
		return &models.UsersCredential{}, http.StatusInternalServerError, err
	}

	return &userCredential, http.StatusOK, nil
}
