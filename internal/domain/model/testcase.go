package model

import (
	"strings"
	"time"
)

// TestCase represents a test case for a problem
type TestCase struct {
	id       int
	input    string
	expected string
	name     string
	timeout  time.Duration
}

// NewTestCase creates a new TestCase instance
func NewTestCase(id int, input, expected string) *TestCase {
	return &TestCase{
		id:       id,
		input:    input,
		expected: expected,
		name:     "",
		timeout:  0, // 0 means use default timeout
	}
}

// NewNamedTestCase creates a new TestCase with a name
func NewNamedTestCase(id int, input, expected, name string) *TestCase {
	return &TestCase{
		id:       id,
		input:    input,
		expected: expected,
		name:     name,
		timeout:  0,
	}
}

// ID returns the test case ID
func (tc *TestCase) ID() int {
	return tc.id
}

// Input returns the input data
func (tc *TestCase) Input() string {
	return tc.input
}

// Expected returns the expected output
func (tc *TestCase) Expected() string {
	return tc.expected
}

// Name returns the test case name
func (tc *TestCase) Name() string {
	return tc.name
}

// Timeout returns the timeout duration
func (tc *TestCase) Timeout() time.Duration {
	return tc.timeout
}

// SetName sets the test case name
func (tc *TestCase) SetName(name string) {
	tc.name = name
}

// SetTimeout sets the timeout duration
func (tc *TestCase) SetTimeout(timeout time.Duration) {
	tc.timeout = timeout
}

// UpdateInput updates the input data
func (tc *TestCase) UpdateInput(input string) {
	tc.input = input
}

// UpdateExpected updates the expected output
func (tc *TestCase) UpdateExpected(expected string) {
	tc.expected = expected
}

// IsValid validates the test case data
func (tc *TestCase) IsValid() bool {
	return tc.id >= 0 && tc.input != ""
}

// HasExpected returns true if the test case has expected output
func (tc *TestCase) HasExpected() bool {
	return tc.expected != ""
}

// HasName returns true if the test case has a name
func (tc *TestCase) HasName() bool {
	return tc.name != ""
}

// HasTimeout returns true if the test case has a custom timeout
func (tc *TestCase) HasTimeout() bool {
	return tc.timeout > 0
}

// InputLines returns the input as lines
func (tc *TestCase) InputLines() []string {
	if tc.input == "" {
		return []string{}
	}
	return strings.Split(strings.TrimRight(tc.input, "\n"), "\n")
}

// ExpectedLines returns the expected output as lines
func (tc *TestCase) ExpectedLines() []string {
	if tc.expected == "" {
		return []string{}
	}
	return strings.Split(strings.TrimRight(tc.expected, "\n"), "\n")
}

// CompareOutput compares the actual output with expected output
func (tc *TestCase) CompareOutput(actual string) bool {
	expectedLines := tc.ExpectedLines()
	actualLines := strings.Split(strings.TrimRight(actual, "\n"), "\n")

	if len(expectedLines) != len(actualLines) {
		return false
	}

	for i, expected := range expectedLines {
		if strings.TrimSpace(expected) != strings.TrimSpace(actualLines[i]) {
			return false
		}
	}

	return true
}

// GetDisplayName returns a display name for the test case
func (tc *TestCase) GetDisplayName() string {
	if tc.HasName() {
		return tc.name
	}
	return tc.GenerateDefaultName()
}

// GenerateDefaultName generates a default name based on the ID
func (tc *TestCase) GenerateDefaultName() string {
	if tc.id == 0 {
		return "sample"
	}
	return "case_" + string(rune('0'+tc.id))
}

// Clone creates a copy of the test case
func (tc *TestCase) Clone() *TestCase {
	return &TestCase{
		id:       tc.id,
		input:    tc.input,
		expected: tc.expected,
		name:     tc.name,
		timeout:  tc.timeout,
	}
}