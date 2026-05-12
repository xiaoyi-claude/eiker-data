package service

import (
	"context"
	"fmt"

	"github.com/xiaoyi-claude/eiker-company-db/model"
)

// SaveFirstLinkLog persists one first-link verification log entry and returns the generated primary key.
func (s *Service) SaveFirstLinkLog(ctx context.Context, req *model.SaveFirstLinkLogRequest) (*model.SaveFirstLinkLogResponse, error) {
	id, err := s.repo.SaveFirstLinkLog(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("save first link log: %w", err)
	}
	return &model.SaveFirstLinkLogResponse{ID: id}, nil
}
