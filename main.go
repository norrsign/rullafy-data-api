package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vanern/goapi/cli"
	"github.com/vanern/goapi/framework/middleware"
	"github.com/vanern/goapi/framework/middleware/auth"

	"github.com/vanern/goapi/framework"
)

func main() {
	os.Setenv("TZ", "UTC")
	now := time.Now()
	fmt.Printf("Starting server at %s\n", now.Format(time.RFC3339))
	// Global logging middleware
	middleware.Use(auth.KeycloakJwt("rullafy-client"))

	//middleware.Use(Test1())
	//middleware.Use(Test2())

	// PUBLIC
	framework.GET("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong health")
	})

	// AUTH but no roles
	framework.GET("/me", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong me")
	}, "seller")

	// AUTH + roles
	framework.GET("/admin", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong admin")
	}, "admin", "ops")

	framework.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong pong")
	})

	cli.Run()
}
