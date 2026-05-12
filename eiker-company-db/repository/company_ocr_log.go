package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/xiaoyi-claude/eiker-company-db/model"
)

// SaveOCRLog inserts a new company_ocr_log row and returns the generated primary key.
func (r *Repository) SaveOCRLog(ctx context.Context, req *model.SaveOCRLogRequest) (int64, error) {
	rawJSON, err := jsonMarshalNullable(req.RawResult)
	if err != nil {
		return 0, fmt.Errorf("marshal raw_result: %w", err)
	}

	now := time.Now().UTC()
	var id int64
	err = r.pool.QueryRow(ctx, `
		INSERT INTO company_ocr_log (
			task_id, model_name, image_path, raw_result,
			credit_code, company_name, legal_rep_name, is_success, fail_reason,
			is_current, trace_id, creater, create_time, updater, update_time
		) VALUES (
			$1,  $2,  $3,  $4,
			$5,  $6,  $7,  $8,  $9,
			TRUE, $10, $11, $12, $11, $12
		) RETURNING id`,
		req.TaskID, req.ModelName, req.ImagePath, rawJSON,
		req.CreditCode, req.CompanyName, req.LegalRepName, req.IsSuccess, req.FailReason,
		req.TraceID, req.Creater, now,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert company_ocr_log: %w", err)
	}
	return id, nil
}
