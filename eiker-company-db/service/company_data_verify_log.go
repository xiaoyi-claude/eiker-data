package service

import (
	"context"
	"fmt"

	"github.com/xiaoyi-claude/eiker-company-db/model"
)

// SaveDataVerifyLog persists one third-party data verification log entry and returns the generated primary key.
func (s *Service) SaveDataVerifyLog(ctx context.Context, req *model.SaveDataVerifyLogRequest) (*model.SaveDataVerifyLogResponse, error) {
	id, err := s.repo.SaveDataVerifyLog(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("save data verify log: %w", err)
	}
	return &model.SaveDataVerifyLogResponse{ID: id}, nil
}
