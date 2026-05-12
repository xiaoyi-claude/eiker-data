# eiker-company-db

企业名称域原子微服务，通过 Dapr gRPC 服务调用暴露 PostgreSQL CRUD 操作，并在企业事实写入成功后通过 Dapr pub/sub 发布 `company-fact-upserted` 事件。

## 架构

```
调用方 ──Dapr sidecar──► eiker-company-db (gRPC :50002)
                                │
                         ┌──────┴──────┐
                       Handler      pub/sub
                         │        (company-fact-upserted)
                       Service
                         │
                      Repository
                         │
                     PostgreSQL (eiker_company)
```

## 数据库表

| 表名 | 说明 |
|------|------|
| company_fact | 企业事实信息（主表） |
| company_task | 企业采集任务 |
| company_ocr_log | OCR 识别日志 |
| company_data_verify_log | 数据校验日志 |
| company_conflict_log | 冲突处理日志 |
| company_first_link_log | 首链核验日志 |

## 环境变量

| 变量 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `DATABASE_URL` | 是 | — | PostgreSQL DSN，例如 `postgres://eiker_company:xxx@localhost:5432/eiker_company` |
| `APP_PORT` | 否 | `50002` | Dapr sidecar 回调的 gRPC 端口 |
| `PUBSUB_NAME` | 否 | `eiker-pubsub` | Dapr pub/sub 组件名称 |

## 构建

```bash
cd eiker-company-db
go mod tidy
go build -o eiker-company-db.exe .
```

## 启动

```bash
set DATABASE_URL=postgres://eiker_company:<password>@localhost:5432/eiker_company
dapr run --app-id eiker-company-db --app-protocol grpc --app-port 50002 -- .\eiker-company-db.exe
```

## Dapr 服务调用方法

| 方法名 | 说明 | 请求类型 | 响应类型 |
|--------|------|----------|----------|
| `upsert-company-fact` | 新增或更新企业事实（按 ID 区分） | `UpsertCompanyFactRequest` | `UpsertCompanyFactResponse` |
| `batch-upsert-company-fact` | 批量新增或更新企业事实 | `BatchUpsertCompanyFactRequest` | `BatchUpsertCompanyFactResponse` |
| `query-company-fact-by-name` | 按企业名称查询事实列表 | `QueryCompanyFactByNameRequest` | `QueryCompanyFactByNameResponse` |
| `query-company-fact-by-code` | 按统一社会信用代码查询事实 | `QueryCompanyFactByCodeRequest` | `QueryCompanyFactByCodeResponse` |
| `query-company-fact-by-ent-code` | 按内部企业编码查询事实列表 | `QueryCompanyFactByEntCodeRequest` | `QueryCompanyFactByEntCodeResponse` |
| `save-ocr-log` | 保存 OCR 识别日志 | `SaveOCRLogRequest` | `SaveOCRLogResponse` |
| `save-data-verify-log` | 保存数据校验日志 | `SaveDataVerifyLogRequest` | `SaveDataVerifyLogResponse` |
| `save-conflict-log` | 保存冲突处理日志 | `SaveConflictLogRequest` | `SaveConflictLogResponse` |
| `save-first-link-log` | 保存首链核验日志 | `SaveFirstLinkLogRequest` | `SaveFirstLinkLogResponse` |

## Pub/Sub 事件

成功执行 `upsert-company-fact` 或 `batch-upsert-company-fact` 后，向 topic `company-fact-upserted` 发布事件：

```json
{
  "action": "create",
  "fact": {
    "id": 1,
    "ent_code": "ENT-xxx",
    "credit_code": "91xxxx",
    "name": "示例企业",
    ...
  }
}
```

`action` 取值：`"create"` 或 `"update"`。

## 测试客户端

```bash
# 先启动服务端（见上方"启动"章节），再运行客户端
cd cmd/client
go run main.go
```

客户端会依次调用全部 9 个 API，并打印响应结果。

## 目录结构

```
eiker-company-db/
├── main.go                        # 程序入口，依赖装配与 Dapr 注册
├── go.mod
├── go.sum
├── sql/
│   └── V1__init.sql               # 建表 DDL（含注释）
├── model/
│   ├── common.go                  # CommonFields 公共字段
│   ├── event.go                   # CompanyFactUpsertedEvent
│   ├── company_fact.go
│   ├── company_task.go
│   ├── company_ocr_log.go
│   ├── company_data_verify_log.go
│   ├── company_conflict_log.go
│   └── company_first_link_log.go
├── repository/
│   ├── repository.go              # Repository 结构体与辅助函数
│   ├── company_fact.go
│   ├── company_ocr_log.go
│   ├── company_data_verify_log.go
│   ├── company_conflict_log.go
│   └── company_first_link_log.go
├── service/
│   ├── service.go                 # Service 结构体与 pub/sub 发布
│   ├── company_fact.go
│   ├── company_ocr_log.go
│   ├── company_data_verify_log.go
│   ├── company_conflict_log.go
│   └── company_first_link_log.go
├── handler/
│   ├── handler.go                 # Handler 结构体与 marshalResponse
│   ├── company_fact.go
│   ├── company_ocr_log.go
│   ├── company_data_verify_log.go
│   ├── company_conflict_log.go
│   └── company_first_link_log.go
└── cmd/
    └── client/
        └── main.go                # 测试客户端
```
