package notification

import (
	"strconv"
	"time"

	"github.com/vesicash/auth-ms/external"
	"github.com/vesicash/auth-ms/external/external_models"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/utility"
	"gorm.io/gorm"
)

func SendWelcomeNotification(logger *utility.Logger, authDb *gorm.DB, accountID int) error {
	var (
		accessToken      = models.AccessToken{}
		outBoundResponse map[string]interface{}
	)
	err := accessToken.GetAccessTokens(authDb)
	if err != nil {
		logger.Error("welcome email", outBoundResponse, err)
		return err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}
	data := external_models.AccountIDModel{AccountId: accountID}
	logger.Info("welcome email", data)
	err = external.SendRequest(logger, "service", "welcome_notification", headers, data, &outBoundResponse)
	if err != nil {
		logger.Error("welcome email", outBoundResponse, err)
		return err
	}
	logger.Info("welcome email", outBoundResponse)

	return nil
}
func SendWelcomeSmsNotification(logger *utility.Logger, authDb *gorm.DB, accountID int) error {
	var (
		accessToken      = models.AccessToken{}
		outBoundResponse map[string]interface{}
	)
	err := accessToken.GetAccessTokens(authDb)
	if err != nil {
		logger.Error("welcome sms", outBoundResponse, err)
		return err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}
	data := external_models.AccountIDModel{AccountId: accountID}
	logger.Info("welcome sms", data)
	err = external.SendRequest(logger, "service", "welcome_sms_notification", headers, data, &outBoundResponse)
	if err != nil {
		logger.Error("welcome sms", outBoundResponse, err)
		return err
	}
	logger.Info("welcome sms", outBoundResponse)

	return nil
}

func SendWelcomePasswordReset(logger *utility.Logger, authDb *gorm.DB, accountID, token int) error {

	var (
		accessToken = models.AccessToken{}
		resetTokens = models.PasswordResetToken{
			AccountID: accountID,
			Token:     strconv.Itoa(token),
			ExpiresAt: strconv.Itoa(int(time.Now().Add(48 * time.Hour).Unix())),
		}
		outBoundResponse map[string]interface{}
	)
	err := accessToken.GetAccessTokens(authDb)
	if err != nil {
		logger.Error("welcome password reset", outBoundResponse, err)
		return err
	}

	err = resetTokens.CreatePasswordResetToken(authDb)
	if err != nil {
		logger.Error("welcome password reset", outBoundResponse, err)
		return err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}
	rToken, _ := strconv.Atoi(resetTokens.Token)
	data := external_models.PhoneEmailVerificationModel{AccountId: resetTokens.AccountID, Token: rToken}
	logger.Info("welcome email", data)
	err = external.SendRequest(logger, "service", "welcome_password_reset_notification", headers, data, &outBoundResponse)
	if err != nil {
		logger.Error("welcome password reset", outBoundResponse, err)
		return err
	}
	logger.Info("welcome password reset", outBoundResponse)

	return nil
}
