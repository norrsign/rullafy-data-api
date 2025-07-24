package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"

	"github.com/vanern/goapi/config"
	"github.com/vanern/goapi/internal/utils"
	"github.com/vanern/goapi/types"
)

// ---------------------------------------------------------------------------
// helpers to extract roles from JWT claims
// ---------------------------------------------------------------------------

func extractSimpleRoles(claims jwt.MapClaims) ([]string, error) {
	raw, ok := claims["roles"].([]interface{})
	if !ok {
		return nil, errors.New("roles claim missing or malformatted")
	}
	var out []string
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out, nil
}

func extractKeycloakRoles(claims jwt.MapClaims, client string) ([]string, error) {
	ra, ok := claims["resource_access"].(map[string]interface{})
	if !ok {
		return nil, errors.New("resource_access claim missing")
	}
	entry, ok := ra[client].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no resource_access entry for client %q", client)
	}
	raw, ok := entry["roles"].([]interface{})
	if !ok {
		return nil, errors.New("roles missing under resource_access")
	}
	var out []string
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out, nil
}

func JWTAuth(client string, requiredRoles ...string) gin.HandlerFunc {

	// Precompute Keyfunc once
	cfg := config.Config.Server
	var keyfunc jwt.Keyfunc
	var isKC = false
	if cfg.JWTRealmURL != "" {
		keyfunc = utils.Keyfunc()
		isKC = true
	} else {
		pub, err := utils.LoadPublicKey(cfg.JWTPublicKey)
		if err != nil {
			logrus.Fatalf("failed to load public key: %v", err)
		}
		keyfunc = func(_ *jwt.Token) (interface{}, error) { return pub, nil }
	}

	return func(c *gin.Context) {

		// 1) Extract the token
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer"))

		// 2) Parse & validate
		claims := jwt.MapClaims{}
		tok, err := jwt.ParseWithClaims(tokenString, claims, keyfunc, jwt.WithLeeway(30*time.Second))
		if err != nil || !tok.Valid {
			logrus.Errorf("JWT parse error: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// 3) Extract subject & roles
		sub, _ := claims["sub"].(string)
		var roles []string
		if isKC {
			roles, err = extractKeycloakRoles(claims, client)
		} else {
			roles, err = extractSimpleRoles(claims)
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		// 4) Enforce requiredRoles (if any)
		//    If the only required role is types.AuthOnly, we skip checks.
		if len(requiredRoles) > 0 {
			allowed := false
			// build a quick lookup of the user's roles
			have := make(map[string]struct{}, len(roles))
			for _, r := range roles {
				have[r] = struct{}{}
			}
			// see if any of the requiredRoles is present
			for _, want := range requiredRoles {
				if _, ok := have[want]; ok {
					allowed = true
					break
				}
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized: missing required roles"})
				return
			}
		}

		// 5) Store AuthUser in context for handlers
		user := &types.AuthUser{ID: sub, Roles: roles, Claims: claims}
		c.Set(types.UserContextKey, user)

		c.Next()
	}
}
