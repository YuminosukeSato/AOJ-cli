package repository

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/entity"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/repository"
	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
	"github.com/YuminosukeSato/AOJ-cli/pkg/logger"
)

// LocalSessionRepository implements SessionRepository for local file storage
type LocalSessionRepository struct {
	configDir string
	logger    *logger.Logger
}

// NewLocalSessionRepository creates a new LocalSessionRepository
func NewLocalSessionRepository(configDir string) repository.SessionRepository {
	return &LocalSessionRepository{
		configDir: configDir,
		logger:    logger.WithGroup("local_session_repository"),
	}
}

// SessionData represents the JSON structure for session storage
type SessionData struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	CreatedAt int64  `json:"created_at"`
	LastUsed  int64  `json:"last_used"`
}

// Save saves a session to local storage
func (r *LocalSessionRepository) Save(ctx context.Context, session *entity.Session) error {
	r.logger.DebugContext(ctx, "saving session", 
		"session_id", session.ID().MaskedString())

	if err := r.ensureConfigDir(); err != nil {
		return cerrors.Wrap(err, "failed to ensure config directory")
	}

	// Convert session to storage format
	data := SessionData{
		ID:        session.ID().String(),
		Username:  session.Username(),
		Token:     session.Token(),
		ExpiresAt: session.ExpiresAt().Unix(),
		CreatedAt: session.CreatedAt().Unix(),
		LastUsed:  session.LastUsed().Unix(),
	}

	// Write to file
	sessionFile := r.getSessionFilePath(session.ID())
	file, err := os.Create(sessionFile)
	if err != nil {
		return cerrors.Wrap(err, "failed to create session file")
	}
	defer func() {
		if err := file.Close(); err != nil {
			r.logger.WarnContext(ctx, "failed to close file", "error", err)
		}
	}()

	if err := json.NewEncoder(file).Encode(data); err != nil {
		return cerrors.Wrap(err, "failed to encode session data")
	}

	// Set file permissions to be readable only by owner
	if err := os.Chmod(sessionFile, 0600); err != nil {
		r.logger.WarnContext(ctx, "failed to set session file permissions", "error", err)
	}

	r.logger.DebugContext(ctx, "session saved successfully", 
		"session_id", session.ID().MaskedString(),
		"file", sessionFile)

	return nil
}

// GetByID retrieves a session by its ID
func (r *LocalSessionRepository) GetByID(ctx context.Context, id model.SessionID) (*entity.Session, error) {
	r.logger.DebugContext(ctx, "getting session by ID", 
		"session_id", id.MaskedString())

	sessionFile := r.getSessionFilePath(id)
	
	// Check if file exists
	if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
		return nil, cerrors.NewAppError(
			cerrors.CodeNotFound,
			"session not found",
			nil,
		)
	}

	// Read and parse session file
	file, err := os.Open(sessionFile)
	if err != nil {
		return nil, cerrors.Wrap(err, "failed to open session file")
	}
	defer func() {
		if err := file.Close(); err != nil {
			r.logger.WarnContext(ctx, "failed to close file", "error", err)
		}
	}()

	var data SessionData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, cerrors.Wrap(err, "failed to decode session data")
	}

	// Convert to entity
	session, err := r.dataToSession(data)
	if err != nil {
		return nil, cerrors.Wrap(err, "failed to convert session data")
	}

	r.logger.DebugContext(ctx, "session retrieved successfully", 
		"session_id", id.MaskedString())

	return session, nil
}

// GetByUsername retrieves the current session for a username
func (r *LocalSessionRepository) GetByUsername(ctx context.Context, username string) (*entity.Session, error) {
	r.logger.DebugContext(ctx, "getting session by username", "username", username)

	// List all session files and find matching username
	sessions, err := r.List(ctx)
	if err != nil {
		return nil, cerrors.Wrap(err, "failed to list sessions")
	}

	for _, session := range sessions {
		if session.Username() == username && session.IsValid() {
			return session, nil
		}
	}

	return nil, cerrors.NewAppError(
		cerrors.CodeNotFound,
		"no valid session found for username",
		nil,
	)
}

// GetCurrent retrieves the current active session
func (r *LocalSessionRepository) GetCurrent(ctx context.Context) (*entity.Session, error) {
	r.logger.DebugContext(ctx, "getting current session")

	currentFile := r.getCurrentSessionFilePath()
	
	// Check if current session file exists
	if _, err := os.Stat(currentFile); os.IsNotExist(err) {
		return nil, cerrors.NewAppError(
			cerrors.CodeNotFound,
			"no current session",
			nil,
		)
	}

	// Read current session ID
	content, err := os.ReadFile(currentFile)
	if err != nil {
		return nil, cerrors.Wrap(err, "failed to read current session file")
	}

	sessionID, err := model.NewSessionID(string(content))
	if err != nil {
		return nil, cerrors.Wrap(err, "invalid session ID in current session file")
	}

	// Get the actual session
	return r.GetByID(ctx, sessionID)
}

// Delete deletes a session by its ID
func (r *LocalSessionRepository) Delete(ctx context.Context, id model.SessionID) error {
	r.logger.DebugContext(ctx, "deleting session", 
		"session_id", id.MaskedString())

	sessionFile := r.getSessionFilePath(id)
	
	if err := os.Remove(sessionFile); err != nil && !os.IsNotExist(err) {
		return cerrors.Wrap(err, "failed to delete session file")
	}

	r.logger.DebugContext(ctx, "session deleted successfully", 
		"session_id", id.MaskedString())

	return nil
}

