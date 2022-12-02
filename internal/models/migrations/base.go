package migrations

import (
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func RunAllMigrations(db postgresql.Databases) {

	// auth migration
	MigrateModels(db.Auth, AuthMigrationModels())

}

func MigrateModels(db *gorm.DB, models []interface{}) {
	_ = db.AutoMigrate(models...)
}
