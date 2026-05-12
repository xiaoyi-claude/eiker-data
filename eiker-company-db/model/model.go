// Package model defines the entity types and all request/response types for
// the eiker-company-db Dapr service invocation interface.
package model

import "time"

// ─────────────────────────────────────────────
// 公共字段
// ─────────────────────────────────────────────

// CommonFields contains the audit and tracking fields that every data entity shares.
type CommonFields struct {
	// ID is the auto-increment primary key.
	ID int64 `json:"id"`
	// IsCurrent indicates whether this record is the current valid version.
	IsCurrent bool `json:"is_current"`
	// TraceID is the distributed-tracing identifier for log correlation.
	TraceID string `json:"trace_id,omitempty"`
	// Owner is the data-owner identifier used for permission control.
	Owner string `json:"owner,omitempty"`
	// Creater is the account that created this record.
	Creater string `json:"creater,omitempty"`
	// CreateTime is the timestamp when the record was inserted.
	CreateTime time.Time `json:"create_time"`
	// Updater is the account that last updated this record.
	Updater string `json:"updater,omitempty"`
	// UpdateTime is the timestamp of the last update.
	UpdateTime time.Time `json:"update_time"`
	// Remark is a free-text annotation.
	Remark string `json:"remark,omitempty"`
}

// ─────────────────────────────────────────────
// 实体类型
// ─────────────────────────────────────────────

// CompanyFact represents one row in the company_fact table.
type CompanyFact struct {
	CommonFields
	// EntCode is the enterprise UUID v4 used as the external identifier.
	EntCode string `json:"ent_code"`
	// CreditCode is the unified social credit code (GB 32100-2015).
	CreditCode string `json:"credit_code"`
	// Name is the enterprise name (4–26 Chinese characters, optionally with brackets).
	Name string `json:"name"`
	// LegalRep is the power-holder full name (2–5 Chinese characters).
	LegalRep string `json:"legal_rep,omitempty"`
	// LocationID references the address_id in eiker-address-db.
	LocationID string `json:"location_id,omitempty"`
	// FactTime is the timestamp when the fact occurred.
	FactTime *time.Time `json:"fact_time,omitempty"`
	// FactRecordTime is the timestamp when the record entered the system.
	FactRecordTime *time.Time `json:"fact_record_time,omitempty"`
	// EntRecordScene is the enterprise on-boarding source:
	// 1=基础拉新 2=用户拉新 3=业务拉新 4=公司拉新.
	EntRecordScene int16 `json:"ent_record_scene"`
	// LegalRepRecordScene is the power-holder on-boarding type:
	// 1=代为确定性主张 2=确定实际主张 3=权力人替换并主张.
	LegalRepRecordScene int16 `json:"legal_rep_record_scene,omitempty"`
	// TaskID is the UUID of the associated task record.
	TaskID string `json:"task_id,omitempty"`
}

// CompanyOCRLog represents one row in the company_ocr_log table.
type CompanyOCRLog struct {
	CommonFields
	// TaskID references the associated task UUID.
	TaskID string `json:"task_id"`
	// ModelName is the OCR model used (paddleocr / easyocr / tesseract).
	ModelName string `json:"model_name"`
	// ImagePath is the path or URL of the source image.
	ImagePath string `json:"image_path,omitempty"`
	// RawResult is the raw JSON output from the OCR model.
	RawResult map[string]interface{} `json:"raw_result,omitempty"`
	// CreditCode is the credit code extracted by OCR.
	CreditCode string `json:"credit_code,omitempty"`
	// CompanyName is the enterprise name extracted by OCR.
	CompanyName string `json:"company_name,omitempty"`
	// LegalRepName is the power-holder name extracted by OCR.
	LegalRepName string `json:"legal_rep_name,omitempty"`
	// IsSuccess indicates whether the extraction succeeded.
	IsSuccess bool `json:"is_success"`
	// FailReason describes the extraction failure cause when IsSuccess is false.
	FailReason string `json:"fail_reason,omitempty"`
}

