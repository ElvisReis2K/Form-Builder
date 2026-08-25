package responses

import "time"

type Response struct {
	ID                    string
	FormID                string
	Answers               map[string]any
	PrivacyAcknowledgedAt time.Time
	SubmittedAt           time.Time
}

type SubmissionInput struct {
	Answers             map[string]any
	PrivacyAcknowledged bool
}
