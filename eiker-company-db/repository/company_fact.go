package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/xiaoyi-claude/eiker-company-db/model"
)

// UpsertCompanyFact inserts a new company_fact row when the credit_code does not exist,
// or updates the matching current row when it does.
// Returns the full resulting entity, the action string ("create" or "update"), and any error.
func (r *Repository) UpsertCompanyFact(ctx context.Context, req *model.UpsertCompanyFactRequest) (*model.CompanyFact, string, error) {
	now := time.Now().UTC()

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
		return r.insertCompanyFact(ctx, req, now)
	}
	return r.updateCompanyFact(ctx, req, existingID, existingEntCode, now)
}

// insertCompanyFact inserts a brand-new company_fact row.
func (r *Repository) insertCompanyFact(ctx context.Context, req *model.UpsertCompanyFactRequest, now time.Time) (*model.CompanyFact, string, error) {
	entCode := req.EntCode
	if entCode == "" {
		entCode = uuid.New().String()
	}
	factRecordTime := now

	var newID int64
	err := r.pool.QueryRow(ctx, `
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
	if err != nil {
		return nil, "", fmt.Errorf("insert company_fact: %w", err)
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

// updateCompanyFact updates an existing company_fact row.
func (r *Repository) updateCompanyFact(ctx context.Context, req *model.UpsertCompanyFactRequest, existingID int64, existingEntCode string, now time.Time) (*model.CompanyFact, string, error) {
	_, err := r.pool.Exec(ctx, `
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
	if err != nil {
		return nil, "", fmt.Errorf("update company_fact: %w", err)
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
	return scanCompanyFactRows(rows)
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
	fact, err := scanCompanyFactRow(row)
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
	fact, err := scanCompanyFactRow(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return fact, err
}

// scanCompanyFactRows iterates pgx.Rows and decodes each row into a CompanyFact.
func scanCompanyFactRows(rows pgx.Rows) ([]model.CompanyFact, error) {
	var facts []model.CompanyFact
	for rows.Next() {
		var f model.CompanyFact
		if err := scanCompanyFactColumns(rows.Scan, &f); err != nil {
			return nil, err
		}
		facts = append(facts, f)
	}
	return facts, rows.Err()
}

// scanCompanyFactRow scans a single pgx.Row into a CompanyFact.
func scanCompanyFactRow(row pgx.Row) (*model.CompanyFact, error) {
	var f model.CompanyFact
	if err := scanCompanyFactColumns(row.Scan, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// scanCompanyFactColumns is the shared column scanner used by both single-row and multi-row scans.
func scanCompanyFactColumns(scan func(...any) error, f *model.CompanyFact) error {
	return scan(
		&f.ID, &f.IsCurrent, &f.TraceID, &f.Owner, &f.Creater,
		&f.CreateTime, &f.Updater, &f.UpdateTime, &f.Remark,
		&f.EntCode, &f.CreditCode, &f.Name, &f.LegalRep, &f.LocationID,
		&f.FactTime, &f.FactRecordTime, &f.EntRecordScene, &f.LegalRepRecordScene, &f.TaskID,
	)
}
