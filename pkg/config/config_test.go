package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.NotNil(t, config)
	assert.NotEmpty(t, config.Login.SessionFile)
	assert.NotEmpty(t, config.Init.TemplateFile)
	assert.Equal(t, "C++17", config.Init.Language)
	assert.True(t, config.Init.FetchTestcases)
	assert.NotEmpty(t, config.Init.DefaultTemplate)
	assert.NotEmpty(t, config.Test.BuildCommand)
	assert.NotEmpty(t, config.Test.RunCommand)
	assert.Greater(t, config.Test.Timeout, 0.0)
	assert.True(t, config.Test.Parallel)
	assert.Equal(t, "main.cpp", config.Submit.SourceFile)
	assert.Equal(t, "C++17", config.Submit.Language)
	assert.True(t, config.Submit.Watch)
}

func TestDefaultLanguages(t *testing.T) {
	languages := DefaultLanguages()

	assert.NotEmpty(t, languages)

	// Test C++17 configuration
	cpp17, exists := languages["cpp17"]
	assert.True(t, exists)
	assert.Equal(t, "cpp", cpp17.Extension)
	assert.Contains(t, cpp17.BuildCommand, "g++")
	assert.Contains(t, cpp17.BuildCommand, "-std=c++17")
	assert.Equal(t, "./a.out", cpp17.RunCommand)
	assert.Equal(t, "C++17", cpp17.AOJLanguageID)

	// Test Python configuration
	python, exists := languages["python"]
	assert.True(t, exists)
	assert.Equal(t, "py", python.Extension)
	assert.Empty(t, python.BuildCommand)
	assert.Contains(t, python.RunCommand, "python3")
	assert.Equal(t, "Python3", python.AOJLanguageID)

	// Test Java configuration
	java, exists := languages["java"]
	assert.True(t, exists)
	assert.Equal(t, "java", java.Extension)
	assert.Contains(t, java.BuildCommand, "javac")
	assert.Equal(t, "java Main", java.RunCommand)
	assert.Equal(t, "Java", java.AOJLanguageID)
}

func TestLoadNonExistentFile(t *testing.T) {
	config, err := Load("/non/existent/file.toml")

	assert.NoError(t, err)
	assert.NotNil(t, config)
	// Should return default config when file doesn't exist
	assert.Equal(t, "C++17", config.Init.Language)
}

func TestSaveAndLoad(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.toml")

	// Create a test configuration
	originalConfig := &Config{
		Login: LoginConfig{
			SessionFile: "/tmp/session.json",
		},
		Init: InitConfig{
			TemplateFile:   "/tmp/template.cpp",
			Language:       "C++23",
			FetchTestcases: false,
		},
		Test: TestConfig{
			BuildCommand: "custom build command",
			RunCommand:   "custom run command",
			Timeout:      5.0,
			Parallel:     false,
		},
		Submit: SubmitConfig{
			SourceFile: "solution.cpp",
			Language:   "C++23",
			Watch:      false,
		},
	}

	// Save configuration
	err := Save(originalConfig, configPath)
	assert.NoError(t, err)

	// Check if file was created
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// Load configuration
	loadedConfig, err := Load(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, loadedConfig)

	// Compare configurations
	assert.Equal(t, originalConfig.Login.SessionFile, loadedConfig.Login.SessionFile)
	assert.Equal(t, originalConfig.Init.Language, loadedConfig.Init.Language)
	assert.Equal(t, originalConfig.Init.FetchTestcases, loadedConfig.Init.FetchTestcases)
	assert.Equal(t, originalConfig.Test.BuildCommand, loadedConfig.Test.BuildCommand)
	assert.Equal(t, originalConfig.Test.Timeout, loadedConfig.Test.Timeout)
	assert.Equal(t, originalConfig.Test.Parallel, loadedConfig.Test.Parallel)
	assert.Equal(t, originalConfig.Submit.SourceFile, loadedConfig.Submit.SourceFile)
	assert.Equal(t, originalConfig.Submit.Language, loadedConfig.Submit.Language)
	assert.Equal(t, originalConfig.Submit.Watch, loadedConfig.Submit.Watch)
}

