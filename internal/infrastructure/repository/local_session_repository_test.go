package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/entity"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
)

func TestLocalSessionRepository_SaveAndGetByID(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	repo := NewLocalSessionRepository(tmpDir)
	ctx := context.Background()

	sessionID := model.MustGenerateSessionID()
	session := entity.NewSessionWithDuration(
		sessionID,
		"testuser",
		"test_token_123",
		24*time.Hour,
	)

	// When - Save
	err := repo.Save(ctx, session)

	// Then
	assert.NoError(t, err)

	// Verify file exists with correct permissions
	sessionFile := filepath.Join(tmpDir, "sessions", sessionID.String())
	info, err := os.Stat(sessionFile)
	assert.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode())

	// When - GetByID
	retrievedSession, err := repo.GetByID(ctx, sessionID)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, retrievedSession)
	assert.Equal(t, session.ID(), retrievedSession.ID())
	assert.Equal(t, session.Username(), retrievedSession.Username())
	assert.Equal(t, session.Token(), retrievedSession.Token())
}

func TestLocalSessionRepository_GetByID_NotFound(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	repo := NewLocalSessionRepository(tmpDir)
	ctx := context.Background()

	nonExistentID := model.MustGenerateSessionID()

	// When
	session, err := repo.GetByID(ctx, nonExistentID)

	// Then
	assert.Error(t, err)
	assert.Nil(t, session)
	assert.True(t, cerrors.IsAppError(err, cerrors.CodeNotFound))
}

func TestLocalSessionRepository_SetCurrentAndGetCurrent(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	repo := NewLocalSessionRepository(tmpDir)
	ctx := context.Background()

	sessionID := model.MustGenerateSessionID()
	session := entity.NewSessionWithDuration(
		sessionID,
		"testuser",
		"test_token_123",
		24*time.Hour,
	)

	// Save the session first
	err := repo.Save(ctx, session)
	assert.NoError(t, err)

	// When - SetCurrent
	err = repo.SetCurrent(ctx, session)
	assert.NoError(t, err)

	// When - GetCurrent
	currentSession, err := repo.GetCurrent(ctx)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, currentSession)
	assert.Equal(t, session.ID(), currentSession.ID())
	assert.Equal(t, session.Username(), currentSession.Username())
}

func TestLocalSessionRepository_GetCurrent_NotFound(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	repo := NewLocalSessionRepository(tmpDir)
	ctx := context.Background()

	// When
	session, err := repo.GetCurrent(ctx)

	// Then
	assert.Error(t, err)
	assert.Nil(t, session)
	assert.True(t, cerrors.IsAppError(err, cerrors.CodeNotFound))
}

