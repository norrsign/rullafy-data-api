package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/norrsign/rullafy-data-api/db"
	"github.com/norrsign/rullafy-data-api/internal/config"
	"github.com/norrsign/rullafy-data-api/internal/handlers"
	"github.com/norrsign/rullafy-data-api/internal/repo"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	// ---- DB ---------------------------------------------------------------
	if err := db.InitDB(cfg.DatabaseURL); err != nil {
		logrus.Fatalf("DB init failed: %v", err)
	}
	q := db.Qrs

	// ---- dependencies -----------------------------------------------------
	userRepo := repo.NewUserRepo(q)

	// ---- Gin engine -------------------------------------------------------
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery()) // proven defaults :contentReference[oaicite:1]{index=1}

	handlers.RegisterUser(router, userRepo)

	// ---- HTTP server with graceful shutdown ------------------------------
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		logrus.Infof("⇢ listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("server shutdown: %v", err)
	}
	logrus.Info("⇢ server stopped")
}
