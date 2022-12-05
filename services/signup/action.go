package signup

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/vesicash/auth-ms/external/notification"
	"github.com/vesicash/auth-ms/external/referral"
	"github.com/vesicash/auth-ms/external/verification"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
	"gorm.io/gorm"
)

func SignupService(req models.CreateUserRequestModel, db postgresql.Databases) (*models.User, int, error) {
	var (
		countryName         = strings.ToLower(req.Country)
		accountType         = req.AccountType
		firstname           = strings.Title(strings.ToLower(req.Firstname))
		lastname            = strings.Title(strings.ToLower(req.Lastname))
		emailAddress        = strings.ToLower(req.EmailAddress)
		phoneNumber         = req.PhoneNumber
		username            = req.Username
		password            = req.Password
		accountID           = 0
		currency            = ""
		countryCode         = ""
		webhookUri          = req.WebhookURI
		paymentGateway      = ""
		disbursementGateway = ""
		businessPercentage  float32
		vesicashPercentage  float32
		processingFee       float32
	)

	if accountType == "" {
		accountType = "individual"
	}

	country := models.Country{Name: countryName}
	code, err := country.FindWithNameOrCode(db.Auth)
	if err != nil {
		return nil, code, err
	}

	password, err = utility.Hash(req.Password)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	currency, countryCode = country.CurrencyCode, country.CountryCode

	accountID, err = GetAccountID(db.Auth)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	user := models.User{
		AccountID:    uint(accountID),
		AccountType:  accountType,
		Firstname:    firstname,
		Lastname:     lastname,
		EmailAddress: emailAddress,
		PhoneNumber:  phoneNumber,
		Username:     username,
		Password:     password,
		BusinessId:   req.BusinessID,
	}

	err = user.CreateUser(db.Auth)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	userProfile := models.UserProfile{
		AccountID: int(user.AccountID),
		Country:   countryCode,
		Currency:  currency,
	}

	err = userProfile.CreateUserProfile(db.Auth)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	userCredential := models.UsersCredential{
		AccountID: int(user.AccountID),
	}

	err = userCredential.CreateUsersCredential(db.Auth)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if user.AccountType == "business" {
		businessProfile := models.BusinessProfile{
			AccountID:       int(user.AccountID),
			BusinessName:    req.BusinessName,
			BusinessType:    req.BusinessType,
			BusinessAddress: req.BusinessAddress,
			Country:         countryCode,
			Currency:        currency,
			Webhook_uri:     webhookUri,
		}

		err = businessProfile.CreateBusinessProfile(db.Auth)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		paymentGateway, disbursementGateway = GetPaymentAndDisbursementGateway(countryCode)
		businessPercentage, vesicashPercentage, processingFee = GetBusinessVesicashAndProcessingFees(req.BusinessType)

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
			return nil, http.StatusInternalServerError, err
		}
	}

	if req.EmailAddress != "" {
		notification.SendWelcomeNotification(db.Auth, int(user.AccountID))

		if req.Password == "" {
			notification.SendWelcomePasswordReset(db.Auth, int(user.AccountID), utility.GetRandomNumbersInRange(100000000, 999999999))
		}

		verification.SendVerificationEmail(db.Auth, int(user.AccountID))
	}

	if req.ReferralCode != "" {
		_, err := referral.CreateReferralRequest(db.Auth, accountID, req.ReferralCode)
		if err != nil {
			fmt.Println(err)
		}

		promo := models.ReferralPromo{ReferralCode: req.ReferralCode}
		_, err = promo.GetReferralPromoByCode(db.Auth)
		if err == nil {
			promo.ActivatePromoCode(db.Auth, int(user.AccountID))
		}

	}

	if req.PhoneNumber != "" {
		notification.SendWelcomeSmsNotification(db.Auth, int(user.AccountID))
		verification.SendVerificationSms(db.Auth, int(user.AccountID))
	}

	return &user, http.StatusCreated, nil
}

func BulkSignupService(req []models.CreateUserRequestModel, db postgresql.Databases) ([]*models.User, int, error) {
	logger := utility.NewLogger()
	newUsers := []*models.User{}
	for _, sData := range req {
		newUser, _, err := SignupService(sData, db)
		if err != nil {
			logger.Error("bulk signup", err)
		} else {
			newUsers = append(newUsers, newUser)
		}

	}
	return newUsers, http.StatusOK, nil

}

func GetAccountID(db *gorm.DB) (int, error) {
	randNum := utility.GetRandomNumbersInRange(1000000000, 9999999999)
	user := models.User{AccountID: uint(randNum)}
	_, err := user.GetUserByAccountID(db)
	if err == nil {
		return GetAccountID(db)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return randNum, nil
	} else {
		return 0, err
	}
}

func GetPaymentAndDisbursementGateway(country string) (payment_gateway, disbursement_gateway string) {
	switch country {
	case "NG":
		payment_gateway, disbursement_gateway = "rave", "rave"
	case "Nigeria":
		payment_gateway, disbursement_gateway = "rave", "rave"
	default:
		payment_gateway, disbursement_gateway = "rave", "rave_momo"
	}
	return
}

func GetBusinessVesicashAndProcessingFees(businessType string) (businessPercentage, vesicashPercentage, processingFee float32) {
	switch businessType {
	case "social_commerce":
		businessPercentage, vesicashPercentage, processingFee = 0, 2.5, 100
	default:
		businessPercentage, vesicashPercentage, processingFee = 0, 2.5, 0
	}
	return
}
