package config

import (
	"os"
)

// Config holds all environment‐based settings.
type Config struct {
	DatabaseURL   string // e.g. postgres://user:pass@host:port/db
	Port          string // e.g. "8080"
	JWTRealmURL   string // if using Keycloak JWKS
	JWTPublicKey  string // path to static PEM file
	JWTClientName string // Keycloak client to extract roles
}

// Load reads from the environment with sensible defaults.
func Load() Config {
	c := Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		Port:          os.Getenv("PORT"),
		JWTRealmURL:   os.Getenv("JWT_REALM_URL"),
		JWTPublicKey:  os.Getenv("JWT_PUBLIC_KEY"),
		JWTClientName: os.Getenv("JWT_CLIENT"),
	}
	if c.DatabaseURL == "" {
		c.DatabaseURL = "postgresql://myuser:mypassword@localhost:5432/mydatabase"
	}
	if c.Port == "" {
		c.Port = "8080"
	}
	return c
}
