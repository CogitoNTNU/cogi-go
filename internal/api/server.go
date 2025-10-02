package api

import (
	"context"
	"database/sql"
	"time"

	routes "github.com/CogitoNTNU/cogi-go/internal/api/router"
	"github.com/CogitoNTNU/cogi-go/internal/config"
	"github.com/CogitoNTNU/cogi-go/internal/handler"
	"github.com/CogitoNTNU/cogi-go/internal/repository/db"
	userRepository "github.com/CogitoNTNU/cogi-go/internal/repository/user"
	"github.com/CogitoNTNU/cogi-go/internal/service"
	"github.com/CogitoNTNU/cogi-go/internal/util/env"
	"github.com/gin-gonic/gin"
	"github.com/mbndr/figlet4go"
	"github.com/sirupsen/logrus"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	healthcheckConfig "github.com/tavsec/gin-healthcheck/config"
)

type Server struct {
	Logger *logrus.Entry
	engine *gin.Engine
	env    *env.EnvConfig
	Ctx    *context.Context
	db     *sql.DB
}

func InitServer() (*Server, error) {
	engine := gin.Default()
	ctx := context.Background()
	queryCtx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer cancel()

	cfg := config.LoadApiConfig()
	e := env.Configure(".env")
	cors := cfg.CorsNew(e)

	if e.Read("ENVIRONMENT") == env.PROD {
		engine.Use(cors)
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	logger := logrus.New().WithField("app", cfg.AppName).WithContext(ctx)

	database, err := db.InitDb(e, logger, queryCtx)

	return &Server{
		Logger: logger,
		engine: engine,
		env:    e,
		Ctx:    &ctx,
		db:     database,
	}, err
}

func (s *Server) Serve() {
	s.Logger.WithContext(*s.Ctx).Info("Initializing server... \n")

	ascii := figlet4go.NewAsciiRender()
	render, _ := ascii.Render("COGI-GO!!")
	s.Logger.Info("\n" + render)

	// Start Database
	s.Logger.WithTime(time.Now()).Info("Starting database...")
	sqlxDb := db.ExportDb(s.db) // TODO: Inject sqlxDb into repositories which are then used by services -> handlers
	defer func() {
		if err := sqlxDb.Close(); err != nil {
			s.Logger.WithError(err).Error("Failed to close sqlx database connection")
			return
		}
	}()

	// Health checks
	s.Logger.WithTime(time.Now()).Info("Performing health checks...")
	sqlCheck := checks.SqlCheck{Sql: s.db}
	healthcheck.New(s.engine, healthcheckConfig.DefaultConfig(), []checks.Check{sqlCheck})

	userRepository := userRepository.NewRepo(sqlxDb, time.Duration(5)*time.Second, s.Logger)
	userService := service.NewUserService(userRepository)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, s.Ctx)

	// Register routes and handlers
	s.Logger.WithTime(time.Now()).Info("Performing health checks...")
	routes.RegisterPublicRoutes(s.engine, userHandler)
	routes.RegisterPrivateRoutes(s.engine)
	routes.RegisterAdminRoutes(s.engine)

	// Start server
	s.Logger.WithContext(*s.Ctx).Info("Initialization Complete! \n")
	s.engine.Run()
}
