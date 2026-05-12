package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dapr/go-sdk/service/common"
	"github.com/xiaoyi-claude/eiker-company-db/model"
)

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
