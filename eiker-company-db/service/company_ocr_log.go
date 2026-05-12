package service

import (
	"context"
	"fmt"

	"github.com/xiaoyi-claude/eiker-company-db/model"
)

// SaveOCRLog persists one OCR processing log entry and returns the generated primary key.
func (s *Service) SaveOCRLog(ctx context.Context, req *model.SaveOCRLogRequest) (*model.SaveOCRLogResponse, error) {
	id, err := s.repo.SaveOCRLog(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("save ocr log: %w", err)
	}
	return &model.SaveOCRLogResponse{ID: id}, nil
}
