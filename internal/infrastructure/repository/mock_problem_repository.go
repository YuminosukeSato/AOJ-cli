package repository

import (
	"context"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/entity"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/repository"
)

// MockProblemRepository is a mock implementation of ProblemRepository for testing
type MockProblemRepository struct{}

// NewMockProblemRepository creates a new mock problem repository
func NewMockProblemRepository() repository.ProblemRepository {
	return &MockProblemRepository{}
}

// GetByID retrieves a problem by its ID
func (r *MockProblemRepository) GetByID(_ context.Context, _ model.ProblemID) (*entity.Problem, error) {
	return nil, nil
}

// GetByIDs retrieves multiple problems by their IDs
func (r *MockProblemRepository) GetByIDs(_ context.Context, _ []model.ProblemID) ([]*entity.Problem, error) {
	return nil, nil
}

// Search searches for problems by criteria
func (r *MockProblemRepository) Search(_ context.Context, _ repository.ProblemSearchCriteria) ([]*entity.Problem, error) {
	return nil, nil
}

// Save saves a problem
func (r *MockProblemRepository) Save(_ context.Context, _ *entity.Problem) error {
	return nil
}

// Delete deletes a problem by its ID
func (r *MockProblemRepository) Delete(_ context.Context, _ model.ProblemID) error {
	return nil
}

// Exists checks if a problem exists
func (r *MockProblemRepository) Exists(_ context.Context, _ model.ProblemID) (bool, error) {
	return true, nil
}

// GetTestCases retrieves test cases for a problem
func (r *MockProblemRepository) GetTestCases(_ context.Context, _ model.ProblemID) ([]model.TestCase, error) {
	// Return empty test cases for now
	// TODO: Fetch actual test cases from AOJ
	return []model.TestCase{}, nil
}

// SaveTestCases saves test cases for a problem
func (r *MockProblemRepository) SaveTestCases(_ context.Context, _ model.ProblemID, _ []model.TestCase) error {
	return nil
}
