package repository

import (
	"context"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/entity"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
)

// ProblemRepository defines the interface for problem data access
type ProblemRepository interface {
	// GetByID retrieves a problem by its ID
	GetByID(ctx context.Context, id model.ProblemID) (*entity.Problem, error)

	// GetByIDs retrieves multiple problems by their IDs
	GetByIDs(ctx context.Context, ids []model.ProblemID) ([]*entity.Problem, error)

	// Search searches for problems by criteria
	Search(ctx context.Context, criteria ProblemSearchCriteria) ([]*entity.Problem, error)

	// Save saves a problem
	Save(ctx context.Context, problem *entity.Problem) error

	// Delete deletes a problem by its ID
	Delete(ctx context.Context, id model.ProblemID) error

	// Exists checks if a problem exists
	Exists(ctx context.Context, id model.ProblemID) (bool, error)

	// GetTestCases retrieves test cases for a problem
	GetTestCases(ctx context.Context, problemID model.ProblemID) ([]model.TestCase, error)

	// SaveTestCases saves test cases for a problem
	SaveTestCases(ctx context.Context, problemID model.ProblemID, testCases []model.TestCase) error
}

// ProblemSearchCriteria defines search criteria for problems
type ProblemSearchCriteria struct {
	Category   string
	Difficulty *int // nil means any difficulty
	Title      string
	Limit      int
	Offset     int
}

// NewProblemSearchCriteria creates a new search criteria with defaults
func NewProblemSearchCriteria() ProblemSearchCriteria {
	return ProblemSearchCriteria{
		Limit: 50,
	}
}

// WithCategory sets the category filter
func (c ProblemSearchCriteria) WithCategory(category string) ProblemSearchCriteria {
	c.Category = category
	return c
}

// WithDifficulty sets the difficulty filter
func (c ProblemSearchCriteria) WithDifficulty(difficulty int) ProblemSearchCriteria {
	c.Difficulty = &difficulty
	return c
}

// WithTitle sets the title filter
func (c ProblemSearchCriteria) WithTitle(title string) ProblemSearchCriteria {
	c.Title = title
	return c
}

// WithLimit sets the limit
func (c ProblemSearchCriteria) WithLimit(limit int) ProblemSearchCriteria {
	c.Limit = limit
	return c
}

// WithOffset sets the offset
func (c ProblemSearchCriteria) WithOffset(offset int) ProblemSearchCriteria {
	c.Offset = offset
	return c
}