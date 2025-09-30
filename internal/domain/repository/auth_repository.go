// Package repository defines interfaces for data access.
package repository

import (
	"context"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/entity"
)

// AuthRepository defines the interface for authentication operations
type AuthRepository interface {
	// Login authenticates a user and returns a session
	Login(ctx context.Context, username, password string) (*entity.Session, error)

	// Logout logs out a user by invalidating their session
	Logout(ctx context.Context, session *entity.Session) error

	// RefreshSession refreshes an existing session
	RefreshSession(ctx context.Context, session *entity.Session) (*entity.Session, error)

	// ValidateSession validates if a session is still active on the server
	ValidateSession(ctx context.Context, session *entity.Session) (bool, error)
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string
	Password string
}

// NewLoginRequest creates a new login request
func NewLoginRequest(username, password string) LoginRequest {
	return LoginRequest{
		Username: username,
		Password: password,
	}
}

// IsValid validates the login request
func (r LoginRequest) IsValid() bool {
	return r.Username != "" && r.Password != ""
}