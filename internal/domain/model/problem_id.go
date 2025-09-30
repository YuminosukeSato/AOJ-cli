// Package model provides value objects for the domain layer.
package model

import (
	"regexp"
	"strings"

	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
)

// ProblemID represents a unique identifier for an AOJ problem
type ProblemID struct {
	value string
}

// Problem ID patterns for AOJ
var (
	// Course problems like ITP1_1_A, ALDS1_1_A
	coursePattern = regexp.MustCompile(`^[A-Z]+\d+_\d+_[A-Z]$`)
	// Volume problems like 0001, 1000
	volumePattern = regexp.MustCompile(`^\d{4}$`)
	// Contest problems like abc123_a, arc456_b
	contestPattern = regexp.MustCompile(`^[a-z]+\d+_[a-z]$`)
)

// NewProblemID creates a new ProblemID
func NewProblemID(value string) (ProblemID, error) {
	if value == "" {
		return ProblemID{}, cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"problem ID cannot be empty",
			nil,
		)
	}

	// Normalize the input
	normalized := strings.TrimSpace(value)
	
	if !isValidProblemIDFormat(normalized) {
		return ProblemID{}, cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"invalid problem ID format",
			cerrors.WithDetail(nil, "expected format: ITP1_1_A, 0001, or abc123_a"),
		)
	}

	return ProblemID{value: normalized}, nil
}

// MustNewProblemID creates a new ProblemID and panics on error
func MustNewProblemID(value string) ProblemID {
	id, err := NewProblemID(value)
	if err != nil {
		panic(err)
	}
	return id
}

// String returns the string representation of the problem ID
func (p ProblemID) String() string {
	return p.value
}

// Value returns the raw value
func (p ProblemID) Value() string {
	return p.value
}

// IsValid returns true if the problem ID is valid
func (p ProblemID) IsValid() bool {
	return p.value != "" && isValidProblemIDFormat(p.value)
}

// IsEmpty returns true if the problem ID is empty
func (p ProblemID) IsEmpty() bool {
	return p.value == ""
}

// Type returns the type of problem ID
func (p ProblemID) Type() string {
	if coursePattern.MatchString(p.value) {
		return "course"
	}
	if volumePattern.MatchString(p.value) {
		return "volume"
	}
	if contestPattern.MatchString(p.value) {
		return "contest"
	}
	return "unknown"
}

// IsCourse returns true if this is a course problem
func (p ProblemID) IsCourse() bool {
	return p.Type() == "course"
}

// IsVolume returns true if this is a volume problem
func (p ProblemID) IsVolume() bool {
	return p.Type() == "volume"
}

// IsContest returns true if this is a contest problem
func (p ProblemID) IsContest() bool {
	return p.Type() == "contest"
}

// Equals compares two problem IDs
func (p ProblemID) Equals(other ProblemID) bool {
	return p.value == other.value
}

// ToDirectoryName returns a directory-safe name for the problem
func (p ProblemID) ToDirectoryName() string {
	return p.value
}

// GetCourseInfo extracts course information for course problems
func (p ProblemID) GetCourseInfo() (course string, chapter int, section int, problem string, ok bool) {
	if !p.IsCourse() {
		return "", 0, 0, "", false
	}

	// Parse ITP1_1_A format
	parts := strings.Split(p.value, "_")
	if len(parts) != 3 {
		return "", 0, 0, "", false
	}

	// Extract course name and number
	courseRe := regexp.MustCompile(`^([A-Z]+)(\d+)$`)
	matches := courseRe.FindStringSubmatch(parts[0])
	if len(matches) != 3 {
		return "", 0, 0, "", false
	}

	courseName := matches[1]
	courseNum := 0
	if len(matches[2]) > 0 {
		courseNum = int(matches[2][0] - '0') // Simple single digit parsing
	}

	// Parse chapter and section
	chapterNum := 0
	if len(parts[1]) > 0 {
		chapterNum = int(parts[1][0] - '0')
	}

	problemLetter := parts[2]

	return courseName, courseNum, chapterNum, problemLetter, true
}

// isValidProblemIDFormat checks if the problem ID matches any valid format
func isValidProblemIDFormat(id string) bool {
	return coursePattern.MatchString(id) ||
		volumePattern.MatchString(id) ||
		contestPattern.MatchString(id)
}