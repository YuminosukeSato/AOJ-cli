// Package usecase implements application business logic.
package usecase

import (
	"context"
	"strings"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/repository"
	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
	"github.com/YuminosukeSato/AOJ-cli/pkg/logger"
)

// LoginUseCase handles user login operations
type LoginUseCase struct {
	authRepo    repository.AuthRepository
	sessionRepo repository.SessionRepository
	logger      *logger.Logger
}

// NewLoginUseCase creates a new LoginUseCase
func NewLoginUseCase(
	authRepo repository.AuthRepository,
	sessionRepo repository.SessionRepository,
) *LoginUseCase {
	return &LoginUseCase{
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
		logger:      logger.WithGroup("login_usecase"),
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string
	Password string
}

// LoginResponse represents a login response
type LoginResponse struct {
	Success   bool   `json:"success"`
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

// Execute executes the login use case
func (uc *LoginUseCase) Execute(ctx context.Context, request LoginRequest) (*LoginResponse, error) {
	uc.logger.InfoContext(ctx, "executing login usecase", "username", request.Username)

	// Validate input
	if err := uc.validateRequest(request); err != nil {
		uc.logger.WarnContext(ctx, "invalid login request", "error", err)
		return nil, err
	}

	// Attempt authentication
	session, err := uc.authRepo.Login(ctx, request.Username, request.Password)
	if err != nil {
		uc.logger.ErrorContext(ctx, "authentication failed", 
			"username", request.Username, 
			"error", err)
		return nil, cerrors.Wrap(err, "authentication failed")
	}

	// Save session locally
	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		uc.logger.ErrorContext(ctx, "failed to save session", 
			"session_id", session.ID().MaskedString(), 
			"error", err)
		return nil, cerrors.Wrap(err, "failed to save session")
	}

	// Set as current session
	if err := uc.sessionRepo.SetCurrent(ctx, session); err != nil {
		uc.logger.ErrorContext(ctx, "failed to set current session", 
			"session_id", session.ID().MaskedString(), 
			"error", err)
		return nil, cerrors.Wrap(err, "failed to set current session")
	}

	uc.logger.InfoContext(ctx, "login successful", 
		"username", request.Username, 
		"session_id", session.ID().MaskedString())

	return &LoginResponse{
		Success:   true,
		Username:  session.Username(),
		SessionID: session.ID().String(),
		Message:   "Login successful",
	}, nil
}

// validateRequest validates the login request
func (uc *LoginUseCase) validateRequest(request LoginRequest) error {
	if strings.TrimSpace(request.Username) == "" {
		return cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"username cannot be empty",
			nil,
		)
	}

	if strings.TrimSpace(request.Password) == "" {
		return cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"password cannot be empty",
			nil,
		)
	}

	return nil
}