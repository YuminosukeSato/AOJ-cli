package repository

import (
	"context"
	"time"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/entity"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
)

// SubmissionRepository defines the interface for submission data access
type SubmissionRepository interface {
	// Submit submits a solution to AOJ
	Submit(ctx context.Context, submission *entity.Submission) error

	// GetByID retrieves a submission by its ID
	GetByID(ctx context.Context, id model.SubmissionID) (*entity.Submission, error)

	// GetByProblemID retrieves submissions for a specific problem
	GetByProblemID(ctx context.Context, problemID model.ProblemID, limit int) ([]*entity.Submission, error)

	// GetRecent retrieves recent submissions
	GetRecent(ctx context.Context, limit int) ([]*entity.Submission, error)

	// GetStatus retrieves the current status of a submission
	GetStatus(ctx context.Context, id model.SubmissionID) (entity.SubmissionStatus, error)

	// WatchStatus watches for status changes of a submission
	WatchStatus(ctx context.Context, id model.SubmissionID, interval time.Duration) (<-chan entity.SubmissionStatus, error)

	// Search searches for submissions by criteria
	Search(ctx context.Context, criteria SubmissionSearchCriteria) ([]*entity.Submission, error)

	// Save saves a submission (for local cache)
	Save(ctx context.Context, submission *entity.Submission) error

	// Delete deletes a submission (for local cache)
	Delete(ctx context.Context, id model.SubmissionID) error

	// Exists checks if a submission exists
	Exists(ctx context.Context, id model.SubmissionID) (bool, error)
}

// SubmissionSearchCriteria defines search criteria for submissions
type SubmissionSearchCriteria struct {
	ProblemID   *model.ProblemID
	Language    string
	Status      *entity.SubmissionStatus
	SubmittedAt *TimeRange
	Limit       int
	Offset      int
}

// TimeRange represents a time range for filtering
type TimeRange struct {
	From *time.Time
	To   *time.Time
}

// NewSubmissionSearchCriteria creates a new search criteria with defaults
func NewSubmissionSearchCriteria() SubmissionSearchCriteria {
	return SubmissionSearchCriteria{
		Limit: 50,
	}
}

// WithProblemID sets the problem ID filter
func (c SubmissionSearchCriteria) WithProblemID(problemID model.ProblemID) SubmissionSearchCriteria {
	c.ProblemID = &problemID
	return c
}

// WithLanguage sets the language filter
func (c SubmissionSearchCriteria) WithLanguage(language string) SubmissionSearchCriteria {
	c.Language = language
	return c
}

// WithStatus sets the status filter
func (c SubmissionSearchCriteria) WithStatus(status entity.SubmissionStatus) SubmissionSearchCriteria {
	c.Status = &status
	return c
}

// WithSubmittedAt sets the submitted time range filter
func (c SubmissionSearchCriteria) WithSubmittedAt(timeRange TimeRange) SubmissionSearchCriteria {
	c.SubmittedAt = &timeRange
	return c
}

// WithLimit sets the limit
func (c SubmissionSearchCriteria) WithLimit(limit int) SubmissionSearchCriteria {
	c.Limit = limit
	return c
}

// WithOffset sets the offset
func (c SubmissionSearchCriteria) WithOffset(offset int) SubmissionSearchCriteria {
	c.Offset = offset
	return c
}

// NewTimeRange creates a new time range
func NewTimeRange(from, to *time.Time) TimeRange {
	return TimeRange{
		From: from,
		To:   to,
	}
}

// Contains checks if the given time is within the range
func (tr TimeRange) Contains(t time.Time) bool {
	if tr.From != nil && t.Before(*tr.From) {
		return false
	}
	if tr.To != nil && t.After(*tr.To) {
		return false
	}
	return true
}