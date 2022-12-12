package auth

import (
	"net/http"
	"strconv"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func UpgradeAccountService(db postgresql.Databases, accountID int, businessType, businessName, webhookUri string) (models.User, int, error) {
	user := models.User{AccountID: uint(accountID)}
	code, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		return user, code, err
	}

	userProfile := models.UserProfile{AccountID: int(user.AccountID)}
	code, err = userProfile.GetByAccountID(db.Auth)
	if err != nil {
		return user, code, err
	}

	countryCode := userProfile.Country
	if countryCode == "" {
		countryCode = "NG"
	}

	country := models.Country{Name: countryCode}
	code, err = country.FindWithNameOrCode(db.Auth)
	if err != nil {
		return user, code, err
	}

	businessProfile := models.BusinessProfile{
		AccountID:    int(user.AccountID),
		BusinessName: businessName,
		BusinessType: businessType,
		Country:      countryCode,
		Currency:     country.CurrencyCode,
		Webhook_uri:  webhookUri,
	}
	code, err = businessProfile.GetByAccountID(db.Auth)
	if err != nil {
		if code == http.StatusInternalServerError {
			return user, code, err
		}

		err := businessProfile.CreateBusinessProfile(db.Auth)
		if err != nil {
			return user, http.StatusInternalServerError, err
		}
	}

	user.AccountType = "business"
	err = user.Update(db.Auth)
	if err != nil {
		return user, http.StatusInternalServerError, err
	}
	paymentGateway, disbursementGateway := GetPaymentAndDisbursementGateway(countryCode)
	businessPercentage, vesicashPercentage, processingFee := GetBusinessVesicashAndProcessingFees(businessType)
	businessCharge := models.BusinessCharge{
		BusinessId:          int(user.AccountID),
		Country:             countryCode,
		BusinessCharge:      strconv.Itoa(int(businessPercentage)),
		VesicashCharge:      strconv.Itoa(int(vesicashPercentage)),
		ProcessingFee:       strconv.Itoa(int(processingFee)),
		DisbursementCharge:  "0",
		PaymentGateway:      paymentGateway,
		DisbursementGateway: disbursementGateway,
	}

	err = businessCharge.CreateBusinessCharge(db.Auth)
	if err != nil {
		return user, http.StatusInternalServerError, err
	}

	accountUpgrade := models.UserAccountUpgrade{AccountID: int(user.AccountID), BusinessType: businessType}
	err = accountUpgrade.CreateUserAccountUpgrade(db.Auth)
	if err != nil {
		return user, http.StatusInternalServerError, err
	}

	return user, http.StatusOK, nil
}
