package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/external/microservice/notification"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
)

func SendOtpService(req models.SendOtpTokenReq, db postgresql.Databases) (int, error) {
	user := models.User{AccountID: uint(req.AccountID)}
	code, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		return code, err
	}

	otp := models.OtpVerification{AccountID: req.AccountID}
	code, err = otp.GetLatestByAccountID(db.Auth)
	if err != nil {
		if code == http.StatusInternalServerError {
			return code, err
		}
	} else {
		err := otp.Delete(db.Auth)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	token := utility.GetRandomNumbersInRange(100000, 999999)
	otp = models.OtpVerification{
		AccountID: req.AccountID,
		OtpToken:  strconv.Itoa(token),
	}
	err = otp.Create(db.Auth)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	notification.SendOtp(db.Auth, req.AccountID, token)

	return http.StatusOK, nil
}

func ValidateOtpService(c *gin.Context, otp string, accountID int, db postgresql.Databases) (interface{}, int, error) {
	var response interface{}
	otpVerification := models.OtpVerification{AccountID: accountID}
	code, err := otpVerification.GetLatestByAccountID(db.Auth)
	if err != nil {
		if code == http.StatusInternalServerError {
			return response, code, err
		}
		return response, code, fmt.Errorf("invalid otp")
	}

	if otpVerification.OtpToken != otp {
		return response, http.StatusBadRequest, fmt.Errorf("invalid token")
	}

	if time.Now().After(otpVerification.ExpiresAt) {
		return response, http.StatusBadRequest, fmt.Errorf("token expired")
	}

	bannedAccount := models.BannedAccount{AccountID: accountID}
	status, err := bannedAccount.CheckByAccountID(db.Auth)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	if status {
		return response, http.StatusBadRequest, fmt.Errorf("this account has been banned")
	}

	TrackUserLogin(c, db, accountID)

	user := models.User{AccountID: uint(accountID)}
	code, err = user.GetUserByAccountID(db.Auth)
	if err != nil {
		return response, code, err
	}

	return LoginResponse(user, db, models.LoginUserRequestModel{})
}
