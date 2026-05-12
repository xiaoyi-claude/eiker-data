package model

// CompanyDataVerifyLog represents one row in the company_data_verify_log table.
// It records the query parameters and response from each third-party data-source call,
// together with the consistency check result.
type CompanyDataVerifyLog struct {
	CommonFields
	// TaskID references the associated task UUID.
	TaskID string `json:"task_id"`
	// SourceName is the data-source name (e.g. "tianyancha").
	SourceName string `json:"source_name"`
	// QueryName is the enterprise name used as the query parameter.
	QueryName string `json:"query_name,omitempty"`
	// ResponseData is the raw JSON response from the data source.
	ResponseData map[string]interface{} `json:"response_data,omitempty"`
	// IsConsistent indicates whether the response matches the input triple (credit code, name, legal rep).
	IsConsistent *bool `json:"is_consistent,omitempty"`
	// IsSuccess indicates whether the API call itself succeeded.
	IsSuccess bool `json:"is_success"`
	// FailReason describes the failure cause when IsSuccess is false.
	FailReason string `json:"fail_reason,omitempty"`
}

// SaveDataVerifyLogRequest is the input payload for the "save-data-verify-log" Dapr method.
type SaveDataVerifyLogRequest struct {
	// TaskID references the associated task UUID.
	TaskID string `json:"task_id"`
	// SourceName is the data-source name.
	SourceName string `json:"source_name"`
	// QueryName is the enterprise name used in the query.
	QueryName string `json:"query_name,omitempty"`
	// ResponseData is the raw JSON from the data source.
	ResponseData map[string]interface{} `json:"response_data,omitempty"`
	// IsConsistent indicates whether the response matches the input triple.
	IsConsistent *bool `json:"is_consistent,omitempty"`
	// IsSuccess indicates whether the API call succeeded.
	IsSuccess bool `json:"is_success"`
	// FailReason is the failure description when IsSuccess is false.
	FailReason string `json:"fail_reason,omitempty"`
	// TraceID is the distributed tracing identifier.
	TraceID string `json:"trace_id,omitempty"`
	// Creater is the operator performing the write.
	Creater string `json:"creater,omitempty"`
}

// SaveDataVerifyLogResponse is the output payload for the "save-data-verify-log" Dapr method.
type SaveDataVerifyLogResponse struct {
	// ID is the auto-generated primary key of the new row.
	ID int64 `json:"id"`
}
