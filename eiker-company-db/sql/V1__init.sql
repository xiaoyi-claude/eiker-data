-- V1__init.sql
-- 企业事实数据服务（eiker-company-db）数据库初始化脚本
-- 包含企业名称域全部数据表及注释

-- ============================================================
-- 1. 企业事实数据表 company_fact
-- ============================================================
CREATE TABLE IF NOT EXISTS company_fact (
    -- 公共字段 --
    id                      BIGSERIAL       PRIMARY KEY,
    is_current              BOOLEAN         NOT NULL DEFAULT TRUE,
    trace_id                VARCHAR(64),
    owner                   VARCHAR(64),
    creater                 VARCHAR(64),
    create_time             TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updater                 VARCHAR(64),
    update_time             TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    remark                  VARCHAR(500),
    -- 业务字段 --
    ent_code                VARCHAR(36)     NOT NULL,
    credit_code             VARCHAR(18)     NOT NULL,
    name                    VARCHAR(100)    NOT NULL,
    legal_rep               VARCHAR(20),
    location_id             VARCHAR(36),
    fact_time               TIMESTAMPTZ,
    fact_record_time        TIMESTAMPTZ,
    ent_record_scene        SMALLINT,
    legal_rep_record_scene  SMALLINT,
    task_id                 VARCHAR(36)
);

COMMENT ON TABLE company_fact IS '企业事实数据表，存储企业名称域核心事实信息';
COMMENT ON COLUMN company_fact.id IS '自增主键';
COMMENT ON COLUMN company_fact.is_current IS '是否当前有效版本（用于多版本数据）';
COMMENT ON COLUMN company_fact.trace_id IS '日志链路追踪ID';
COMMENT ON COLUMN company_fact.owner IS '数据拥有人（用于数据权限控制）';
COMMENT ON COLUMN company_fact.creater IS '创建人';
COMMENT ON COLUMN company_fact.create_time IS '创建时间（DB insert 时自动写入）';
COMMENT ON COLUMN company_fact.updater IS '修改人';
COMMENT ON COLUMN company_fact.update_time IS '修改时间（DB update 时自动写入）';
COMMENT ON COLUMN company_fact.remark IS '备注';
COMMENT ON COLUMN company_fact.ent_code IS '企业 UUID v4（对外标识，公域使用）';
COMMENT ON COLUMN company_fact.credit_code IS '统一社会信用代码（GB 32100-2015）';
COMMENT ON COLUMN company_fact.name IS '企业名称（4-26字，汉字或汉字+括号）';
COMMENT ON COLUMN company_fact.legal_rep IS '权力人姓名（2-5字，纯汉字）';
COMMENT ON COLUMN company_fact.location_id IS '关联 eiker-address-db 中的 address_id';
COMMENT ON COLUMN company_fact.fact_time IS '事实发生时间';
COMMENT ON COLUMN company_fact.fact_record_time IS '入巢时间';
COMMENT ON COLUMN company_fact.ent_record_scene IS '企业新增场景：1=基础拉新；2=用户拉新；3=业务拉新；4=公司拉新';
COMMENT ON COLUMN company_fact.legal_rep_record_scene IS '权力人新增场景：1=代为确定性主张；2=确定实际主张；3=权力人替换并主张';
COMMENT ON COLUMN company_fact.task_id IS '关联任务单 UUID';

CREATE UNIQUE INDEX IF NOT EXISTS uidx_company_fact_credit_code ON company_fact (credit_code) WHERE is_current = TRUE;
CREATE UNIQUE INDEX IF NOT EXISTS uidx_company_fact_ent_code ON company_fact (ent_code) WHERE is_current = TRUE;
CREATE INDEX IF NOT EXISTS idx_company_fact_name ON company_fact (name);
CREATE INDEX IF NOT EXISTS idx_company_fact_task_id ON company_fact (task_id);

-- ============================================================
-- 2. 任务单表 company_task
-- ============================================================
CREATE TABLE IF NOT EXISTS company_task (
    -- 公共字段 --
    id                      BIGSERIAL       PRIMARY KEY,
    is_current              BOOLEAN         NOT NULL DEFAULT TRUE,
    trace_id                VARCHAR(64),
    owner                   VARCHAR(64),
    creater                 VARCHAR(64),
    create_time             TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updater                 VARCHAR(64),
    update_time             TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    remark                  VARCHAR(500),
    -- 业务字段 --
    task_id                 VARCHAR(36)     NOT NULL,
    task_type               SMALLINT        NOT NULL,
    task_status             SMALLINT        NOT NULL DEFAULT 0,
    input_source            VARCHAR(500),
    payload                 JSONB,
    result                  JSONB,
    workflow_instance_id    VARCHAR(64),
    retry_count             SMALLINT        NOT NULL DEFAULT 0,
    error_msg               TEXT
);

