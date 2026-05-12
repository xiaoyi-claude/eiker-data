package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dapr/go-sdk/service/common"
	"github.com/xiaoyi-claude/eiker-company-db/model"
)

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
