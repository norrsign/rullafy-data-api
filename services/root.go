package services

import (
	"context"
	"fmt"
	"os"

	"github.com/norrsign/rullafy-data-api/db"
)

// Init wires up everything services need. Call from main().
func Init() error {
	dsn := os.Getenv("DATABASE_URL") // e.g. postgres://user:pass@localhost:5432/app
	if dsn == "" {
		dsn = "postgresql://myuser:mypassword@localhost:5432/mydatabase"
	}
	return db.InitDB(dsn)
}

// Example business logic -------------------------------------------------

// SeedAuthors shows how service code interacts **without** context/pool.
func SeedAuthors() error {
	_, err := db.Qrs.CreateUser(context.Background(),
		db.User{
			Job: "Fucker",
		})

	if err != nil {
		return fmt.Errorf("seeding user: %w", err)
	}
	return nil
}
