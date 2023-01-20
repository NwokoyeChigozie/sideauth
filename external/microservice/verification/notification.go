package verification

import (
	"github.com/vesicash/auth-ms/external"
	"github.com/vesicash/auth-ms/external/external_models"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/utility"
	"gorm.io/gorm"
)

func SendVerificationEmail(logger *utility.Logger, authDb *gorm.DB, accountID int) error {
	var (
		accessToken      = models.AccessToken{}
		outBoundResponse map[string]interface{}
	)
	err := accessToken.GetAccessTokens(authDb)
	if err != nil {
		logger.Info("verification email", outBoundResponse, err)
		return err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}
	data := map[string]interface{}{"account_id": accountID}
	logger.Info("verification email", data)
	err = external.SendRequest(logger, "service", "verification_email", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("verification email", outBoundResponse, err)
		return err
	}
	logger.Info("verification email", outBoundResponse)

	return nil
}

func SendVerificationSms(logger *utility.Logger, authDb *gorm.DB, accountID int) error {
	var (
		accessToken      = models.AccessToken{}
		outBoundResponse map[string]interface{}
	)
	err := accessToken.GetAccessTokens(authDb)
	if err != nil {
		logger.Info("verification email", outBoundResponse, err)
		return err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}
	data := external_models.AccountIDModel{AccountId: accountID}
	logger.Info("welcome email", data)
	err = external.SendRequest(logger, "service", "verification_sms", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("verification email", outBoundResponse, err)
		return err
	}
	logger.Info("verification email", outBoundResponse)

	return nil
}
