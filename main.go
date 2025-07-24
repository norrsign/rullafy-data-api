package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norrsign/rullafy-data-api/db"
	"github.com/sirupsen/logrus"

	"github.com/vanern/goapi/cli"
	"github.com/vanern/goapi/framework"
	"github.com/vanern/goapi/framework/middleware/auth"
	"github.com/vanern/goapi/types"
)

func serverAfterConfigHook() error {
	os.Setenv("TZ", "UTC")

	// initialize DB (requires $DATABASE_URL)
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		logrus.Fatal("DATABASE_URL environment variable is required")
	}
	if err := db.InitDB(connStr); err != nil {
		logrus.Fatalf("failed to init DB: %v", err)
	}

	// start time log
	now := time.Now()
	fmt.Printf("Starting server at %s\n", now.Format(time.RFC3339))

	// register routes
	gn := framework.RegisterGin(":8080")

	// health check
	gn.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong\n")
	})

	// protected endpoint
	gr := gn.Group("/protected")
	gr.Use(auth.JWTAuth("rullafy-client"))
	gr.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "protected\n")
	})

	return nil
}

func tokenAfterConfigHook() error {
	return nil
}
func main() {
	// start service (blocks until shutdown)

	cli.AddRunHook(types.TokenAfterConfigHook, tokenAfterConfigHook)
	cli.AddRunHook(types.ServerAfterConfigHook, serverAfterConfigHook)
	cli.Run()
	// once HTTP servers are down, close DB pool
	db.CloseDB()
}