func TestLoadInvalidToml(t *testing.T) {
	// Create a temporary file with invalid TOML
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.toml")

	invalidToml := `
[login
session_file = "incomplete toml
`

	err := os.WriteFile(configPath, []byte(invalidToml), 0644)
	assert.NoError(t, err)

	// Try to load invalid TOML
	config, err := Load(configPath)
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "failed to decode config file")
}

func TestGetConfigDir(t *testing.T) {
	configDir, err := GetConfigDir()
	assert.NoError(t, err)
	assert.NotEmpty(t, configDir)
	assert.True(t, strings.HasSuffix(configDir, ".aoj-cli"))
}

func TestGetConfigPath(t *testing.T) {
	configPath, err := GetConfigPath()
	assert.NoError(t, err)
	assert.NotEmpty(t, configPath)
	assert.True(t, strings.HasSuffix(configPath, "config.toml"))
}

func TestEnsureConfigDir(t *testing.T) {
	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer func() {
		_ = os.Setenv("HOME", originalHome)
	}()

	// Set temporary home directory
	tmpDir := t.TempDir()
	_ = os.Setenv("HOME", tmpDir)

	err := EnsureConfigDir()
	assert.NoError(t, err)

	// Check if directory was created
	configDir := filepath.Join(tmpDir, ".aoj-cli")
	info, err := os.Stat(configDir)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestGetLanguageConfig(t *testing.T) {
	t.Run("Existing language", func(t *testing.T) {
		lang, exists := GetLanguageConfig("cpp17")
		assert.True(t, exists)
		assert.Equal(t, "cpp", lang.Extension)
		assert.Equal(t, "C++17", lang.AOJLanguageID)
	})

	t.Run("Non-existing language", func(t *testing.T) {
		lang, exists := GetLanguageConfig("nonexistent")
		assert.False(t, exists)
		assert.Empty(t, lang.Extension)
	})
}

func TestValidateConfig(t *testing.T) {
	t.Run("Valid config", func(t *testing.T) {
		config := DefaultConfig()
		err := ValidateConfig(config)
		assert.NoError(t, err)
	})

	t.Run("Nil config", func(t *testing.T) {
		err := ValidateConfig(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config is nil")
	})

	t.Run("Invalid timeout", func(t *testing.T) {
		config := DefaultConfig()
		config.Test.Timeout = -1.0
		err := ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test timeout must be positive")
	})

	t.Run("Empty source file", func(t *testing.T) {
		config := DefaultConfig()
		config.Submit.SourceFile = ""
		err := ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "submit source file cannot be empty")
	})

	t.Run("Empty language", func(t *testing.T) {
		config := DefaultConfig()
		config.Submit.Language = ""
		err := ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "submit language cannot be empty")
	})

	t.Run("Zero timeout", func(t *testing.T) {
		config := DefaultConfig()
		config.Test.Timeout = 0.0
		err := ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test timeout must be positive")
	})
}

func TestDefaultTemplate(t *testing.T) {
	config := DefaultConfig()
	template := config.Init.DefaultTemplate

	assert.NotEmpty(t, template)
	assert.Contains(t, template, "#include <iostream>")
	assert.Contains(t, template, "int main()")
	assert.Contains(t, template, "using namespace std;")
	assert.Contains(t, template, "TODO: Implement solution")
}

func TestLoadDefaultAndSaveDefault(t *testing.T) {
	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer func() {
		_ = os.Setenv("HOME", originalHome)
	}()

	// Set temporary home directory
	tmpDir := t.TempDir()
	_ = os.Setenv("HOME", tmpDir)

	// Load default config (should create default since file doesn't exist)
	config, err := LoadDefault()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Modify the config
	config.Init.Language = "Python3"
	config.Test.Timeout = 10.0

	// Save the modified config
	err = SaveDefault(config)
	assert.NoError(t, err)

	// Load again and verify changes were saved
	loadedConfig, err := LoadDefault()
	assert.NoError(t, err)
	assert.Equal(t, "Python3", loadedConfig.Init.Language)
	assert.Equal(t, 10.0, loadedConfig.Test.Timeout)
}