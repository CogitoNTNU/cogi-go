package api

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	routes "github.com/CogitoNTNU/cogi-go/internal/api/router"
	"github.com/CogitoNTNU/cogi-go/internal/config"
	"github.com/CogitoNTNU/cogi-go/internal/handler"
	"github.com/CogitoNTNU/cogi-go/internal/repository/db"
	projectRepository "github.com/CogitoNTNU/cogi-go/internal/repository/project"
	sponsorRepository "github.com/CogitoNTNU/cogi-go/internal/repository/sponsor"
	tempApplicationRepository "github.com/CogitoNTNU/cogi-go/internal/repository/tempApplication"
	userRepository "github.com/CogitoNTNU/cogi-go/internal/repository/user"
	"github.com/CogitoNTNU/cogi-go/internal/service"
	"github.com/CogitoNTNU/cogi-go/internal/util/env"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/mbndr/figlet4go"
	"github.com/sirupsen/logrus"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	healthcheckConfig "github.com/tavsec/gin-healthcheck/config"
)

type Server struct {
	cfg           *config.Api
	Logger        *logrus.Entry
	engine        *gin.Engine
	Env           *env.EnvConfig
	Ctx           *context.Context
	db            *sql.DB
	jwtMiddleware *jwt.GinJWTMiddleware
}

func InitServer() (*Server, error) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	// engine.RedirectTrailingSlash = false
	// engine.RedirectFixedPath = false

	cfg := config.LoadApiConfig()
	e := env.Configure(".env")

	cors := cfg.CorsNew()
	envVal, err := e.Read("ENVIRONMENT")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to read ENVIRONMENT from env")
	}
	if envVal == env.PROD {
		engine.Use(cors)
	}

	ctx := context.Background()
	queryCtx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer cancel()

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
		Env:           e,
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
	sqlxDB := db.ExportDb(s.db) // TODO: Inject sqlxDb into repositories which are then used by services -> handlers
	defer func() {
		if err := sqlxDB.Close(); err != nil {
			s.Logger.WithError(err).Error("Failed to close sqlx database connection")
			return
		}
	}()

	s.Logger.WithTime(time.Now()).Info("Performing health checks...")
	sqlCheck := checks.SqlCheck{Sql: s.db}
	err := healthcheck.New(s.engine, healthcheckConfig.DefaultConfig(), []checks.Check{sqlCheck})
	if err != nil {
		s.Logger.WithError(err).Error("Error serving on main port")
		return
	}

	s.Logger.WithTime(time.Now()).Info("Registering routes...")
	handlers := routerHandlers(sqlxDB, s.Logger, s.Ctx, s.Env)

	routes.RegisterPublicRoutes(s.engine, handlers)
	routes.RegisterPrivateRoutes(s.engine, s.jwtMiddleware, handlers)
	routes.RegisterAdminRoutes(s.engine)

	err = s.engine.Run(fmt.Sprintf(":%d", s.cfg.MainPort))
	s.Logger.WithContext(*s.Ctx).Info("Initialization Complete! \n")
	if err != nil {
		s.Logger.WithError(err).Error("Error serving on main port")
		return
	}
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

func routerHandlers(sqlxDB *sqlx.DB, logger *logrus.Entry, ctx *context.Context, env *env.EnvConfig) *routes.Handlers {
	queryTimeoutLimit := time.Duration(5) * time.Second
	userRepository := userRepository.NewRepo(sqlxDB, queryTimeoutLimit, logger)
	userService := service.NewUserService(userRepository, logger)
	userHandler := handler.NewUserHandler(userService, ctx)

	projectRepository := projectRepository.NewRepo(sqlxDB, queryTimeoutLimit, logger)
	projectService := service.NewProjectService(projectRepository, logger)
	projectHandler := handler.NewProjectHandler(projectService, ctx)

	sponsorRepository := sponsorRepository.NewRepo(sqlxDB, queryTimeoutLimit, logger)
	sponsorService := service.NewSponsorService(sponsorRepository, logger)
	sponsorHandler := handler.NewSponsorHandler(sponsorService, ctx)

	tempApplicationRepository := tempApplicationRepository.NewRepo(sqlxDB, queryTimeoutLimit, logger)
	templateApplicationService := service.NewTempApplicationService(tempApplicationRepository, logger)

	// TODO: Inject env into gin
	tempApplicationHandler := handler.NewTempApplicationHandler(templateApplicationService, ctx, env)

	return &routes.Handlers{
		User:            userHandler,
		Project:         projectHandler,
		Sponsor:         sponsorHandler,
		TempApplication: tempApplicationHandler,
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
