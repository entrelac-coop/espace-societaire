package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Middleware(key []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			authorization := c.GetHeader("Authorization")
			if authorization == "" {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "The Authorization header or the token query is required for this route.",
					"code":  "authorization-header-missing",
				})
				return
			}

			if !strings.HasPrefix(authorization, "Bearer ") {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "The Authorization header must start with 'Bearer '.",
					"code":  "authorization-header-missing",
				})
				return
			}

			token = authorization[7:]
		}

		parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})
		if parsedToken.Valid {
			if claims, ok := parsedToken.Claims.(*CustomClaims); ok {
				c.Set("userID", claims.UserID)
				c.Set("admin", claims.Admin)
				sentry.ConfigureScope(func(scope *sentry.Scope) {
					scope.SetUser(sentry.User{ID: claims.UserID})
				})
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "The token is invalid.",
					"code":  "token-invalid",
				})
			}
		} else if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "The token is invalid.",
				"code":  "token-invalid",
			})
			return
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "The token has expired.",
				"code":  "token-expired",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error.",
				"code":  "internal",
			})
			return
		}
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		admin := c.GetBool("admin")

		if !admin {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "You must be admin to call this route.",
				"code":  "not-admin",
			})
			return
		}

		c.Next()
	}
}
