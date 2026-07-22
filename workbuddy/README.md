# WorkBuddy 插件

[WorkBuddy / CodeBuddy]（Tencent CodeBuddy `copilot.tencent.com`）CPA 插件。

## 功能特性

- **OAuth 登录**：通过 CodeBuddy 网页授权流程获取 access token，多账号按 `workbuddy-<uid>.json` 独立保存
- **动态模型拉取**：登录后自动从 `copilot.tencent.com/console/enterprises/personal/models` 拉取可用模型列表（5 分钟缓存 + 硬编码 fallback）
- **Executor 转发**：流式 + 非流式 chat completions
  - 后端始终返回 SSE，非流式请求自动聚合为单个 `chat.completion` JSON
  - `cleanChunk` 清洗空 `function_call` / `tool_calls`，防客户端 truncated
  - 跨协议入口（claude/gemini/codex）自动补 SSE `data:` framing
  - alias 反解：`ExecutorRequest.Model` 可能是别名，用 `model.static`/`model.for_auth` 缓存的 Host 别名表反解后发上游
  - `tool_choice` 归一化：OpenAI object 类型 → CodeBuddy 上游 string 类型（v0.3.15+）
- **Usage 上报**：插件 executor 路径直接 `usage.PublishRecord`（宿主不对插件 executor 自动记用量），覆盖成功/上游 4xx/网络错误三条出口
- **每日自动签到**：09:00 / 21:00 本地时区自动签到，可在面板开关（运行时开关，不持久化到 config.yaml）
- **积分与套餐面板**：查看账号昵称、积分余额、套餐余量、用量进度、签到状态；并发拉取 + billing 5xx 重试 + stale-while-error 缓存
- **OAuth 模型别名 / 禁用**：由 CPA 原生 `oauth-model-alias` / `oauth-excluded-models` 管理，插件只返回完整上游列表

## 安装方法

1. 将 `workbuddy.so` 复制到 CPA 插件目录：

```bash
cp workbuddy.so /path/to/cliproxyapi/plugins/workbuddy.so
```

2. 在 CPA `config.yaml` 中启用：

```yaml
plugins:
  enabled: true
  dir: "plugins"
  configs:
    workbuddy:
      enabled: true
```

3. 重启 CPA 服务。

## 构建

需要 Go 1.26.0+ 和 CGO：

```bash
cd workbuddy
make build VERSION=0.3.18
# 产物：dist/workbuddy.so
```

### 测试

```bash
make test    # go test ./...
make vet     # go vet ./...
make check   # test + vet
```

## 配置说明

### 自动签到

插件默认启用自动签到。可通过管理面板开关，或调用管理 API：

```bash
curl -X POST -H "Authorization: Bearer $CPA_MGMT_KEY" \
  -H "Content-Type: application/json" \
  -d '{"enabled":true}' \
  http://127.0.0.1:12888/v0/management/plugins/workbuddy/checkin/config
```

> **限制**：`checkin_auto` 是运行时开关，CPA 重启后恢复 config 默认值。宿主无 plugin-config 写回调。

### OAuth 模型别名 / 禁用

在 CPA `config.yaml` 中配置（由宿主原生处理）：

```yaml
oauth-excluded-models:
  workbuddy:
    - hy3

oauth-model-alias:
  workbuddy:
    - name: deepseek-v4-pro
      alias: workbuddy-dsv4-pro
      fork: true
```

## 管理面板

```
http://<cpa-host>/v0/resource/plugins/workbuddy/panel
```

- 查看所有 WorkBuddy 账号的积分、套餐、用量进度
- 手动签到 / 全部签到
- 开启/关闭自动签到
- 刷新缓存

面板鉴权（三通道，无硬编码 key）：
1. 从 CPA 主面板 iframe 嵌入时自动读 `localStorage["cli-proxy-auth"]`
2. URL `?key=` 参数（读后 `replaceState` 清除）
3. 手动输入（存 `sessionStorage`）

## 管理 API

| 接口 | 方法 | 说明 |
|---|---|---|
| `/v0/management/plugins/workbuddy/accounts` | GET | 列出账号、积分、签到状态 |
| `/v0/management/plugins/workbuddy/refresh` | POST | 强制刷新积分缓存 |
| `/v0/management/plugins/workbuddy/checkin` | POST | 手动签到（单账号 `auth_index` 或全部） |
| `/v0/management/plugins/workbuddy/checkin/config` | POST | 设置自动签到开关 |
| `/v0/management/plugins/workbuddy/credits` | GET | 实时查询积分（单账号 `auth_index` 或全部） |

## 能力声明

| 能力 | 状态 |
|---|---|
| AuthProvider | ✅ OAuth login/poll/refresh/parse |
| ModelProvider | ✅ 动态模型 + 静态 fallback |
| Executor | ✅ 流式 + 非流式 + tool_choice 归一 |
| ManagementAPI | ✅ 面板 + 管理 API |
| FrontendAuthProvider | ❌ 不声明（依赖标准 OAuth 流程） |

## 文件说明

| 文件 | 说明 |
|---|---|
| `main.go` | 插件主入口：OAuth/模型/executor/usage |
| `management.go` | 管理 API、签到调度、面板后端 |
| `panel.html` | 管理面板前端 |
| `*_test.go` | 单测：cleanchunk / expiry / credits / billing_retry / toolchoice |
| `Makefile` | 构建/测试/检查 |
| `go.mod` / `go.sum` | Go 模块依赖 |

## 注意事项

- 插件需要 CPA 管理密钥才能访问管理 API；从 CPA 主面板嵌入时自动读取
- 多个 WorkBuddy 账号登录时按 `workbuddy-<uid>.json` 命名保存
- 模型列表通过上游 API 动态获取，需要账号已登录且 token 有效
- `checkin_auto` 开关不持久化到 config.yaml（宿主限制）
- 插件不自行处理 `oauth-model-alias` / `oauth-excluded-models`（宿主原生处理）
