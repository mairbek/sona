package db

import (
	"context"
	"testing"
	"time"
)

func TestServiceWithMigrations(t *testing.T) {
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

	// Create service
	service := NewService(pool)
	defer service.Close()

	// Test creating a user
	user, err := service.CreateUser(ctx, "testuser")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be non-zero")
	}
	if user.Name != "testuser" {
		t.Errorf("Expected user name to be 'testuser', got '%s'", user.Name)
	}

	// Test getting user by ID
	retrievedUser, err := service.GetUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, retrievedUser.ID)
	}
	if retrievedUser.Name != user.Name {
		t.Errorf("Expected user name '%s', got '%s'", user.Name, retrievedUser.Name)
	}

	// Test getting user by name
	retrievedByName, err := service.GetUserByName(ctx, "testuser")
	if err != nil {
		t.Fatalf("Failed to get user by name: %v", err)
	}

	if retrievedByName.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, retrievedByName.ID)
	}

	// Test listing users
	users, err := service.ListUsers(ctx)
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

	// Test updating user
	updatedUser, err := service.UpdateUser(ctx, user.ID, "updateduser")
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	if updatedUser.Name != "updateduser" {
		t.Errorf("Expected updated name 'updateduser', got '%s'", updatedUser.Name)
	}

	// Test unique constraint
	_, err = service.CreateUser(ctx, "updateduser")
	if err == nil {
		t.Error("Expected error when creating user with duplicate name")
	}

	// Test deleting user
	err = service.DeleteUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify user was deleted
	_, err = service.GetUser(ctx, user.ID)
	if err == nil {
		t.Error("Expected error when getting deleted user")
	}
}
