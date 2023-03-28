package main

import (
	"github.com/vesicash/auth-ms/internal/config"
	"github.com/vesicash/auth-ms/internal/models/migrations"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"

	"github.com/vesicash/auth-ms/utility"

	"github.com/go-playground/validator/v10"
	"github.com/vesicash/auth-ms/pkg/router"
)

var (
	g errgroup.Group
)

func main() {
	logger := utility.NewLogger() //Warning !!!!! Do not recreate this action anywhere on the app

	configuration := config.Setup(logger, "./app")

	postgresql.ConnectToDatabases(logger, configuration.Databases)
	validatorRef := validator.New()
	db := postgresql.Connection()

	if configuration.Databases.Migrate {
		migrations.RunAllMigrations(db)
	}

	r := router.Setup(logger, validatorRef, db, &configuration.App)
	rM := router.SetupMetrics(&configuration.App)

	server := &http.Server{
		Addr:    ":" + configuration.Server.Port,
		Handler: r,
	}

	metricsServer := &http.Server{
		Addr:    ":" + configuration.Server.MetricsPort,
		Handler: rM,
	}

	g.Go(func() error {
		utility.LogAndPrint(logger, "Server is starting at 127.0.0.1:%s", configuration.Server.Port)
		return server.ListenAndServe()
	})

	g.Go(func() error {
		utility.LogAndPrint(logger, "Metrics Server is starting at 127.0.0.1:%s", "8015")
		return metricsServer.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
