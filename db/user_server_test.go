package db

import (
	"context"
	"testing"
	"time"

	dbgen "sona/db/gen"
	proto "sona/gen"

	"connectrpc.com/connect"
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

	queries := dbgen.New(pool)

	// Create service
	service := NewUserServer(queries)

	// Test creating a user
	resp, err := service.CreateUser(ctx, connect.NewRequest(&proto.CreateUserRequest{
		Name: "testuser",
	}))
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	user := resp.Msg

	if user.Id == 0 {
		t.Error("Expected user ID to be non-zero")
	}
	if user.Name != "testuser" {
		t.Errorf("Expected user name to be 'testuser', got '%s'", user.Name)
	}

	// Test getting user by ID
	getUserResp, err := service.GetUser(ctx, connect.NewRequest(&proto.GetUserRequest{
		Id: user.Id,
	}))
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	retrievedUser := getUserResp.Msg
	if retrievedUser.Id != user.Id {
		t.Errorf("Expected user ID %d, got %d", user.Id, retrievedUser.Id)
	}
	if retrievedUser.Name != user.Name {
		t.Errorf("Expected user name '%s', got '%s'", user.Name, retrievedUser.Name)
	}

	// Test listing users
	listUsersResp, err := service.ListUsers(ctx, connect.NewRequest(&proto.ListUsersRequest{}))
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}
	users := listUsersResp.Msg.Users
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

}
