package verification

import (
	"github.com/vesicash/auth-ms/external"
	"github.com/vesicash/auth-ms/external/external_models"
	"github.com/vesicash/auth-ms/utility"
	"gorm.io/gorm"
)

func GetVerifications(logger *utility.Logger, authDb *gorm.DB, accountID int, token string) ([]external_models.Verification, error) {
	var (
		outBoundResponse external_models.GetVerifications
	)
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}
	data := external_models.AccountIDModel{AccountId: accountID}
	logger.Info("get verifications", data)
	err := external.SendRequest(logger, "service", "get_verifications", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("get verifications", outBoundResponse, err)
		return outBoundResponse.Data, err
	}
	logger.Info("get verifications", outBoundResponse)

	return outBoundResponse.Data, nil
}
