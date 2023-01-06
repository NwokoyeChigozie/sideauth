package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/vesicash/auth-ms/external/microservice/notification"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
)

func UpdatePassword(db postgresql.Databases, accountID int, oldPassword, newPassword string) (int, error) {
	user := models.User{AccountID: uint(accountID)}
	code, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		return code, err
	}

	if !utility.CompareHash(oldPassword, user.Password) {
		return http.StatusBadRequest, fmt.Errorf("incorrect password")
	}

	password, err := utility.Hash(newPassword)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	user.Password = password
	err = user.Update(db.Auth)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, err
}

func RequestPasswordResetService(logger *utility.Logger, db postgresql.Databases, email, phoneNumber string) (int, int, error) {
	if email == "" && phoneNumber == "" {
		return 0, http.StatusBadRequest, fmt.Errorf("provide either email address or phone number")
	}

	user := models.User{EmailAddress: email, PhoneNumber: phoneNumber}
	code, err := user.GetUserByUsernameEmailOrPhone(db.Auth)
	if err != nil {
		return 0, code, err
	}

	token := utility.GetRandomNumbersInRange(100000000, 999999999)
	resetToken := models.PasswordResetToken{
		AccountID: int(user.AccountID),
		Token:     strconv.Itoa(token),
		ExpiresAt: strconv.Itoa(int(time.Now().Add(30 * time.Minute).Unix())),
	}
	err = resetToken.CreatePasswordResetToken(db.Auth)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	if email != "" {
		notification.SendEmailPasswordReset(logger, db.Auth, int(user.AccountID), token)
	}

	if phoneNumber != "" {
		notification.SendPhonePasswordReset(logger, db.Auth, int(user.AccountID), token)
	}

	return int(user.AccountID), http.StatusOK, nil

}

func UpdatePasswordWithTokenService(logger *utility.Logger, db postgresql.Databases, accountID int, token int, password string) (int, error) {
	if password == "" {
		return http.StatusBadRequest, fmt.Errorf("password is empty")
	}

	if token == 0 {
		return http.StatusBadRequest, fmt.Errorf("invalid token")
	}

	user := models.User{AccountID: uint(accountID)}
	code, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		return code, err
	}

	resetToken := models.PasswordResetToken{AccountID: int(user.AccountID), Token: strconv.Itoa(token)}
	code, err = resetToken.GetLatestByAccountIDAndToken(db.Auth)
	if err != nil {
		if code == http.StatusInternalServerError {
			return code, err
		}
		return code, fmt.Errorf("invalid token")
	}

	if resetToken.ExpiresAt == "" {
		return http.StatusBadRequest, fmt.Errorf("invalid token")
	}

	parseInt, err := strconv.Atoi(resetToken.ExpiresAt)
	if err != nil {
		return http.StatusBadRequest, err
	}

	unixTimeUTC := time.Unix(int64(parseInt), 0)
	if time.Now().After(unixTimeUTC) {
		return http.StatusBadRequest, fmt.Errorf("expired token")
	}

	newPassword, err := utility.Hash(password)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	user.Password = newPassword
	err = user.Update(db.Auth)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	notification.SendEmailPasswordDoneReset(logger, db.Auth, int(user.AccountID))
	notification.SendPhonePasswordDoneReset(logger, db.Auth, int(user.AccountID))

	err = resetToken.DeletePasswordResetToken(db.Auth)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, err
}
