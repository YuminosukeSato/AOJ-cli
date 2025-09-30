package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/entity"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
)

// MockAuthRepository is a mock implementation of AuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) Login(ctx context.Context, username, password string) (*entity.Session, error) {
	args := m.Called(ctx, username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Session), args.Error(1)
}

func (m *MockAuthRepository) Logout(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockAuthRepository) RefreshSession(ctx context.Context, session *entity.Session) (*entity.Session, error) {
	args := m.Called(ctx, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Session), args.Error(1)
}

func (m *MockAuthRepository) ValidateSession(ctx context.Context, session *entity.Session) (bool, error) {
	args := m.Called(ctx, session)
	return args.Bool(0), args.Error(1)
}

// MockSessionRepository is a mock implementation of SessionRepository
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Save(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) GetByID(ctx context.Context, id model.SessionID) (*entity.Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Session), args.Error(1)
}

func (m *MockSessionRepository) GetByUsername(ctx context.Context, username string) (*entity.Session, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Session), args.Error(1)
}

func (m *MockSessionRepository) GetCurrent(ctx context.Context) (*entity.Session, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Session), args.Error(1)
}

func (m *MockSessionRepository) Delete(ctx context.Context, id model.SessionID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteByUsername(ctx context.Context, username string) error {
	args := m.Called(ctx, username)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockSessionRepository) Exists(ctx context.Context, id model.SessionID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockSessionRepository) IsValid(ctx context.Context, id model.SessionID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockSessionRepository) SetCurrent(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) ClearCurrent(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockSessionRepository) List(ctx context.Context) ([]*entity.Session, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Session), args.Error(1)
}

// TDD: Red - First failing test
func TestLoginUseCase_Execute_ShouldFailWithEmptyUsername(t *testing.T) {
	// Given
	mockAuthRepo := &MockAuthRepository{}
	mockSessionRepo := &MockSessionRepository{}
	usecase := NewLoginUseCase(mockAuthRepo, mockSessionRepo)
	
	ctx := context.Background()
	request := LoginRequest{
		Username: "", // Empty username should fail
		Password: "password123",
	}

	// When
	response, err := usecase.Execute(ctx, request)

	// Then
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, cerrors.IsAppError(err, cerrors.CodeInvalidInput))
	
	// Verify no repository calls were made
	mockAuthRepo.AssertNotCalled(t, "Login")
	mockSessionRepo.AssertNotCalled(t, "Save")
}

func TestLoginUseCase_Execute_ShouldFailWithEmptyPassword(t *testing.T) {
	// Given
	mockAuthRepo := &MockAuthRepository{}
	mockSessionRepo := &MockSessionRepository{}
	usecase := NewLoginUseCase(mockAuthRepo, mockSessionRepo)
	
	ctx := context.Background()
	request := LoginRequest{
		Username: "testuser",
		Password: "", // Empty password should fail
	}

	// When
	response, err := usecase.Execute(ctx, request)

	// Then
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, cerrors.IsAppError(err, cerrors.CodeInvalidInput))
	
	// Verify no repository calls were made
	mockAuthRepo.AssertNotCalled(t, "Login")
	mockSessionRepo.AssertNotCalled(t, "Save")
}

func TestLoginUseCase_Execute_ShouldSucceedWithValidCredentials(t *testing.T) {
	// Given
	mockAuthRepo := &MockAuthRepository{}
	mockSessionRepo := &MockSessionRepository{}
	usecase := NewLoginUseCase(mockAuthRepo, mockSessionRepo)
	
	ctx := context.Background()
	request := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	// Setup expected session
	sessionID := model.MustGenerateSessionID()
	expectedSession := entity.NewSessionWithDuration(
		sessionID,
		"testuser",
		"session_token_123",
		24*time.Hour,
	)

	// Setup mock expectations
	mockAuthRepo.On("Login", ctx, "testuser", "password123").Return(expectedSession, nil)
	mockSessionRepo.On("Save", ctx, expectedSession).Return(nil)
	mockSessionRepo.On("SetCurrent", ctx, expectedSession).Return(nil)

	// When
	response, err := usecase.Execute(ctx, request)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Success)
	assert.Equal(t, "testuser", response.Username)
	assert.NotEmpty(t, response.SessionID)
	assert.NotEmpty(t, response.Message)
	
	// Verify all expected calls were made
	mockAuthRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestLoginUseCase_Execute_ShouldFailWithInvalidCredentials(t *testing.T) {
	// Given
	mockAuthRepo := &MockAuthRepository{}
	mockSessionRepo := &MockSessionRepository{}
	usecase := NewLoginUseCase(mockAuthRepo, mockSessionRepo)
	
	ctx := context.Background()
	request := LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	// Setup mock expectations
	authError := cerrors.NewAppError(
		cerrors.CodeUnauthorized,
		"invalid credentials",
		nil,
	)
	mockAuthRepo.On("Login", ctx, "testuser", "wrongpassword").Return(nil, authError)

	// When
	response, err := usecase.Execute(ctx, request)

	// Then
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, cerrors.IsAppError(err, cerrors.CodeUnauthorized))
	
	// Verify auth was called but session save was not
	mockAuthRepo.AssertExpectations(t)
	mockSessionRepo.AssertNotCalled(t, "Save")
	mockSessionRepo.AssertNotCalled(t, "SetCurrent")
}

func TestLoginUseCase_Execute_ShouldFailWhenSessionSaveFails(t *testing.T) {
	// Given
	mockAuthRepo := &MockAuthRepository{}
	mockSessionRepo := &MockSessionRepository{}
	usecase := NewLoginUseCase(mockAuthRepo, mockSessionRepo)
	
	ctx := context.Background()
	request := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	// Setup expected session
	sessionID := model.MustGenerateSessionID()
	expectedSession := entity.NewSessionWithDuration(
		sessionID,
		"testuser",
		"session_token_123",
		24*time.Hour,
	)

	// Setup mock expectations
	mockAuthRepo.On("Login", ctx, "testuser", "password123").Return(expectedSession, nil)
	
	saveError := cerrors.NewAppError(
		cerrors.CodeInternalServer,
		"failed to save session",
		nil,
	)
	mockSessionRepo.On("Save", ctx, expectedSession).Return(saveError)

	// When
	response, err := usecase.Execute(ctx, request)

	// Then
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, cerrors.IsAppError(err, cerrors.CodeInternalServer))
	
	// Verify auth was called but SetCurrent was not
	mockAuthRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
	mockSessionRepo.AssertNotCalled(t, "SetCurrent")
}