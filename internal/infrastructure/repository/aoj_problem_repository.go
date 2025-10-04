// Package repository implements the data access layer.
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/entity"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/repository"
	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
	"github.com/YuminosukeSato/AOJ-cli/pkg/logger"
)

// AOJProblemRepository implements ProblemRepository for AOJ API
type AOJProblemRepository struct {
	baseURL    string
	httpClient *http.Client
	logger     *logger.Logger
}

// NewAOJProblemRepository creates a new AOJProblemRepository
func NewAOJProblemRepository(baseURL string) repository.ProblemRepository {
	return &AOJProblemRepository{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger.WithGroup("aoj_problem_repository"),
	}
}

// TestCaseResponse represents a single test case in the API response
type TestCaseResponse struct {
	Serial int    `json:"serial"`
	In     string `json:"in"`
	Out    string `json:"out"`
}

// GetByID retrieves a problem by its ID
func (r *AOJProblemRepository) GetByID(_ context.Context, _ model.ProblemID) (*entity.Problem, error) {
	return nil, cerrors.New("GetByID not implemented")
}

// GetByIDs retrieves multiple problems by their IDs
func (r *AOJProblemRepository) GetByIDs(_ context.Context, _ []model.ProblemID) ([]*entity.Problem, error) {
	return nil, cerrors.New("GetByIDs not implemented")
}

// Search searches for problems by criteria
func (r *AOJProblemRepository) Search(_ context.Context, _ repository.ProblemSearchCriteria) ([]*entity.Problem, error) {
	return nil, cerrors.New("Search not implemented")
}

// Save saves a problem
func (r *AOJProblemRepository) Save(_ context.Context, _ *entity.Problem) error {
	return cerrors.New("Save not implemented")
}

// Delete deletes a problem by its ID
func (r *AOJProblemRepository) Delete(_ context.Context, _ model.ProblemID) error {
	return cerrors.New("Delete not implemented")
}

// Exists checks if a problem exists
func (r *AOJProblemRepository) Exists(_ context.Context, _ model.ProblemID) (bool, error) {
	return false, cerrors.New("Exists not implemented")
}

// GetTestCases retrieves test cases for a problem from AOJ API
// AOJ API requires fetching test cases one by one by serial number
func (r *AOJProblemRepository) GetTestCases(ctx context.Context, problemID model.ProblemID) ([]model.TestCase, error) {
	r.logger.InfoContext(ctx, "fetching test cases from AOJ", "problem_id", problemID.String())

	testCases := make([]model.TestCase, 0)

	// Fetch test cases sequentially until we get a 404
	// Most problems have 1-20 test cases
	const maxTestCases = 100
	for serial := 1; serial <= maxTestCases; serial++ {
		testCase, found, err := r.fetchSingleTestCase(ctx, problemID, serial)
		if err != nil {
			return nil, err
		}
		if !found {
			// No more test cases available
			break
		}
		testCases = append(testCases, *testCase)
	}

	r.logger.InfoContext(ctx, "successfully fetched test cases", "count", len(testCases))
	return testCases, nil
}

// fetchSingleTestCase fetches a single test case by serial number
// Returns (testCase, found, error)
func (r *AOJProblemRepository) fetchSingleTestCase(ctx context.Context, problemID model.ProblemID, serial int) (*model.TestCase, bool, error) {
	// AOJ test cases are available at https://judgedat.u-aizu.ac.jp/testcases/{problemId}/{serial}
	url := fmt.Sprintf("https://judgedat.u-aizu.ac.jp/testcases/%s/%d", problemID.String(), serial)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, false, cerrors.Wrap(err, "failed to create HTTP request")
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.logger.ErrorContext(ctx, "HTTP request failed", "error", err, "serial", serial)
		return nil, false, cerrors.NewAppError(
			cerrors.CodeNetworkError,
			"failed to connect to AOJ",
			err,
		)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			r.logger.WarnContext(ctx, "failed to close response body", "error", err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		var apiTC TestCaseResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiTC); err != nil {
			return nil, false, cerrors.Wrap(err, "failed to decode test case response")
		}
		tc := model.NewTestCase(apiTC.Serial, apiTC.In, apiTC.Out)
		return tc, true, nil
	case http.StatusNotFound:
		// No more test cases available
		return nil, false, nil
	case http.StatusBadRequest:
		return nil, false, cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"invalid problem ID format",
			nil,
		)
	case http.StatusInternalServerError:
		return nil, false, cerrors.NewAppError(
			cerrors.CodeServiceUnavailable,
			"AOJ server error",
			nil,
		)
	default:
		r.logger.ErrorContext(ctx, "unexpected response status", "status", resp.StatusCode, "serial", serial)
		return nil, false, cerrors.NewAppError(
			cerrors.CodeInternalServer,
			"unexpected response from AOJ",
			cerrors.WithDetail(nil, "status_code: "+resp.Status),
		)
	}
}

// SaveTestCases saves test cases for a problem
func (r *AOJProblemRepository) SaveTestCases(_ context.Context, _ model.ProblemID, _ []model.TestCase) error {
	return cerrors.New("SaveTestCases not implemented")
}
