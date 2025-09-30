package cerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New("test error")
	assert.EqualError(t, err, "test error")
}

func TestErrorf(t *testing.T) {
	err := Errorf("test error: %s", "details")
	assert.EqualError(t, err, "test error: details")
}

func TestWrap(t *testing.T) {
	originalErr := New("original error")
	wrappedErr := Wrap(originalErr, "wrapped")

	assert.Contains(t, wrappedErr.Error(), "wrapped")
	assert.Contains(t, wrappedErr.Error(), "original error")
	assert.True(t, Is(wrappedErr, originalErr))
}

func TestWrapf(t *testing.T) {
	originalErr := New("original error")
	wrappedErr := Wrapf(originalErr, "wrapped: %s", "details")

	assert.Contains(t, wrappedErr.Error(), "wrapped: details")
	assert.Contains(t, wrappedErr.Error(), "original error")
	assert.True(t, Is(wrappedErr, originalErr))
}

func TestIs(t *testing.T) {
	baseErr := New("base error")
	wrappedErr := Wrap(baseErr, "wrapped")

	assert.True(t, Is(wrappedErr, baseErr))
	assert.False(t, Is(wrappedErr, New("different error")))
}

func TestAs(t *testing.T) {
	appErr := &AppError{
		Code:    CodeNotFound,
		Message: "resource not found",
		Err:     New("underlying error"),
	}

	wrappedErr := Wrap(appErr, "wrapped app error")

	var target *AppError
	assert.True(t, As(wrappedErr, &target))
	assert.Equal(t, CodeNotFound, target.Code)
	assert.Equal(t, "resource not found", target.Message)
}

func TestAppError(t *testing.T) {
	t.Run("AppError without underlying error", func(t *testing.T) {
		appErr := NewAppError(CodeInvalidInput, "invalid data", nil)

		assert.Equal(t, CodeInvalidInput, appErr.Code)
		assert.Equal(t, "invalid data", appErr.Message)
		assert.Equal(t, "invalid data", appErr.Error())
		assert.Nil(t, appErr.Unwrap())
	})

	t.Run("AppError with underlying error", func(t *testing.T) {
		underlyingErr := New("underlying problem")
		appErr := NewAppError(CodeInternalServer, "server error", underlyingErr)

		assert.Equal(t, CodeInternalServer, appErr.Code)
		assert.Equal(t, "server error", appErr.Message)
		assert.Equal(t, "server error: underlying problem", appErr.Error())
		assert.Equal(t, underlyingErr, appErr.Unwrap())
	})
}

func TestIsAppError(t *testing.T) {
	appErr := NewAppError(CodeNotFound, "not found", nil)
	wrappedErr := Wrap(appErr, "wrapped")

	assert.True(t, IsAppError(appErr, CodeNotFound))
	assert.True(t, IsAppError(wrappedErr, CodeNotFound))
	assert.False(t, IsAppError(appErr, CodeInvalidInput))
	assert.False(t, IsAppError(New("regular error"), CodeNotFound))
}

func TestGetErrorCode(t *testing.T) {
	t.Run("AppError", func(t *testing.T) {
		appErr := NewAppError(CodeTimeout, "timeout occurred", nil)
		code := GetErrorCode(appErr)
		assert.Equal(t, CodeTimeout, code)
	})

	t.Run("Wrapped AppError", func(t *testing.T) {
		appErr := NewAppError(CodeUnauthorized, "auth failed", nil)
		wrappedErr := Wrap(appErr, "wrapped")
		code := GetErrorCode(wrappedErr)
		assert.Equal(t, CodeUnauthorized, code)
	})

	t.Run("Regular error", func(t *testing.T) {
		regularErr := New("regular error")
		code := GetErrorCode(regularErr)
		assert.Empty(t, code)
	})
}

func TestCommonErrors(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{"ErrNotFound", ErrNotFound},
		{"ErrInvalidInput", ErrInvalidInput},
		{"ErrUnauthorized", ErrUnauthorized},
		{"ErrForbidden", ErrForbidden},
		{"ErrConflict", ErrConflict},
		{"ErrInternalServer", ErrInternalServer},
		{"ErrServiceUnavailable", ErrServiceUnavailable},
		{"ErrTimeout", ErrTimeout},
		{"ErrNetworkError", ErrNetworkError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Error(t, tc.err)
			assert.NotEmpty(t, tc.err.Error())
		})
	}
}

func TestWithDetail(t *testing.T) {
	baseErr := New("base error")
	detailedErr := WithDetail(baseErr, "request_id: req-123")

	// Check that the error still behaves as expected
	assert.True(t, Is(detailedErr, baseErr))
	assert.Error(t, detailedErr)
}

func TestWithHint(t *testing.T) {
	baseErr := New("authentication failed")
	hintedErr := WithHint(baseErr, "please check your credentials")

	// Check that the error still behaves as expected
	assert.True(t, Is(hintedErr, baseErr))
	assert.Error(t, hintedErr)
}