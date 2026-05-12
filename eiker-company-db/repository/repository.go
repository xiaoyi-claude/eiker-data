// Package repository provides PostgreSQL data-access operations for the
// company domain using pgx/v5 connection pooling.
package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xiaoyi-claude/eiker-company-db/model"
)

// Repository provides all PostgreSQL CRUD operations for the eiker-company-db service.
type Repository struct {
	// pool is the pgx connection pool shared by all operations.
	pool *pgxpool.Pool
}

// New constructs a Repository backed by the given connection pool.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// ─────────────────────────────────────────────
// company_fact
// ─────────────────────────────────────────────

// UpsertCompanyFact inserts a new company_fact row when the credit_code does not exist,
// or updates the existing current row when it does.
// Returns the resulting entity, the action string ("create" or "update"), and any error.
func (r *Repository) UpsertCompanyFact(ctx context.Context, req *model.UpsertCompanyFactRequest) (*model.CompanyFact, string, error) {
	now := time.Now().UTC()

	// Check for an existing current record by credit_code.
	var existingID int64
	var existingEntCode string
	err := r.pool.QueryRow(ctx,
		`SELECT id, ent_code FROM company_fact WHERE credit_code = $1 AND is_current = TRUE LIMIT 1`,
		req.CreditCode,
	).Scan(&existingID, &existingEntCode)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, "", fmt.Errorf("query existing company_fact: %w", err)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		// No existing record – insert a new one.
		entCode := req.EntCode
		if entCode == "" {
			entCode = uuid.New().String()
		}
		factRecordTime := now

		var newID int64
		insertErr := r.pool.QueryRow(ctx, `
			INSERT INTO company_fact (
				ent_code, credit_code, name, legal_rep, location_id,
				fact_time, fact_record_time, ent_record_scene, legal_rep_record_scene,
				task_id, is_current, trace_id, owner, creater, create_time, updater, update_time, remark
			) VALUES (
				$1,  $2,  $3,  $4,  $5,
				$6,  $7,  $8,  $9,
				$10, TRUE, $11, $12, $13, $14, $13, $14, $15
			) RETURNING id`,
			entCode, req.CreditCode, req.Name, req.LegalRep, req.LocationID,
			req.FactTime, factRecordTime, req.EntRecordScene, req.LegalRepRecordScene,
			req.TaskID, req.TraceID, req.Owner, req.Creater, now, req.Remark,
		).Scan(&newID)
		if insertErr != nil {
			return nil, "", fmt.Errorf("insert company_fact: %w", insertErr)
		}

		fact := &model.CompanyFact{
			CommonFields: model.CommonFields{
				ID: newID, IsCurrent: true,
				TraceID: req.TraceID, Owner: req.Owner, Creater: req.Creater,
				CreateTime: now, Updater: req.Creater, UpdateTime: now, Remark: req.Remark,
			},
			EntCode: entCode, CreditCode: req.CreditCode, Name: req.Name,
			LegalRep: req.LegalRep, LocationID: req.LocationID,
			FactTime: req.FactTime, FactRecordTime: &factRecordTime,
			EntRecordScene: req.EntRecordScene, LegalRepRecordScene: req.LegalRepRecordScene,
			TaskID: req.TaskID,
		}
		return fact, "create", nil
	}

	// Existing record found – update it.
	_, updateErr := r.pool.Exec(ctx, `
		UPDATE company_fact SET
			name                    = $1,
			legal_rep               = $2,
			location_id             = $3,
			fact_time               = $4,
			ent_record_scene        = $5,
			legal_rep_record_scene  = $6,
			task_id                 = $7,
			updater                 = $8,
			update_time             = $9,
			remark                  = $10
		WHERE id = $11`,
		req.Name, req.LegalRep, req.LocationID,
		req.FactTime, req.EntRecordScene, req.LegalRepRecordScene,
		req.TaskID, req.Creater, now, req.Remark,
		existingID,
	)
	if updateErr != nil {
		return nil, "", fmt.Errorf("update company_fact: %w", updateErr)
	}

	fact := &model.CompanyFact{
		CommonFields: model.CommonFields{
			ID: existingID, IsCurrent: true,
			Updater: req.Creater, UpdateTime: now,
		},
		EntCode: existingEntCode, CreditCode: req.CreditCode, Name: req.Name,
		LegalRep: req.LegalRep, LocationID: req.LocationID,
		FactTime: req.FactTime, EntRecordScene: req.EntRecordScene,
		LegalRepRecordScene: req.LegalRepRecordScene, TaskID: req.TaskID,
	}
	return fact, "update", nil
}

