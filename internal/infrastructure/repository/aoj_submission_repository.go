// Package repository implements the data access layer.
package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/entity"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
	"github.com/YuminosukeSato/AOJ-cli/internal/domain/repository"
	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
	"github.com/YuminosukeSato/AOJ-cli/pkg/logger"
)

// AOJSubmissionRepository implements SubmissionRepository for AOJ API
type AOJSubmissionRepository struct {
	baseURL    string
	httpClient *http.Client
	logger     *logger.Logger
}

// NewAOJSubmissionRepository creates a new AOJSubmissionRepository
func NewAOJSubmissionRepository(baseURL string) repository.SubmissionRepository {
	return &AOJSubmissionRepository{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger.WithGroup("aoj_submission_repository"),
	}
}

// SubmitRequest represents the JSON payload for submission
type SubmitRequest struct {
	ProblemID  string `json:"problemId"`
	Language   string `json:"language"`
	SourceCode string `json:"sourceCode"`
}

// SubmitResponse represents the JSON response from submission
type SubmitResponse struct {
	SubmissionID    string `json:"submissionId"`
	ProblemID       string `json:"problemId"`
	Status          string `json:"status"`
	SubmittedAt     int64  `json:"submittedAt"`
	JudgeType       string `json:"judgeType"`
	Score           int    `json:"score"`
	ExecutionTime   int    `json:"cpuTime"`
	ExecutionMemory int    `json:"memory"`
	Message         string `json:"message"`
}

// Submit submits a solution to AOJ
func (r *AOJSubmissionRepository) Submit(ctx context.Context, submission *entity.Submission) error {
	r.logger.InfoContext(ctx, "submitting solution to AOJ",
		"problem_id", submission.ProblemID().String(),
		"language", submission.Language())

	// Prepare request payload
	submitReq := SubmitRequest{
		ProblemID:  submission.ProblemID().String(),
		Language:   r.normalizeLanguage(submission.Language()),
		SourceCode: submission.SourceCode(),
	}

	payload, err := json.Marshal(submitReq)
	if err != nil {
		return cerrors.Wrap(err, "failed to marshal submission request")
	}

	// Create HTTP request
	// Note: The exact endpoint needs to be verified with AOJ API documentation
	url := r.baseURL + "/submissions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return cerrors.Wrap(err, "failed to create HTTP request")
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	// Execute request
	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.logger.ErrorContext(ctx, "HTTP request failed", "error", err)
		return cerrors.NewAppError(
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

	// Handle different response status codes
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return r.parseSubmitResponse(ctx, resp, submission)
	case http.StatusUnauthorized:
		return cerrors.NewAppError(
			cerrors.CodeUnauthorized,
			"authentication required. Please login first",
			nil,
		)
	case http.StatusBadRequest:
		return cerrors.NewAppError(
			cerrors.CodeInvalidInput,
			"invalid submission request",
			nil,
		)
	case http.StatusInternalServerError:
		return cerrors.NewAppError(
			cerrors.CodeServiceUnavailable,
			"AOJ server error",
			nil,
		)
	default:
		r.logger.ErrorContext(ctx, "unexpected response status", "status", resp.StatusCode)
		return cerrors.NewAppError(
			cerrors.CodeInternalServer,
			"unexpected response from AOJ",
			cerrors.WithDetail(nil, "status_code: "+resp.Status),
		)
	}
}

// parseSubmitResponse parses the submission response
func (r *AOJSubmissionRepository) parseSubmitResponse(_ context.Context, resp *http.Response, submission *entity.Submission) error {
	var submitResp SubmitResponse
	if err := json.NewDecoder(resp.Body).Decode(&submitResp); err != nil {
		return cerrors.Wrap(err, "failed to decode submission response")
	}

	// Update submission with response data
	status := r.mapSubmissionStatus(submitResp.Status)
	submission.UpdateResult(
		status,
		submitResp.Score,
		time.Duration(submitResp.ExecutionTime)*time.Millisecond,
		int64(submitResp.ExecutionMemory),
		submitResp.Message,
	)

	return nil
}

