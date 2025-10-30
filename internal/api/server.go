package api

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	routes "github.com/CogitoNTNU/cogi-go/internal/api/router"
	authpkg "github.com/CogitoNTNU/cogi-go/internal/auth"
	"github.com/CogitoNTNU/cogi-go/internal/config"
	"github.com/CogitoNTNU/cogi-go/internal/handler"
	"github.com/CogitoNTNU/cogi-go/internal/repository/db"
	projectRepository "github.com/CogitoNTNU/cogi-go/internal/repository/project"
	userRepository "github.com/CogitoNTNU/cogi-go/internal/repository/user"
	"github.com/CogitoNTNU/cogi-go/internal/service"
	"github.com/CogitoNTNU/cogi-go/internal/util/env"
	"github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v5"
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
	env           *env.EnvConfig
	Ctx           *context.Context
	db            *sql.DB
	jwtMiddleware *jwt.GinJWTMiddleware
}

func InitServer() (*Server, error) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	ctx := context.Background()
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &Server{
		cfg:           cfg,
		Logger:        logger,
		engine:        engine,
		env:           e,
		Ctx:           &ctx,
		db:            database,
		jwtMiddleware: authMiddleware,
	}, nil
}

func (s *Server) Serve() {
	s.Logger.WithContext(*s.Ctx).Info("Initializing server...")

	render := renderAscii(s.cfg.AppName)
	s.Logger.Info("\n\n" + render + "\n\n")

	s.Logger.WithTime(time.Now()).Info("Starting database...")
	sqlxDb := db.ExportDb(s.db)

	// Initialize repositories
	userRepo := userRepository.NewRepo(sqlxDb, 5*time.Second, s.Logger)
	projectRepo := projectRepository.NewRepo(sqlxDb, 5*time.Second, s.Logger)

	// Initialize services
	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(projectRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, s.Ctx)
	projectHandler := handler.NewProjectHandler(projectService, s.Ctx)

	// Initialize authentication
	authMW, err := authpkg.NewAuthMiddleware(s.env, s.Logger, userRepo, s.Ctx)
	if err != nil {
		s.Logger.WithError(err).Fatal("failed to init auth middleware")
		return
	}

	googleLogin := authpkg.GoogleLoginHandler(s.env, s.Logger, authMW)

	// Bundle handlers
	handlers := &routes.Handlers{
		User:        userHandler,
		Project:     projectHandler,
		AuthMW:      authMW,
		GoogleLogin: googleLogin,
	}
	
	// Health checks
	s.Logger.WithTime(time.Now()).Info("Registering health checks...")
	sqlCheck := checks.SqlCheck{Sql: s.db}
	healthcheck.New(s.engine, healthcheckConfig.DefaultConfig(), []checks.Check{sqlCheck})

	// Register routes
	s.Logger.WithTime(time.Now()).Info("Registering routes...")
	routes.RegisterPublicRoutes(s.engine, handlers)
	routes.RegisterPrivateRoutes(s.engine, authMW, handlers)
	routes.RegisterAdminRoutes(s.engine)

	// Start the server
	s.Logger.Infof("Server running on port %d", s.cfg.MainPort)
	defer func() {
		if err := sqlxDb.Close(); err != nil {
			s.Logger.WithError(err).Error("Failed to close sqlx database connection")
		}
	}()

	if err := s.engine.Run(fmt.Sprintf(":%d", s.cfg.MainPort)); err != nil {
		s.Logger.WithError(err).Fatal("Failed to start server")
	}
}

func initParams() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:      "example zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: 24 * time.Hour,
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