// CompanyDataVerifyLog represents one row in the company_data_verify_log table.
type CompanyDataVerifyLog struct {
	CommonFields
	// TaskID references the associated task UUID.
	TaskID string `json:"task_id"`
	// SourceName is the name of the third-party data source (e.g. tianyancha).
	SourceName string `json:"source_name"`
	// QueryName is the enterprise name used as the query parameter.
	QueryName string `json:"query_name,omitempty"`
	// ResponseData is the raw JSON response from the data source.
	ResponseData map[string]interface{} `json:"response_data,omitempty"`
	// IsConsistent indicates whether the response matches the three key fields.
	IsConsistent *bool `json:"is_consistent,omitempty"`
	// IsSuccess indicates whether the API call succeeded.
	IsSuccess bool `json:"is_success"`
	// FailReason describes the failure cause when IsSuccess is false.
	FailReason string `json:"fail_reason,omitempty"`
}

// CompanyConflictLog represents one row in the company_conflict_log table.
type CompanyConflictLog struct {
	CommonFields
	// TaskID references the associated task UUID.
	TaskID string `json:"task_id"`
	// ConflictType is the conflict category: 1=同名不同企 2=同企不同名.
	ConflictType int16 `json:"conflict_type"`
	// NewCreditCode is the credit code of the incoming record.
	NewCreditCode string `json:"new_credit_code,omitempty"`
	// NewName is the enterprise name of the incoming record.
	NewName string `json:"new_name,omitempty"`
	// ExistingEntCode is the UUID of the existing enterprise.
	ExistingEntCode string `json:"existing_ent_code,omitempty"`
	// ExistingCreditCode is the credit code of the existing record.
	ExistingCreditCode string `json:"existing_credit_code,omitempty"`
	// ExistingName is the name of the existing enterprise.
	ExistingName string `json:"existing_name,omitempty"`
	// Resolution is the handling decision:
	// 1=全部保留独立存在 2=复用UUID挂载新名称 3=拒绝拉新.
	Resolution int16 `json:"resolution,omitempty"`
	// ResolutionReason explains the resolution decision.
	ResolutionReason string `json:"resolution_reason,omitempty"`
}

// CompanyFirstLinkLog represents one row in the company_first_link_log table.
type CompanyFirstLinkLog struct {
	CommonFields
	// TaskID references the associated task UUID.
	TaskID string `json:"task_id"`
	// FactID references company_fact.id.
	FactID int64 `json:"fact_id"`
	// VerifierType is the verifier category:
	// 1=系统自验 2=用户确认 3=业务方确认.
	VerifierType int16 `json:"verifier_type"`
	// ExpectedData is the snapshot taken before writing, used for comparison.
	ExpectedData map[string]interface{} `json:"expected_data,omitempty"`
	// ActualData is the data read back after writing.
	ActualData map[string]interface{} `json:"actual_data,omitempty"`
	// IsConsistent indicates whether expected and actual data match.
	IsConsistent *bool `json:"is_consistent,omitempty"`
	// VerifyStatus is the verification outcome: 0=待确认 1=已通过 2=已拒绝.
	VerifyStatus int16 `json:"verify_status"`
	// VerifiedAt is the timestamp when the verification decision was made.
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}

// ─────────────────────────────────────────────
// upsert-company-fact
// ─────────────────────────────────────────────

