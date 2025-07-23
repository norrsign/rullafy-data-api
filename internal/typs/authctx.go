package typs

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// key used both in gin.Context and standard context
const UserContextKey = "auth-user"

type AuthUser struct {
	ID     string
	Roles  []string
	Claims jwt.MapClaims
}

/* ---------------- helpers ---------------- */

// From std-lib context (e.g. inside repo)
func UserFromContext(ctx context.Context) (*AuthUser, bool) {
	v := ctx.Value(UserContextKey)
	u, ok := v.(*AuthUser)
	return u, ok
}

// From gin.Context (inside handlers / middleware)
func GetAuthUser(c *gin.Context) (*AuthUser, bool) {
	v, ok := c.Get(UserContextKey)
	if !ok {
		return nil, false
	}
	u, ok := v.(*AuthUser)
	return u, ok
}
