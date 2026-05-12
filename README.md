# eiker-data

atomic/eiker-data/eiker-company-db/
├── go.mod                   # 模块声明，依赖 dapr/go-sdk v1.11, pgx/v5 v5.5.5, uuid v1.6                                                                                                                       
├── go.sum
├── main.go                  # 入口：初始化 pgxpool + Dapr 客户端，注册 9 个 handler
├── model/model.go           # 所有实体类型 + 9 个 API 的 Request/Response + pub/sub 事件
├── repository/repository.go # pgx/v5 CRUD（upsert/query company_fact + 4 种日志写入）
├── service/service.go       # 业务逻辑层 + upsert 后 publish company-fact-upserted
├── handler/handler.go       # 9 个 Dapr service invocation handler（JSON 解/序列化）
└── sql/V1__init.sql         # 6 张表建表 DDL（含全量字段注释 + 索引）