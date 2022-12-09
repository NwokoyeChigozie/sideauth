package notification

import (
	"github.com/vesicash/auth-ms/external"
	"github.com/vesicash/auth-ms/external/external_models"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/utility"
	"gorm.io/gorm"
)

func SendEmailPasswordReset(authDb *gorm.DB, accountID, token int) error {
	logger := utility.NewLogger()
	var (
		accessToken      = models.AccessToken{}
		outBoundResponse map[string]interface{}
	)
	err := accessToken.GetAccessTokens(authDb)
	if err != nil {
		logger.Info("email password reset", outBoundResponse, err)
		return err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	data := external_models.PhoneEmailVerificationModel{AccountId: accountID, Token: token}
	logger.Info("email password reset", data)
	err = external.SendRequest(logger, "service", "email_password_reset_notification", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("welcome password reset", outBoundResponse, err)
		return err
	}
	logger.Info("email password reset", outBoundResponse)

	return nil
}
func SendPhonePasswordReset(authDb *gorm.DB, accountID, token int) error {
	logger := utility.NewLogger()
	var (
		accessToken      = models.AccessToken{}
		outBoundResponse map[string]interface{}
	)
	err := accessToken.GetAccessTokens(authDb)
	if err != nil {
		logger.Info("phone password reset", outBoundResponse, err)
		return err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	data := external_models.PhoneEmailVerificationModel{AccountId: accountID, Token: token}
	logger.Info("phone password reset", data)
	err = external.SendRequest(logger, "service", "phone_password_reset_notification", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("welcome password reset", outBoundResponse, err)
		return err
	}
	logger.Info("phone password reset", outBoundResponse)

	return nil
}

func SendEmailPasswordDoneReset(authDb *gorm.DB, accountID int) error {
	logger := utility.NewLogger()
	var (
		accessToken      = models.AccessToken{}
		outBoundResponse map[string]interface{}
	)
	err := accessToken.GetAccessTokens(authDb)
	if err != nil {
		logger.Info("email password reset done", outBoundResponse, err)
		return err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	data := external_models.AccountIDModel{AccountId: accountID}
	logger.Info("email password reset done", data)
	err = external.SendRequest(logger, "service", "email_password_reset_done_notification", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("welcome password reset done", outBoundResponse, err)
		return err
	}
	logger.Info("email password reset done", outBoundResponse)

	return nil
}

func SendPhonePasswordDoneReset(authDb *gorm.DB, accountID int) error {
	logger := utility.NewLogger()
	var (
		accessToken      = models.AccessToken{}
		outBoundResponse map[string]interface{}
	)
	err := accessToken.GetAccessTokens(authDb)
	if err != nil {
		logger.Info("phone password reset done", outBoundResponse, err)
		return err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	data := external_models.AccountIDModel{AccountId: accountID}
	logger.Info("phone password reset done", data)
	err = external.SendRequest(logger, "service", "phone_password_reset_done_notification", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("welcome password reset done", outBoundResponse, err)
		return err
	}
	logger.Info("phone password reset done", outBoundResponse)

	return nil
}