COMMENT ON TABLE company_task IS '企业拉新任务单表，由 business 层 eiker-company-update 管理，atomic 层提供 CRUD';
COMMENT ON COLUMN company_task.id IS '自增主键';
COMMENT ON COLUMN company_task.is_current IS '是否当前有效版本';
COMMENT ON COLUMN company_task.trace_id IS '日志链路追踪ID';
COMMENT ON COLUMN company_task.owner IS '数据拥有人';
COMMENT ON COLUMN company_task.creater IS '创建人';
COMMENT ON COLUMN company_task.create_time IS '创建时间';
COMMENT ON COLUMN company_task.updater IS '修改人';
COMMENT ON COLUMN company_task.update_time IS '修改时间';
COMMENT ON COLUMN company_task.remark IS '备注';
COMMENT ON COLUMN company_task.task_id IS '任务单 UUID';
COMMENT ON COLUMN company_task.task_type IS '任务类型：1=基础拉新；2=用户拉新；3=业务拉新；4=公司拉新';
COMMENT ON COLUMN company_task.task_status IS '任务状态：0=PENDING；1=IN_PROGRESS；2=SUCCESS；3=FAILED';
COMMENT ON COLUMN company_task.input_source IS '输入来源标识（文件路径/图片路径/数据标识）';
COMMENT ON COLUMN company_task.payload IS '采信来源、凭证、处理规则（JSONB）';
COMMENT ON COLUMN company_task.result IS '处理结果（JSONB）';
COMMENT ON COLUMN company_task.workflow_instance_id IS 'Dapr Workflow 实例ID';
COMMENT ON COLUMN company_task.retry_count IS '重试次数';
COMMENT ON COLUMN company_task.error_msg IS '失败原因';

CREATE UNIQUE INDEX IF NOT EXISTS uidx_company_task_task_id ON company_task (task_id);
CREATE INDEX IF NOT EXISTS idx_company_task_status ON company_task (task_status);

-- ============================================================
-- 3. OCR 处理日志表 company_ocr_log
-- ============================================================
CREATE TABLE IF NOT EXISTS company_ocr_log (
    -- 公共字段 --
    id                      BIGSERIAL       PRIMARY KEY,
    is_current              BOOLEAN         NOT NULL DEFAULT TRUE,
    trace_id                VARCHAR(64),
    owner                   VARCHAR(64),
    creater                 VARCHAR(64),
    create_time             TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updater                 VARCHAR(64),
    update_time             TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    remark                  VARCHAR(500),
    -- 业务字段 --
    task_id                 VARCHAR(36)     NOT NULL,
    model_name              VARCHAR(50)     NOT NULL,
    image_path              VARCHAR(500),
    raw_result              JSONB,
    credit_code             VARCHAR(18),
    company_name            VARCHAR(100),
    legal_rep_name          VARCHAR(20),
    is_success              BOOLEAN         NOT NULL DEFAULT FALSE,
    fail_reason             VARCHAR(200)
);

COMMENT ON TABLE company_ocr_log IS 'OCR 处理日志表，记录每次 OCR 识别的原始输出和提取结果';
COMMENT ON COLUMN company_ocr_log.id IS '自增主键';
COMMENT ON COLUMN company_ocr_log.is_current IS '是否当前有效版本';
COMMENT ON COLUMN company_ocr_log.trace_id IS '日志链路追踪ID';
COMMENT ON COLUMN company_ocr_log.owner IS '数据拥有人';
COMMENT ON COLUMN company_ocr_log.creater IS '创建人';
COMMENT ON COLUMN company_ocr_log.create_time IS '创建时间';
COMMENT ON COLUMN company_ocr_log.updater IS '修改人';
COMMENT ON COLUMN company_ocr_log.update_time IS '修改时间';
COMMENT ON COLUMN company_ocr_log.remark IS '备注';
COMMENT ON COLUMN company_ocr_log.task_id IS '关联任务单 UUID';
COMMENT ON COLUMN company_ocr_log.model_name IS 'OCR 模型名称（paddleocr / easyocr / tesseract）';
COMMENT ON COLUMN company_ocr_log.image_path IS '图片路径或 URL';
COMMENT ON COLUMN company_ocr_log.raw_result IS 'OCR 原始输出（JSONB）';
COMMENT ON COLUMN company_ocr_log.credit_code IS '提取的统一社会信用代码';
COMMENT ON COLUMN company_ocr_log.company_name IS '提取的企业名称';
COMMENT ON COLUMN company_ocr_log.legal_rep_name IS '提取的权力人姓名';
COMMENT ON COLUMN company_ocr_log.is_success IS '提取是否成功';
COMMENT ON COLUMN company_ocr_log.fail_reason IS '失败原因';

