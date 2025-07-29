package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer represents a PostgreSQL test container
type PostgresContainer struct {
	Container testcontainers.Container
	Host      string
	Port      string
	User      string
	Password  string
	Database  string
}

// NewPostgresContainer creates and starts a new PostgreSQL test container
func NewPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithTmpfs(map[string]string{
			"/var/lib/postgresql/data": "rw,nosuid,nodev,noexec,relatime,size=1g",
		}),
		testcontainers.WithCmdArgs(
			"-c", "fsync=off",
			"-c", "synchronous_commit=off",
			"-c", "full_page_writes=off",
		),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	return &PostgresContainer{
		Container: postgresContainer,
		Host:      host,
		Port:      port.Port(),
		User:      "testuser",
		Password:  "testpass",
		Database:  "testdb",
	}, nil
}

// GetConnectionString returns the connection string for the PostgreSQL container
func (pc *PostgresContainer) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pc.User, pc.Password, pc.Host, pc.Port, pc.Database)
}

// Connect establishes a connection to the PostgreSQL container using pgx
func (pc *PostgresContainer) Connect(ctx context.Context) (*pgxpool.Pool, error) {
	connString := pc.GetConnectionString()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection with retry
	var pingErr error
	for i := 0; i < 3; i++ {
		pingErr = pool.Ping(ctx)
		if pingErr == nil {
			break
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if pingErr != nil {
		return nil, fmt.Errorf("failed to ping database after retries: %w", pingErr)
	}

	log.Printf("Successfully connected to PostgreSQL container at %s:%s", pc.Host, pc.Port)
	return pool, nil
}

// Close stops and removes the PostgreSQL container
func (pc *PostgresContainer) Close(ctx context.Context) error {
	if pc.Container != nil {
		return pc.Container.Terminate(ctx)
	}
	return nil
}
