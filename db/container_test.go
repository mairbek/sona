package db

import (
	"context"
	"testing"
	"time"
)

func TestPostgresContainer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create and start PostgreSQL container
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

	// Test: Create a table
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS test_users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Test: Insert data
	_, err = pool.Exec(ctx, `
		INSERT INTO test_users (name, email) VALUES ($1, $2)
	`, "Test User", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Test: Query data
	rows, err := pool.Query(ctx, "SELECT id, name, email FROM test_users WHERE email = $1", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to query data: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}

		// Verify the data
		if name != "Test User" || email != "test@example.com" {
			t.Errorf("Expected name='Test User', email='test@example.com', got name='%s', email='%s'", name, email)
		}
	} else {
		t.Error("Expected to find one row, but found none")
	}

	t.Log("PostgreSQL container test completed successfully!")
}

func TestPostgresContainerConnectionString(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	container, err := NewPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}
	defer container.Close(ctx)

	connString := container.GetConnectionString()
	expectedPrefix := "postgres://testuser:testpass@"

	if len(connString) < len(expectedPrefix) || connString[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Connection string should start with '%s', got: %s", expectedPrefix, connString)
	}

	t.Logf("Connection string: %s", connString)
}
