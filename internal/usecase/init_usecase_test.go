package usecase_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/entity"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/repository"
	"github.com/YuminosukeSato/AOJ-cli/internal/usecase"
)

// MockProblemRepository is a mock implementation of ProblemRepository
type MockProblemRepository struct {
	testCases []model.TestCase
	getError  error
	saveError error
}

func (m *MockProblemRepository) GetByID(_ context.Context, _ model.ProblemID) (*entity.Problem, error) {
	return nil, nil
}

func (m *MockProblemRepository) GetByIDs(_ context.Context, _ []model.ProblemID) ([]*entity.Problem, error) {
	return nil, nil
}

func (m *MockProblemRepository) Search(_ context.Context, _ repository.ProblemSearchCriteria) ([]*entity.Problem, error) {
	return nil, nil
}

func (m *MockProblemRepository) Save(_ context.Context, _ *entity.Problem) error {
	return m.saveError
}

func (m *MockProblemRepository) Delete(_ context.Context, _ model.ProblemID) error {
	return nil
}

func (m *MockProblemRepository) Exists(_ context.Context, _ model.ProblemID) (bool, error) {
	return true, nil
}

func (m *MockProblemRepository) GetTestCases(_ context.Context, _ model.ProblemID) ([]model.TestCase, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return m.testCases, nil
}

func (m *MockProblemRepository) SaveTestCases(_ context.Context, _ model.ProblemID, _ []model.TestCase) error {
	return m.saveError
}

func TestInitUseCase_Execute_EmptyProblemID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mockRepo := &MockProblemRepository{}
	uc := usecase.NewInitUseCase(mockRepo)

	err := uc.Execute(ctx, "")
	if err == nil {
		t.Error("expected error for empty problem ID, got nil")
	}
}

func TestInitUseCase_Execute_Success(t *testing.T) {
	t.Parallel()

	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	ctx := context.Background()
	mockRepo := &MockProblemRepository{
		testCases: []model.TestCase{
			*model.NewTestCase(1, "5\n", "5\n"),
		},
	}
	uc := usecase.NewInitUseCase(mockRepo)

	problemID := "ALDS1_1_A"
	err := uc.Execute(ctx, problemID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// ディレクトリが作成されたか確認
	if _, err := os.Stat(problemID); os.IsNotExist(err) {
		t.Errorf("problem directory was not created")
	}

	// main.goが作成されたか確認
	mainFile := filepath.Join(problemID, "main.go")
	if _, err := os.Stat(mainFile); os.IsNotExist(err) {
		t.Errorf("main.go was not created")
	}

	// testディレクトリが作成されたか確認
	testDir := filepath.Join(problemID, "test")
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Errorf("test directory was not created")
	}
}