// UpsertCompanyFactRequest is the input payload for the "upsert-company-fact" method.
type UpsertCompanyFactRequest struct {
	// EntCode is the enterprise UUID; leave empty to generate a new one on insert.
	EntCode string `json:"ent_code,omitempty"`
	// CreditCode is the unique social credit code used as the upsert key.
	CreditCode string `json:"credit_code"`
	// Name is the enterprise name.
	Name string `json:"name"`
	// LegalRep is the power-holder name.
	LegalRep string `json:"legal_rep,omitempty"`
	// LocationID references an address in eiker-address-db.
	LocationID string `json:"location_id,omitempty"`
	// FactTime is the time the fact occurred.
	FactTime *time.Time `json:"fact_time,omitempty"`
	// EntRecordScene is the on-boarding source type.
	EntRecordScene int16 `json:"ent_record_scene"`
	// LegalRepRecordScene is the power-holder on-boarding type.
	LegalRepRecordScene int16 `json:"legal_rep_record_scene,omitempty"`
	// TaskID is the associated task UUID.
	TaskID string `json:"task_id,omitempty"`
	// TraceID is the distributed tracing identifier.
	TraceID string `json:"trace_id,omitempty"`
	// Owner is the data owner.
	Owner string `json:"owner,omitempty"`
	// Creater is the operator performing the write.
	Creater string `json:"creater,omitempty"`
	// Remark is an optional free-text annotation.
	Remark string `json:"remark,omitempty"`
}

// UpsertCompanyFactResponse is the output payload for the "upsert-company-fact" method.
type UpsertCompanyFactResponse struct {
	// EntCode is the enterprise UUID of the upserted record.
	EntCode string `json:"ent_code"`
	// FactID is the primary key of the upserted company_fact row.
	FactID int64 `json:"fact_id"`
	// Action is "create" for a new insert or "update" for an existing record.
	Action string `json:"action"`
}

// ─────────────────────────────────────────────
// batch-upsert-company-fact
// ─────────────────────────────────────────────

// BatchUpsertCompanyFactRequest is the input payload for the "batch-upsert-company-fact" method.
type BatchUpsertCompanyFactRequest struct {
	// Items contains the list of fact records to upsert.
	Items []UpsertCompanyFactRequest `json:"items"`
}

// BatchUpsertCompanyFactResponse is the output payload for the "batch-upsert-company-fact" method.
type BatchUpsertCompanyFactResponse struct {
	// Results contains the outcome for each input item in the same order.
	Results []UpsertCompanyFactResponse `json:"results"`
	// SuccessCount is the number of items that were successfully upserted.
	SuccessCount int `json:"success_count"`
	// FailCount is the number of items that failed.
	FailCount int `json:"fail_count"`
}

// ─────────────────────────────────────────────
// query-company-fact-by-name
// ─────────────────────────────────────────────

// QueryCompanyFactByNameRequest is the input payload for the "query-company-fact-by-name" method.
type QueryCompanyFactByNameRequest struct {
	// Name is the enterprise name to search for (exact match).
	Name string `json:"name"`
}

// QueryCompanyFactByNameResponse is the output payload for the "query-company-fact-by-name" method.
type QueryCompanyFactByNameResponse struct {
	// Items contains all current records matching the given name.
	Items []CompanyFact `json:"items"`
}

// ─────────────────────────────────────────────
// query-company-fact-by-code
// ─────────────────────────────────────────────

// QueryCompanyFactByCodeRequest is the input payload for the "query-company-fact-by-code" method.
type QueryCompanyFactByCodeRequest struct {
	// CreditCode is the unified social credit code to look up.
	CreditCode string `json:"credit_code"`
}

// QueryCompanyFactByCodeResponse is the output payload for the "query-company-fact-by-code" method.
type QueryCompanyFactByCodeResponse struct {
	// Item is the current company_fact record for the given credit code, or nil if not found.
	Item *CompanyFact `json:"item"`
}

// ─────────────────────────────────────────────
// query-company-fact-by-ent-code
// ─────────────────────────────────────────────

// QueryCompanyFactByEntCodeRequest is the input payload for the "query-company-fact-by-ent-code" method.
type QueryCompanyFactByEntCodeRequest struct {
	// EntCode is the enterprise UUID to look up.
	EntCode string `json:"ent_code"`
}

// QueryCompanyFactByEntCodeResponse is the output payload for the "query-company-fact-by-ent-code" method.
type QueryCompanyFactByEntCodeResponse struct {
	// Item is the current company_fact record for the given ent_code, or nil if not found.
	Item *CompanyFact `json:"item"`
}

