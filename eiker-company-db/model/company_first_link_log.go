package model

import "time"

// CompanyFirstLinkLog represents one row in the company_first_link_log table.
// It captures the before/after snapshot for a first-link verification step and
// records the final verification decision.
type CompanyFirstLinkLog struct {
	CommonFields
	// TaskID references the associated task UUID.
	TaskID string `json:"task_id"`
	// FactID references company_fact.id.
	FactID int64 `json:"fact_id"`
	// VerifierType is the verifier category:
	// 1=系统自验（基础/公司拉新） 2=用户确认 3=业务方确认.
	VerifierType int16 `json:"verifier_type"`
	// ExpectedData is the snapshot taken before writing, used for comparison.
	ExpectedData map[string]interface{} `json:"expected_data,omitempty"`
	// ActualData is the data read back after the write.
	ActualData map[string]interface{} `json:"actual_data,omitempty"`
	// IsConsistent indicates whether expected and actual data match.
	IsConsistent *bool `json:"is_consistent,omitempty"`
	// VerifyStatus is the verification outcome: 0=待确认 1=已通过 2=已拒绝.
	VerifyStatus int16 `json:"verify_status"`
	// VerifiedAt is the timestamp when the verification decision was made.
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}

// SaveFirstLinkLogRequest is the input payload for the "save-first-link-log" Dapr method.
type SaveFirstLinkLogRequest struct {
	// TaskID references the associated task UUID.
	TaskID string `json:"task_id"`
	// FactID references company_fact.id.
	FactID int64 `json:"fact_id"`
	// VerifierType is the verifier category: 1=系统自验 2=用户确认 3=业务方确认.
	VerifierType int16 `json:"verifier_type"`
	// ExpectedData is the snapshot before writing.
	ExpectedData map[string]interface{} `json:"expected_data,omitempty"`
	// ActualData is the data read back after writing.
	ActualData map[string]interface{} `json:"actual_data,omitempty"`
	// IsConsistent indicates whether expected and actual match.
	IsConsistent *bool `json:"is_consistent,omitempty"`
	// VerifyStatus is the verification outcome: 0=待确认 1=已通过 2=已拒绝.
	VerifyStatus int16 `json:"verify_status"`
	// VerifiedAt is the time the decision was made.
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	// TraceID is the distributed tracing identifier.
	TraceID string `json:"trace_id,omitempty"`
	// Creater is the operator performing the write.
	Creater string `json:"creater,omitempty"`
}

// SaveFirstLinkLogResponse is the output payload for the "save-first-link-log" Dapr method.
type SaveFirstLinkLogResponse struct {
	// ID is the auto-generated primary key of the new row.
	ID int64 `json:"id"`
}
