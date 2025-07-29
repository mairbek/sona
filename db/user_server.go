package db

import (
	"context"
	"fmt"

	dbgen "sona/db/gen"
	proto "sona/gen"
	"sona/gen/sonav1connect"

	"connectrpc.com/connect"
)

// UserServer provides user operations using sqlc-generated code
type UserServer struct {
	sonav1connect.UnimplementedUserServiceHandler

	Queries *dbgen.Queries
}

var _ sonav1connect.UserServiceHandler = (*UserServer)(nil)

// NewUserServer creates a new user service
func NewUserServer(queries *dbgen.Queries) *UserServer {
	return &UserServer{
		Queries: queries,
	}
}

// CreateUser creates a new user
// func (UnimplementedUserServiceHandler) CreateUser(context.Context, *connect.Request[gen.CreateUserRequest]) (*connect.Response[gen.User], error) {
func (s *UserServer) CreateUser(ctx context.Context, req *connect.Request[proto.CreateUserRequest]) (*connect.Response[proto.User], error) {
	user, err := s.Queries.CreateUser(ctx, req.Msg.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return connect.NewResponse(&proto.User{
		Id:   user.ID,
		Name: user.Name,
	}), nil
}

// GetUser retrieves a user by ID
func (s *UserServer) GetUser(ctx context.Context, req *connect.Request[proto.GetUserRequest]) (*connect.Response[proto.User], error) {
	id := req.Msg.Id
	user, err := s.Queries.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return connect.NewResponse(&proto.User{
		Id:   user.ID,
		Name: user.Name,
	}), nil
}

// ListUsers retrieves all users
func (s *UserServer) ListUsers(ctx context.Context, req *connect.Request[proto.ListUsersRequest]) (*connect.Response[proto.ListUsersResponse], error) {
	users, err := s.Queries.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert to slice of pointers
	result := make([]*proto.User, len(users))
	for i, user := range users {
		result[i] = &proto.User{
			Id:   user.ID,
			Name: user.Name,
		}
	}
	return connect.NewResponse(&proto.ListUsersResponse{
		Users: result,
	}), nil
}
