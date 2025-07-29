package db_test

import (
	"context"
	"testing"
	"time"

	dbgen "sona/db/gen"
	"sona/db"
	"sona/dbstub"
	proto "sona/gen"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get test pool from global container
	pool, err := dbstub.GetTestPool(ctx)
	assert.NoError(t, err, "Failed to get test pool")
	defer pool.Close()

	// Create queries and service
	queries := dbgen.New(pool)
	service := db.NewUserServer(queries)

	// Test creating a user
	resp, err := service.CreateUser(ctx, connect.NewRequest(&proto.CreateUserRequest{
		Name: "testuser",
	}))
	assert.NoError(t, err, "Failed to create user")

	user := resp.Msg
	assert.NotZero(t, user.Id, "Expected user ID to be non-zero")
	assert.Equal(t, "testuser", user.Name, "Expected user name to be 'testuser'")
}

func TestGetUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get test pool from global container
	pool, err := dbstub.GetTestPool(ctx)
	assert.NoError(t, err, "Failed to get test pool")
	defer pool.Close()

	// Create queries and service
	queries := dbgen.New(pool)
	service := db.NewUserServer(queries)

	// First create a user
	createResp, err := service.CreateUser(ctx, connect.NewRequest(&proto.CreateUserRequest{
		Name: "getusertest",
	}))
	assert.NoError(t, err, "Failed to create user for get test")
	createdUser := createResp.Msg

	// Test getting user by ID
	getUserResp, err := service.GetUser(ctx, connect.NewRequest(&proto.GetUserRequest{
		Id: createdUser.Id,
	}))
	assert.NoError(t, err, "Failed to get user")

	retrievedUser := getUserResp.Msg
	assert.Equal(t, createdUser.Id, retrievedUser.Id, "User ID should match")
	assert.Equal(t, createdUser.Name, retrievedUser.Name, "User name should match")
}

func TestListUsers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get test pool from global container
	pool, err := dbstub.GetTestPool(ctx)
	assert.NoError(t, err, "Failed to get test pool")
	defer pool.Close()

	// Create queries and service
	queries := dbgen.New(pool)
	service := db.NewUserServer(queries)

	// Create a user first
	_, err = service.CreateUser(ctx, connect.NewRequest(&proto.CreateUserRequest{
		Name: "listusertest",
	}))
	assert.NoError(t, err, "Failed to create user for list test")

	// Test listing users
	listUsersResp, err := service.ListUsers(ctx, connect.NewRequest(&proto.ListUsersRequest{}))
	assert.NoError(t, err, "Failed to list users")

	users := listUsersResp.Msg.Users
	assert.GreaterOrEqual(t, len(users), 1, "Should have at least 1 user")
} 