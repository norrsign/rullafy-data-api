package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/you/rullafy-data-api/internal/config"
	"github.com/you/rullafy-data-api/internal/db"
	"github.com/you/rullafy-data-api/internal/handlers"
	"github.com/you/rullafy-data-api/internal/repo"
)

func main() {
	_ = godotenv.Load() // optional

	cfg := config.Load()

	// Init DB once
	if err := db.InitDB(cfg.DatabaseURL); err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	// Create all repos
	qr := db.Qrs
	userRepo := repo.NewUserRepo(qr)
	companyRepo := repo.NewCompanyRepo(qr)
	productRepo := repo.NewProductRepo(qr)

	// Gin setup
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Register our routes
	handlers.RegisterUser(router, userRepo)
	handlers.RegisterCompany(router, companyRepo)
	handlers.RegisterProduct(router, productRepo)

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	go func() {
		log.Printf("⇢  listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen failed: %v", err)
		}
	}()

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("⇢  shutting down…")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown failed: %v", err)
	}
	log.Println("⇢  server stopped")
}
