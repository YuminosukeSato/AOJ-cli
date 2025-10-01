// Package usecase implements application business logic.
package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/repository"
	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
	"github.com/YuminosukeSato/AOJ-cli/pkg/logger"
)

// InitUseCase handles problem initialization operations
type InitUseCase struct {
	problemRepo repository.ProblemRepository
	logger      *logger.Logger
}

// NewInitUseCase creates a new InitUseCase
func NewInitUseCase(problemRepo repository.ProblemRepository) *InitUseCase {
	return &InitUseCase{
		problemRepo: problemRepo,
		logger:      logger.WithGroup("init_usecase"),
	}
}

// Execute executes the init use case
func (uc *InitUseCase) Execute(ctx context.Context, problemID string) error {
	uc.logger.InfoContext(ctx, "initializing problem directory", "problem_id", problemID)

	// Validate input
	if strings.TrimSpace(problemID) == "" {
		return cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"problem ID cannot be empty",
			nil,
		)
	}

	// Create ProblemID value object
	pid, err := model.NewProblemID(problemID)
	if err != nil {
		return cerrors.Wrap(err, "invalid problem ID")
	}

	// Create problem directory
	if err := os.MkdirAll(problemID, 0755); err != nil {
		return cerrors.Wrap(err, "failed to create problem directory")
	}

	// Get test cases from repository
	testCases, err := uc.problemRepo.GetTestCases(ctx, pid)
	if err != nil {
		uc.logger.WarnContext(ctx, "failed to get test cases, continuing with empty test cases", "error", err)
		testCases = []model.TestCase{}
	}

	// Create test directory and save test cases
	testDir := filepath.Join(problemID, "test")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return cerrors.Wrap(err, "failed to create test directory")
	}

	// Save test cases
	for i, tc := range testCases {
		inputFile := filepath.Join(testDir, fmt.Sprintf("sample-%d.in", i+1))
		outputFile := filepath.Join(testDir, fmt.Sprintf("sample-%d.out", i+1))

		if err := os.WriteFile(inputFile, []byte(tc.Input()), 0644); err != nil {
			return cerrors.Wrap(err, fmt.Sprintf("failed to write test input file %s", inputFile))
		}

		if err := os.WriteFile(outputFile, []byte(tc.Expected()), 0644); err != nil {
			return cerrors.Wrap(err, fmt.Sprintf("failed to write test output file %s", outputFile))
		}
	}

	// Create main.go template
	mainTemplate := `package main

import (
	"fmt"
)

func main() {
	// TODO: Implement solution for %s
	fmt.Println("Hello, AOJ!")
}
`
	mainContent := fmt.Sprintf(mainTemplate, problemID)
	mainFile := filepath.Join(problemID, "main.go")
	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		return cerrors.Wrap(err, "failed to create main.go")
	}

	uc.logger.InfoContext(ctx, "successfully initialized problem directory", "problem_id", problemID)
	return nil
}
