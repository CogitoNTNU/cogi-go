package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	gojwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"

	"github.com/CogitoNTNU/cogi-go/internal/util/env"
	userRepository "github.com/CogitoNTNU/cogi-go/internal/repository/user"
)

const identityKey = "userId"

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type googleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

type Identity struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

func NewAuthMiddleware(
	e *env.EnvConfig,
	logger *logrus.Entry,
	userRepo *userRepository.Repo,
	appCtx *context.Context,
) (*jwt.GinJWTMiddleware, error) {
	secret := e.Read("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	timeoutMin := e.ReadIntDefault("JWT_TIMEOUT_MIN", 15)
	maxRefreshHours := e.ReadIntDefault("JWT_MAX_REFRESH_HOURS", 24)

	mw, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "cogi-go",
		Key:         []byte(secret),
		Timeout:     time.Duration(timeoutMin) * time.Minute,
		MaxRefresh:  time.Duration(maxRefreshHours) * time.Hour,
		TimeFunc:    time.Now,
		IdentityKey: identityKey,

		Authenticator: func(c *gin.Context) (any, error) {
			var req loginRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			ctx := context.Background()
			if appCtx != nil {
				ctx = *appCtx
			}

			user, berr := userRepo.GetUserByEmail(&ctx, req.Email)
			if berr != nil {
				return nil, errors.New("invalid credentials")
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
				return nil, errors.New("invalid credentials")
			}

			return &Identity{
				UserID: user.Id.String(),
				Email:  user.Email,
			}, nil
		},

		PayloadFunc: func(data any) gojwt.MapClaims {
			if id, ok := data.(*Identity); ok {
				return gojwt.MapClaims{
					identityKey: id.UserID,
					"email":     id.Email,
				}
			}
			return gojwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) any {
			claims := jwt.ExtractClaims(c)
			uid, _ := claims[identityKey].(string)
			email, _ := claims["email"].(string)
			return &Identity{UserID: uid, Email: email}
		},

		Authorizer: func(c *gin.Context, data any) bool {
			_, ok := data.(*Identity)
			return ok
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"error": message})
		},
	})
	if err != nil {
		return nil, err
	}

	if err := mw.MiddlewareInit(); err != nil {
		return nil, err
	}

	return mw, nil
}

func GoogleLoginHandler(
	e *env.EnvConfig,
	logger *logrus.Entry,
	auth *jwt.GinJWTMiddleware,
) gin.HandlerFunc {
	audience := e.Read("GOOGLE_CLIENT_ID")

	return func(c *gin.Context) {
		var req googleLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		payload, err := idtoken.Validate(c, req.IDToken, audience)
		if err != nil {
			logger.WithError(err).Warn("google id token validation failed")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid google token"})
			return
		}

		userID, _ := payload.Claims["sub"].(string)
		email, _ := payload.Claims["email"].(string)

		tokenPair, err := auth.TokenGenerator(&Identity{
			UserID: userID,
			Email:  email,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
			return
		}

		setAuthCookie(c, tokenPair.AccessToken)

		c.JSON(http.StatusOK, gin.H{
			"access_token":  tokenPair.AccessToken,
			"refresh_token": tokenPair.RefreshToken,
			"expires_at":    tokenPair.ExpiresAt,
			"token_type":    tokenPair.TokenType,
		})
	}
}

func setAuthCookie(c *gin.Context, token string) {
	
	maxAge := 86400
	secure := false
	httpOnly := true

	c.SetCookie("jwt", token, maxAge, "/", "", secure, httpOnly)
}