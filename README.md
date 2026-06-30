# Prom Lens Backend

Prometheus 规则与告警通知管理后端：提供认证鉴权、审计日志、Prometheus/Thanos 规则与采集目标配置（同步 K8s ConfigMap）、Lark 告警通知通道及 Alertmanager receiver/route 自动同步。

## 技术栈

- **语言**: Go 1.24+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL（库名默认 `prom_lens`，表前缀 `prom_lens_`）
- **认证**: JWT（Access Token + Refresh Token，HS512）
- **Kubernetes**: client-go（规则 / 采集目标 / Alertmanager 配置同步）

## 项目结构

```
cmd/server/                 # 程序入口
config/
  config.example.yaml       # 配置模板（复制为 config.yaml 后修改）
  migrations/               # 数据库迁移 SQL（手工执行）
internal/
  apps/
    authn/                  # 登录、改密、用户资料、用户管理
    audit/                  # 登录审计、操作审计
    prometheus/             # 规则组、规则、记录、采集目标；K8s 正向同步与反向导入
    alerting/               # 通知通道、Lark 转发、Alertmanager 同步
  core/                     # 配置、数据库、日志、分页、中间件
  pkg/di/                   # 依赖注入
  routes/                   # 路由注册
```

## 功能模块

### 认证（`/api/v1/authn`）

- `POST /login` — 登录
- `POST /refresh-token` — 刷新 Token
- `POST /me/change-password` — 修改密码
- `GET /me/profile` — 当前用户资料

**用户管理（需 JWT + 超级管理员）**

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/users` | 用户列表（`page`/`size`/`sortBy`/`sortOrder`/`keyword`） |
| GET | `/users/:userID` | 用户详情 |
| POST | `/users` | 创建用户 |
| PATCH | `/users/:userID` | 更新用户 |
| DELETE | `/users/:userID` | 删除用户 |
| POST | `/users/:userID/reset-password` | 重置密码 |

创建用户 body 示例：

```json
{
  "username": "ops",
  "password": "secret123",
  "nickname": "运维",
  "email": "ops@example.com",
  "phone": "13800138000",
  "isActive": true,
  "isSuperuser": false
}
```

### 审计（`/api/v1/audit`）

- `GET /auth-logs` — 登录审计
- `GET /logs` — 操作审计

### Prometheus 配置（`/api/v1/prometheus`）

- **规则组** `groups`：类型为 `ALERTING RULES` 或 `ALERTING RECORDS`（含空格）
- **告警规则** `groups/:groupID/rules`、**记录规则** `groups/:groupID/records`
- **采集目标组** `target-groups`、**采集目标** `target-groups/:groupID/targets`
- **反向导入** `POST /sync/import-rules` — 从 ConfigMap 导入规则（仅新增）

CRUD 后自动将规则/采集配置同步到 K8s ConfigMap（见 `config.yaml` 中 `prometheus.rule` / `prometheus.target`）。

规则 `annotations` 支持 `summary`、`description` 及 `extraAnnotations`（合并写入 YAML）。

更多 curl 示例见 [internal/apps/prometheus/README.md](internal/apps/prometheus/README.md)。

### 告警通知（`/api/v1/alerting`）

**管理接口（需 JWT）**

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/webhooks` | 通道列表 |
| GET | `/webhooks/:webhookID` | 通道详情 |
| POST | `/webhooks` | 创建通道（可选 route matchers，自动同步 AM） |
| PATCH | `/webhooks/:webhookID` | 更新通道 |
| DELETE | `/webhooks/:webhookID` | 删除通道（仅 superuser） |
| POST | `/webhooks/verify` | 验证 Lark Webhook（不落库） |
| POST | `/sync/receivers` | 全量同步 Alertmanager receivers/routes |

**回调接口（免 JWT）**

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/webhook/:channel` | Alertmanager 回调，需 `Authorization: Bearer <callbackToken>` |

#### 告警链路

```
Prometheus/Thanos Ruler → Alertmanager（route 匹配）
    → POST Prom Lens /api/v1/alerting/webhook/:channel（Bearer）
    → Lark 卡片通知
```

创建通道时：

1. 写入 DB（`webhook` + 可选 `route` / `matchers`）
2. 自动生成 `callbackToken`
3. 同步 Alertmanager ConfigMap：`receivers` + `route.routes`（`matchers` 格式，非 deprecated `match`）
4. 创建时若 AM 同步失败且已配置同步，**回滚**通道记录

#### 创建通道请求示例

```json
{
  "name": "hk-test",
  "url": "https://open.larksuite.com/open-apis/bot/v2/hook/xxx",
  "description": "测试通道",
  "enabled": true,
  "route": {
    "priority": 10,
    "continue": false,
    "matchers": [
      { "label": "env", "operator": "=", "value": "prod" },
      { "label": "team", "operator": "=", "value": "ops" }
    ]
  }
}
```

响应含 `data.callbackToken` 与 `amSync`（同步结果）。

#### Alertmanager receiver 配置（自动同步结果）

```yaml
receivers:
  - name: hk-test
    webhook_configs:
      - url: http://<prom-lens-host>:8081/api/v1/alerting/webhook/hk-test
        send_resolved: true
        http_config:
          bearer_token: <callbackToken>