// DeleteByUsername deletes all sessions for a username
func (r *LocalSessionRepository) DeleteByUsername(ctx context.Context, username string) error {
	r.logger.DebugContext(ctx, "deleting sessions by username", "username", username)

	sessions, err := r.List(ctx)
	if err != nil {
		return cerrors.Wrap(err, "failed to list sessions")
	}

	deleted := 0
	for _, session := range sessions {
		if session.Username() == username {
			if err := r.Delete(ctx, session.ID()); err != nil {
				r.logger.WarnContext(ctx, "failed to delete session", 
					"session_id", session.ID().MaskedString(), 
					"error", err)
			} else {
				deleted++
			}
		}
	}

	r.logger.DebugContext(ctx, "sessions deleted by username", 
		"username", username, 
		"deleted_count", deleted)

	return nil
}

// DeleteExpired deletes all expired sessions
func (r *LocalSessionRepository) DeleteExpired(ctx context.Context) error {
	r.logger.DebugContext(ctx, "deleting expired sessions")

	sessions, err := r.List(ctx)
	if err != nil {
		return cerrors.Wrap(err, "failed to list sessions")
	}

	deleted := 0
	for _, session := range sessions {
		if session.IsExpired() {
			if err := r.Delete(ctx, session.ID()); err != nil {
				r.logger.WarnContext(ctx, "failed to delete expired session", 
					"session_id", session.ID().MaskedString(), 
					"error", err)
			} else {
				deleted++
			}
		}
	}

	r.logger.DebugContext(ctx, "expired sessions deleted", 
		"deleted_count", deleted)

	return nil
}

// Exists checks if a session exists
func (r *LocalSessionRepository) Exists(_ context.Context, id model.SessionID) (bool, error) {
	sessionFile := r.getSessionFilePath(id)
	_, err := os.Stat(sessionFile)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, cerrors.Wrap(err, "failed to check session file")
	}
	return true, nil
}

// IsValid checks if a session is valid (exists and not expired)
func (r *LocalSessionRepository) IsValid(ctx context.Context, id model.SessionID) (bool, error) {
	session, err := r.GetByID(ctx, id)
	if err != nil {
		if cerrors.IsAppError(err, cerrors.CodeNotFound) {
			return false, nil
		}
		return false, err
	}
	return session.IsValid(), nil
}

// SetCurrent sets the current active session
func (r *LocalSessionRepository) SetCurrent(ctx context.Context, session *entity.Session) error {
	r.logger.DebugContext(ctx, "setting current session", 
		"session_id", session.ID().MaskedString())

	if err := r.ensureConfigDir(); err != nil {
		return cerrors.Wrap(err, "failed to ensure config directory")
	}

	currentFile := r.getCurrentSessionFilePath()
	
	if err := os.WriteFile(currentFile, []byte(session.ID().String()), 0600); err != nil {
		return cerrors.Wrap(err, "failed to write current session file")
	}

	r.logger.DebugContext(ctx, "current session set successfully", 
		"session_id", session.ID().MaskedString())

	return nil
}

// ClearCurrent clears the current active session
func (r *LocalSessionRepository) ClearCurrent(ctx context.Context) error {
	r.logger.DebugContext(ctx, "clearing current session")

	currentFile := r.getCurrentSessionFilePath()
	
	if err := os.Remove(currentFile); err != nil && !os.IsNotExist(err) {
		return cerrors.Wrap(err, "failed to remove current session file")
	}

	r.logger.DebugContext(ctx, "current session cleared successfully")

	return nil
}

// List lists all sessions
func (r *LocalSessionRepository) List(ctx context.Context) ([]*entity.Session, error) {
	r.logger.DebugContext(ctx, "listing all sessions")

	sessionsDir := r.getSessionsDir()
	
	// Check if sessions directory exists
	if _, err := os.Stat(sessionsDir); os.IsNotExist(err) {
		return []*entity.Session{}, nil
	}

	// Read directory
	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		return nil, cerrors.Wrap(err, "failed to read sessions directory")
	}

	var sessions []*entity.Session
	for _, entry := range entries {
		if entry.IsDir() || !entry.Type().IsRegular() {
			continue
		}

		// Try to parse as session ID
		sessionID, err := model.NewSessionID(entry.Name())
		if err != nil {
			r.logger.WarnContext(ctx, "invalid session file name", 
				"filename", entry.Name())
			continue
		}

		session, err := r.GetByID(ctx, sessionID)
		if err != nil {
			r.logger.WarnContext(ctx, "failed to load session", 
				"session_id", sessionID.MaskedString(), 
				"error", err)
			continue
		}

		sessions = append(sessions, session)
	}

	r.logger.DebugContext(ctx, "sessions listed successfully", 
		"count", len(sessions))

	return sessions, nil
}

// Helper methods

func (r *LocalSessionRepository) ensureConfigDir() error {
	return os.MkdirAll(r.getSessionsDir(), 0755)
}

func (r *LocalSessionRepository) getSessionsDir() string {
	return filepath.Join(r.configDir, "sessions")
}

func (r *LocalSessionRepository) getSessionFilePath(id model.SessionID) string {
	return filepath.Join(r.getSessionsDir(), id.String())
}

func (r *LocalSessionRepository) getCurrentSessionFilePath() string {
	return filepath.Join(r.configDir, "current_session")
}

func (r *LocalSessionRepository) dataToSession(data SessionData) (*entity.Session, error) {
	sessionID, err := model.NewSessionID(data.ID)
	if err != nil {
		return nil, err
	}
	// Use reflection or factory method to create session with all fields
	// This is a simplified version - in practice, you might need a more sophisticated approach
	
	session := entity.NewSession(
		sessionID,
		data.Username,
		data.Token,
		time.Unix(data.ExpiresAt, 0),
	)

	// Update timestamps
	session.UpdateLastUsedAt(time.Unix(data.LastUsed, 0))

	return session, nil
}