package repository

import (
	"context"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/entity"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
)

// SessionRepository defines the interface for session data access
type SessionRepository interface {
	// Save saves a session
	Save(ctx context.Context, session *entity.Session) error

	// GetByID retrieves a session by its ID
	GetByID(ctx context.Context, id model.SessionID) (*entity.Session, error)

	// GetByUsername retrieves the current session for a username
	GetByUsername(ctx context.Context, username string) (*entity.Session, error)

	// GetCurrent retrieves the current active session
	GetCurrent(ctx context.Context) (*entity.Session, error)

	// Delete deletes a session by its ID
	Delete(ctx context.Context, id model.SessionID) error

	// DeleteByUsername deletes all sessions for a username
	DeleteByUsername(ctx context.Context, username string) error

	// DeleteExpired deletes all expired sessions
	DeleteExpired(ctx context.Context) error

	// Exists checks if a session exists
	Exists(ctx context.Context, id model.SessionID) (bool, error)

	// IsValid checks if a session is valid (exists and not expired)
	IsValid(ctx context.Context, id model.SessionID) (bool, error)

	// SetCurrent sets the current active session
	SetCurrent(ctx context.Context, session *entity.Session) error

	// ClearCurrent clears the current active session
	ClearCurrent(ctx context.Context) error

	// List lists all sessions (for admin purposes)
	List(ctx context.Context) ([]*entity.Session, error)
}