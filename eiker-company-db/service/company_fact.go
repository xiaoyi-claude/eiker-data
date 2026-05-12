package service

import (
	"context"
	"fmt"
	"log"

	"github.com/xiaoyi-claude/eiker-company-db/model"
)

// UpsertCompanyFact inserts or updates a company_fact record and publishes a
// "company-fact-upserted" event containing the full fact on success.
func (s *Service) UpsertCompanyFact(ctx context.Context, req *model.UpsertCompanyFactRequest) (*model.UpsertCompanyFactResponse, error) {
	fact, action, err := s.repo.UpsertCompanyFact(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("upsert company fact: %w", err)
	}

	s.publishFactUpserted(ctx, fact, action)

	return &model.UpsertCompanyFactResponse{
		EntCode: fact.EntCode,
		FactID:  fact.ID,
		Action:  action,
	}, nil
}

// BatchUpsertCompanyFact upserts multiple company_fact records sequentially.
// A failure on one item is recorded in the response but does not abort the remaining items.
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
			resp.Results = append(resp.Results, model.UpsertCompanyFactResponse{})
			continue
		}
		resp.Results = append(resp.Results, *r)
		resp.SuccessCount++
	}

	return resp, nil
}

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
