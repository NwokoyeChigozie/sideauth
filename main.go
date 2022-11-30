package main

import (
	"log"

	"github.com/vesicash/auth-ms/internal/config"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"

	"github.com/go-playground/validator/v10"
	"github.com/vesicash/auth-ms/utility"

	"github.com/vesicash/auth-ms/pkg/router"
)

func init() {
	config.Setup()
	postgresql.ConnectToDatabases()

}

func main() {
	//Load config
	logger := utility.NewLogger()
	getConfig := config.GetConfig()
	validatorRef := validator.New()
	r := router.Setup(validatorRef, logger)

	utility.LogAndPrint("Server is starting at 127.0.0.1:%s", getConfig.Server.Port)
	log.Fatal(r.Run(":" + getConfig.Server.Port))
}
