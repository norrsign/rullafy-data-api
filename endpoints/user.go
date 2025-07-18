package endpoints

import (
	"net/http"
	"strings"

	"github.com/norrsign/rullafy-data-api/repositories"
	"github.com/vanern/goapi/framework"
	"github.com/vanern/goapi/framework/httpx"
	"github.com/vanern/goapi/typs"
)

// RegisterUser wires CRUD paths for /users/**
func RegisterUser(ur *repositories.UserRepo) {

	parseID := func(r *http.Request) (string, error) {
		// URL is .../users/{id}.  Grab last segment.
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		return parts[len(parts)-1], nil
	}

	framework.GET("/users", httpx.MakeList(ur), typs.AuthOnly)
	framework.POST("/users", httpx.MakeCreate(ur))
	framework.GET("/users/{id}", httpx.MakeGet(ur, parseID))
	//framework.PUT("/users/{id}", httpx.MakeUpdate(ur, parseID))
	//framework.DELETE("/users/{id}", httpx.MakeDelete(ur, parseID))
} //
