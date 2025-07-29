package db

import (
	"context"
	"testing"
	"time"
)

func TestMigrations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create PostgreSQL container
	container, err := NewPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}
	defer container.Close(ctx)

	// Connect to the database
	pool, err := container.Connect(ctx)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Run migrations
	err = RunMigrationsForContainer(ctx, container)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Verify the users table was created
	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query users table: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 users, got %d", count)
	}

	// Test inserting a user
	_, err = pool.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", "testuser")
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// Verify the user was inserted
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE name = $1", "testuser").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query for inserted user: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 user, got %d", count)
	}

	// Test unique constraint
	_, err = pool.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", "testuser")
	if err == nil {
		t.Error("Expected error when inserting duplicate name, got none")
	}
}
