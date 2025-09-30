// Package config provides configuration management utilities.
package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
	"github.com/YuminosukeSato/AOJ-cli/pkg/logger"
)

// Config represents the application configuration
type Config struct {
	Login  LoginConfig  `toml:"login"`
	Init   InitConfig   `toml:"init"`
	Test   TestConfig   `toml:"test"`
	Submit SubmitConfig `toml:"submit"`
}

// LoginConfig holds login-related configuration
type LoginConfig struct {
	SessionFile string `toml:"session_file"`
}

// InitConfig holds init command configuration
type InitConfig struct {
	TemplateFile    string `toml:"template_file"`
	Language        string `toml:"language"`
	FetchTestcases  bool   `toml:"fetch_testcases"`
	DefaultTemplate string `toml:"default_template"`
}

// TestConfig holds test command configuration
type TestConfig struct {
	BuildCommand string  `toml:"build_command"`
	RunCommand   string  `toml:"run_command"`
	Timeout      float64 `toml:"timeout"`
	Parallel     bool    `toml:"parallel"`
}

// SubmitConfig holds submit command configuration
type SubmitConfig struct {
	SourceFile string `toml:"source_file"`
	Language   string `toml:"language"`
	Watch      bool   `toml:"watch"`
}

// LanguageConfig represents language-specific configuration
type LanguageConfig struct {
	Extension    string `toml:"extension"`
	BuildCommand string `toml:"build_command"`
	RunCommand   string `toml:"run_command"`
	AOJLanguageID string `toml:"aoj_language_id"`
}

// Languages holds all language configurations
type Languages map[string]LanguageConfig

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	aojDir := filepath.Join(homeDir, ".aoj-cli")

	return &Config{
		Login: LoginConfig{
			SessionFile: filepath.Join(aojDir, "session.json"),
		},
		Init: InitConfig{
			TemplateFile:    filepath.Join(aojDir, "template.cpp"),
			Language:        "C++17",
			FetchTestcases:  true,
			DefaultTemplate: defaultCppTemplate,
		},
		Test: TestConfig{
			BuildCommand: "g++ -std=c++17 -O2 -o a.out main.cpp",
			RunCommand:   "./a.out",
			Timeout:      2.0,
			Parallel:     true,
		},
		Submit: SubmitConfig{
			SourceFile: "main.cpp",
			Language:   "C++17",
			Watch:      true,
		},
	}
}

// DefaultLanguages returns the default language configurations
func DefaultLanguages() Languages {
	return Languages{
		"cpp17": {
			Extension:     "cpp",
			BuildCommand:  "g++ -std=c++17 -O2 -o a.out {file}",
			RunCommand:    "./a.out",
			AOJLanguageID: "C++17",
		},
		"cpp23": {
			Extension:     "cpp",
			BuildCommand:  "g++ -std=c++23 -O2 -o a.out {file}",
			RunCommand:    "./a.out",
			AOJLanguageID: "C++23",
		},
		"python": {
			Extension:     "py",
			BuildCommand:  "",
			RunCommand:    "python3 {file}",
			AOJLanguageID: "Python3",
		},
		"java": {
			Extension:     "java",
			BuildCommand:  "javac {file}",
			RunCommand:    "java Main",
			AOJLanguageID: "Java",
		},
		"go": {
			Extension:     "go",
			BuildCommand:  "go build -o main {file}",
			RunCommand:    "./main",
			AOJLanguageID: "Go",
		},
	}
}

const defaultCppTemplate = `#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include <map>
#include <set>
#include <queue>
#include <stack>
#include <cmath>
#include <climits>

using namespace std;

int main() {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);
    
    // TODO: Implement solution
    
    return 0;
}
`

// Load loads configuration from the specified file
func Load(filePath string) (*Config, error) {
	config := DefaultConfig()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Debug("config file not found, using defaults", "path", filePath)
		return config, nil
	}

	if _, err := toml.DecodeFile(filePath, config); err != nil {
		return nil, cerrors.Wrap(err, "failed to decode config file")
	}

	logger.Debug("config loaded successfully", "path", filePath)
	return config, nil
}

// Save saves configuration to the specified file
func Save(config *Config, filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return cerrors.Wrap(err, "failed to create config directory")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return cerrors.Wrap(err, "failed to create config file")
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Warn("failed to close config file", "error", err)
		}
	}()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return cerrors.Wrap(err, "failed to encode config")
	}

	logger.Debug("config saved successfully", "path", filePath)
	return nil
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", cerrors.Wrap(err, "failed to get user home directory")
	}

	configDir := filepath.Join(homeDir, ".aoj-cli")
	return configDir, nil
}

// GetConfigPath returns the default configuration file path
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.toml"), nil
}

// EnsureConfigDir ensures the configuration directory exists
func EnsureConfigDir() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return cerrors.Wrap(err, "failed to create config directory")
	}

	return nil
}

// LoadDefault loads configuration from the default location
func LoadDefault() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	return Load(configPath)
}

// SaveDefault saves configuration to the default location
func SaveDefault(config *Config) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	return Save(config, configPath)
}

// GetLanguageConfig returns the configuration for the specified language
func GetLanguageConfig(langName string) (LanguageConfig, bool) {
	languages := DefaultLanguages()
	lang, exists := languages[langName]
	return lang, exists
}

// ValidateConfig validates the configuration
func ValidateConfig(config *Config) error {
	if config == nil {
		return cerrors.New("config is nil")
	}

	if config.Test.Timeout <= 0 {
		return cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"test timeout must be positive",
			nil,
		)
	}

	if config.Submit.SourceFile == "" {
		return cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"submit source file cannot be empty",
			nil,
		)
	}

	if config.Submit.Language == "" {
		return cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"submit language cannot be empty",
			nil,
		)
	}

	return nil
}