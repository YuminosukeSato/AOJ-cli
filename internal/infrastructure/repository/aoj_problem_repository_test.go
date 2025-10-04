package repository

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
	"github.com/YuminosukeSato/AOJ-cli/pkg/cerrors"
)

func TestAOJProblemRepository_GetTestCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		problemID      string
		serverResponse string
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name:      "successful response with test cases",
			problemID: "ITP1_1_A",
			serverResponse: `[
				{"serial": 1, "in": "input1", "out": "output1"},
				{"serial": 2, "in": "input2", "out": "output2"}
			]`,
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name:           "not found - returns empty test cases",
			problemID:      "ITP1_9_Z",
			serverResponse: "",
			serverStatus:   http.StatusNotFound,
			wantErr:        false,
			wantCount:      0,
		},
		{
			name:           "bad request",
			problemID:      "ITP1_8_C",
			serverResponse: "",
			serverStatus:   http.StatusBadRequest,
			wantErr:        true,
			wantCount:      0,
		},
		{
			name:           "server error",
			problemID:      "ITP1_1_A",
			serverResponse: "",
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
			wantCount:      0,
		},
		{
			name:           "invalid JSON response",
			problemID:      "ITP1_1_A",
			serverResponse: `invalid json`,
			serverStatus:   http.StatusOK,
			wantErr:        true,
			wantCount:      0,
		},
		{
			name:           "empty array",
			problemID:      "ITP1_1_A",
			serverResponse: `[]`,
			serverStatus:   http.StatusOK,
			wantErr:        false,
			wantCount:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				expectedPath := "/testcases/samples/" + tt.problemID
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
				if tt.serverResponse != "" {
					_, _ = w.Write([]byte(tt.serverResponse))
				}
			}))
			defer server.Close()

			// Create repository with mock server URL
			repo := NewAOJProblemRepository(server.URL)

			// Create problem ID
			pid, err := model.NewProblemID(tt.problemID)
			if err != nil {
				t.Fatalf("failed to create problem ID: %v", err)
			}

			// Execute test
			ctx := context.Background()
			testCases, err := repo.GetTestCases(ctx, pid)

			// Verify error
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTestCases() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify test case count
			if len(testCases) != tt.wantCount {
				t.Errorf("GetTestCases() got %d test cases, want %d", len(testCases), tt.wantCount)
			}

			// Verify test case content for successful case
			if !tt.wantErr && tt.wantCount > 0 {
				if testCases[0].Input() != "input1" {
					t.Errorf("first test case input = %v, want %v", testCases[0].Input(), "input1")
				}
				if testCases[0].Expected() != "output1" {
					t.Errorf("first test case expected = %v, want %v", testCases[0].Expected(), "output1")
				}
			}
		})
	}
}

func TestAOJProblemRepository_GetTestCases_NetworkError(t *testing.T) {
	t.Parallel()

	// Create repository with invalid URL to simulate network error
	repo := NewAOJProblemRepository("http://invalid-url-that-does-not-exist.local")

	pid, err := model.NewProblemID("ITP1_1_A")
	if err != nil {
		t.Fatalf("failed to create problem ID: %v", err)
	}

	ctx := context.Background()
	_, err = repo.GetTestCases(ctx, pid)

	if err == nil {
		t.Error("expected error for network failure, got nil")
	}

	// Verify it's a network error
	var appErr *cerrors.AppError
	if cerrors.As(err, &appErr) {
		if appErr.Code != cerrors.CodeNetworkError {
			t.Errorf("expected network error code, got %v", appErr.Code)
		}
	}
}

func TestAOJProblemRepository_NotImplementedMethods(t *testing.T) {
	t.Parallel()

	repo := NewAOJProblemRepository("http://example.com")
	ctx := context.Background()

	pid, _ := model.NewProblemID("TEST")

	t.Run("GetByID", func(t *testing.T) {
		_, err := repo.GetByID(ctx, pid)
		if err == nil {
			t.Error("expected error for GetByID, got nil")
		}
	})

	t.Run("Exists", func(t *testing.T) {
		_, err := repo.Exists(ctx, pid)
		if err == nil {
			t.Error("expected error for Exists, got nil")
		}
	})
}
