// Package main is a test client that exercises all 9 Dapr service invocation
// methods exposed by eiker-company-db.
//
// Prerequisites:
//  1. eiker-company-db is running under its Dapr sidecar (dapr run ...).
//  2. The Dapr sidecar's gRPC port is reachable at DAPR_GRPC_PORT (default 50001).
//
// Run:
//
//	set DAPR_GRPC_PORT=50001
//	go run ./cmd/client/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/xiaoyi-claude/eiker-company-db/model"
)

const (
	appID = "eiker-company-db"
	verb  = "post"
)

func main() {
	client, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("create Dapr client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// ── 1. upsert-company-fact (create) ──────────────────────────────────────
	fmt.Println("=== 1. upsert-company-fact (create) ===")
	upsertReq := model.UpsertCompanyFactRequest{
		CreditCode:     "91440300MA5DTEST01",
		Name:           "测试科技有限公司",
		LegalRep:       "张三",
		EntRecordScene: 1,
		TraceID:        "trace-client-001",
		Creater:        "client-test",
	}
	var upsertResp model.UpsertCompanyFactResponse
	mustInvoke(ctx, client, "upsert-company-fact", upsertReq, &upsertResp)
	fmt.Printf("  → action=%s  ent_code=%s  fact_id=%d\n\n", upsertResp.Action, upsertResp.EntCode, upsertResp.FactID)

	entCode := upsertResp.EntCode

	// ── 2. upsert-company-fact (update) ──────────────────────────────────────
	fmt.Println("=== 2. upsert-company-fact (update) ===")
	updateReq := model.UpsertCompanyFactRequest{
		EntCode:        entCode,
		CreditCode:     "91440300MA5DTEST01",
		Name:           "测试科技有限公司（已更新）",
		LegalRep:       "李四",
		EntRecordScene: 1,
		TraceID:        "trace-client-002",
		Creater:        "client-test",
	}
	var updateResp model.UpsertCompanyFactResponse
	mustInvoke(ctx, client, "upsert-company-fact", updateReq, &updateResp)
	fmt.Printf("  → action=%s  ent_code=%s  fact_id=%d\n\n", updateResp.Action, updateResp.EntCode, updateResp.FactID)

	// ── 3. batch-upsert-company-fact ─────────────────────────────────────────
	fmt.Println("=== 3. batch-upsert-company-fact ===")
	batchReq := model.BatchUpsertCompanyFactRequest{
		Items: []model.UpsertCompanyFactRequest{
			{
				CreditCode:     "91440300MA5DBATCH1",
				Name:           "批量企业甲有限公司",
				EntRecordScene: 2,
				Creater:        "client-test",
			},
			{
				CreditCode:     "91440300MA5DBATCH2",
				Name:           "批量企业乙有限公司",
				EntRecordScene: 2,
				Creater:        "client-test",
			},
		},
	}
	var batchResp model.BatchUpsertCompanyFactResponse
	mustInvoke(ctx, client, "batch-upsert-company-fact", batchReq, &batchResp)
	fmt.Printf("  → success=%d  fail=%d\n", batchResp.SuccessCount, batchResp.FailCount)
	for i, r := range batchResp.Results {
		fmt.Printf("     [%d] action=%s  ent_code=%s  fact_id=%d\n", i, r.Action, r.EntCode, r.FactID)
	}
	fmt.Println()

	// ── 4. query-company-fact-by-name ─────────────────────────────────────────
	fmt.Println("=== 4. query-company-fact-by-name ===")
	var byNameResp model.QueryCompanyFactByNameResponse
	mustInvoke(ctx, client, "query-company-fact-by-name",
		model.QueryCompanyFactByNameRequest{Name: "测试科技有限公司（已更新）"},
		&byNameResp)
	fmt.Printf("  → found %d items\n", len(byNameResp.Items))
	for _, f := range byNameResp.Items {
		fmt.Printf("     ent_code=%s  name=%s  legal_rep=%s\n", f.EntCode, f.Name, f.LegalRep)
	}
	fmt.Println()

	// ── 5. query-company-fact-by-code ────────────────────────────────────────
	fmt.Println("=== 5. query-company-fact-by-code ===")
	var byCodeResp model.QueryCompanyFactByCodeResponse
	mustInvoke(ctx, client, "query-company-fact-by-code",
		model.QueryCompanyFactByCodeRequest{CreditCode: "91440300MA5DTEST01"},
		&byCodeResp)
	if byCodeResp.Item != nil {
		fmt.Printf("  → ent_code=%s  name=%s\n\n", byCodeResp.Item.EntCode, byCodeResp.Item.Name)
	} else {
		fmt.Println("  → not found\n")
	}

	// ── 6. query-company-fact-by-ent-code ────────────────────────────────────
	fmt.Println("=== 6. query-company-fact-by-ent-code ===")
	var byEntCodeResp model.QueryCompanyFactByEntCodeResponse
	mustInvoke(ctx, client, "query-company-fact-by-ent-code",
		model.QueryCompanyFactByEntCodeRequest{EntCode: entCode},
		&byEntCodeResp)
	if byEntCodeResp.Item != nil {
		fmt.Printf("  → id=%d  ent_code=%s  name=%s\n\n", byEntCodeResp.Item.ID, byEntCodeResp.Item.EntCode, byEntCodeResp.Item.Name)
	} else {
		fmt.Println("  → not found\n")
	}

	// ── 7. save-ocr-log ───────────────────────────────────────────────────────
	fmt.Println("=== 7. save-ocr-log ===")
	var ocrResp model.SaveOCRLogResponse
	mustInvoke(ctx, client, "save-ocr-log", model.SaveOCRLogRequest{
		TaskID:      "task-ocr-001",
		ModelName:   "paddleocr",
		ImagePath:   "/images/sample.jpg",
		CreditCode:  "91440300MA5DTEST01",
		CompanyName: "测试科技有限公司",
		IsSuccess:   true,
		TraceID:     "trace-client-007",
		Creater:     "client-test",
		RawResult:   map[string]interface{}{"confidence": 0.98},
	}, &ocrResp)
	fmt.Printf("  → id=%d\n\n", ocrResp.ID)

	// ── 8. save-data-verify-log ───────────────────────────────────────────────
	fmt.Println("=== 8. save-data-verify-log ===")
	isConsistent := true
	var dvResp model.SaveDataVerifyLogResponse
	mustInvoke(ctx, client, "save-data-verify-log", model.SaveDataVerifyLogRequest{
		TaskID:       "task-verify-001",
		SourceName:   "tianyancha",
		QueryName:    "测试科技有限公司",
		IsConsistent: &isConsistent,
		IsSuccess:    true,
		TraceID:      "trace-client-008",
		Creater:      "client-test",
		ResponseData: map[string]interface{}{"status": "ok"},
	}, &dvResp)
	fmt.Printf("  → id=%d\n\n", dvResp.ID)

	// ── 9. save-conflict-log ──────────────────────────────────────────────────
	fmt.Println("=== 9. save-conflict-log ===")
	var conflictResp model.SaveConflictLogResponse
	mustInvoke(ctx, client, "save-conflict-log", model.SaveConflictLogRequest{
		TaskID:           "task-conflict-001",
		ConflictType:     1,
		NewCreditCode:    "91440300MA5DNEW001",
		NewName:          "同名企业新公司",
		ExistingEntCode:  entCode,
		ExistingName:     "测试科技有限公司（已更新）",
		Resolution:       1,
		ResolutionReason: "两家独立企业，保留各自记录",
		TraceID:          "trace-client-009",
		Creater:          "client-test",
	}, &conflictResp)
	fmt.Printf("  → id=%d\n\n", conflictResp.ID)

	// ── 10. save-first-link-log ───────────────────────────────────────────────
	fmt.Println("=== 10. save-first-link-log ===")
	now := time.Now()
	isOk := true
	var firstLinkResp model.SaveFirstLinkLogResponse
	mustInvoke(ctx, client, "save-first-link-log", model.SaveFirstLinkLogRequest{
		TaskID:       "task-firstlink-001",
		FactID:       upsertResp.FactID,
		VerifierType: 1,
		ExpectedData: map[string]interface{}{"name": "测试科技有限公司"},
		ActualData:   map[string]interface{}{"name": "测试科技有限公司（已更新）"},
		IsConsistent: &isOk,
		VerifyStatus: 1,
		VerifiedAt:   &now,
		TraceID:      "trace-client-010",
		Creater:      "client-test",
	}, &firstLinkResp)
	fmt.Printf("  → id=%d\n\n", firstLinkResp.ID)

	fmt.Println("=== all 9 APIs called successfully ===")
}

// mustInvoke marshals req, calls the named Dapr method on appID, and unmarshals the response into out.
func mustInvoke(ctx context.Context, client dapr.Client, method string, req, out interface{}) {
	data, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("[%s] marshal request: %v", method, err)
	}

	resp, err := client.InvokeMethodWithContent(ctx, appID, method, verb, &dapr.DataContent{
		ContentType: "application/json",
		Data:        data,
	})
	if err != nil {
		log.Fatalf("[%s] invoke: %v", method, err)
	}

	if err := json.Unmarshal(resp, out); err != nil {
		log.Fatalf("[%s] unmarshal response: %v", method, err)
	}
}
