package model

// CompanyOCRLog represents one row in the company_ocr_log table.
// It records the raw output and extracted triple (credit code, name, legal rep)
// for each OCR recognition attempt.
type CompanyOCRLog struct {
	CommonFields
	// TaskID references the associated task UUID.
	TaskID string `json:"task_id"`
	// ModelName is the OCR engine used (paddleocr / easyocr / tesseract).
	ModelName string `json:"model_name"`
	// ImagePath is the path or URL of the source image.
	ImagePath string `json:"image_path,omitempty"`
	// RawResult is the raw JSON output from the OCR engine.
	RawResult map[string]interface{} `json:"raw_result,omitempty"`
	// CreditCode is the unified social credit code extracted by OCR.
	CreditCode string `json:"credit_code,omitempty"`
	// CompanyName is the enterprise name extracted by OCR.
	CompanyName string `json:"company_name,omitempty"`
	// LegalRepName is the power-holder name extracted by OCR.
	LegalRepName string `json:"legal_rep_name,omitempty"`
	// IsSuccess indicates whether the triple extraction succeeded.
	IsSuccess bool `json:"is_success"`
	// FailReason describes the failure cause when IsSuccess is false.
	FailReason string `json:"fail_reason,omitempty"`
}

// SaveOCRLogRequest is the input payload for the "save-ocr-log" Dapr method.
type SaveOCRLogRequest struct {
	// TaskID references the associated task UUID.
	TaskID string `json:"task_id"`
	// ModelName is the OCR engine used.
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

// SaveOCRLogResponse is the output payload for the "save-ocr-log" Dapr method.
type SaveOCRLogResponse struct {
	// ID is the auto-generated primary key of the new row.
	ID int64 `json:"id"`
}
