# go-code 项目分析文档

## 项目概述

go-code 是一个 AI 编程助手命令行工具，灵感来自 Claude Code，使用 Go 语言重写。

**GitHub**: https://github.com/lucaswangdev/go-code

## 技术架构

### 项目结构

```
go-code/
├── cmd/go-code/main.go          # 程序入口
├── internal/
│   ├── agent/agent.go           # Agent 核心逻辑
│   ├── config/config.go         # 配置管理
│   ├── context/context.go       # 上下文管理（token 压缩）
│   ├── llm/
│   │   ├── client.go            # LLM API 客户端
│   │   ├── toolcall.go          # ToolCall 数据结构
│   │   └── pricing.go           # 模型定价
│   ├── prompt/prompt.go         # 系统提示词
│   ├── session/session.go       # 会话持久化
│   └── tools/
│       ├── base.go              # Tool 接口定义
│       ├── registry.go          # Tool 注册表
│       ├── bash.go              # bash 工具
│       ├── read.go              # read_file 工具
│       ├── write.go             # write_file 工具
│       ├── edit.go              # edit_file 工具
│       ├── glob.go              # glob 工具
│       ├── grep.go              # grep 工具
│       └── agent.go             # agent（子代理）工具
├── .github/workflows/release.yml # GitHub Actions 发布流程
├── README.md
└── CLAUDE.md
```

## 核心组件

### 1. Agent (agent.go)

Agent 是核心编排器，负责：
- 管理对话消息历史
- 调用 LLM 并处理响应
- 执行 Tool（函数调用）
- 支持并行执行多个 Tool
- 最多 50 轮 Tool 调用循环

```go
type Agent struct {
    LLM          *llm.LLM       // LLM 客户端
    Tools        []tools.Tool   // 可用工具列表
    Messages     []map[string]interface{} // 对话历史
    Context      *context.ContextManager  // 上下文管理器
    MaxRounds    int            // 最大循环次数 (50)
    SystemPrompt string         // 系统提示词
}
```

**工作流程**：
1. 用户输入 → 添加到 Messages
2. 调用 Context.MaybeCompress 压缩过长上下文
3. 构造 Messages（system + user + history）
4. 调用 LLM.Chat 获取响应
5. 如果没有 ToolCall → 返回内容给用户
6. 如果有 ToolCall → 执行工具，结果添加到 Messages
7. 重复步骤 3-6 直到无 ToolCall 或达到最大轮数

### 2. LLM 客户端 (llm/client.go)

封装 OpenAI 兼容 API 调用：
- 支持自定义 BaseURL（兼容 MiniMax、DeepSeek、Ollama 等）
- 支持 Tool/Function Calling
- 记录 Token 使用量
- 调试日志写入 `~/.go-code/debug.log`

```go
type LLM struct {
    Model                 string
    APIKey                string
    BaseURL               string
    Temperature           float64
    MaxTokens             int
    TotalPromptTokens     int
    TotalCompletionTokens int
    client                *openai.Client
}
```

### 3. Tools (tools/)

实现 CLI 编程辅助工具：

| 工具 | 功能 |
|------|------|
| bash | 执行 Shell 命令 |
| read_file | 读取文件内容 |
| write_file | 创建/重写文件 |
| edit_file | 精确修改文件（old/new 字符串匹配） |
| glob | 文件模式匹配 |
| grep | 内容搜索 |
| agent | 子代理（递归调用） |

每个 Tool 实现以下接口：
```go
type Tool interface {
    Name() string
    Description() string
    Parameters() map[string]interface{}
    Execute(args map[string]interface{}) (string, error)
    Schema() map[string]interface{}
}
```

### 4. Context 管理 (context/context.go)

当对话历史过长时自动压缩：
- Token 估算：约 `len(text) / 3`
- 阈值设置：
  - snipAt: 50% maxTokens → 裁剪 Tool 输出
  - summarizeAt: 70% → 待实现
  - collapseAt: 90% → 待实现

当前实现：仅裁剪超长 Tool 输出（保留前后各 3 行）

### 5. Session 管理 (session/session.go)

会话持久化到 `~/.go-code/sessions/`：
- JSON 格式存储
- 自动生成会话 ID：`session_YYYYMMDD_HHMMSS_uuid`
- 支持恢复指定会话

### 6. Prompt (prompt/prompt.go)

动态生成系统提示词，包含：
- 工具列表
- 工作目录
- 操作系统信息
- 8 条核心规则

## 配置与部署

### 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| OPENAI_API_KEY | API 密钥 | - |
| OPENAI_BASE_URL | API 端点 | - |
| OPENAI_API_BASE | API 端点（备用） | - |
| CORECODER_API_KEY | API 密钥（备用） | - |
| CORECODER_MODEL | 模型名称 | gpt-4o |
| CORECODER_MAX_TOKENS | 最大生成 Token 数 | 4096 |
| CORECODER_TEMPERATURE | 采样温度 | 0.0 |
| CORECODER_MAX_CONTEXT | 最大上下文 Token | 128000 |

### 使用示例

```bash
# MiniMax
export OPENAI_API_KEY="sk-..."
export OPENAI_API_BASE="https://api.minimaxi.com/v1"
export CORECODER_MODEL="MiniMax-M2.7"

# 交互模式
go-code

# 单次请求
go-code -p "Hello, write a hello world in Go"

# 恢复会话
go-code -r session_20260121_120000_abc123
```

### 交互命令

- `/help` - 显示帮助
- `/model` - 显示当前模型
- `/model <name>` - 切换模型
- `/tokens` - 显示 Token 用量
- `/save` - 保存当前会话
- `/sessions` - 列出所有会话
- `/reset` - 清空对话历史
- `quit` / `exit` - 退出

## 发布流程

1. 更新 `README.md` 中的版本号
2. 提交 README 更改
3. 创建 Git Tag：`git tag v0.x.x`
4. 推送 Tag：`git push origin v0.x.x`
5. GitHub Actions 自动构建并发布

## 待改进项

1. **Context 压缩** - summarizeAt 和 collapseAt 逻辑未实现
2. **Diff 追踪** - `/diff` 命令显示空（未实现）
3. **价格估算** - pricing.go 中模型定价数据可能不完整
4. **Error Handling** - 部分错误处理可以更健壮
5. **测试** - 缺少单元测试和集成测试

## 依赖

- `github.com/sashabaranov/go-openai` - OpenAI API 客户端
- `github.com/google/uuid` - UUID 生成