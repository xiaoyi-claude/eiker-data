package service

import (
	"context"
	"fmt"

	"github.com/xiaoyi-claude/eiker-company-db/model"
)

// SaveConflictLog persists one conflict-handling log entry and returns the generated primary key.
func (s *Service) SaveConflictLog(ctx context.Context, req *model.SaveConflictLogRequest) (*model.SaveConflictLogResponse, error) {
	id, err := s.repo.SaveConflictLog(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("save conflict log: %w", err)
	}
	return &model.SaveConflictLogResponse{ID: id}, nil
}