// QueryCompanyFactByName returns all current company_fact rows whose name exactly matches the given value.
func (r *Repository) QueryCompanyFactByName(ctx context.Context, name string) ([]model.CompanyFact, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, is_current, trace_id, owner, creater, create_time, updater, update_time, remark,
		        ent_code, credit_code, name, legal_rep, location_id,
		        fact_time, fact_record_time, ent_record_scene, legal_rep_record_scene, task_id
		 FROM company_fact WHERE name = $1 AND is_current = TRUE`,
		name,
	)
	if err != nil {
		return nil, fmt.Errorf("query company_fact by name: %w", err)
	}
	defer rows.Close()

	return scanCompanyFacts(rows)
}

// QueryCompanyFactByCode returns the current company_fact row for the given credit_code, or nil if not found.
func (r *Repository) QueryCompanyFactByCode(ctx context.Context, creditCode string) (*model.CompanyFact, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, is_current, trace_id, owner, creater, create_time, updater, update_time, remark,
		        ent_code, credit_code, name, legal_rep, location_id,
		        fact_time, fact_record_time, ent_record_scene, legal_rep_record_scene, task_id
		 FROM company_fact WHERE credit_code = $1 AND is_current = TRUE LIMIT 1`,
		creditCode,
	)
	fact, err := scanCompanyFact(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return fact, err
}

// QueryCompanyFactByEntCode returns the current company_fact row for the given ent_code, or nil if not found.
func (r *Repository) QueryCompanyFactByEntCode(ctx context.Context, entCode string) (*model.CompanyFact, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, is_current, trace_id, owner, creater, create_time, updater, update_time, remark,
		        ent_code, credit_code, name, legal_rep, location_id,
		        fact_time, fact_record_time, ent_record_scene, legal_rep_record_scene, task_id
		 FROM company_fact WHERE ent_code = $1 AND is_current = TRUE LIMIT 1`,
		entCode,
	)
	fact, err := scanCompanyFact(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return fact, err
}

// ─────────────────────────────────────────────
// company_ocr_log
// ─────────────────────────────────────────────

// SaveOCRLog inserts a new row into company_ocr_log and returns the generated primary key.
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

// ─────────────────────────────────────────────
// company_data_verify_log
// ─────────────────────────────────────────────

// SaveDataVerifyLog inserts a new row into company_data_verify_log and returns the generated primary key.
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

// ─────────────────────────────────────────────
// company_conflict_log
// ─────────────────────────────────────────────

// SaveConflictLog inserts a new row into company_conflict_log and returns the generated primary key.
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

// ─────────────────────────────────────────────
// company_first_link_log
// ─────────────────────────────────────────────

// SaveFirstLinkLog inserts a new row into company_first_link_log and returns the generated primary key.
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

// ─────────────────────────────────────────────
// 内部辅助函数
// ─────────────────────────────────────────────

// scanCompanyFacts iterates pgx.Rows and decodes each row into a CompanyFact.
func scanCompanyFacts(rows pgx.Rows) ([]model.CompanyFact, error) {
	var facts []model.CompanyFact
	for rows.Next() {
		var f model.CompanyFact
		if err := scanCompanyFactRow(rows.Scan, &f); err != nil {
			return nil, err
		}
		facts = append(facts, f)
	}
	return facts, rows.Err()
}

// scanCompanyFact scans a single pgx.Row into a CompanyFact.
func scanCompanyFact(row pgx.Row) (*model.CompanyFact, error) {
	var f model.CompanyFact
	if err := scanCompanyFactRow(row.Scan, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// scanCompanyFactRow is the shared column scanner for both pgx.Row and pgx.Rows.
func scanCompanyFactRow(scan func(...any) error, f *model.CompanyFact) error {
	return scan(
		&f.ID, &f.IsCurrent, &f.TraceID, &f.Owner, &f.Creater,
		&f.CreateTime, &f.Updater, &f.UpdateTime, &f.Remark,
		&f.EntCode, &f.CreditCode, &f.Name, &f.LegalRep, &f.LocationID,
		&f.FactTime, &f.FactRecordTime, &f.EntRecordScene, &f.LegalRepRecordScene, &f.TaskID,
	)
}

// jsonMarshalNullable returns nil when m is empty, otherwise the JSON-encoded bytes.
// This prevents storing an explicit "null" in JSONB columns when no data is provided.
func jsonMarshalNullable(m map[string]interface{}) ([]byte, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return json.Marshal(m)
}
