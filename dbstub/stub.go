package dbstub

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"sona/db"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/golangmigrator"
)

var (
	container      *db.PostgresContainer
	containerOnce  sync.Once
	containerErr   error
	migrationsPath string
)

// GetTestContainer returns the global test container, initializing it if needed
func GetTestContainer() (*db.PostgresContainer, error) {
	containerOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		container, containerErr = db.NewPostgresContainer(ctx)
		if containerErr != nil {
			log.Printf("Failed to create test container: %v", containerErr)
			return
		}

		log.Println("Test container initialized successfully")
	})

	return container, containerErr
}

// TestDBPool returns a connection pool to an isolated test database
// This function should be called from within a test function
func TestDBPool(ctx context.Context, t testing.TB) (*pgxpool.Pool, error) {
	container, err := GetTestContainer()
	if err != nil {
		return nil, err
	}

	// Create pgtestdb config using the container's connection details
	config := pgtestdb.Config{
		DriverName: "pgx",
		User:       container.User,
		Password:   container.Password,
		Host:       container.Host,
		Port:       container.Port,
		Database:   container.Database,
		Options:    "sslmode=disable",
	}

	migrator := golangmigrator.New(migrationsPath)
	testConfig := pgtestdb.Custom(t, config, migrator)

	// Connect to the test database using the returned config
	pool, err := pgxpool.New(ctx, testConfig.URL())
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func init() {
	// Discover migrations path once
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Warning: Failed to get current working directory: %v", err)
		return
	}

	// Find the project root by looking for go.mod
	projectRoot := cwd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			log.Printf("Warning: Could not find project root (go.mod not found)")
			return
		}
		projectRoot = parent
	}

	migrationsPath = filepath.Join(projectRoot, "db", "migrations")

}
