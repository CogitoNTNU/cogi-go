package api

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"strconv"
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
	mail "github.com/xhit/go-simple-mail/v2"
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

	// TODO: FIX HACKY CORS CONFIG
	engine.RedirectTrailingSlash = false
	engine.RedirectFixedPath = false
	corscfg := cors.DefaultConfig()
	corscfg.AllowOrigins = []string{"http://localhost:3000", "https://cogito-ntnu.no", "*"}
	corscfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corscfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corscfg.AllowCredentials = true
	corscfg.MaxAge = 12 * time.Hour
	engine.Use(cors.New(corscfg))
	// cors := cfg.CorsNew(e)
	// envVal, err := e.Read("ENVIRONMENT")
	// if err != nil {
	// 	logrus.WithError(err).Fatal("Failed to read ENVIRONMENT from env")
	// }
	// if envVal == env.PROD {
	// }
	// engine.Use(cors)

	ctx := context.Background()
	cfg := config.LoadApiConfig()
	e := env.Configure(".env")

	logger := logrus.New().WithField("app", cfg.AppName).WithContext(ctx)

	// cors := cfg.CorsNew()
	// envVal, err := e.Read("ENVIRONMENT")
	// if err != nil {
	// 	logger.Fatalf("Failed to read ENVIRONMENT from env")
	// }
	// if envVal == env.PROD {
	// 	engine.Use(cors)
	// }

	// WARNING: THIS IS NOT IDEAL
	engine.Use(CORSMiddleware())

	queryCtx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer cancel()

	authMiddleware, err := jwt.New(initParams())
	if err != nil {
		logger.Fatalf("JWT Error: %s", err.Error())
	}

	err = authMiddleware.MiddlewareInit()
	if err != nil {
		logger.Fatalf("authMiddleware.MiddlewareInit() Error: %s", err.Error())
	}

	database, err := db.InitDb(e, logger, queryCtx)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %s", err.Error())
	}

	smtpServer, err := initSMTPServer(e)
	if err != nil {
		logger.Fatalf("Failed to initialize SMTP client: %s", err.Error())
	}
	engine.Use(SMTPMiddleware(smtpServer))

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
	sqlxDB := db.ExportDb(s.db)
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

	// TODO: Use Gin Context for env
	tempApplicationHandler := handler.NewTempApplicationHandler(templateApplicationService, ctx, env)

	return &routes.Handlers{
		User:            userHandler,
		Project:         projectHandler,
		Sponsor:         sponsorHandler,
		TempApplication: tempApplicationHandler,
	}
}

func SMTPMiddleware(server *mail.SMTPServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("smtp_server", server)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func initSMTPServer(e *env.EnvConfig) (*mail.SMTPServer, error) {
	server := mail.NewSMTPClient()

	serverHost, err := e.Read("SMTP_HOST")
	if err != nil {
		return nil, err
	}
	server.Host = serverHost

	serverPort, err := e.Read("SMTP_PORT")
	if err != nil {
		return nil, err
	}
	server.Port, err = strconv.Atoi(serverPort)
	if err != nil {
		return nil, err
	}

	serverUser, err := e.Read("SMTP_USER")
	if err != nil {
		return nil, err
	}
	server.Username = serverUser

	serverPassword, err := e.Read("SMTP_PASSWORD")
	if err != nil {
		return nil, err
	}
	server.Password = serverPassword

	serverEncryption, err := e.Read("SMTP_ENCRYPTION")
	if err != nil {
		return nil, err
	}
	serverEncryptionInt, err := strconv.Atoi(serverEncryption)
	if err != nil {
		return nil, err
	}
	server.Encryption = mail.Encryption(serverEncryptionInt)
	server.KeepAlive = true
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return server, nil
}

func renderAscii(input string) string {
	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()

	options.FontName = "speed"
	ascii.LoadFont("./misc/fonts/speed.flf")

	render, _ := ascii.RenderOpts(input, options)

	return render
}
