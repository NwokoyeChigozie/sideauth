package migrations

import (
	"github.com/vesicash/auth-ms/internal/models"
)

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.User{},
		models.Country{},
		models.UserProfile{},
		models.AccessToken{},
		models.BusinessProfile{},
		models.BusinessCharge{},
		models.PasswordResetToken{},
		models.ReferralPromo{},
		models.UsersCredential{},
	}
}
