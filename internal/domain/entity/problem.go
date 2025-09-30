// Package entity provides domain entities for the AOJ CLI application.
package entity

import (
	"time"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
)

// Problem represents an AOJ problem
type Problem struct {
	id          model.ProblemID
	title       string
	description string
	timeLimit   time.Duration
	memoryLimit int64 // in KB
	category    string
	difficulty  int
	testCases   []model.TestCase
	createdAt   time.Time
	updatedAt   time.Time
}

// NewProblem creates a new Problem instance
func NewProblem(
	id model.ProblemID,
	title, description string,
	timeLimit time.Duration,
	memoryLimit int64,
	category string,
	difficulty int,
) *Problem {
	now := time.Now()
	return &Problem{
		id:          id,
		title:       title,
		description: description,
		timeLimit:   timeLimit,
		memoryLimit: memoryLimit,
		category:    category,
		difficulty:  difficulty,
		testCases:   make([]model.TestCase, 0),
		createdAt:   now,
		updatedAt:   now,
	}
}

// ID returns the problem ID
func (p *Problem) ID() model.ProblemID {
	return p.id
}

// Title returns the problem title
func (p *Problem) Title() string {
	return p.title
}

// Description returns the problem description
func (p *Problem) Description() string {
	return p.description
}

// TimeLimit returns the time limit
func (p *Problem) TimeLimit() time.Duration {
	return p.timeLimit
}

// MemoryLimit returns the memory limit in KB
func (p *Problem) MemoryLimit() int64 {
	return p.memoryLimit
}

// Category returns the problem category
func (p *Problem) Category() string {
	return p.category
}

// Difficulty returns the difficulty level
func (p *Problem) Difficulty() int {
	return p.difficulty
}

// TestCases returns the test cases
func (p *Problem) TestCases() []model.TestCase {
	// Return a copy to prevent external modification
	result := make([]model.TestCase, len(p.testCases))
	copy(result, p.testCases)
	return result
}

// CreatedAt returns the creation time
func (p *Problem) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns the last update time
func (p *Problem) UpdatedAt() time.Time {
	return p.updatedAt
}

// AddTestCase adds a test case to the problem
func (p *Problem) AddTestCase(testCase model.TestCase) {
	p.testCases = append(p.testCases, testCase)
	p.updatedAt = time.Now()
}

// SetTestCases sets all test cases for the problem
func (p *Problem) SetTestCases(testCases []model.TestCase) {
	p.testCases = make([]model.TestCase, len(testCases))
	copy(p.testCases, testCases)
	p.updatedAt = time.Now()
}

// UpdateTitle updates the problem title
func (p *Problem) UpdateTitle(title string) {
	p.title = title
	p.updatedAt = time.Now()
}

// UpdateDescription updates the problem description
func (p *Problem) UpdateDescription(description string) {
	p.description = description
	p.updatedAt = time.Now()
}

// HasTestCases returns true if the problem has test cases
func (p *Problem) HasTestCases() bool {
	return len(p.testCases) > 0
}

// TestCaseCount returns the number of test cases
func (p *Problem) TestCaseCount() int {
	return len(p.testCases)
}

// IsValid validates the problem data
func (p *Problem) IsValid() bool {
	return p.id.IsValid() &&
		p.title != "" &&
		p.timeLimit > 0 &&
		p.memoryLimit > 0 &&
		p.difficulty >= 0
}