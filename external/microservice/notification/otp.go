package notification

import (
	"github.com/vesicash/auth-ms/external"
	"github.com/vesicash/auth-ms/external/external_models"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/utility"
	"gorm.io/gorm"
)

func SendOtp(logger *utility.Logger, authDb *gorm.DB, accountID, token int) error {
	var (
		accessToken      = models.AccessToken{}
		outBoundResponse map[string]interface{}
	)
	err := accessToken.GetAccessTokens(authDb)
	if err != nil {
		logger.Error("send otp", outBoundResponse, err)
		return err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}
	data := external_models.SendOtpModel{AccountId: accountID, OtpToken: token}
	logger.Info("send otp", data)
	err = external.SendRequest(logger, "service", "send_otp_notification", headers, data, &outBoundResponse)
	if err != nil {
		logger.Error("send otp", outBoundResponse, err)
		return err
	}
	logger.Info("send otp", outBoundResponse)

	return nil
}
