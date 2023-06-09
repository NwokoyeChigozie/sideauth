package migrations

import (
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func RunAllMigrations(db postgresql.Databases) {

	// auth migration
	MigrateModels(db.Auth, AuthMigrationModels())

	// add countries
	models.AddCountriesIfNotExist(db.Auth)

}

func MigrateModels(db *gorm.DB, models []interface{}) {
	_ = db.AutoMigrate(models...)
}
