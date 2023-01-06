package verification

import (
	"github.com/vesicash/auth-ms/external"
	"github.com/vesicash/auth-ms/external/external_models"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/utility"
	"gorm.io/gorm"
)

func GetVerifications(logger *utility.Logger, authDb *gorm.DB, accountID int) ([]external_models.Verifications, error) {
	var (
		accessToken      = models.AccessToken{}
		outBoundResponse external_models.GetVerifications
	)
	err := accessToken.GetAccessTokens(authDb)
	if err != nil {
		logger.Info("get verifications", outBoundResponse, err)
		return outBoundResponse.Data, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}
	data := external_models.AccountIDModel{AccountId: accountID}
	logger.Info("get verifications", data)
	err = external.SendRequest(logger, "service", "get_verifications", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("get verifications", outBoundResponse, err)
		return outBoundResponse.Data, err
	}
	logger.Info("get verifications", outBoundResponse)

	return outBoundResponse.Data, nil
}
