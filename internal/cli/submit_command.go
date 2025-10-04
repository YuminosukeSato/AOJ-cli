// Package cli provides command-line interface functionality for the AOJ CLI.
package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YuminosukeSato/AOJ-cli/internal/usecase"
	"github.com/YuminosukeSato/AOJ-cli/pkg/logger"
)

// SubmitCommand represents the submit command
type SubmitCommand struct {
	submitUseCase *usecase.SubmitUseCase
	logger        *logger.Logger
}

// NewSubmitCommand creates a new submit command
func NewSubmitCommand(submitUseCase *usecase.SubmitUseCase) *SubmitCommand {
	return &SubmitCommand{
		submitUseCase: submitUseCase,
		logger:        logger.WithGroup("submit_command"),
	}
}

// Command returns the cobra command for submit
func (c *SubmitCommand) Command() *cobra.Command {
	var (
		problemID string
		filePath  string
		language  string
	)

	cmd := &cobra.Command{
		Use:   "submit",
		Short: "Submit a solution to AOJ",
		Long: `Submit a solution to AOJ for the current problem.

By default, this command:
- Uses the current directory name as the problem ID
- Submits the main.go file
- Auto-detects the language from the file extension

Examples:
  # Submit main.go in current directory (problem ID from directory name)
  aoj submit

  # Submit a specific file
  aoj submit --file solution.cpp

  # Submit with explicit problem ID
  aoj submit --problem-id ITP1_1_A

  # Submit with explicit language
  aoj submit --language C++17`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, problemID, filePath, language)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&problemID, "problem-id", "p", "", "Problem ID (default: current directory name)")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Source file to submit (default: main.go)")
	cmd.Flags().StringVarP(&language, "language", "l", "", "Programming language (default: auto-detect from extension)")

	return cmd
}

// run executes the submit command
func (c *SubmitCommand) run(cmd *cobra.Command, problemID, filePath, language string) error {
	ctx := cmd.Context()

	c.logger.InfoContext(ctx, "executing submit command",
		"problem_id", problemID,
		"file_path", filePath,
		"language", language)

	// Prepare options
	opts := usecase.SubmitOptions{
		ProblemID: problemID,
		FilePath:  filePath,
		Language:  language,
	}

	// Execute use case
	submission, err := c.submitUseCase.Execute(ctx, opts)
	if err != nil {
		c.logger.ErrorContext(ctx, "submission failed", "error", err)
		return fmt.Errorf("submission failed: %w", err)
	}

	// Display result
	fmt.Printf("Successfully submitted solution!\n")
	fmt.Printf("Problem ID: %s\n", submission.ProblemID().String())
	fmt.Printf("Language: %s\n", submission.Language())
	fmt.Printf("Status: %s\n", submission.Status())
	fmt.Printf("Submission ID: %s\n", submission.ID().String())

	if submission.IsAccepted() {
		fmt.Printf("\n\u001b[32m✓ Accepted!\u001b[0m\n")
	} else if submission.HasError() {
		fmt.Printf("\n\u001b[31m✗ %s\u001b[0m\n", submission.Status())
		if submission.Message() != "" {
			fmt.Printf("Message: %s\n", submission.Message())
		}
	}

	return nil
}
