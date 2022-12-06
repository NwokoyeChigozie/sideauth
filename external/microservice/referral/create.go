package referral

import (
	"github.com/vesicash/auth-ms/external"
	"github.com/vesicash/auth-ms/external/external_models"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/utility"
	"gorm.io/gorm"
)

func CreateReferralRequest(authDb *gorm.DB, accountID int, code string) (interface{}, error) {
	logger := utility.NewLogger()
	var (
		accessToken      = models.AccessToken{}
		outBoundResponse map[string]interface{}
	)
	err := accessToken.GetAccessTokens(authDb)
	if err != nil {
		logger.Info("create_referral email", outBoundResponse, err)
		return outBoundResponse, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}
	data := external_models.ReferralCreateModel{AccountId: accountID, ReferralCode: code}
	logger.Info("welcome email", data)
	err = external.SendRequest(logger, "service", "create_referral", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("create_referral email", outBoundResponse, err)
		return outBoundResponse, err
	}
	logger.Info("create_referral email", outBoundResponse)

	return outBoundResponse, nil

}
