package dbstub

import (
	"context"
	"log"
	"sync"
	"time"

	"sona/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	container     *db.PostgresContainer
	containerOnce sync.Once
	containerErr  error
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

		// Run migrations
		pool, err := container.Connect(ctx)
		if err != nil {
			containerErr = err
			return
		}
		defer pool.Close()

		// Run migrations using the existing db package function
		err = db.RunMigrationsForContainer(ctx, container)
		if err != nil {
			containerErr = err
			return
		}

		log.Println("Test container initialized successfully")
	})

	return container, containerErr
}

// GetTestPool returns a connection pool to the test database
func GetTestPool(ctx context.Context) (*pgxpool.Pool, error) {
	container, err := GetTestContainer()
	if err != nil {
		return nil, err
	}

	return container.Connect(ctx)
}

func init() {
	// Initialize the test container when the package is imported
	_, err := GetTestContainer()
	if err != nil {
		log.Printf("Warning: Failed to initialize test container: %v", err)
	}
}
