package main

import (
	"fmt"
	"os"
	"time"

	"github.com/norrsign/rullafy-data-api/db"
	"github.com/norrsign/rullafy-data-api/endpoints"
	"github.com/vanern/goapi/cli"
	"github.com/vanern/goapi/framework/middleware"
	"github.com/vanern/goapi/framework/middleware/auth"
)

func main() {
	os.Setenv("TZ", "UTC")

	db.InitDB("postgresql://myuser:mypassword@localhost:5432/mydatabase")

	now := time.Now()
	fmt.Printf("Starting server at %s\n", now.Format(time.RFC3339))
	// Global logging middleware
	middleware.Use(auth.KeycloakJwt("rullafy-client"))
	endpoints.Init()
	//middleware.Use(Test1())
	//middleware.Use(Test2())

	cli.Run()
}
