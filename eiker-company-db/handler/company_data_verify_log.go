package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dapr/go-sdk/service/common"
	"github.com/xiaoyi-claude/eiker-company-db/model"
)

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
