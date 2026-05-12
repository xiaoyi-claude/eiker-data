package model

// CompanyConflictLog represents one row in the company_conflict_log table.
// It records each name/credit-code conflict encountered during on-boarding
// and the resolution decision taken.
type CompanyConflictLog struct {
	CommonFields
	// TaskID references the associated task UUID.
	TaskID string `json:"task_id"`
	// ConflictType is the conflict category: 1=同名不同企 2=同企不同名.
	ConflictType int16 `json:"conflict_type"`
	// NewCreditCode is the credit code of the incoming (new) record.
	NewCreditCode string `json:"new_credit_code,omitempty"`
	// NewName is the enterprise name of the incoming record.
	NewName string `json:"new_name,omitempty"`
	// ExistingEntCode is the UUID of the already-stored enterprise.
	ExistingEntCode string `json:"existing_ent_code,omitempty"`
	// ExistingCreditCode is the credit code of the already-stored record.
	ExistingCreditCode string `json:"existing_credit_code,omitempty"`
	// ExistingName is the name of the already-stored enterprise.
	ExistingName string `json:"existing_name,omitempty"`
	// Resolution is the handling decision:
	// 1=全部保留独立存在 2=复用UUID挂载新名称 3=拒绝拉新.
	Resolution int16 `json:"resolution,omitempty"`
	// ResolutionReason explains the resolution decision in free text.
	ResolutionReason string `json:"resolution_reason,omitempty"`
}

// SaveConflictLogRequest is the input payload for the "save-conflict-log" Dapr method.
type SaveConflictLogRequest struct {
	// TaskID references the associated task UUID.
	TaskID string `json:"task_id"`
	// ConflictType is the conflict category: 1=同名不同企 2=同企不同名.
	ConflictType int16 `json:"conflict_type"`
	// NewCreditCode is the credit code of the incoming record.
	NewCreditCode string `json:"new_credit_code,omitempty"`
	// NewName is the enterprise name of the incoming record.
	NewName string `json:"new_name,omitempty"`
	// ExistingEntCode is the UUID of the already-stored enterprise.
	ExistingEntCode string `json:"existing_ent_code,omitempty"`
	// ExistingCreditCode is the credit code of the already-stored record.
	ExistingCreditCode string `json:"existing_credit_code,omitempty"`
	// ExistingName is the name of the already-stored enterprise.
	ExistingName string `json:"existing_name,omitempty"`
	// Resolution is the handling decision: 1=全部保留 2=复用UUID 3=拒绝.
	Resolution int16 `json:"resolution,omitempty"`
	// ResolutionReason explains the decision.
	ResolutionReason string `json:"resolution_reason,omitempty"`
	// TraceID is the distributed tracing identifier.
	TraceID string `json:"trace_id,omitempty"`
	// Creater is the operator performing the write.
	Creater string `json:"creater,omitempty"`
}

// SaveConflictLogResponse is the output payload for the "save-conflict-log" Dapr method.
type SaveConflictLogResponse struct {
	// ID is the auto-generated primary key of the new row.
	ID int64 `json:"id"`
}
