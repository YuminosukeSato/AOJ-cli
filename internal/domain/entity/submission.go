package entity

import (
	"time"

	"github.com/YuminosukeSato/AOJ-cli/internal/domain/model"
)

// SubmissionStatus represents the status of a submission
type SubmissionStatus string

// Submission status constants
const (
	StatusPending       SubmissionStatus = "PENDING"
	StatusJudging       SubmissionStatus = "JUDGING"
	StatusAccepted      SubmissionStatus = "ACCEPTED"
	StatusWrongAnswer   SubmissionStatus = "WRONG_ANSWER"
	StatusTimeLimitExceeded SubmissionStatus = "TIME_LIMIT_EXCEEDED"
	StatusMemoryLimitExceeded SubmissionStatus = "MEMORY_LIMIT_EXCEEDED"
	StatusRuntimeError  SubmissionStatus = "RUNTIME_ERROR"
	StatusCompileError  SubmissionStatus = "COMPILE_ERROR"
	StatusPresentationError SubmissionStatus = "PRESENTATION_ERROR"
	StatusOutputLimitExceeded SubmissionStatus = "OUTPUT_LIMIT_EXCEEDED"
	StatusInternalError SubmissionStatus = "INTERNAL_ERROR"
)

// IsSuccess returns true if the status indicates success
func (s SubmissionStatus) IsSuccess() bool {
	return s == StatusAccepted
}

// IsError returns true if the status indicates an error
func (s SubmissionStatus) IsError() bool {
	return s != StatusPending && s != StatusJudging && s != StatusAccepted
}

// IsFinal returns true if the status is final (not pending or judging)
func (s SubmissionStatus) IsFinal() bool {
	return s != StatusPending && s != StatusJudging
}

// Submission represents a code submission to AOJ
type Submission struct {
	id         model.SubmissionID
	problemID  model.ProblemID
	language   string
	sourceCode string
	status     SubmissionStatus
	score      int
	time       time.Duration
	memory     int64 // in KB
	message    string
	submittedAt time.Time
	judgedAt   *time.Time
}

// NewSubmission creates a new Submission instance
func NewSubmission(
	id model.SubmissionID,
	problemID model.ProblemID,
	language, sourceCode string,
) *Submission {
	return &Submission{
		id:         id,
		problemID:  problemID,
		language:   language,
		sourceCode: sourceCode,
		status:     StatusPending,
		score:      0,
		time:       0,
		memory:     0,
		message:    "",
		submittedAt: time.Now(),
		judgedAt:   nil,
	}
}

// ID returns the submission ID
func (s *Submission) ID() model.SubmissionID {
	return s.id
}

// ProblemID returns the problem ID
func (s *Submission) ProblemID() model.ProblemID {
	return s.problemID
}

// Language returns the programming language
func (s *Submission) Language() string {
	return s.language
}

// SourceCode returns the source code
func (s *Submission) SourceCode() string {
	return s.sourceCode
}

// Status returns the submission status
func (s *Submission) Status() SubmissionStatus {
	return s.status
}

// Score returns the score
func (s *Submission) Score() int {
	return s.score
}

// Time returns the execution time
func (s *Submission) Time() time.Duration {
	return s.time
}

// Memory returns the memory usage in KB
func (s *Submission) Memory() int64 {
	return s.memory
}

// Message returns the judge message
func (s *Submission) Message() string {
	return s.message
}

// SubmittedAt returns the submission time
func (s *Submission) SubmittedAt() time.Time {
	return s.submittedAt
}

// JudgedAt returns the judge time (nil if not judged yet)
func (s *Submission) JudgedAt() *time.Time {
	if s.judgedAt == nil {
		return nil
	}
	judgedTime := *s.judgedAt
	return &judgedTime
}

// UpdateStatus updates the submission status
func (s *Submission) UpdateStatus(status SubmissionStatus) {
	s.status = status
	if status.IsFinal() && s.judgedAt == nil {
		now := time.Now()
		s.judgedAt = &now
	}
}

// UpdateResult updates the submission result
func (s *Submission) UpdateResult(
	status SubmissionStatus,
	score int,
	execTime time.Duration,
	memory int64,
	message string,
) {
	s.status = status
	s.score = score
	s.time = execTime
	s.memory = memory
	s.message = message
	
	if status.IsFinal() && s.judgedAt == nil {
		now := time.Now()
		s.judgedAt = &now
	}
}

// IsJudged returns true if the submission has been judged
func (s *Submission) IsJudged() bool {
	return s.judgedAt != nil
}

// IsAccepted returns true if the submission was accepted
func (s *Submission) IsAccepted() bool {
	return s.status == StatusAccepted
}

// HasError returns true if the submission has an error
func (s *Submission) HasError() bool {
	return s.status.IsError()
}

// IsPending returns true if the submission is pending
func (s *Submission) IsPending() bool {
	return s.status == StatusPending || s.status == StatusJudging
}

// GetJudgeDuration returns the duration from submission to judgment
func (s *Submission) GetJudgeDuration() *time.Duration {
	if s.judgedAt == nil {
		return nil
	}
	duration := s.judgedAt.Sub(s.submittedAt)
	return &duration
}

// IsValid validates the submission data
func (s *Submission) IsValid() bool {
	return s.id.IsValid() &&
		s.problemID.IsValid() &&
		s.language != "" &&
		s.sourceCode != ""
}

// Clone creates a copy of the submission
func (s *Submission) Clone() *Submission {
	clone := &Submission{
		id:         s.id,
		problemID:  s.problemID,
		language:   s.language,
		sourceCode: s.sourceCode,
		status:     s.status,
		score:      s.score,
		time:       s.time,
		memory:     s.memory,
		message:    s.message,
		submittedAt: s.submittedAt,
		judgedAt:   nil,
	}
	
	if s.judgedAt != nil {
		judgedTime := *s.judgedAt
		clone.judgedAt = &judgedTime
	}
	
	return clone
}