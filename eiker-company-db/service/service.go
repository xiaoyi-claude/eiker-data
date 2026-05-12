// Package service implements the business logic for the eiker-company-db service.
// After a successful company_fact upsert it publishes a "company-fact-upserted"
// event to the configured Dapr pub/sub component so that downstream consumers
// (e.g. eiker-company-es) can synchronise their indexes.
package service

import (
	"context"
	"fmt"
	"log"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/xiaoyi-claude/eiker-company-db/model"
	"github.com/xiaoyi-claude/eiker-company-db/repository"
)

// Service coordinates data access via Repository and event publishing via the Dapr client.
type Service struct {
	// repo provides all PostgreSQL CRUD operations.
	repo *repository.Repository
	// daprClient is the Dapr SDK client used for pub/sub publishing.
	daprClient dapr.Client
	// pubsubName is the name of the Dapr pub/sub component (e.g. "eiker-pubsub").
	pubsubName string
}

// New constructs a Service with the given repository, Dapr client, and pub/sub component name.
func New(repo *repository.Repository, daprClient dapr.Client, pubsubName string) *Service {
	return &Service{
		repo:       repo,
		daprClient: daprClient,
		pubsubName: pubsubName,
	}
}

// ─────────────────────────────────────────────
// company_fact 写操作
// ─────────────────────────────────────────────

// UpsertCompanyFact inserts or updates a company_fact record and publishes a
// "company-fact-upserted" event on success.
func (s *Service) UpsertCompanyFact(ctx context.Context, req *model.UpsertCompanyFactRequest) (*model.UpsertCompanyFactResponse, error) {
	fact, action, err := s.repo.UpsertCompanyFact(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("upsert company fact: %w", err)
	}

	s.publishFactUpserted(ctx, fact.EntCode, fact.ID, action)

	return &model.UpsertCompanyFactResponse{
		EntCode: fact.EntCode,
		FactID:  fact.ID,
		Action:  action,
	}, nil
}

// BatchUpsertCompanyFact upserts multiple company_fact records sequentially.
// Each item is processed independently; a failure on one item is recorded in the
// response but does not abort the remaining items.
func (s *Service) BatchUpsertCompanyFact(ctx context.Context, req *model.BatchUpsertCompanyFactRequest) (*model.BatchUpsertCompanyFactResponse, error) {
	resp := &model.BatchUpsertCompanyFactResponse{
		Results: make([]model.UpsertCompanyFactResponse, 0, len(req.Items)),
	}

	for i := range req.Items {
		item := &req.Items[i]
		r, err := s.UpsertCompanyFact(ctx, item)
		if err != nil {
			log.Printf("batch-upsert item[%d] credit_code=%q failed: %v", i, item.CreditCode, err)
			resp.FailCount++
			// Append a zero-value placeholder so the result slice stays index-aligned.
			resp.Results = append(resp.Results, model.UpsertCompanyFactResponse{})
			continue
		}
		resp.Results = append(resp.Results, *r)
		resp.SuccessCount++
	}

	return resp, nil
}

// ─────────────────────────────────────────────
// company_fact 查操作
// ─────────────────────────────────────────────

// QueryCompanyFactByName returns all current records whose name exactly matches the given value.
func (s *Service) QueryCompanyFactByName(ctx context.Context, req *model.QueryCompanyFactByNameRequest) (*model.QueryCompanyFactByNameResponse, error) {
	facts, err := s.repo.QueryCompanyFactByName(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("query company fact by name: %w", err)
	}
	if facts == nil {
		facts = []model.CompanyFact{}
	}
	return &model.QueryCompanyFactByNameResponse{Items: facts}, nil
}

// QueryCompanyFactByCode returns the current record for the given credit_code, or nil Item if not found.
func (s *Service) QueryCompanyFactByCode(ctx context.Context, req *model.QueryCompanyFactByCodeRequest) (*model.QueryCompanyFactByCodeResponse, error) {
	fact, err := s.repo.QueryCompanyFactByCode(ctx, req.CreditCode)
	if err != nil {
		return nil, fmt.Errorf("query company fact by code: %w", err)
	}
	return &model.QueryCompanyFactByCodeResponse{Item: fact}, nil
}

// QueryCompanyFactByEntCode returns the current record for the given ent_code, or nil Item if not found.
func (s *Service) QueryCompanyFactByEntCode(ctx context.Context, req *model.QueryCompanyFactByEntCodeRequest) (*model.QueryCompanyFactByEntCodeResponse, error) {
	fact, err := s.repo.QueryCompanyFactByEntCode(ctx, req.EntCode)
	if err != nil {
		return nil, fmt.Errorf("query company fact by ent_code: %w", err)
	}
	return &model.QueryCompanyFactByEntCodeResponse{Item: fact}, nil
}

// ─────────────────────────────────────────────
// 日志写操作
// ─────────────────────────────────────────────

// SaveOCRLog persists one OCR processing log entry and returns the generated ID.
func (s *Service) SaveOCRLog(ctx context.Context, req *model.SaveOCRLogRequest) (*model.SaveOCRLogResponse, error) {
	id, err := s.repo.SaveOCRLog(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("save ocr log: %w", err)
	}
	return &model.SaveOCRLogResponse{ID: id}, nil
}

// SaveDataVerifyLog persists one third-party data verification log entry and returns the generated ID.
func (s *Service) SaveDataVerifyLog(ctx context.Context, req *model.SaveDataVerifyLogRequest) (*model.SaveDataVerifyLogResponse, error) {
	id, err := s.repo.SaveDataVerifyLog(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("save data verify log: %w", err)
	}
	return &model.SaveDataVerifyLogResponse{ID: id}, nil
}

// SaveConflictLog persists one conflict-handling log entry and returns the generated ID.
func (s *Service) SaveConflictLog(ctx context.Context, req *model.SaveConflictLogRequest) (*model.SaveConflictLogResponse, error) {
	id, err := s.repo.SaveConflictLog(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("save conflict log: %w", err)
	}
	return &model.SaveConflictLogResponse{ID: id}, nil
}

// SaveFirstLinkLog persists one first-link verification log entry and returns the generated ID.
func (s *Service) SaveFirstLinkLog(ctx context.Context, req *model.SaveFirstLinkLogRequest) (*model.SaveFirstLinkLogResponse, error) {
	id, err := s.repo.SaveFirstLinkLog(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("save first link log: %w", err)
	}
	return &model.SaveFirstLinkLogResponse{ID: id}, nil
}

// ─────────────────────────────────────────────
// 内部辅助
// ─────────────────────────────────────────────

// publishFactUpserted publishes a "company-fact-upserted" event to the Dapr pub/sub component.
// Publish failures are logged but do not propagate – the write is already committed.
func (s *Service) publishFactUpserted(ctx context.Context, entCode string, factID int64, action string) {
	event := &model.CompanyFactUpsertedEvent{
		EntCode: entCode,
		FactID:  factID,
		Action:  action,
	}
	if err := s.daprClient.PublishEvent(ctx, s.pubsubName, "company-fact-upserted", event); err != nil {
		log.Printf("warn: publish company-fact-upserted (ent_code=%s fact_id=%d): %v", entCode, factID, err)
	}
}
