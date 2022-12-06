package otp

import (
	"net/http"
	"strconv"

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
