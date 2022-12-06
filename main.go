package main

import (
	"log"

	"github.com/vesicash/auth-ms/external/microservice/verification"
	"github.com/vesicash/auth-ms/internal/config"
	"github.com/vesicash/auth-ms/internal/models/migrations"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"

	"github.com/vesicash/auth-ms/utility"

	"github.com/go-playground/validator/v10"
	"github.com/vesicash/auth-ms/pkg/router"
)

func init() {
	config.Setup()
	postgresql.ConnectToDatabases()

}

func main() {
	//Load config
	getConfig := config.GetConfig()
	validatorRef := validator.New()
	db := postgresql.Connection()

	if getConfig.Databases.Migrate {
		migrations.RunAllMigrations(db)
	}

	verification.GetVerifications(db.Auth, 3245681126)

	r := router.Setup(validatorRef, db)

	utility.LogAndPrint("Server is starting at 127.0.0.1:%s", getConfig.Server.Port)
	log.Fatal(r.Run(":" + getConfig.Server.Port))
}
