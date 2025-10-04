// Package usecase implements application business logic.
package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/entity"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/repository"
	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
	"github.com/YuminosukeSato/AOJ-cli/pkg/logger"
)

// SubmitUseCase handles solution submission operations
type SubmitUseCase struct {
	submissionRepo repository.SubmissionRepository
	sessionRepo    repository.SessionRepository
	logger         *logger.Logger
}

// NewSubmitUseCase creates a new SubmitUseCase
func NewSubmitUseCase(
	submissionRepo repository.SubmissionRepository,
	sessionRepo repository.SessionRepository,
) *SubmitUseCase {
	return &SubmitUseCase{
		submissionRepo: submissionRepo,
		sessionRepo:    sessionRepo,
		logger:         logger.WithGroup("submit_usecase"),
	}
}

// SubmitOptions contains options for submission
type SubmitOptions struct {
	ProblemID string // Optional: explicit problem ID (defaults to directory name)
	FilePath  string // Optional: source file path (defaults to main.go)
	Language  string // Optional: language (defaults to auto-detect from extension)
}

// Execute executes the submit use case
func (uc *SubmitUseCase) Execute(ctx context.Context, opts SubmitOptions) (*entity.Submission, error) {
	uc.logger.InfoContext(ctx, "starting submission", "options", fmt.Sprintf("%+v", opts))

	// Determine problem ID
	problemID, err := uc.determineProblemID(opts.ProblemID)
	if err != nil {
		return nil, err
	}
	uc.logger.InfoContext(ctx, "determined problem ID", "problem_id", problemID.String())

	// Determine source file path
	filePath := opts.FilePath
	if filePath == "" {
		filePath = "main.go" // Default
	}

	// Read source code
	sourceCode, err := os.ReadFile(filePath)
	if err != nil {
		return nil, cerrors.Wrap(err, fmt.Sprintf("failed to read source file: %s", filePath))
	}
	uc.logger.InfoContext(ctx, "read source file", "file_path", filePath, "size", len(sourceCode))

	// Determine language
	language := opts.Language
	if language == "" {
		language = uc.detectLanguage(filePath)
	}
	uc.logger.InfoContext(ctx, "determined language", "language", language)

	// Get current session
	session, err := uc.sessionRepo.GetCurrent(ctx)
	if err != nil {
		return nil, cerrors.Wrap(err, "failed to get current session")
	}

	if session == nil {
		return nil, cerrors.NewAppError(
			cerrors.CodeUnauthorized,
			"no active session found. Please login first with 'aoj login'",
			nil,
		)
	}

	if session.IsExpired() {
		return nil, cerrors.NewAppError(
			cerrors.CodeUnauthorized,
			"session has expired. Please login again with 'aoj login'",
			nil,
		)
	}

	// Generate submission ID
	submissionID, err := model.GenerateSubmissionID()
	if err != nil {
		return nil, cerrors.Wrap(err, "failed to generate submission ID")
	}

	// Create submission entity
	submission := entity.NewSubmission(
		submissionID,
		problemID,
		language,
		string(sourceCode),
	)

	// Submit to AOJ
	if err := uc.submissionRepo.Submit(ctx, submission); err != nil {
		uc.logger.ErrorContext(ctx, "submission failed", "error", err)
		return nil, cerrors.Wrap(err, "failed to submit solution")
	}

	uc.logger.InfoContext(ctx, "submission successful",
		"submission_id", submissionID.String(),
		"problem_id", problemID.String())

	return submission, nil
}

// determineProblemID determines the problem ID from options or current directory
func (uc *SubmitUseCase) determineProblemID(explicitID string) (model.ProblemID, error) {
	if explicitID != "" {
		return model.NewProblemID(explicitID)
	}

	// Get current directory name
	cwd, err := os.Getwd()
	if err != nil {
		return model.ProblemID{}, cerrors.Wrap(err, "failed to get current directory")
	}

	dirName := filepath.Base(cwd)

	// Try to parse directory name as problem ID
	problemID, err := model.NewProblemID(dirName)
	if err != nil {
		return model.ProblemID{}, cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			fmt.Sprintf("could not determine problem ID from directory name '%s'. Please specify --problem-id", dirName),
			err,
		)
	}

	return problemID, nil
}

// detectLanguage detects the language from file extension
func (uc *SubmitUseCase) detectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	languageMap := map[string]string{
		".c":     "C",
		".cpp":   "C++14",
		".cc":    "C++14",
		".cxx":   "C++14",
		".c++":   "C++14",
		".java":  "JAVA",
		".py":    "Python3",
		".rb":    "Ruby",
		".go":    "Go",
		".js":    "JavaScript",
		".cs":    "C#",
		".php":   "PHP",
		".d":     "D",
		".rs":    "Rust",
		".kt":    "Kotlin",
		".scala": "Scala",
	}

	if lang, ok := languageMap[ext]; ok {
		return lang
	}

	// Default to C++ if unknown
	return "C++14"
}
