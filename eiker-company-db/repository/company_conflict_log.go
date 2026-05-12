package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/xiaoyi-claude/eiker-company-db/model"
)

// SaveConflictLog inserts a new company_conflict_log row and returns the generated primary key.
func (r *Repository) SaveConflictLog(ctx context.Context, req *model.SaveConflictLogRequest) (int64, error) {
	now := time.Now().UTC()
	var id int64
	err := r.pool.QueryRow(ctx, `
		INSERT INTO company_conflict_log (
			task_id, conflict_type, new_credit_code, new_name,
			existing_ent_code, existing_credit_code, existing_name,
			resolution, resolution_reason,
			is_current, trace_id, creater, create_time, updater, update_time
		) VALUES (
			$1,  $2,  $3,  $4,
			$5,  $6,  $7,
			$8,  $9,
			TRUE, $10, $11, $12, $11, $12
		) RETURNING id`,
		req.TaskID, req.ConflictType, req.NewCreditCode, req.NewName,
		req.ExistingEntCode, req.ExistingCreditCode, req.ExistingName,
		req.Resolution, req.ResolutionReason,
		req.TraceID, req.Creater, now,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert company_conflict_log: %w", err)
	}
	return id, nil
}
