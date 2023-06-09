package migrations

import (
	"github.com/vesicash/auth-ms/internal/models"
)

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.AccessToken{},
		models.Authorize{},
		models.BankDetail{},
		models.Bank{},
		models.BannedAccount{},
		models.BusinessCharge{},
		models.BusinessProfile{},
		models.BusinessType{},
		models.ContactUs{},
		models.Country{},
		models.EscrowCharge{},
		models.OtpVerification{},
		models.PasswordResetToken{},
		models.ReferralPromo{},
		models.UserAccountUpgrade{},
		models.UserProfile{},
		models.UserTracking{},
		models.User{},
		models.UsersCredential{},
		models.WalletBalance{},
		models.WalletHistory{},
		models.WalletTransaction{},
	}
}
