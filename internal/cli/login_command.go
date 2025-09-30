// Package cli provides command-line interface functionality for the AOJ CLI.
package cli

import (
	"context"
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/YuminosukeSato/AOJ-cli/internal/usecase"
	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
	"github.com/YuminosukeSato/AOJ-cli/pkg/logger"
)

// LoginCommand represents the login command
type LoginCommand struct {
	loginUseCase *usecase.LoginUseCase
	logger       *logger.Logger
}

// NewLoginCommand creates a new login command
func NewLoginCommand(loginUseCase *usecase.LoginUseCase) *LoginCommand {
	return &LoginCommand{
		loginUseCase: loginUseCase,
		logger:       logger.WithGroup("login_command"),
	}
}

// Command returns the cobra command for login
func (c *LoginCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to AOJ",
		Long:  "Authenticate with Aizu Online Judge and save session locally",
		RunE:  c.run,
	}

	return cmd
}

// run executes the login command
func (c *LoginCommand) run(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	c.logger.InfoContext(ctx, "starting login command")

	// Get username from user input
	username, err := c.promptUsername()
	if err != nil {
		return cerrors.Wrap(err, "failed to get username")
	}

	// Get password from user input (hidden)
	password, err := c.promptPassword()
	if err != nil {
		return cerrors.Wrap(err, "failed to get password")
	}

	// Execute login use case
	request := usecase.LoginRequest{
		Username: username,
		Password: password,
	}

	response, err := c.loginUseCase.Execute(ctx, request)
	if err != nil {
		return c.handleLoginError(err)
	}

	// Display success message
	c.displaySuccessMessage(response)

	c.logger.InfoContext(ctx, "login command completed successfully", 
		"username", response.Username)

	return nil
}

// promptUsername prompts the user for their username
func (c *LoginCommand) promptUsername() (string, error) {
	fmt.Print("Username: ")
	
	var username string
	_, err := fmt.Scanln(&username)
	if err != nil {
		return "", cerrors.Wrap(err, "failed to read username")
	}

	if username == "" {
		return "", cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"username cannot be empty",
			nil,
		)
	}

	return username, nil
}

// promptPassword prompts the user for their password (hidden input)
func (c *LoginCommand) promptPassword() (string, error) {
	fmt.Print("Password: ")
	
	// Read password without echoing to terminal
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", cerrors.Wrap(err, "failed to read password")
	}
	
	// Print newline after password input
	fmt.Println()

	password := string(passwordBytes)
	if password == "" {
		return "", cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"password cannot be empty",
			nil,
		)
	}

	return password, nil
}

// handleLoginError handles different types of login errors
func (c *LoginCommand) handleLoginError(err error) error {
	c.logger.ErrorContext(context.Background(), "login failed", "error", err)

	// Extract error code for user-friendly messages
	if cerrors.IsAppError(err, cerrors.CodeUnauthorized) {
		fmt.Println("❌ Login failed: Invalid username or password")
		return nil // Don't return error to avoid double error output
	}

	if cerrors.IsAppError(err, cerrors.CodeNetworkError) {
		fmt.Println("❌ Login failed: Unable to connect to AOJ. Please check your internet connection.")
		return nil
	}

	if cerrors.IsAppError(err, cerrors.CodeServiceUnavailable) {
		fmt.Println("❌ Login failed: AOJ service is currently unavailable. Please try again later.")
		return nil
	}

	if cerrors.IsAppError(err, cerrors.CodeInvalidInput) {
		fmt.Printf("❌ Login failed: %s\n", err.Error())
		return nil
	}

	// Generic error
	fmt.Printf("❌ Login failed: %s\n", err.Error())
	return nil
}

// displaySuccessMessage displays a success message to the user
func (c *LoginCommand) displaySuccessMessage(response *usecase.LoginResponse) {
	fmt.Println("✅ Login successful!")
	fmt.Printf("Logged in as: %s\n", response.Username)
	fmt.Printf("Session ID: %s\n", response.SessionID[:8]+"...")
	fmt.Println("You can now use AOJ CLI commands.")
}