// ─────────────────────────────────────────────
// save-ocr-log
// ─────────────────────────────────────────────

// SaveOCRLogRequest is the input payload for the "save-ocr-log" method.
type SaveOCRLogRequest struct {
	// TaskID references the associated task UUID.
	TaskID string `json:"task_id"`
	// ModelName is the OCR model (paddleocr / easyocr / tesseract).
	ModelName string `json:"model_name"`
	// ImagePath is the image path or URL.
	ImagePath string `json:"image_path,omitempty"`
	// RawResult is the raw JSON output from the OCR engine.
	RawResult map[string]interface{} `json:"raw_result,omitempty"`
	// CreditCode is the extracted credit code.
	CreditCode string `json:"credit_code,omitempty"`
	// CompanyName is the extracted enterprise name.
	CompanyName string `json:"company_name,omitempty"`
	// LegalRepName is the extracted power-holder name.
	LegalRepName string `json:"legal_rep_name,omitempty"`
	// IsSuccess indicates whether extraction succeeded.
	IsSuccess bool `json:"is_success"`
	// FailReason is the failure description when IsSuccess is false.
	FailReason string `json:"fail_reason,omitempty"`
	// TraceID is the distributed tracing identifier.
	TraceID string `json:"trace_id,omitempty"`
	// Creater is the operator performing the write.
	Creater string `json:"creater,omitempty"`
}

// SaveOCRLogResponse is the output payload for the "save-ocr-log" method.
type SaveOCRLogResponse struct {
	// ID is the auto-generated primary key of the new row.
	ID int64 `json:"id"`
}

// ─────────────────────────────────────────────
// save-data-verify-log
// ─────────────────────────────────────────────

// SaveDataVerifyLogRequest is the input payload for the "save-data-verify-log" method.
type SaveDataVerifyLogRequest struct {
	// TaskID references the associated task UUID.
	TaskID string `json:"task_id"`
	// SourceName is the data-source name (e.g. tianyancha).
	SourceName string `json:"source_name"`
	// QueryName is the enterprise name used in the query.
	QueryName string `json:"query_name,omitempty"`
	// ResponseData is the raw JSON from the data source.
	ResponseData map[string]interface{} `json:"response_data,omitempty"`
	// IsConsistent indicates whether the response matches the input three key fields.
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

// SaveDataVerifyLogResponse is the output payload for the "save-data-verify-log" method.
type SaveDataVerifyLogResponse struct {
	// ID is the auto-generated primary key of the new row.
	ID int64 `json:"id"`
}

// ─────────────────────────────────────────────
// save-conflict-log
// ─────────────────────────────────────────────

// SaveConflictLogRequest is the input payload for the "save-conflict-log" method.
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

// SaveConflictLogResponse is the output payload for the "save-conflict-log" method.
type SaveConflictLogResponse struct {
	// ID is the auto-generated primary key of the new row.
	ID int64 `json:"id"`
}

// ─────────────────────────────────────────────
// save-first-link-log
// ─────────────────────────────────────────────

// SaveFirstLinkLogRequest is the input payload for the "save-first-link-log" method.
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

// SaveFirstLinkLogResponse is the output payload for the "save-first-link-log" method.
type SaveFirstLinkLogResponse struct {
	// ID is the auto-generated primary key of the new row.
	ID int64 `json:"id"`
}

// ─────────────────────────────────────────────
// Dapr pub/sub 事件 Payload
// ─────────────────────────────────────────────

// CompanyFactUpsertedEvent is published to the "company-fact-upserted" topic
// after a successful upsert of company_fact.
type CompanyFactUpsertedEvent struct {
	// EntCode is the enterprise UUID.
	EntCode string `json:"ent_code"`
	// FactID is the primary key of the upserted row.
	FactID int64 `json:"fact_id"`
	// Action is "create" or "update".
	Action string `json:"action"`
}
