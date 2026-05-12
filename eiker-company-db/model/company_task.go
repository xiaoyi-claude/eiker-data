package model

// CompanyTask represents one row in the company_task table.
// The task lifecycle is managed by the business layer (eiker-company-update);
// this atomic service only exposes CRUD when called via Dapr service invocation.
type CompanyTask struct {
	CommonFields
	// TaskID is the task UUID assigned by the business layer.
	TaskID string `json:"task_id"`
	// TaskType is the on-boarding source: 1=基础拉新 2=用户拉新 3=业务拉新 4=公司拉新.
	TaskType int16 `json:"task_type"`
	// TaskStatus is the task state: 0=PENDING 1=IN_PROGRESS 2=SUCCESS 3=FAILED.
	TaskStatus int16 `json:"task_status"`
	// InputSource identifies the input origin (file path, image URL, data key, etc.).
	InputSource string `json:"input_source,omitempty"`
	// Payload stores adoption source, credentials, and processing rules as free-form JSON.
	Payload map[string]interface{} `json:"payload,omitempty"`
	// Result stores the processing outcome as free-form JSON.
	Result map[string]interface{} `json:"result,omitempty"`
	// WorkflowInstanceID is the Dapr Workflow instance ID for this task.
	WorkflowInstanceID string `json:"workflow_instance_id,omitempty"`
	// RetryCount is the number of times this task has been retried.
	RetryCount int16 `json:"retry_count"`
	// ErrorMsg describes the failure reason when TaskStatus is FAILED.
	ErrorMsg string `json:"error_msg,omitempty"`
}