// normalizeLanguage normalizes language names for AOJ API
func (r *AOJSubmissionRepository) normalizeLanguage(lang string) string {
	// Map common language names to AOJ's expected format
	languageMap := map[string]string{
		"C":          "C",
		"C++":        "C++14",
		"C++14":      "C++14",
		"C++17":      "C++17",
		"Java":       "JAVA",
		"JAVA":       "JAVA",
		"Python":     "Python3",
		"Python3":    "Python3",
		"Ruby":       "Ruby",
		"Go":         "Go",
		"JavaScript": "JavaScript",
		"C#":         "C#",
		"PHP":        "PHP",
		"D":          "D",
		"Rust":       "Rust",
		"Kotlin":     "Kotlin",
		"Scala":      "Scala",
	}

	if normalized, ok := languageMap[lang]; ok {
		return normalized
	}

	return lang
}

// mapSubmissionStatus maps AOJ status to our domain status
func (r *AOJSubmissionRepository) mapSubmissionStatus(aojStatus string) entity.SubmissionStatus {
	statusMap := map[string]entity.SubmissionStatus{
		"PENDING":               entity.StatusPending,
		"JUDGING":               entity.StatusJudging,
		"ACCEPTED":              entity.StatusAccepted,
		"WRONG_ANSWER":          entity.StatusWrongAnswer,
		"TIME_LIMIT_EXCEEDED":   entity.StatusTimeLimitExceeded,
		"MEMORY_LIMIT_EXCEEDED": entity.StatusMemoryLimitExceeded,
		"RUNTIME_ERROR":         entity.StatusRuntimeError,
		"COMPILE_ERROR":         entity.StatusCompileError,
		"PRESENTATION_ERROR":    entity.StatusPresentationError,
		"OUTPUT_LIMIT_EXCEEDED": entity.StatusOutputLimitExceeded,
		"INTERNAL_ERROR":        entity.StatusInternalError,
	}

	if status, ok := statusMap[aojStatus]; ok {
		return status
	}

	return entity.StatusPending
}

// Not implemented methods - return errors

func (r *AOJSubmissionRepository) GetByID(_ context.Context, _ model.SubmissionID) (*entity.Submission, error) {
	return nil, cerrors.New("GetByID not implemented")
}

func (r *AOJSubmissionRepository) GetByProblemID(_ context.Context, _ model.ProblemID, _ int) ([]*entity.Submission, error) {
	return nil, cerrors.New("GetByProblemID not implemented")
}

func (r *AOJSubmissionRepository) GetRecent(_ context.Context, _ int) ([]*entity.Submission, error) {
	return nil, cerrors.New("GetRecent not implemented")
}

func (r *AOJSubmissionRepository) GetStatus(_ context.Context, _ model.SubmissionID) (entity.SubmissionStatus, error) {
	return "", cerrors.New("GetStatus not implemented")
}

func (r *AOJSubmissionRepository) WatchStatus(_ context.Context, _ model.SubmissionID, _ time.Duration) (<-chan entity.SubmissionStatus, error) {
	return nil, cerrors.New("WatchStatus not implemented")
}

func (r *AOJSubmissionRepository) Search(_ context.Context, _ repository.SubmissionSearchCriteria) ([]*entity.Submission, error) {
	return nil, cerrors.New("Search not implemented")
}

func (r *AOJSubmissionRepository) Save(_ context.Context, _ *entity.Submission) error {
	return cerrors.New("Save not implemented")
}

func (r *AOJSubmissionRepository) Delete(_ context.Context, _ model.SubmissionID) error {
	return cerrors.New("Delete not implemented")
}

func (r *AOJSubmissionRepository) Exists(_ context.Context, _ model.SubmissionID) (bool, error) {
	return false, cerrors.New("Exists not implemented")
}
