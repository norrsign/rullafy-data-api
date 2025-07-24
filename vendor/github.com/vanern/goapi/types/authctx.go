package types

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// key under which we store *AuthUser in Gin context
const UserContextKey = "auth-user"

// AuthUser holds the user identity & roles extracted from the JWT
type AuthUser struct {
	ID     string
	Roles  []string
	Claims jwt.MapClaims
}

// GetAuthUser retrieves the *AuthUser from `c *gin.Context`.
func GetAuthUser(c *gin.Context) (*AuthUser, bool) {
	ui, ok := c.Get(UserContextKey)
	user, ok2 := ui.(*AuthUser)
	return user, ok && ok2
}
