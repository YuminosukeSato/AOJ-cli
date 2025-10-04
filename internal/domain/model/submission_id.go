package model

import (
	"strconv"
	"strings"
	"time"

	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
)

// SubmissionID represents a unique identifier for a submission
type SubmissionID struct {
	value string
}

// NewSubmissionID creates a new SubmissionID
func NewSubmissionID(value string) (SubmissionID, error) {
	if value == "" {
		return SubmissionID{}, cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"submission ID cannot be empty",
			nil,
		)
	}

	normalized := strings.TrimSpace(value)

	if !isValidSubmissionIDFormat(normalized) {
		return SubmissionID{}, cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"invalid submission ID format",
			cerrors.WithDetail(nil, "expected numeric ID"),
		)
	}

	return SubmissionID{value: normalized}, nil
}

// NewSubmissionIDFromInt creates a new SubmissionID from an integer
func NewSubmissionIDFromInt(id int64) SubmissionID {
	return SubmissionID{value: strconv.FormatInt(id, 10)}
}

// GenerateSubmissionID generates a new submission ID
// Note: In a real implementation, this would be assigned by the server
// For now, we use a temporary client-side ID based on timestamp
func GenerateSubmissionID() (SubmissionID, error) {
	// Use Unix nano timestamp as temporary ID
	timestamp := time.Now().UnixNano()
	return NewSubmissionIDFromInt(timestamp), nil
}

// MustNewSubmissionID creates a new SubmissionID and panics on error
func MustNewSubmissionID(value string) SubmissionID {
	id, err := NewSubmissionID(value)
	if err != nil {
		panic(err)
	}
	return id
}

// String returns the string representation of the submission ID
func (s SubmissionID) String() string {
	return s.value
}

// Value returns the raw value
func (s SubmissionID) Value() string {
	return s.value
}

// IsValid returns true if the submission ID is valid
func (s SubmissionID) IsValid() bool {
	return s.value != "" && isValidSubmissionIDFormat(s.value)
}

// IsEmpty returns true if the submission ID is empty
func (s SubmissionID) IsEmpty() bool {
	return s.value == ""
}

// ToInt64 converts the submission ID to int64
func (s SubmissionID) ToInt64() (int64, error) {
	return strconv.ParseInt(s.value, 10, 64)
}

// MustToInt64 converts the submission ID to int64 and panics on error
func (s SubmissionID) MustToInt64() int64 {
	id, err := s.ToInt64()
	if err != nil {
		panic(err)
	}
	return id
}

// Equals compares two submission IDs
func (s SubmissionID) Equals(other SubmissionID) bool {
	return s.value == other.value
}

// Compare compares two submission IDs numerically
// Returns -1 if s < other, 0 if s == other, 1 if s > other
func (s SubmissionID) Compare(other SubmissionID) int {
	sInt, sErr := s.ToInt64()
	otherInt, otherErr := other.ToInt64()

	if sErr != nil || otherErr != nil {
		// Fallback to string comparison
		if s.value < other.value {
			return -1
		} else if s.value > other.value {
			return 1
		}
		return 0
	}

	if sInt < otherInt {
		return -1
	} else if sInt > otherInt {
		return 1
	}
	return 0
}

// IsGreaterThan returns true if this submission ID is greater than other
func (s SubmissionID) IsGreaterThan(other SubmissionID) bool {
	return s.Compare(other) > 0
}

// IsLessThan returns true if this submission ID is less than other
func (s SubmissionID) IsLessThan(other SubmissionID) bool {
	return s.Compare(other) < 0
}

// isValidSubmissionIDFormat checks if the submission ID is a valid numeric format
func isValidSubmissionIDFormat(id string) bool {
	if id == "" {
		return false
	}

	_, err := strconv.ParseInt(id, 10, 64)
	return err == nil
}
