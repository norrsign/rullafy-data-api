package db

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var (
	pool    *pgxpool.Pool // internal connection pool
	Qrs     *Queries      // shared sqlc handle
	initErr error
	once    sync.Once
)

// Init sets up the global connection pool exactly once.
// Call this early from main() or a services.Init() helper.
func InitDB(connString string) error {
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		pool, initErr = pgxpool.New(ctx, connString)
		if initErr != nil {
			return
		}
		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			initErr = err
			return
		}

		Qrs = New(pool) // sqlc handle
		registerShutdownHook()
	})
	return initErr
}

// Close manually closes the pool. Normally called by the shutdown hook.
func close() {
	if pool != nil {
		pool.Close()
		logrus.Info("database pool closed")
	}
}

// internal graceful-shutdown hook
func registerShutdownHook() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-c
		logrus.Infof("signal %s received - shutting down db pool...", sig)
		close()
	}()
}
