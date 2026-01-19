package main

import (
	"context"
	"log"

	"github.com/CogitoNTNU/cogi-go/internal/config"
	"github.com/CogitoNTNU/cogi-go/internal/repository/db"
	_ "github.com/CogitoNTNU/cogi-go/internal/repository/db"
	"github.com/CogitoNTNU/cogi-go/internal/util/env"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadApiConfig()

	logger := logrus.New().WithField("app", cfg.AppName).WithContext(ctx)

	logger.Info("Migration has started \n")

	env := env.Configure(".env")

	database, err := db.InitDb(env, logger, ctx)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to db")
	}

	defer database.Close()

	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create postgres driver:", err)
	}

	migrate, err := migrate.NewWithDatabaseInstance(
		"file://internal/repository/db/migrations",
		"postgres", driver)

	if err != nil {
		log.Fatal(err)
	}
	if err := migrate.Up(); err != nil {
		log.Fatal(err)
	}
}
