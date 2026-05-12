package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/xiaoyi-claude/eiker-company-db/model"
)

// SaveFirstLinkLog inserts a new company_first_link_log row and returns the generated primary key.
func (r *Repository) SaveFirstLinkLog(ctx context.Context, req *model.SaveFirstLinkLogRequest) (int64, error) {
	expectedJSON, err := jsonMarshalNullable(req.ExpectedData)
	if err != nil {
		return 0, fmt.Errorf("marshal expected_data: %w", err)
	}
	actualJSON, err := jsonMarshalNullable(req.ActualData)
	if err != nil {
		return 0, fmt.Errorf("marshal actual_data: %w", err)
	}

	now := time.Now().UTC()
	var id int64
	err = r.pool.QueryRow(ctx, `
		INSERT INTO company_first_link_log (
			task_id, fact_id, verifier_type,
			expected_data, actual_data, is_consistent,
			verify_status, verified_at,
			is_current, trace_id, creater, create_time, updater, update_time
		) VALUES (
			$1,  $2,  $3,
			$4,  $5,  $6,
			$7,  $8,
			TRUE, $9,  $10, $11, $10, $11
		) RETURNING id`,
		req.TaskID, req.FactID, req.VerifierType,
		expectedJSON, actualJSON, req.IsConsistent,
		req.VerifyStatus, req.VerifiedAt,
		req.TraceID, req.Creater, now,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert company_first_link_log: %w", err)
	}
	return id, nil
}
