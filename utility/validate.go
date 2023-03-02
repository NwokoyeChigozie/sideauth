package utility

import (
	"net/mail"
	"os"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

func EmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func PhoneValid(phone string) (string, bool) {
	parsed, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return strings.ReplaceAll(phone, " ", ""), false
	}

	if !phonenumbers.IsValidNumber(parsed) {
		return strings.ReplaceAll(phone, " ", ""), false
	}

	formattedNum := phonenumbers.Format(parsed, phonenumbers.NATIONAL)

	return strings.ReplaceAll(formattedNum, " ", ""), true
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
