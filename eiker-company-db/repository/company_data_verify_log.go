package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/xiaoyi-claude/eiker-company-db/model"
)

// SaveDataVerifyLog inserts a new company_data_verify_log row and returns the generated primary key.
func (r *Repository) SaveDataVerifyLog(ctx context.Context, req *model.SaveDataVerifyLogRequest) (int64, error) {
	responseJSON, err := jsonMarshalNullable(req.ResponseData)
	if err != nil {
		return 0, fmt.Errorf("marshal response_data: %w", err)
	}

	now := time.Now().UTC()
	var id int64
	err = r.pool.QueryRow(ctx, `
		INSERT INTO company_data_verify_log (
			task_id, source_name, query_name, response_data,
			is_consistent, is_success, fail_reason,
			is_current, trace_id, creater, create_time, updater, update_time
		) VALUES (
			$1,  $2,  $3,  $4,
			$5,  $6,  $7,
			TRUE, $8,  $9,  $10, $9, $10
		) RETURNING id`,
		req.TaskID, req.SourceName, req.QueryName, responseJSON,
		req.IsConsistent, req.IsSuccess, req.FailReason,
		req.TraceID, req.Creater, now,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert company_data_verify_log: %w", err)
	}
	return id, nil
}
