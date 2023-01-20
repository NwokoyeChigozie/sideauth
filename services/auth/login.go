package auth

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/external/microservice/verification"
	"github.com/vesicash/auth-ms/external/thirdparty"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/middleware"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
)

func LoginService(c *gin.Context, logger *utility.Logger, req models.LoginUserRequestModel, db postgresql.Databases) (map[string]interface{}, int, error) {
	var (
		responseData = gin.H{}
	)

	if req.EmailAddress == "" && req.PhoneNumber == "" && req.Username == "" {
		return responseData, http.StatusBadRequest, fmt.Errorf("provide either username, email_address, or phone_number")
	}
	user := models.User{Username: req.Username, EmailAddress: req.EmailAddress, PhoneNumber: req.PhoneNumber}
	code, err := user.GetUserByUsernameEmailOrPhone(db.Auth)
	if err != nil {
		if code == http.StatusBadRequest {
			return responseData, code, fmt.Errorf("invalid login details")
		}
		return responseData, code, err
	}

	bannedAccount := models.BannedAccount{AccountID: int(user.AccountID)}
	status, err := bannedAccount.CheckByAccountID(db.Auth)
	if err != nil {
		return responseData, http.StatusInternalServerError, err
	}

	if status {
		return responseData, http.StatusBadRequest, fmt.Errorf("this account has been banned")
	}

	if !utility.CompareHash(req.Password, user.Password) {
		return responseData, http.StatusBadRequest, fmt.Errorf("invalid login details")
	}

	TrackUserLogin(c, logger, db, int(user.AccountID))

	return LoginResponse(logger, user, db, req)
}

func PhoneOtpLogin(c *gin.Context, logger *utility.Logger, phoneNumber string, db postgresql.Databases) (int, int, error) {
	var (
		accountID int
	)

	user := models.User{PhoneNumber: phoneNumber}
	code, err := user.GetUserByUsernameEmailOrPhone(db.Auth)
	if err != nil {
		if code == http.StatusBadRequest {
			return accountID, code, fmt.Errorf("user with phone number not found")
		}
		return accountID, code, err
	}
	accountID = int(user.AccountID)

	bannedAccount := models.BannedAccount{AccountID: int(user.AccountID)}
	status, err := bannedAccount.CheckByAccountID(db.Auth)
	if err != nil {
		return accountID, http.StatusInternalServerError, err
	}

	if status {
		return accountID, http.StatusBadRequest, fmt.Errorf("this account has been banned")
	}

	otpReq := models.SendOtpTokenReq{AccountID: int(user.AccountID)}
	SendOtpService(logger, otpReq, db)

	TrackUserLogin(c, logger, db, int(user.AccountID))
	return accountID, http.StatusOK, nil
}

func TrackUserLogin(c *gin.Context, logger *utility.Logger, db postgresql.Databases, accountID int) error {
	var (
		ipAddress    = c.ClientIP()
		browser      = c.Request.UserAgent()
		geo_location = ""
	)

	data, err := thirdparty.GetIpInfo(logger, ipAddress)
	if err != nil {
		return err
	}
	// outBoundResponse["geoplugin_regionName"], outBoundResponse["geoplugin_countryName"]
	if data["geoplugin_regionName"] != nil {
		geo_location = data["geoplugin_regionName"].(string)
	}

	if data["geoplugin_countryName"] != nil {
		if geo_location != "" {
			geo_location += ", " + data["geoplugin_countryName"].(string)
		} else {
			geo_location = data["geoplugin_countryName"].(string)
		}
	}

	userTracking := models.UserTracking{
		AccountID: accountID,
		IpAddress: ipAddress,
		Browser:   browser,
		Location:  geo_location,
	}
	err = userTracking.CreateUserTracking(db.Auth)
	if err != nil {
		return err
	}
	return nil
}

func LoginResponse(logger *utility.Logger, user models.User, db postgresql.Databases, req models.LoginUserRequestModel) (map[string]interface{}, int, error) {
	var (
		responseData = gin.H{}
	)

	token, err := middleware.CreateToken(user, false)
	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error creating token: " + err.Error())
	}

	user.LoginAccessToken = token.AccessToken
	user.LoginAccessTokenExpiresIn = strconv.Itoa(int(token.AtExpiresTime.Unix()))
	err = user.Update(db.Auth)
	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}

	verifications, _ := verification.GetVerifications(logger, db.Auth, int(user.AccountID), user.LoginAccessToken)

	tracking := models.UserTracking{AccountID: int(user.AccountID)}
	trackings, err := tracking.GetAllByAccountID(db.Auth)
	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error getting tracking records: " + err.Error())
	}

	trackingCount := len(trackings)

	userProfile := models.UserProfile{AccountID: int(user.AccountID)}
	code, err := userProfile.GetByAccountID(db.Auth)
	if code == http.StatusInternalServerError {
		return responseData, code, fmt.Errorf("error getting user profile: " + err.Error())
	}

	businessProfile := models.BusinessProfile{AccountID: int(user.AccountID)}
	code, err = businessProfile.GetByAccountID(db.Auth)
	if code == http.StatusInternalServerError {
		return responseData, code, fmt.Errorf("error getting business profile: " + err.Error())
	}

	bankDetail := models.BankDetail{AccountID: int(user.AccountID)}
	bankDetails, err := bankDetail.GetAllByAccountID(db.Auth)
	if err != nil {
		return responseData, code, fmt.Errorf("error getting bank details: " + err.Error())
	}

	escrowCharge := models.EscrowCharge{BusinessID: int(user.AccountID)}
	code, err = escrowCharge.GetByBusinessID(db.Auth)
	if code == http.StatusInternalServerError {
		return responseData, code, fmt.Errorf("error getting escrow charge: " + err.Error())
	}

	businessCharge := models.BusinessCharge{}
	businessCharges, err := businessCharge.GetAllByBusinessID(db.Auth)
	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error getting business charges: " + err.Error())
	}

	if req.PhoneNumber != "" {
		otpReq := models.SendOtpTokenReq{AccountID: int(user.AccountID)}
		SendOtpService(logger, otpReq, db)
	}

	return gin.H{
		"token_type":   "auth",
		"expires_in":   token.AtExpiresTime,
		"access_token": token.AccessToken,
		"user":         user,
		"login_count":  trackingCount,
		"profile": gin.H{
			"business":         businessProfile,
			"user":             userProfile,
			"bank_details":     bankDetails,
			"escrow_charge":    escrowCharge,
			"business_charges": businessCharges,
			"verifications":    verifications,
		},
	}, http.StatusOK, nil
}
