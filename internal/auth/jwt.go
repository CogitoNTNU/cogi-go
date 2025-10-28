package auth
import(
	"context"
	"errors"
	"net/http"
	"time"
jwt"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
	"github.com/CogitoNTNU/cogi-go/internal/util/env"
userRepository "github.com/CogitoNTNU/cogi-go/internal/repository/user"
)
const identityKey = "userId"
type loginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type googleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}
type Identity struct {
	UserID string `json:"userId"`
	Email stirng `json:"email"`
}
func NewAuthMiddleware(e *env.EnvConfig, logger *logrus.Entry, userRepo *userRepository.Repo, appCtx *context.Context) (*jwt.GinJWTMiddleware, error) {
	secret := e.Read("JWT_SECRET")
	if secret == "" {
		return nilm errors.New("JWT_SECRET not set")
	}
timeoutMin := e.ReadIntDefault("JWT_TIMEOUT_MIN", 15)
	maxRedreshHours := e.ReadIntDefault("JWT_MAX_REFRESH_HOURS", 24) 
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm: "cogi-go",
		Key: []byte(secret),
		Timeout: time.Duration(timeoutMin) * time.Minute,
		MaxRefresh: time.Duration(maxRedreshHours) * time.Hour,
		identityKey: identityKey,
		SendCookie: true,
		CookieHTTPOnly: true,
		CookieMaxAge: timeoutMin * 60,
		CookieName: "access_token",
		TokenLookup: "cookie_access_token, header:Authorization, query:token",
		TokenHeadName: "Bearer",
		TimeFunc: time.Now,
		Authenticator: func(c *gin.Context) (any, error) {
			var req loginRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
user, uerr := userRepo.getUserByEmail(appCtx, req.Email)
			if uerr != nil {
				return nil, errors.New("invalid credentials")
			}
if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
	return nil, errors.New("invalid credentials")

}
return &Identity{UserID: user.ID, Email: user.Email}, nil
		},
		PayloadFunc: func(data any) jwt.MapClaims {
			if id, ok := data.(*Identity); ok {
				return jwt.MapClaims{
					identityKey: id.UserID,
					"email": id.email
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler : func(c *gin.Context) any {
			claims := jwt.ExtractClaims(c)
			uid,_:=claims[identityKey].(string)
			email,_:=claims["email"].(string)
			return &Identity{UserID: uid, Email: email}
		},
		Authorizator: func(data any, c *gin.Context) bool {
			-, ok := data.(*Identity)
			return ok
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"error": message})
		},
		
	})
if err != nil {
	return nil, err
}
return authMiddleware, nil
}
func GoogleLoginHandler(e *env.EnvConfig, logger *logrus.Entry, auth *jwt.GinJWTMiddleware) gin.HandlerFunc {
	audience := e.Read("GOOGLE_CLIENT_ID")
	return func(c *gin.Contect) {
		var req googleLoginRequest
		if err := c.ShouldBindJSON(&req); err !=nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
payload, err := idtoken.Validate(c, req.IDToken, audience)
		if err != nil {
			logger.WithError(err).Warn("google id token validation failed")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid google token"})
			return
		}
userID, _ := pauload.Claims["sub"].(string)
		email, _ := payload.Claims["email"].(string)
		token, expire, err := auth.TokenGeneratoer(&Identity{UserID: userID, Email: email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
			return
		}
		auth.SetCookie(c, token, expire)
		c.JSON(http.StatusOK, gin.H{"token": token, "expire": expire})
	}
}

