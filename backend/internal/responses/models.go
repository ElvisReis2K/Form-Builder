package responses

import "time"

type Response struct {
	ID          string
	FormID      string
	Answers     map[string]any
	SubmittedAt time.Time
}

type SubmissionInput struct {
	Answers map[string]any
}
