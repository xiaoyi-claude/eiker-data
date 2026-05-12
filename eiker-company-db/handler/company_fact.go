package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dapr/go-sdk/service/common"
	"github.com/xiaoyi-claude/eiker-company-db/model"
)

// UpsertCompanyFact is the Dapr handler for the "upsert-company-fact" method.
// Input:  JSON-encoded model.UpsertCompanyFactRequest
// Output: JSON-encoded model.UpsertCompanyFactResponse
func (h *Handler) UpsertCompanyFact(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	var req model.UpsertCompanyFactRequest
	if err := json.Unmarshal(in.Data, &req); err != nil {
		return nil, fmt.Errorf("upsert-company-fact: invalid payload: %w", err)
	}
	resp, err := h.svc.UpsertCompanyFact(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("upsert-company-fact: %w", err)
	}
	return marshalResponse(resp)
}

// BatchUpsertCompanyFact is the Dapr handler for the "batch-upsert-company-fact" method.
// Input:  JSON-encoded model.BatchUpsertCompanyFactRequest
// Output: JSON-encoded model.BatchUpsertCompanyFactResponse
func (h *Handler) BatchUpsertCompanyFact(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	var req model.BatchUpsertCompanyFactRequest
	if err := json.Unmarshal(in.Data, &req); err != nil {
		return nil, fmt.Errorf("batch-upsert-company-fact: invalid payload: %w", err)
	}
	resp, err := h.svc.BatchUpsertCompanyFact(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("batch-upsert-company-fact: %w", err)
	}
	return marshalResponse(resp)
}

// QueryCompanyFactByName is the Dapr handler for the "query-company-fact-by-name" method.
// Input:  JSON-encoded model.QueryCompanyFactByNameRequest
// Output: JSON-encoded model.QueryCompanyFactByNameResponse
func (h *Handler) QueryCompanyFactByName(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	var req model.QueryCompanyFactByNameRequest
	if err := json.Unmarshal(in.Data, &req); err != nil {
		return nil, fmt.Errorf("query-company-fact-by-name: invalid payload: %w", err)
	}
	resp, err := h.svc.QueryCompanyFactByName(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("query-company-fact-by-name: %w", err)
	}
	return marshalResponse(resp)
}

// QueryCompanyFactByCode is the Dapr handler for the "query-company-fact-by-code" method.
// Input:  JSON-encoded model.QueryCompanyFactByCodeRequest
// Output: JSON-encoded model.QueryCompanyFactByCodeResponse
func (h *Handler) QueryCompanyFactByCode(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	var req model.QueryCompanyFactByCodeRequest
	if err := json.Unmarshal(in.Data, &req); err != nil {
		return nil, fmt.Errorf("query-company-fact-by-code: invalid payload: %w", err)
	}
	resp, err := h.svc.QueryCompanyFactByCode(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("query-company-fact-by-code: %w", err)
	}
	return marshalResponse(resp)
}

// QueryCompanyFactByEntCode is the Dapr handler for the "query-company-fact-by-ent-code" method.
// Input:  JSON-encoded model.QueryCompanyFactByEntCodeRequest
// Output: JSON-encoded model.QueryCompanyFactByEntCodeResponse
func (h *Handler) QueryCompanyFactByEntCode(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	var req model.QueryCompanyFactByEntCodeRequest
	if err := json.Unmarshal(in.Data, &req); err != nil {
		return nil, fmt.Errorf("query-company-fact-by-ent-code: invalid payload: %w", err)
	}
	resp, err := h.svc.QueryCompanyFactByEntCode(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("query-company-fact-by-ent-code: %w", err)
	}
	return marshalResponse(resp)
}
