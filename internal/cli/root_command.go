package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/YuminosukeSato/AOJ-cli/pkg/logger"
)

// RootCommand represents the root command
type RootCommand struct {
	logger *logger.Logger
}

// NewRootCommand creates a new root command
func NewRootCommand() *RootCommand {
	return &RootCommand{
		logger: logger.WithGroup("root_command"),
	}
}

// Command returns the root cobra command
func (c *RootCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aoj",
		Short: "AOJ CLI - Aizu Online Judge command line tool",
		Long: `AOJ CLI is a command line tool for Aizu Online Judge.
It provides functionality similar to atcoder-cli for AOJ.

Features:
- Login to AOJ and manage sessions
- Initialize problem directories with test cases
- Run tests locally
- Submit solutions to AOJ`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// Setup context for the command
			ctx := context.Background()
			cmd.SetContext(ctx)
			return nil
		},
	}

	// Add global flags
	cmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	cmd.PersistentFlags().BoolP("quiet", "q", false, "quiet output")

	return cmd
}

// AddSubcommands adds all subcommands to the root command
func (c *RootCommand) AddSubcommands(cmd *cobra.Command, commands ...*cobra.Command) {
	cmd.AddCommand(commands...)
}

// Execute executes the root command
func (c *RootCommand) Execute(cmd *cobra.Command) error {
	return cmd.Execute()
}

// HandleError handles command execution errors
func (c *RootCommand) HandleError(err error) {
	if err != nil {
		c.logger.Error("command execution failed", "error", err)
		os.Exit(1)
	}
}