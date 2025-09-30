package entity

import (
	"time"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
)

// Session represents an AOJ login session
type Session struct {
	id        model.SessionID
	username  string
	token     string
	expiresAt time.Time
	createdAt time.Time
	lastUsed  time.Time
}

// NewSession creates a new Session instance
func NewSession(
	id model.SessionID,
	username, token string,
	expiresAt time.Time,
) *Session {
	now := time.Now()
	return &Session{
		id:        id,
		username:  username,
		token:     token,
		expiresAt: expiresAt,
		createdAt: now,
		lastUsed:  now,
	}
}

// NewSessionWithDuration creates a new Session with a duration from now
func NewSessionWithDuration(
	id model.SessionID,
	username, token string,
	duration time.Duration,
) *Session {
	now := time.Now()
	return &Session{
		id:        id,
		username:  username,
		token:     token,
		expiresAt: now.Add(duration),
		createdAt: now,
		lastUsed:  now,
	}
}

// ID returns the session ID
func (s *Session) ID() model.SessionID {
	return s.id
}

// Username returns the username
func (s *Session) Username() string {
	return s.username
}

// Token returns the session token
func (s *Session) Token() string {
	return s.token
}

// ExpiresAt returns the expiration time
func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

// CreatedAt returns the creation time
func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}

// LastUsed returns the last used time
func (s *Session) LastUsed() time.Time {
	return s.lastUsed
}

// IsValid returns true if the session is valid
func (s *Session) IsValid() bool {
	return s.id.IsValid() &&
		s.username != "" &&
		s.token != "" &&
		!s.IsExpired()
}

// IsExpired returns true if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.expiresAt)
}

// IsExpiredAt returns true if the session is expired at the given time
func (s *Session) IsExpiredAt(t time.Time) bool {
	return t.After(s.expiresAt)
}

// TimeUntilExpiry returns the duration until the session expires
func (s *Session) TimeUntilExpiry() time.Duration {
	return time.Until(s.expiresAt)
}

// RemainingTime returns the remaining time for the session
func (s *Session) RemainingTime() time.Duration {
	remaining := s.TimeUntilExpiry()
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Age returns the age of the session
func (s *Session) Age() time.Duration {
	return time.Since(s.createdAt)
}

// TimeSinceLastUse returns the time since the session was last used
func (s *Session) TimeSinceLastUse() time.Duration {
	return time.Since(s.lastUsed)
}

// UpdateLastUsed updates the last used time to now
func (s *Session) UpdateLastUsed() {
	s.lastUsed = time.Now()
}

// UpdateLastUsedAt updates the last used time to the specified time
func (s *Session) UpdateLastUsedAt(t time.Time) {
	s.lastUsed = t
}

// Refresh extends the session expiration time
func (s *Session) Refresh(duration time.Duration) {
	s.expiresAt = time.Now().Add(duration)
	s.UpdateLastUsed()
}

// RefreshFromNow extends the session expiration time from the current expiration
func (s *Session) RefreshFromNow(duration time.Duration) {
	now := time.Now()
	if s.expiresAt.Before(now) {
		s.expiresAt = now.Add(duration)
	} else {
		s.expiresAt = s.expiresAt.Add(duration)
	}
	s.UpdateLastUsed()
}

// UpdateToken updates the session token
func (s *Session) UpdateToken(token string) {
	s.token = token
	s.UpdateLastUsed()
}

// Clone creates a copy of the session
func (s *Session) Clone() *Session {
	return &Session{
		id:        s.id,
		username:  s.username,
		token:     s.token,
		expiresAt: s.expiresAt,
		createdAt: s.createdAt,
		lastUsed:  s.lastUsed,
	}
}

// ToMap converts the session to a map for serialization
func (s *Session) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         s.id.String(),
		"username":   s.username,
		"token":      s.token,
		"expires_at": s.expiresAt.Unix(),
		"created_at": s.createdAt.Unix(),
		"last_used":  s.lastUsed.Unix(),
	}
}

// FromMap creates a session from a map
func FromMap(data map[string]interface{}) (*Session, error) {
	id, err := model.NewSessionID(data["id"].(string))
	if err != nil {
		return nil, err
	}

	username := data["username"].(string)
	token := data["token"].(string)
	
	expiresAt := time.Unix(int64(data["expires_at"].(float64)), 0)
	createdAt := time.Unix(int64(data["created_at"].(float64)), 0)
	lastUsed := time.Unix(int64(data["last_used"].(float64)), 0)

	session := &Session{
		id:        id,
		username:  username,
		token:     token,
		expiresAt: expiresAt,
		createdAt: createdAt,
		lastUsed:  lastUsed,
	}

	return session, nil
}