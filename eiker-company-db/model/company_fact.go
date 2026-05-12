package model

import "time"

// ─────────────────────────────────────────────
// 实体
// ─────────────────────────────────────────────

// CompanyFact represents one row in the company_fact table.
// It stores the core fact data for an enterprise in the company name domain.
type CompanyFact struct {
	CommonFields
	// EntCode is the enterprise UUID v4 used as the external identifier (public domain).
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

// ─────────────────────────────────────────────
// upsert-company-fact
// ─────────────────────────────────────────────

// UpsertCompanyFactRequest is the input payload for the "upsert-company-fact" Dapr method.
type UpsertCompanyFactRequest struct {
	// EntCode is the enterprise UUID; leave empty to auto-generate on insert.
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

// UpsertCompanyFactResponse is the output payload for the "upsert-company-fact" Dapr method.
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

// BatchUpsertCompanyFactRequest is the input payload for the "batch-upsert-company-fact" Dapr method.
type BatchUpsertCompanyFactRequest struct {
	// Items contains the list of fact records to upsert in order.
	Items []UpsertCompanyFactRequest `json:"items"`
}

// BatchUpsertCompanyFactResponse is the output payload for the "batch-upsert-company-fact" Dapr method.
type BatchUpsertCompanyFactResponse struct {
	// Results contains the outcome for each input item, index-aligned with the request.
	Results []UpsertCompanyFactResponse `json:"results"`
	// SuccessCount is the number of items successfully upserted.
	SuccessCount int `json:"success_count"`
	// FailCount is the number of items that failed.
	FailCount int `json:"fail_count"`
}

// ─────────────────────────────────────────────
// query-company-fact-by-name
// ─────────────────────────────────────────────

// QueryCompanyFactByNameRequest is the input payload for the "query-company-fact-by-name" Dapr method.
type QueryCompanyFactByNameRequest struct {
	// Name is the enterprise name to search for (exact match against current records).
	Name string `json:"name"`
}

// QueryCompanyFactByNameResponse is the output payload for the "query-company-fact-by-name" Dapr method.
type QueryCompanyFactByNameResponse struct {
	// Items contains all current records whose name matches exactly.
	Items []CompanyFact `json:"items"`
}

// ─────────────────────────────────────────────
// query-company-fact-by-code
// ─────────────────────────────────────────────

// QueryCompanyFactByCodeRequest is the input payload for the "query-company-fact-by-code" Dapr method.
type QueryCompanyFactByCodeRequest struct {
	// CreditCode is the unified social credit code to look up.
	CreditCode string `json:"credit_code"`
}

// QueryCompanyFactByCodeResponse is the output payload for the "query-company-fact-by-code" Dapr method.
type QueryCompanyFactByCodeResponse struct {
	// Item is the current record for the given credit code, or nil if not found.
	Item *CompanyFact `json:"item"`
}

// ─────────────────────────────────────────────
// query-company-fact-by-ent-code
// ─────────────────────────────────────────────

// QueryCompanyFactByEntCodeRequest is the input payload for the "query-company-fact-by-ent-code" Dapr method.
type QueryCompanyFactByEntCodeRequest struct {
	// EntCode is the enterprise UUID to look up.
	EntCode string `json:"ent_code"`
}

// QueryCompanyFactByEntCodeResponse is the output payload for the "query-company-fact-by-ent-code" Dapr method.
type QueryCompanyFactByEntCodeResponse struct {
	// Item is the current record for the given ent_code, or nil if not found.
	Item *CompanyFact `json:"item"`
}
