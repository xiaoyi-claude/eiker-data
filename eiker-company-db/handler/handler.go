// Package handler exposes all eiker-company-db business operations as Dapr
// service invocation handlers.  Each handler decodes the incoming JSON payload,
// delegates to the Service layer, and encodes the response back as JSON.
package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dapr/go-sdk/service/common"
	"github.com/xiaoyi-claude/eiker-company-db/model"
	"github.com/xiaoyi-claude/eiker-company-db/service"
)

// Handler bridges the Dapr gRPC transport layer with the Service layer.
type Handler struct {
	// svc is the business-logic service that all handlers delegate to.
	svc *service.Service
}

// New constructs a Handler wrapping the given Service.
func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// ─────────────────────────────────────────────
// company_fact 写操作
// ─────────────────────────────────────────────

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

// ─────────────────────────────────────────────
// company_fact 查操作
// ─────────────────────────────────────────────

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

// ─────────────────────────────────────────────
// 日志写操作
// ─────────────────────────────────────────────

// SaveOCRLog is the Dapr handler for the "save-ocr-log" method.
// Input:  JSON-encoded model.SaveOCRLogRequest
// Output: JSON-encoded model.SaveOCRLogResponse
func (h *Handler) SaveOCRLog(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	var req model.SaveOCRLogRequest
	if err := json.Unmarshal(in.Data, &req); err != nil {
		return nil, fmt.Errorf("save-ocr-log: invalid payload: %w", err)
	}
	resp, err := h.svc.SaveOCRLog(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("save-ocr-log: %w", err)
	}
	return marshalResponse(resp)
}

// SaveDataVerifyLog is the Dapr handler for the "save-data-verify-log" method.
// Input:  JSON-encoded model.SaveDataVerifyLogRequest
// Output: JSON-encoded model.SaveDataVerifyLogResponse
func (h *Handler) SaveDataVerifyLog(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	var req model.SaveDataVerifyLogRequest
	if err := json.Unmarshal(in.Data, &req); err != nil {
		return nil, fmt.Errorf("save-data-verify-log: invalid payload: %w", err)
	}
	resp, err := h.svc.SaveDataVerifyLog(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("save-data-verify-log: %w", err)
	}
	return marshalResponse(resp)
}

// SaveConflictLog is the Dapr handler for the "save-conflict-log" method.
// Input:  JSON-encoded model.SaveConflictLogRequest
// Output: JSON-encoded model.SaveConflictLogResponse
func (h *Handler) SaveConflictLog(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	var req model.SaveConflictLogRequest
	if err := json.Unmarshal(in.Data, &req); err != nil {
		return nil, fmt.Errorf("save-conflict-log: invalid payload: %w", err)
	}
	resp, err := h.svc.SaveConflictLog(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("save-conflict-log: %w", err)
	}
	return marshalResponse(resp)
}

// SaveFirstLinkLog is the Dapr handler for the "save-first-link-log" method.
// Input:  JSON-encoded model.SaveFirstLinkLogRequest
// Output: JSON-encoded model.SaveFirstLinkLogResponse
func (h *Handler) SaveFirstLinkLog(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	var req model.SaveFirstLinkLogRequest
	if err := json.Unmarshal(in.Data, &req); err != nil {
		return nil, fmt.Errorf("save-first-link-log: invalid payload: %w", err)
	}
	resp, err := h.svc.SaveFirstLinkLog(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("save-first-link-log: %w", err)
	}
	return marshalResponse(resp)
}

// ─────────────────────────────────────────────
// 内部辅助
// ─────────────────────────────────────────────

// marshalResponse serialises v to JSON and wraps it in a Dapr Content object.
func marshalResponse(v interface{}) (*common.Content, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal response: %w", err)
	}
	return &common.Content{ContentType: "application/json", Data: data}, nil
}
