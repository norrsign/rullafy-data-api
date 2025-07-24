// internal/cmd/server.go
package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vanern/goapi/config"
	"github.com/vanern/goapi/internal/server"
	"github.com/vanern/goapi/internal/utils"
	"github.com/vanern/goapi/types"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Server management commands",
	}
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the HTTP service",
		Run:   startServer,
		PreRunE: func(cmd *cobra.Command, args []string) error {

			afterConfigHook := hooks[types.ServerAfterConfigHook]
			if afterConfigHook != nil {
				return afterConfigHook()
			}
			return nil
		},
	}
)

func init() {
	RootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(startCmd)

	// Bind flags to Viper keys
	startCmd.Flags().BoolP("verbose", "v", false,
		"enable verbose logging (env: GOAPI_START_VERBOSE)")
	viper.BindPFlag("start_verbose", startCmd.Flags().Lookup("verbose"))
	viper.SetDefault("start_verbose", false)

	startCmd.Flags().String("jwt-public-key", "",
		"path to RSA public key for JWT validation (env: GOAPI_JWT_PUBLIC_KEY)")
	viper.BindPFlag("jwt_public_key", startCmd.Flags().Lookup("jwt-public-key"))

	startCmd.Flags().String("jwt-realm-url", "",
		"Keycloak realm base URL (env: GOAPI_JWT_REALM_URL)")
	viper.BindPFlag("jwt_realm_url", startCmd.Flags().Lookup("jwt-realm-url"))

	startCmd.Flags().Duration("jwt-key-refresh-interval", time.Minute,
		"JWKS refresh interval (env: GOAPI_JWT_KEY_REFRESH_INTERVAL)")
	viper.BindPFlag("jwt_key_refresh_interval", startCmd.Flags().Lookup("jwt-key-refresh-interval"))
	viper.SetDefault("jwt_key_refresh_interval", time.Minute)
}

func startServer(cmd *cobra.Command, args []string) {
	// Honor verbose flag
	if config.Config.Server.Start.Verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	// Initialize JWKS if using Keycloak
	if config.Config.Server.JWTRealmURL != "" {
		if err := utils.InitJWKs(
			config.Config.Server.JWTRealmURL,
			config.Config.Server.JWTKeyRefreshInterval,
		); err != nil {
			logrus.Fatalf("failed to initialize JWKS: %v", err)
		}
		logrus.Infof("JWKS initialized from realm %s (refresh every %s)",
			config.Config.Server.JWTRealmURL,
			config.Config.Server.JWTKeyRefreshInterval,
		)
	} else if config.Config.Server.JWTPublicKey == "" {
		logrus.Fatal("either --jwt-public-key or --jwt-realm-url must be provided")
	}

	runWithGracefulShutdown()
}

// runWithGracefulShutdown spins up one http.Server per registered Gin,
// then waits for SIGINT/SIGTERM to shut them down gracefully.
func runWithGracefulShutdown() {
	var servers []*http.Server

	// Build and start each server concurrently
	for _, srv := range server.Gins {
		httpSrv := &http.Server{
			Addr:              srv.Address,
			Handler:           srv.Gin,
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
		}
		servers = append(servers, httpSrv)

		go func(s *http.Server) {
			logrus.Infof("Listening on %s", s.Addr)
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("listen %s failed: %v", s.Addr, err)
			}
		}(httpSrv)
	}

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logrus.Infof("signal %s received, shutting down...", sig)

	// Graceful shutdown with 15s timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for _, s := range servers {
		if err := s.Shutdown(ctx); err != nil {
			logrus.Errorf("shutdown of %s failed: %v", s.Addr, err)
		}
	}
	logrus.Info("all servers stopped")
}