CREATE INDEX IF NOT EXISTS idx_company_ocr_log_task_id ON company_ocr_log (task_id);

-- ============================================================
-- 4. 第三方数据源验证日志表 company_data_verify_log
-- ============================================================
CREATE TABLE IF NOT EXISTS company_data_verify_log (
    -- 公共字段 --
    id                      BIGSERIAL       PRIMARY KEY,
    is_current              BOOLEAN         NOT NULL DEFAULT TRUE,
    trace_id                VARCHAR(64),
    owner                   VARCHAR(64),
    creater                 VARCHAR(64),
    create_time             TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updater                 VARCHAR(64),
    update_time             TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    remark                  VARCHAR(500),
    -- 业务字段 --
    task_id                 VARCHAR(36)     NOT NULL,
    source_name             VARCHAR(50)     NOT NULL,
    query_name              VARCHAR(100),
    response_data           JSONB,
    is_consistent           BOOLEAN,
    is_success              BOOLEAN         NOT NULL DEFAULT FALSE,
    fail_reason             VARCHAR(200)
);

COMMENT ON TABLE company_data_verify_log IS '第三方数据源验证日志表，记录天眼查等数据源的查询与比对结果';
COMMENT ON COLUMN company_data_verify_log.id IS '自增主键';
COMMENT ON COLUMN company_data_verify_log.is_current IS '是否当前有效版本';
COMMENT ON COLUMN company_data_verify_log.trace_id IS '日志链路追踪ID';
COMMENT ON COLUMN company_data_verify_log.owner IS '数据拥有人';
COMMENT ON COLUMN company_data_verify_log.creater IS '创建人';
COMMENT ON COLUMN company_data_verify_log.create_time IS '创建时间';
COMMENT ON COLUMN company_data_verify_log.updater IS '修改人';
COMMENT ON COLUMN company_data_verify_log.update_time IS '修改时间';
COMMENT ON COLUMN company_data_verify_log.remark IS '备注';
COMMENT ON COLUMN company_data_verify_log.task_id IS '关联任务单 UUID';
COMMENT ON COLUMN company_data_verify_log.source_name IS '数据源名称（tianyancha 等）';
COMMENT ON COLUMN company_data_verify_log.query_name IS '查询企业名称';
COMMENT ON COLUMN company_data_verify_log.response_data IS '数据源返回结果（JSONB）';
COMMENT ON COLUMN company_data_verify_log.is_consistent IS '与输入三要素是否一致';
COMMENT ON COLUMN company_data_verify_log.is_success IS 'API 调用是否成功';
COMMENT ON COLUMN company_data_verify_log.fail_reason IS '失败原因';

CREATE INDEX IF NOT EXISTS idx_company_data_verify_log_task_id ON company_data_verify_log (task_id);

-- ============================================================
-- 5. 冲突处理日志表 company_conflict_log
-- ============================================================
CREATE TABLE IF NOT EXISTS company_conflict_log (
    -- 公共字段 --
    id                      BIGSERIAL       PRIMARY KEY,
    is_current              BOOLEAN         NOT NULL DEFAULT TRUE,
    trace_id                VARCHAR(64),
    owner                   VARCHAR(64),
    creater                 VARCHAR(64),
    create_time             TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updater                 VARCHAR(64),
    update_time             TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    remark                  VARCHAR(500),
    -- 业务字段 --
    task_id                 VARCHAR(36)     NOT NULL,
    conflict_type           SMALLINT        NOT NULL,
    new_credit_code         VARCHAR(18),
    new_name                VARCHAR(100),
    existing_ent_code       VARCHAR(36),
    existing_credit_code    VARCHAR(18),
    existing_name           VARCHAR(100),
    resolution              SMALLINT,
    resolution_reason       VARCHAR(200)
);

