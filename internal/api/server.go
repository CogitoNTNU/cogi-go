package api

import (
	"context"
	"database/sql"
	"fmt"
	routes "github.com/CogitoNTNU/cogi-go/internal/api/router"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/CogitoNTNU/cogi-go/internal/config"
	"github.com/CogitoNTNU/cogi-go/internal/handler"
	"github.com/CogitoNTNU/cogi-go/internal/repository/db"
	projectRepository "github.com/CogitoNTNU/cogi-go/internal/repository/project"
	userRepository "github.com/CogitoNTNU/cogi-go/internal/repository/user"
	"github.com/CogitoNTNU/cogi-go/internal/service"
	"github.com/CogitoNTNU/cogi-go/internal/util/env"
	"github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/mbndr/figlet4go"
	"github.com/sirupsen/logrus"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	healthcheckConfig "github.com/tavsec/gin-healthcheck/config"
	"time"
)

type Server struct {
	cfg           *config.Api
	Logger        *logrus.Entry
	engine        *gin.Engine
	env           *env.EnvConfig
	Ctx           *context.Context
	db            *sql.DB
	jwtMiddleware *jwt.GinJWTMiddleware
}

func InitServer() (*Server, error) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	ctx := context.Background()
	queryCtx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer cancel()

	cfg := config.LoadApiConfig()
	e := env.Configure(".env")
	cors := cfg.CorsNew(e)

	if e.Read("ENVIRONMENT") == env.PROD {
		engine.Use(cors)
	}

	logger := logrus.New().WithField("app", cfg.AppName).WithContext(ctx)

	authMiddleware, err := jwt.New(initParams())
	if err != nil {
		logger.Infof("JWT Error: %s", err.Error())
	}

	err = authMiddleware.MiddlewareInit()
	if err != nil {
		logger.Infof("authMiddleware.MiddlewareInit() Error: %s", err.Error())
	}

	database, err := db.InitDb(e, logger, queryCtx)

	return &Server{
		cfg:           cfg,
		Logger:        logger,
		engine:        engine,
		env:           e,
		Ctx:           &ctx,
		db:            database,
		jwtMiddleware: authMiddleware,
	}, err
}

func (s *Server) Serve() {
	s.Logger.WithContext(*s.Ctx).Info("Initializing server... \n")

	render := renderAscii(s.cfg.AppName)
	s.Logger.Info("\n \n" + render + "\n \n")

	s.Logger.WithTime(time.Now()).Info("Starting database...")
	sqlxDb := db.ExportDb(s.db) // TODO: Inject sqlxDb into repositories which are then used by services -> handlers
	defer func() {
		if err := sqlxDb.Close(); err != nil {
			s.Logger.WithError(err).Error("Failed to close sqlx database connection")
			return
		}
	}()

	s.Logger.WithTime(time.Now()).Info("Performing health checks...")
	sqlCheck := checks.SqlCheck{Sql: s.db}
	healthcheck.New(s.engine, healthcheckConfig.DefaultConfig(), []checks.Check{sqlCheck})

	s.Logger.WithTime(time.Now()).Info("Registering routes...")

	// User endpoints
	userRepository := userRepository.NewRepo(sqlxDb, time.Duration(5)*time.Second, s.Logger)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService, s.Ctx)

	// Project endpoints
	projectRepository := projectRepository.NewRepo(sqlxDb, time.Duration(5)*time.Second, s.Logger)
	projectService := service.NewProjectService(projectRepository)
	projectHandler := handler.NewProjectHandler(projectService, s.Ctx)

	handlers := &routes.Handlers{
		User:    userHandler,
		Project: projectHandler,
	}

	routes.RegisterPublicRoutes(s.engine, handlers)
	routes.RegisterPrivateRoutes(s.engine, s.jwtMiddleware, handlers)
	routes.RegisterAdminRoutes(s.engine)

	s.Logger.WithContext(*s.Ctx).Info("Initialization Complete! \n")
	s.engine.Run(fmt.Sprintf(":%d", s.cfg.MainPort))
}

func initParams() *jwt.GinJWTMiddleware {

	return &jwt.GinJWTMiddleware{
		Realm:      "example zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		PayloadFunc: func(data any) gojwt.MapClaims {
			return gojwt.MapClaims{
				"user_id": data,
			}
		},
	}
}

func renderAscii(input string) string {
	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()

	options.FontName = "speed"
	ascii.LoadFont("./misc/fonts/speed.flf")

	render, _ := ascii.RenderOpts(input, options)

	return render
}
