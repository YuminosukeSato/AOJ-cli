package model

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
)

// SessionID represents a unique identifier for a session
type SessionID struct {
	value string
}

// Session ID format: hexadecimal string (32-128 characters)
var sessionIDPattern = regexp.MustCompile(`^[a-fA-F0-9]{32,128}$`)

// NewSessionID creates a new SessionID
func NewSessionID(value string) (SessionID, error) {
	if value == "" {
		return SessionID{}, cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"session ID cannot be empty",
			nil,
		)
	}

	normalized := strings.TrimSpace(value)
	
	if !isValidSessionIDFormat(normalized) {
		return SessionID{}, cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"invalid session ID format",
			cerrors.WithDetail(nil, "expected hexadecimal string (32-128 characters)"),
		)
	}

	return SessionID{value: normalized}, nil
}

// GenerateSessionID generates a new random session ID
func GenerateSessionID() (SessionID, error) {
	bytes := make([]byte, 32) // 32 bytes = 64 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return SessionID{}, cerrors.Wrap(err, "failed to generate random session ID")
	}
	
	value := hex.EncodeToString(bytes)
	return SessionID{value: value}, nil
}

// MustGenerateSessionID generates a new random session ID and panics on error
func MustGenerateSessionID() SessionID {
	id, err := GenerateSessionID()
	if err != nil {
		panic(err)
	}
	return id
}

// MustNewSessionID creates a new SessionID and panics on error
func MustNewSessionID(value string) SessionID {
	id, err := NewSessionID(value)
	if err != nil {
		panic(err)
	}
	return id
}

// String returns the string representation of the session ID
func (s SessionID) String() string {
	return s.value
}

// Value returns the raw value
func (s SessionID) Value() string {
	return s.value
}

// IsValid returns true if the session ID is valid
func (s SessionID) IsValid() bool {
	return s.value != "" && isValidSessionIDFormat(s.value)
}

// IsEmpty returns true if the session ID is empty
func (s SessionID) IsEmpty() bool {
	return s.value == ""
}

// Equals compares two session IDs
func (s SessionID) Equals(other SessionID) bool {
	return s.value == other.value
}

// Bytes returns the session ID as bytes (decoded from hex)
func (s SessionID) Bytes() ([]byte, error) {
	return hex.DecodeString(s.value)
}

// MustBytes returns the session ID as bytes and panics on error
func (s SessionID) MustBytes() []byte {
	bytes, err := s.Bytes()
	if err != nil {
		panic(err)
	}
	return bytes
}

// Length returns the length of the session ID string
func (s SessionID) Length() int {
	return len(s.value)
}

// IsSecure returns true if the session ID has sufficient length for security
func (s SessionID) IsSecure() bool {
	return len(s.value) >= 64 // At least 32 bytes (64 hex chars)
}

// Truncate returns a truncated version of the session ID for logging
func (s SessionID) Truncate(length int) string {
	if length <= 0 {
		return ""
	}
	if length >= len(s.value) {
		return s.value
	}
	return s.value[:length] + "..."
}

// MaskedString returns a masked version of the session ID for logging
func (s SessionID) MaskedString() string {
	if len(s.value) <= 8 {
		return strings.Repeat("*", len(s.value))
	}
	return s.value[:4] + strings.Repeat("*", len(s.value)-8) + s.value[len(s.value)-4:]
}

// isValidSessionIDFormat checks if the session ID is a valid hexadecimal format
func isValidSessionIDFormat(id string) bool {
	return sessionIDPattern.MatchString(id)
}