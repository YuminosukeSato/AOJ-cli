// Package repository implements the data access layer.
package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/entity"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/repository"
	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
	"github.com/YuminosukeSato/AOJ-cli/pkg/logger"
)

// AOJAuthRepository implements AuthRepository for AOJ API
type AOJAuthRepository struct {
	baseURL    string
	httpClient *http.Client
	logger     *logger.Logger
}

// NewAOJAuthRepository creates a new AOJAuthRepository
func NewAOJAuthRepository(baseURL string) repository.AuthRepository {
	return &AOJAuthRepository{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger.WithGroup("aoj_auth_repository"),
	}
}

// LoginRequest represents the JSON payload for AOJ login
type LoginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// LoginResponse represents the JSON response from AOJ login
type LoginResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SessionID string `json:"sessionId"`
	Token     string `json:"token"`
}

// Login authenticates a user with AOJ and returns a session
func (r *AOJAuthRepository) Login(ctx context.Context, username, password string) (*entity.Session, error) {
	r.logger.InfoContext(ctx, "attempting AOJ login", "username", username)

	// Prepare request payload
	loginReq := LoginRequest{
		ID:       username,
		Password: password,
	}

	payload, err := json.Marshal(loginReq)
	if err != nil {
		return nil, cerrors.Wrap(err, "failed to marshal login request")
	}

	// Create HTTP request
	url := r.baseURL + "/session"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, cerrors.Wrap(err, "failed to create HTTP request")
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	// Execute request
	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.logger.ErrorContext(ctx, "HTTP request failed", "error", err)
		return nil, cerrors.NewAppError(
			cerrors.CodeNetworkError,
			"failed to connect to AOJ",
			err,
		)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			r.logger.WarnContext(ctx, "failed to close response body", "error", err)
		}
	}()

	// Handle different response status codes
	switch resp.StatusCode {
	case http.StatusOK:
		return r.parseLoginResponse(ctx, resp)
	case http.StatusUnauthorized:
		r.logger.WarnContext(ctx, "authentication failed", "username", username)
		return nil, cerrors.NewAppError(
			cerrors.CodeUnauthorized,
			"invalid username or password",
			nil,
		)
	case http.StatusBadRequest:
		return nil, cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"invalid login request format",
			nil,
		)
	case http.StatusInternalServerError:
		return nil, cerrors.NewAppError(
			cerrors.CodeServiceUnavailable,
			"AOJ server error",
			nil,
		)
	default:
		r.logger.ErrorContext(ctx, "unexpected response status", "status", resp.StatusCode)
		return nil, cerrors.NewAppError(
			cerrors.CodeInternalServer,
			"unexpected response from AOJ",
			cerrors.WithDetail(nil, "status_code: "+resp.Status),
		)
	}
}

// parseLoginResponse parses the successful login response
func (r *AOJAuthRepository) parseLoginResponse(ctx context.Context, resp *http.Response) (*entity.Session, error) {
	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, cerrors.Wrap(err, "failed to decode login response")
	}

	// Generate session ID
	sessionID, err := model.GenerateSessionID()
	if err != nil {
		return nil, cerrors.Wrap(err, "failed to generate session ID")
	}

	// Create session entity
	session := entity.NewSessionWithDuration(
		sessionID,
		loginResp.ID,
		loginResp.Token,
		24*time.Hour, // AOJ sessions typically last 24 hours
	)

	r.logger.InfoContext(ctx, "login successful", 
		"username", loginResp.ID,
		"session_id", sessionID.MaskedString())

	return session, nil
}

// Logout logs out a user by invalidating their session
func (r *AOJAuthRepository) Logout(ctx context.Context, session *entity.Session) error {
	r.logger.InfoContext(ctx, "attempting AOJ logout", 
		"session_id", session.ID().MaskedString())

	// Create HTTP request to delete session
	url := r.baseURL + "/session"
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return cerrors.Wrap(err, "failed to create logout request")
	}

	// Add session token to request (implementation depends on AOJ API)
	req.Header.Set("Authorization", "Bearer "+session.Token())

	// Execute request
	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.logger.ErrorContext(ctx, "logout request failed", "error", err)
		return cerrors.NewAppError(
			cerrors.CodeNetworkError,
			"failed to connect to AOJ for logout",
			err,
		)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			r.logger.WarnContext(ctx, "failed to close response body", "error", err)
		}
	}()

	// Handle response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		r.logger.WarnContext(ctx, "logout request returned unexpected status", 
			"status", resp.StatusCode)
		// Don't return error for logout - best effort
	}

	r.logger.InfoContext(ctx, "logout completed", 
		"session_id", session.ID().MaskedString())

	return nil
}

// RefreshSession refreshes an existing session
func (r *AOJAuthRepository) RefreshSession(ctx context.Context, session *entity.Session) (*entity.Session, error) {
	r.logger.InfoContext(ctx, "refreshing session", 
		"session_id", session.ID().MaskedString())

	// For AOJ, we might need to validate the current session and extend it
	// Implementation depends on AOJ API capabilities
	isValid, err := r.ValidateSession(ctx, session)
	if err != nil {
		return nil, cerrors.Wrap(err, "failed to validate session for refresh")
	}

	if !isValid {
		return nil, cerrors.NewAppError(
			cerrors.CodeUnauthorized,
			"session is no longer valid",
			nil,
		)
	}

	// Create new session with extended duration
	newSessionID, err := model.GenerateSessionID()
	if err != nil {
		return nil, cerrors.Wrap(err, "failed to generate new session ID")
	}

	refreshedSession := entity.NewSessionWithDuration(
		newSessionID,
		session.Username(),
		session.Token(), // Keep the same token or get a new one from AOJ
		24*time.Hour,
	)

	r.logger.InfoContext(ctx, "session refreshed", 
		"old_session_id", session.ID().MaskedString(),
		"new_session_id", newSessionID.MaskedString())

	return refreshedSession, nil
}

// ValidateSession validates if a session is still active on the server
func (r *AOJAuthRepository) ValidateSession(ctx context.Context, session *entity.Session) (bool, error) {
	r.logger.DebugContext(ctx, "validating session", 
		"session_id", session.ID().MaskedString())

	// Check if session is expired locally first
	if session.IsExpired() {
		r.logger.DebugContext(ctx, "session is locally expired")
		return false, nil
	}

	// Make a lightweight request to AOJ to validate the session
	// This could be a request to get user info or any authenticated endpoint
	url := r.baseURL + "/user/" + session.Username()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, cerrors.Wrap(err, "failed to create validation request")
	}

	req.Header.Set("Authorization", "Bearer "+session.Token())

	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.logger.ErrorContext(ctx, "session validation request failed", "error", err)
		return false, cerrors.NewAppError(
			cerrors.CodeNetworkError,
			"failed to validate session with AOJ",
			err,
		)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			r.logger.WarnContext(ctx, "failed to close response body", "error", err)
		}
	}()

	isValid := resp.StatusCode == http.StatusOK

	r.logger.DebugContext(ctx, "session validation completed", 
		"session_id", session.ID().MaskedString(),
		"is_valid", isValid)

	return isValid, nil
}