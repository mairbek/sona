package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"sona/db/gen"
)

// Service provides database operations using sqlc-generated code
type Service struct {
	*gen.Queries
	pool *pgxpool.Pool
}

// NewService creates a new database service
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{
		Queries: gen.New(pool),
		pool:    pool,
	}
}

// Close closes the database connection pool
func (s *Service) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, name string) (*gen.User, error) {
	user, err := s.Queries.CreateUser(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &user, nil
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(ctx context.Context, id int32) (*gen.User, error) {
	user, err := s.Queries.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByName retrieves a user by name
func (s *Service) GetUserByName(ctx context.Context, name string) (*gen.User, error) {
	user, err := s.Queries.GetUserByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by name: %w", err)
	}
	return &user, nil
}

// ListUsers retrieves all users
func (s *Service) ListUsers(ctx context.Context) ([]*gen.User, error) {
	users, err := s.Queries.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert to slice of pointers
	result := make([]*gen.User, len(users))
	for i, user := range users {
		result[i] = &user
	}
	return result, nil
}

// UpdateUser updates a user's name
func (s *Service) UpdateUser(ctx context.Context, id int32, name string) (*gen.User, error) {
	user, err := s.Queries.UpdateUser(ctx, gen.UpdateUserParams{
		ID:   id,
		Name: name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return &user, nil
}

// DeleteUser deletes a user by ID
func (s *Service) DeleteUser(ctx context.Context, id int32) error {
	err := s.Queries.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