func TestLocalSessionRepository_Delete(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	repo := NewLocalSessionRepository(tmpDir)
	ctx := context.Background()

	sessionID := model.MustGenerateSessionID()
	session := entity.NewSessionWithDuration(
		sessionID,
		"testuser",
		"test_token_123",
		24*time.Hour,
	)

	// Save session first
	err := repo.Save(ctx, session)
	assert.NoError(t, err)

	// Verify it exists
	exists, err := repo.Exists(ctx, sessionID)
	assert.NoError(t, err)
	assert.True(t, exists)

	// When - Delete
	err = repo.Delete(ctx, sessionID)

	// Then
	assert.NoError(t, err)

	// Verify it no longer exists
	exists, err = repo.Exists(ctx, sessionID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestLocalSessionRepository_ClearCurrent(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	repo := NewLocalSessionRepository(tmpDir)
	ctx := context.Background()

	sessionID := model.MustGenerateSessionID()
	session := entity.NewSessionWithDuration(
		sessionID,
		"testuser",
		"test_token_123",
		24*time.Hour,
	)

	// Save and set as current
	err := repo.Save(ctx, session)
	assert.NoError(t, err)
	err = repo.SetCurrent(ctx, session)
	assert.NoError(t, err)

	// Verify current exists
	currentSession, err := repo.GetCurrent(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, currentSession)

	// When - ClearCurrent
	err = repo.ClearCurrent(ctx)

	// Then
	assert.NoError(t, err)

	// Verify no current session
	currentSession, err = repo.GetCurrent(ctx)
	assert.Error(t, err)
	assert.Nil(t, currentSession)
	assert.True(t, cerrors.IsAppError(err, cerrors.CodeNotFound))
}

func TestLocalSessionRepository_GetByUsername(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	repo := NewLocalSessionRepository(tmpDir)
	ctx := context.Background()

	// Create multiple sessions for different users
	session1 := entity.NewSessionWithDuration(
		model.MustGenerateSessionID(),
		"user1",
		"token1",
		24*time.Hour,
	)
	session2 := entity.NewSessionWithDuration(
		model.MustGenerateSessionID(),
		"user2",
		"token2",
		24*time.Hour,
	)

	// Save both sessions
	err := repo.Save(ctx, session1)
	assert.NoError(t, err)
	err = repo.Save(ctx, session2)
	assert.NoError(t, err)

	// When - GetByUsername
	retrievedSession, err := repo.GetByUsername(ctx, "user1")

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, retrievedSession)
	assert.Equal(t, "user1", retrievedSession.Username())
	assert.Equal(t, session1.ID(), retrievedSession.ID())
}

func TestLocalSessionRepository_DeleteByUsername(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	repo := NewLocalSessionRepository(tmpDir)
	ctx := context.Background()

	// Create multiple sessions for the same user
	session1 := entity.NewSessionWithDuration(
		model.MustGenerateSessionID(),
		"testuser",
		"token1",
		24*time.Hour,
	)
	session2 := entity.NewSessionWithDuration(
		model.MustGenerateSessionID(),
		"testuser",
		"token2",
		24*time.Hour,
	)
	session3 := entity.NewSessionWithDuration(
		model.MustGenerateSessionID(),
		"otheruser",
		"token3",
		24*time.Hour,
	)

	// Save all sessions
	err := repo.Save(ctx, session1)
	assert.NoError(t, err)
	err = repo.Save(ctx, session2)
	assert.NoError(t, err)
	err = repo.Save(ctx, session3)
	assert.NoError(t, err)

	// When - DeleteByUsername
	err = repo.DeleteByUsername(ctx, "testuser")

	// Then
	assert.NoError(t, err)

	// Verify testuser sessions are deleted
	_, err = repo.GetByID(ctx, session1.ID())
	assert.Error(t, err)
	assert.True(t, cerrors.IsAppError(err, cerrors.CodeNotFound))

	_, err = repo.GetByID(ctx, session2.ID())
	assert.Error(t, err)
	assert.True(t, cerrors.IsAppError(err, cerrors.CodeNotFound))

	// Verify otheruser session still exists
	retrievedSession, err := repo.GetByID(ctx, session3.ID())
	assert.NoError(t, err)
	assert.Equal(t, "otheruser", retrievedSession.Username())
}

func TestLocalSessionRepository_DeleteExpired(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	repo := NewLocalSessionRepository(tmpDir)
	ctx := context.Background()

	// Create expired and valid sessions
	expiredSession := entity.NewSessionWithDuration(
		model.MustGenerateSessionID(),
		"user1",
		"token1",
		-time.Hour, // Expired 1 hour ago
	)
	validSession := entity.NewSessionWithDuration(
		model.MustGenerateSessionID(),
		"user2",
		"token2",
		24*time.Hour, // Valid for 24 hours
	)

	// Save both sessions
	err := repo.Save(ctx, expiredSession)
	assert.NoError(t, err)
	err = repo.Save(ctx, validSession)
	assert.NoError(t, err)

	// When - DeleteExpired
	err = repo.DeleteExpired(ctx)

	// Then
	assert.NoError(t, err)

	// Verify expired session is deleted
	_, err = repo.GetByID(ctx, expiredSession.ID())
	assert.Error(t, err)
	assert.True(t, cerrors.IsAppError(err, cerrors.CodeNotFound))

	// Verify valid session still exists
	retrievedSession, err := repo.GetByID(ctx, validSession.ID())
	assert.NoError(t, err)
	assert.Equal(t, "user2", retrievedSession.Username())
}

func TestLocalSessionRepository_List(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	repo := NewLocalSessionRepository(tmpDir)
	ctx := context.Background()

	// Create multiple sessions
	session1 := entity.NewSessionWithDuration(
		model.MustGenerateSessionID(),
		"user1",
		"token1",
		24*time.Hour,
	)
	session2 := entity.NewSessionWithDuration(
		model.MustGenerateSessionID(),
		"user2",
		"token2",
		24*time.Hour,
	)

	// Save sessions
	err := repo.Save(ctx, session1)
	assert.NoError(t, err)
	err = repo.Save(ctx, session2)
	assert.NoError(t, err)

	// When - List
	sessions, err := repo.List(ctx)

	// Then
	assert.NoError(t, err)
	assert.Len(t, sessions, 2)

	// Verify sessions are present (order might vary)
	usernames := make(map[string]bool)
	for _, session := range sessions {
		usernames[session.Username()] = true
	}
	assert.True(t, usernames["user1"])
	assert.True(t, usernames["user2"])
}

func TestLocalSessionRepository_IsValid(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	repo := NewLocalSessionRepository(tmpDir)
	ctx := context.Background()

	// Create valid and expired sessions
	validSession := entity.NewSessionWithDuration(
		model.MustGenerateSessionID(),
		"user1",
		"token1",
		24*time.Hour,
	)
	expiredSession := entity.NewSessionWithDuration(
		model.MustGenerateSessionID(),
		"user2",
		"token2",
		-time.Hour, // Expired
	)

	// Save sessions
	err := repo.Save(ctx, validSession)
	assert.NoError(t, err)
	err = repo.Save(ctx, expiredSession)
	assert.NoError(t, err)

	// When & Then - Valid session
	isValid, err := repo.IsValid(ctx, validSession.ID())
	assert.NoError(t, err)
	assert.True(t, isValid)

	// When & Then - Expired session
	isValid, err = repo.IsValid(ctx, expiredSession.ID())
	assert.NoError(t, err)
	assert.False(t, isValid)

	// When & Then - Non-existent session
	nonExistentID := model.MustGenerateSessionID()
	isValid, err = repo.IsValid(ctx, nonExistentID)
	assert.NoError(t, err)
	assert.False(t, isValid)
}