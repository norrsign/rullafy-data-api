package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/norrsign/rullafy-data-api/internal/config"
	"github.com/norrsign/rullafy-data-api/internal/typs"
	"github.com/norrsign/rullafy-data-api/internal/utils"
	"github.com/sirupsen/logrus"
)

func init() {
	// allow loading .env for local dev
	_ = godotenv.Load()
}

// JWTAuth must be applied on any route that needs a signed-in user.
// It injects *typs.AuthUser into Gin’s context.
func JWTAuth() gin.HandlerFunc {
	cfg := config.Load()

	// choose keyfunc: JWKS vs. static PEM
	var keyfunc jwt.Keyfunc
	if cfg.JWTRealmURL != "" {
		keyfunc = utils.Keyfunc() // from your existing utils
	} else {
		pub, err := utils.LoadPublicKey(cfg.JWTPublicKey)
		if err != nil {
			logrus.Fatalf("loading public key: %v", err)
		}
		keyfunc = func(*jwt.Token) (interface{}, error) { return pub, nil }
	}

	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer"))

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, keyfunc, jwt.WithLeeway(30*time.Second))
		if err != nil || !token.Valid {
			logrus.Errorf("JWT parse error: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// extract subject
		sub, _ := claims["sub"].(string)

		// extract roles (simple array-based or Keycloak style)
		var roles []string
		if ra, ok := claims["roles"].([]interface{}); ok {
			for _, v := range ra {
				if s, ok := v.(string); ok {
					roles = append(roles, s)
				}
			}
		} else if ra, ok := claims["resource_access"].(map[string]interface{}); ok {
			if rc, ok := ra[cfg.JWTClientName].(map[string]interface{}); ok {
				if rr, ok := rc["roles"].([]interface{}); ok {
					for _, v := range rr {
						if s, ok := v.(string); ok {
							roles = append(roles, s)
						}
					}
				}
			}
		}

		user := &typs.AuthUser{
			ID:     sub,
			Roles:  roles,
			Claims: claims,
		}

		// store into Gin context
		c.Set(typs.UserContextKey, user)
		c.Next()
	}
}
