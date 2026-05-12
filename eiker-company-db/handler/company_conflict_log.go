package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dapr/go-sdk/service/common"
	"github.com/xiaoyi-claude/eiker-company-db/model"
)

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
