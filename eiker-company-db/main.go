// eiker-company-db is a Dapr-enabled atomic microservice that provides PostgreSQL
// CRUD operations for the company name domain and publishes a "company-fact-upserted"
// event via Dapr pub/sub after each successful fact write.
//
// Environment variables:
//
//	DATABASE_URL  – PostgreSQL DSN, e.g. "postgres://user:pass@localhost:5432/eiker_company"
//	APP_PORT      – gRPC port the Dapr sidecar will call into (default: 50002)
//	PUBSUB_NAME   – Dapr pub/sub component name (default: "eiker-pubsub")
//
// Start with the Dapr sidecar:
//
//	dapr run --app-id eiker-company-db --app-protocol grpc --app-port 50002 -- ./eiker-company-db

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	dapr "github.com/dapr/go-sdk/client"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xiaoyi-claude/eiker-company-db/handler"
	"github.com/xiaoyi-claude/eiker-company-db/repository"
	"github.com/xiaoyi-claude/eiker-company-db/service"
)

func main() {
	// ── 配置读取 ──────────────────────────────
	databaseURL := mustEnv("DATABASE_URL")
	port := envOrDefault("APP_PORT", "50002")
	pubsubName := envOrDefault("PUBSUB_NAME", "eiker-pubsub")

	// ── PostgreSQL 连接池 ──────────────────────
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("failed to create pgx pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}
	log.Printf("connected to PostgreSQL")

	// ── Dapr 客户端（用于 pub/sub 发布） ──────
	daprClient, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("failed to create Dapr client: %v", err)
	}
	defer daprClient.Close()

	// ── 依赖装配 ──────────────────────────────
	repo := repository.New(pool)
	svc := service.New(repo, daprClient, pubsubName)
	h := handler.New(svc)

	// ── Dapr gRPC 服务注册 ────────────────────
	s, err := daprd.NewService(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to create Dapr gRPC service on port %s: %v", port, err)
	}

	if err := s.AddServiceInvocationHandler("upsert-company-fact", h.UpsertCompanyFact); err != nil {
		log.Fatalf("failed to register upsert-company-fact: %v", err)
	}
	if err := s.AddServiceInvocationHandler("batch-upsert-company-fact", h.BatchUpsertCompanyFact); err != nil {
		log.Fatalf("failed to register batch-upsert-company-fact: %v", err)
	}
	if err := s.AddServiceInvocationHandler("query-company-fact-by-name", h.QueryCompanyFactByName); err != nil {
		log.Fatalf("failed to register query-company-fact-by-name: %v", err)
	}
	if err := s.AddServiceInvocationHandler("query-company-fact-by-code", h.QueryCompanyFactByCode); err != nil {
		log.Fatalf("failed to register query-company-fact-by-code: %v", err)
	}
	if err := s.AddServiceInvocationHandler("query-company-fact-by-ent-code", h.QueryCompanyFactByEntCode); err != nil {
		log.Fatalf("failed to register query-company-fact-by-ent-code: %v", err)
	}
	if err := s.AddServiceInvocationHandler("save-ocr-log", h.SaveOCRLog); err != nil {
		log.Fatalf("failed to register save-ocr-log: %v", err)
	}
	if err := s.AddServiceInvocationHandler("save-data-verify-log", h.SaveDataVerifyLog); err != nil {
		log.Fatalf("failed to register save-data-verify-log: %v", err)
	}
	if err := s.AddServiceInvocationHandler("save-conflict-log", h.SaveConflictLog); err != nil {
		log.Fatalf("failed to register save-conflict-log: %v", err)
	}
	if err := s.AddServiceInvocationHandler("save-first-link-log", h.SaveFirstLinkLog); err != nil {
		log.Fatalf("failed to register save-first-link-log: %v", err)
	}

	log.Printf("eiker-company-db listening on :%s (Dapr gRPC, pubsub=%s)", port, pubsubName)
	if err := s.Start(); err != nil {
		log.Fatalf("service stopped with error: %v", err)
	}
}

// mustEnv returns the value of the named environment variable or terminates the process.
func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %q is not set", key)
	}
	return v
}

// envOrDefault returns the value of the named environment variable,
// or defaultVal when the variable is empty or unset.
func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
