package main

import (
	"context"

	"github.com/CogitoNTNU/cogi-go/internal/config"
	"github.com/CogitoNTNU/cogi-go/internal/repository/db"
	"github.com/CogitoNTNU/cogi-go/internal/util/env"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()

	ctx := context.Background()

	cfg := config.LoadApiConfig()

	logger := logrus.New().WithField("app", cfg.AppName).WithContext(ctx)

	logger.Info("Starting server \n")

	env := env.Configure(".env")

	_, err := db.InitDb(env, logger, ctx)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to db")
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run()
}
