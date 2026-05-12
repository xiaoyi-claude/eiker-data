// Package handler exposes all eiker-company-db operations as Dapr service invocation handlers.
// Each handler decodes the incoming JSON payload, delegates to the Service layer,
// and encodes the response back as JSON.
package handler

import (
	"encoding/json"
	"fmt"

	"github.com/dapr/go-sdk/service/common"
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

// marshalResponse serialises v to JSON and wraps it in a Dapr Content object.
func marshalResponse(v interface{}) (*common.Content, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal response: %w", err)
	}
	return &common.Content{ContentType: "application/json", Data: data}, nil
}
