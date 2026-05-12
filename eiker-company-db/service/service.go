// Package service implements the business logic for the eiker-company-db service.
// After every successful company_fact upsert it publishes a "company-fact-upserted"
// event containing the full CompanyFact to the configured Dapr pub/sub component.
package service

import (
	"context"
	"log"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/xiaoyi-claude/eiker-company-db/model"
	"github.com/xiaoyi-claude/eiker-company-db/repository"
)

// Service coordinates data access via Repository and event publishing via the Dapr client.
type Service struct {
	// repo provides all PostgreSQL CRUD operations.
	repo *repository.Repository
	// daprClient is the Dapr SDK client used for pub/sub event publishing.
	daprClient dapr.Client
	// pubsubName is the Dapr pub/sub component name (e.g. "eiker-pubsub").
	pubsubName string
}

// New constructs a Service with the given repository, Dapr client, and pub/sub component name.
func New(repo *repository.Repository, daprClient dapr.Client, pubsubName string) *Service {
	return &Service{
		repo:       repo,
		daprClient: daprClient,
		pubsubName: pubsubName,
	}
}

// publishFactUpserted publishes a "company-fact-upserted" event to the Dapr pub/sub component.
// The full CompanyFact is included in the payload so downstream consumers need no additional DB query.
// Publish failures are only logged; the DB write is already committed at this point.
func (s *Service) publishFactUpserted(ctx context.Context, fact *model.CompanyFact, action string) {
	event := &model.CompanyFactUpsertedEvent{
		Action: action,
		Fact:   fact,
	}
	if err := s.daprClient.PublishEvent(ctx, s.pubsubName, "company-fact-upserted", event); err != nil {
		log.Printf("warn: publish company-fact-upserted (ent_code=%s fact_id=%d): %v",
			fact.EntCode, fact.ID, err)
	}
}
