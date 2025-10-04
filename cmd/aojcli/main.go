// Package main provides the entry point for the AOJ CLI application.
package main

import (
	"os"

	"github.com/YuminosukeSato/AOJ-cli/internal/cli"
	"github.com/YuminosukeSato/AOJ-cli/internal/infrastructure/repository"
	"github.com/YuminosukeSato/AOJ-cli/internal/usecase"
	"github.com/YuminosukeSato/AOJ-cli/pkg/config"
	"github.com/YuminosukeSato/AOJ-cli/pkg/logger"
)

const (
	aojBaseURL = "https://judgeapi.u-aizu.ac.jp"
)

func main() {
	// Initialize logger
	logConfig := logger.Config{
		Level:  logger.LevelInfo,
		Format: logger.FormatText,
		Output: os.Stderr,
	}
	logger.SetGlobal(logger.New(logConfig))

	// Load configuration
	configDir, err := config.GetConfigDir()
	if err != nil {
		logger.Error("failed to get config directory", "error", err)
		os.Exit(1)
	}

	// Ensure config directory exists
	if err := config.EnsureConfigDir(); err != nil {
		logger.Error("failed to ensure config directory", "error", err)
		os.Exit(1)
	}

	// Initialize dependencies
	dependencies := initializeDependencies(configDir)

	// Create root command
	rootCmd := cli.NewRootCommand()
	rootCommand := rootCmd.Command()

	// Create and add login command
	loginCmd := cli.NewLoginCommand(dependencies.LoginUseCase)
	loginCommand := loginCmd.Command()

	// Create and add init command
	initCmd := cli.NewInitCommand(dependencies.InitUseCase)
	initCommand := initCmd.Command()

	// Create and add submit command
	submitCmd := cli.NewSubmitCommand(dependencies.SubmitUseCase)
	submitCommand := submitCmd.Command()

	// Add subcommands to root
	rootCmd.AddSubcommands(rootCommand, loginCommand, initCommand, submitCommand)

	// Execute root command
	err = rootCmd.Execute(rootCommand)
	rootCmd.HandleError(err)
}

// Dependencies holds all application dependencies
type Dependencies struct {
	LoginUseCase  *usecase.LoginUseCase
	InitUseCase   *usecase.InitUseCase
	SubmitUseCase *usecase.SubmitUseCase
}

// initializeDependencies initializes all application dependencies
func initializeDependencies(configDir string) *Dependencies {
	// Initialize repositories
	authRepo := repository.NewAOJAuthRepository(aojBaseURL)
	sessionRepo := repository.NewLocalSessionRepository(configDir)
	problemRepo := repository.NewAOJProblemRepository(aojBaseURL)
	submissionRepo := repository.NewAOJSubmissionRepository(aojBaseURL)

	// Initialize use cases
	loginUseCase := usecase.NewLoginUseCase(authRepo, sessionRepo)
	initUseCase := usecase.NewInitUseCase(problemRepo)
	submitUseCase := usecase.NewSubmitUseCase(submissionRepo, sessionRepo)

	return &Dependencies{
		LoginUseCase:  loginUseCase,
		InitUseCase:   initUseCase,
		SubmitUseCase: submitUseCase,
	}
}
