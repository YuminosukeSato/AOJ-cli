// Package cli provides command-line interface functionality for the AOJ CLI.
package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YuminosukeSato/AOJ-cli/internal/usecase"
	"github.com/YuminosukeSato/AOJ-cli/pkg/logger"
)

// InitCommand represents the init command
type InitCommand struct {
	initUseCase *usecase.InitUseCase
	logger      *logger.Logger
}

// NewInitCommand creates a new init command
func NewInitCommand(initUseCase *usecase.InitUseCase) *InitCommand {
	return &InitCommand{
		initUseCase: initUseCase,
		logger:      logger.WithGroup("init_command"),
	}
}

// Command returns the cobra command for init
func (c *InitCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <problem-id>",
		Short: "Initialize a problem directory",
		Long: `Initialize a new problem directory with the given problem ID.
This command will:
- Create a directory named after the problem ID
- Download test cases from AOJ
- Generate solution template files`,
		Args: cobra.ExactArgs(1),
		RunE: c.run,
	}

	return cmd
}

// run executes the init command
func (c *InitCommand) run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	problemID := args[0]

	c.logger.InfoContext(ctx, "initializing problem directory", "problem_id", problemID)

	// Execute the use case
	if err := c.initUseCase.Execute(ctx, problemID); err != nil {
		c.logger.ErrorContext(ctx, "failed to initialize problem", "problem_id", problemID, "error", err)
		return fmt.Errorf("failed to initialize problem %s: %w", problemID, err)
	}

	c.logger.InfoContext(ctx, "successfully initialized problem directory", "problem_id", problemID)
	fmt.Printf("Successfully initialized problem: %s\n", problemID)
	return nil
}