route:
  routes:
    - receiver: hk-test
      matchers:
        - env="prod"
        - team="ops"
      continue: false
```

#### 直接投递 Alertmanager 测试（不经过规则）

```bash
NOW=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
kubectl exec -n monitoring deploy/alertmanager -- wget -qO- \
  --header='Content-Type: application/json' \
  --post-data='[{"labels":{"alertname":"DirectAMRouteTest","env":"prod","team":"ops"},"annotations":{"summary":"路由测试"},"startsAt":"'"${NOW}"'"}]' \
  http://localhost:9093/api/v2/alerts
```

#### 模拟 Prom Lens 回调

```bash
curl -X POST "http://<prom-lens-host>:8081/api/v1/alerting/webhook/hk-test" \
  -H "Authorization: Bearer <callbackToken>" \
  -H "Content-Type: application/json" \
  -d '{"status":"firing","receiver":"hk-test","alerts":[{"status":"firing","labels":{"alertname":"Test"},"annotations":{"summary":"测试"},"startsAt":"2026-01-01T00:00:00Z"}]}'
```

## 权限说明

- 所有 `/api/v1` 接口默认需 JWT（白名单除外）
- **DELETE** 请求仅 **superuser**（JWT 中 `su=true`）可执行
- Alertmanager 回调仅校验 Bearer `callbackToken`，不走 JWT

## 健康检查

| 路径 | 说明 | 鉴权 |
|------|------|------|
| `GET /health` | 存活探针 | 免鉴权 |
| `GET /readyz` | 就绪探针（含 DB ping） | 免鉴权 |

## 配置

复制模板并修改本地配置（`config/config.yaml` 已加入 `.gitignore`，不会提交密钥）：

```bash
cp config/config.example.yaml config/config.yaml
```

默认读取 `config/config.yaml`，启动参数 `-c` 可指定其他路径。

环境变量前缀：`PROM_LENS_`，嵌套用下划线，例如：

```bash
export PROM_LENS_SERVER_PORT=8081
export PROM_LENS_DATABASE_PASSWORD=your-password
export PROM_LENS_APP_JWT_SECRET=your-secret-at-least-16-chars
export PROM_LENS_BASE_URL=http://your-host:8081
```

主要配置项：

```yaml
server:
  port: 8081

base_url: "http://127.0.0.1:8081"   # 对外访问地址，AM 回调用；环境变量 PROM_LENS_BASE_URL

app:
  jwt_secret: "<生产环境请用环境变量注入>"

auth:
  access_token_expires: 1h
  refresh_token_expires: 24h
  whitelist:                    # 免 JWT 路径前缀
    - /api/v1/authn/login
    - /api/v1/authn/refresh-token
    - /api/v1/alerting/webhook/
    - /health
    - /readyz

database:
  host: "127.0.0.1"
  dbname: "prom_lens"
  # ...

prometheus:
  target:
    namespace: monitoring-test
    configmap: prometheus-targets
  rule:
    namespace: thanos
    configmap: thanos-ruler-config

alerting:
  lark_timeout: 10s
  alertmanager:
    namespace: monitoring
    configmap: alertmanager-config
    config_key: alertmanager.yml
    default_receiver: ""                            # 可选，写入根 route.receiver
```

`base_url`（`PROM_LENS_BASE_URL`）必须为 Alertmanager Pod **能访问**的地址（建议 K8s Service 或稳定内网 IP，不要用 `localhost`）。

## 数据库迁移

迁移脚本位于 `config/migrations/`，需**手工执行**（无内置 migrate 工具）。

### 新环境（推荐顺序）

```bash
mysql ... < config/migrations/00_database.sql
mysql ... prom_lens < config/migrations/authn_init.sql
mysql ... prom_lens < config/migrations/audit_init.sql
mysql ... prom_lens < config/migrations/prometheus_init.sql
mysql ... prom_lens < config/migrations/alerting_init.sql
mysql ... prom_lens < config/migrations/authn_admin_seed.sql   # 默认 admin，上线后务必改密
```

各模块 `*_init.sql` 为完整建表脚本；`alerting_init.sql` 已含 `callback_token`、`route`、`route_matcher`。

### 已有环境升级

执行 `upgrade_legacy.sql`，按当前库结构跳过已完成的步骤（重复执行可能报错，属正常）。

```bash
mysql ... prom_lens < config/migrations/upgrade_legacy.sql
```

## 运行

```bash
# 依赖：MySQL 已就绪并完成迁移
go run ./cmd/server -c config/config.yaml
```

默认监听 `http://localhost:8081`，API 前缀 `/api/v1`。

```bash
# 编译
go build -o bin/prom-lens-backend ./cmd/server
```

## 日志

- 规则同步：`[prom-sync]`
- 规则导入：`[prom-import]`
- Alertmanager 同步：`[am-sync]`
- Lark 通知：`[alert-notify]`

日志目录默认 `logs/`（`config.yaml` 中 `logging.log_dir`）。

## 模块说明

Go module 名为 `prom-lens-backend`，与仓库名一致。
