package db

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norrsign/rullafy-data-api/db/models"
	"github.com/sirupsen/logrus"
)

var (
	pool    *pgxpool.Pool
	Qrs     *models.Queries
	initErr error
	once    sync.Once
)

// InitDB sets up the global connection pool exactly once.
// Call this early in application startup.
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

		Qrs = models.New(pool)
	})
	return initErr
}

// CloseDB closes the database connection pool.
func CloseDB() {
	if pool != nil {
		pool.Close()
		logrus.Info("database pool closed")
	}
}
