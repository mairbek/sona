package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs all pending migrations on the database using connection string
func RunMigrations(connString, migrationsPath string) error {
	// Create migrate instance using connection string
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		connString)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// RunMigrationsFromProjectRoot runs migrations using the default migrations path
func RunMigrationsFromProjectRoot(connString string) error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Construct the migrations path - use absolute path
	migrationsPath := filepath.Join(cwd, "migrations")

	// Verify the migrations directory exists
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", migrationsPath)
	}

	return RunMigrations(connString, migrationsPath)
}

// RunMigrationsForContainer runs migrations for a PostgreSQL container
func RunMigrationsForContainer(ctx context.Context, container *PostgresContainer) error {
	connString := container.GetConnectionString()
	return RunMigrationsFromProjectRoot(connString)
}