COMMENT ON TABLE company_conflict_log IS '冲突处理日志表，记录拉新过程中企业名称/信用代码冲突及处理方式';
COMMENT ON COLUMN company_conflict_log.id IS '自增主键';
COMMENT ON COLUMN company_conflict_log.is_current IS '是否当前有效版本';
COMMENT ON COLUMN company_conflict_log.trace_id IS '日志链路追踪ID';
COMMENT ON COLUMN company_conflict_log.owner IS '数据拥有人';
COMMENT ON COLUMN company_conflict_log.creater IS '创建人';
COMMENT ON COLUMN company_conflict_log.create_time IS '创建时间';
COMMENT ON COLUMN company_conflict_log.updater IS '修改人';
COMMENT ON COLUMN company_conflict_log.update_time IS '修改时间';
COMMENT ON COLUMN company_conflict_log.remark IS '备注';
COMMENT ON COLUMN company_conflict_log.task_id IS '关联任务单 UUID';
COMMENT ON COLUMN company_conflict_log.conflict_type IS '冲突类型：1=同名不同企；2=同企不同名';
COMMENT ON COLUMN company_conflict_log.new_credit_code IS '新记录统一社会信用代码';
COMMENT ON COLUMN company_conflict_log.new_name IS '新记录企业名称';
COMMENT ON COLUMN company_conflict_log.existing_ent_code IS '已有企业 UUID';
COMMENT ON COLUMN company_conflict_log.existing_credit_code IS '已有统一社会信用代码';
COMMENT ON COLUMN company_conflict_log.existing_name IS '已有企业名称';
COMMENT ON COLUMN company_conflict_log.resolution IS '处理方式：1=全部保留独立存在；2=复用UUID挂载新名称；3=拒绝拉新';
COMMENT ON COLUMN company_conflict_log.resolution_reason IS '处理原因';

CREATE INDEX IF NOT EXISTS idx_company_conflict_log_task_id ON company_conflict_log (task_id);

-- ============================================================
-- 6. 首支验证日志表 company_first_link_log
-- ============================================================
CREATE TABLE IF NOT EXISTS company_first_link_log (
    -- 公共字段 --
    id                      BIGSERIAL       PRIMARY KEY,
    is_current              BOOLEAN         NOT NULL DEFAULT TRUE,
    trace_id                VARCHAR(64),
    owner                   VARCHAR(64),
    creater                 VARCHAR(64),
    create_time             TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updater                 VARCHAR(64),
    update_time             TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    remark                  VARCHAR(500),
    -- 业务字段 --
    task_id                 VARCHAR(36)     NOT NULL,
    fact_id                 BIGINT          NOT NULL,
    verifier_type           SMALLINT        NOT NULL,
    expected_data           JSONB,
    actual_data             JSONB,
    is_consistent           BOOLEAN,
    verify_status           SMALLINT        NOT NULL DEFAULT 0,
    verified_at             TIMESTAMPTZ
);

COMMENT ON TABLE company_first_link_log IS '首支验证日志表，记录企业事实数据写入前后的一致性验证结果';
COMMENT ON COLUMN company_first_link_log.id IS '自增主键';
COMMENT ON COLUMN company_first_link_log.is_current IS '是否当前有效版本';
COMMENT ON COLUMN company_first_link_log.trace_id IS '日志链路追踪ID';
COMMENT ON COLUMN company_first_link_log.owner IS '数据拥有人';
COMMENT ON COLUMN company_first_link_log.creater IS '创建人';
COMMENT ON COLUMN company_first_link_log.create_time IS '创建时间';
COMMENT ON COLUMN company_first_link_log.updater IS '修改人';
COMMENT ON COLUMN company_first_link_log.update_time IS '修改时间';
COMMENT ON COLUMN company_first_link_log.remark IS '备注';
COMMENT ON COLUMN company_first_link_log.task_id IS '关联任务单 UUID';
COMMENT ON COLUMN company_first_link_log.fact_id IS '关联 company_fact.id';
COMMENT ON COLUMN company_first_link_log.verifier_type IS '验证人类型：1=系统自验（基础/公司拉新）；2=用户确认；3=业务方确认';
COMMENT ON COLUMN company_first_link_log.expected_data IS '写入前信息，用于比对（JSONB）';
COMMENT ON COLUMN company_first_link_log.actual_data IS '写入后信息（JSONB）';
COMMENT ON COLUMN company_first_link_log.is_consistent IS '写入前后是否一致';
COMMENT ON COLUMN company_first_link_log.verify_status IS '验证状态：0=待确认；1=已通过；2=已拒绝';
COMMENT ON COLUMN company_first_link_log.verified_at IS '确认时间';

CREATE INDEX IF NOT EXISTS idx_company_first_link_log_task_id ON company_first_link_log (task_id);
CREATE INDEX IF NOT EXISTS idx_company_first_link_log_fact_id ON company_first_link_log (fact_id);
