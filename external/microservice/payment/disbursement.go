package payment

import (
	"strconv"

	"github.com/vesicash/auth-ms/external"
	"github.com/vesicash/auth-ms/external/external_models"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/utility"
	"gorm.io/gorm"
)

func GetDisbursement(authDb *gorm.DB, accountID int) ([]external_models.Disbursements, error) {
	logger := utility.NewLogger()
	var (
		accessToken      = models.AccessToken{}
		outBoundResponse external_models.GetDisbursement
	)
	err := accessToken.GetAccessTokens(authDb)
	if err != nil {
		logger.Info("get disbursements", outBoundResponse, err)
		return outBoundResponse.Data, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}
	data := external_models.AccountIDModel{AccountId: accountID}
	logger.Info("get disbursements", data)
	err = external.SendRequest(logger, "service", "get_disbursements", headers, data, &outBoundResponse, "/"+strconv.Itoa(accountID))
	if err != nil {
		logger.Info("get disbursements", outBoundResponse, err)
		return outBoundResponse.Data, err
	}
	logger.Info("get disbursements", outBoundResponse)

	return outBoundResponse.Data, nil
}
