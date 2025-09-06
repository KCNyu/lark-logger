# lark-logger

[![CI](https://github.com/KCNyu/lark-logger/workflows/CI/badge.svg)](https://github.com/KCNyu/lark-logger/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/KCNyu/lark-logger)](https://goreportcard.com/report/github.com/KCNyu/lark-logger)
[![GoDoc](https://godoc.org/github.com/KCNyu/lark-logger?status.svg)](https://godoc.org/github.com/KCNyu/lark-logger)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[English](../README.md) | 中文文档

一个优雅的Go SDK，用于向Lark（飞书）机器人发送结构化的日志消息，支持美观的交互式卡片格式。

## ✨ 特性

- 🎨 **美观的卡片格式** - 支持emoji、颜色模板和Markdown格式
- 📊 **结构化日志** - 支持`map[string]interface{}`参数，智能格式化显示
- 🚦 **多级别日志** - 支持Info、Warn、Error三种级别，每种级别都有独特的样式和颜色
- 🔄 **重试机制** - 内置重试逻辑，确保消息可靠送达
- ⚙️ **类型安全配置** - 使用函数式选项模式，提供类型安全的配置
- 🎨 **自定义标题** - 支持自定义主标题和格式化副标题
- 📋 **智能字段显示** - 清晰的KV表格布局，良好的对齐和可读性
- 🔧 **增强布局** - 专业的两列布局，交替行颜色
- 📱 **移动端优化** - 响应式设计，合适的填充和字体大小
- 🎨 **视觉层次** - 标题和内容的不同字体大小和颜色
- 🚀 **简化API** - 支持Infof/Warnf/Errorf方法，键值对参数
- 🔧 **配置显示** - 可选的2x2网格布局显示logger配置（默认隐藏）
- 🎯 **环境管理** - 智能环境变量处理，支持本地/CI/测试模式
- 🧪 **完整测试覆盖** - 包含单元测试和集成测试
- 🚀 **CI/CD就绪** - 包含GitHub Actions工作流和GoReleaser配置

## 📦 安装

```bash
go get github.com/KCNyu/lark-logger
```

## 🚀 快速开始

### 环境配置

在使用lark-logger之前，需要配置环境变量：

#### 方式1：使用.env文件（推荐）

1. 复制示例环境文件：
```bash
cp env.example .env.local
```

2. 编辑`.env.local`文件，填入您的实际webhook URL：
```bash
# 生产环境使用（发送真实消息到Lark）
LARK_WEBHOOK_URL=https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url
LARK_TEST_MODE=false

# 测试/开发环境（使用测试webhook，不发送真实消息）
# LARK_WEBHOOK_URL=https://test.webhook.url
# LARK_TEST_MODE=true
```

3. 加载环境变量：
```bash
source scripts/load-env.sh
```

**注意**：`.env.local`文件会被git自动忽略，确保您的webhook URL安全。

#### 方式2：直接导出

```bash
# 生产环境使用（发送真实消息到Lark）
export LARK_WEBHOOK_URL="https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url"
export LARK_TEST_MODE="false"

# 测试/开发环境（使用测试webhook，不发送真实消息）
export LARK_TEST_MODE="true"
export LARK_WEBHOOK_URL="https://test.webhook.url"
```

#### CI/测试配置

对于CI和测试环境，以下环境变量会自动设置：
- `LARK_TEST_MODE=true`
- `LARK_WEBHOOK_URL=https://test.webhook.url`

### 基本使用

```go
package main

import (
    "github.com/KCNyu/lark-logger"
)

func main() {
    // 从环境变量获取webhook URL
    webhookURL := larklogger.GetWebhookURL()

    // 使用webhook URL和选项创建logger
    logger := larklogger.New(
        webhookURL,
        larklogger.WithEnv("production"),
        larklogger.WithTitle("系统监控"),
    )

    // 使用传统map格式发送日志
    logger.Info("API网关初始化成功", map[string]interface{}{
        "port":     8080,
        "version":  "2.1.0",
        "features": []string{"authentication", "rate_limiting", "metrics"},
    })

    // 使用简化的Infof/Warnf/Errorf格式发送日志
    logger.Infof("服务健康检查", "status", "healthy", "response_time", "120ms", "uptime", "2h30m")
    logger.Warnf("内存使用率接近阈值", "usage", "87%", "threshold", "85%", "recommendation", "考虑水平扩展")
    logger.Errorf("数据库连接池耗尽", "error", "连接超时30秒", "retry_count", 3, "pool_size", 20)

    // 创建启用配置显示的logger
    configLogger := larklogger.New(
        webhookURL,
        larklogger.WithService("config-demo"),
        larklogger.WithEnv("production"),
        larklogger.WithHostname("server-01"),
        larklogger.WithTitle("配置演示"),
        larklogger.WithShowConfig(true), // 启用配置区块
    )

    // 这将显示配置区块，使用2x2网格布局
    configLogger.Info("配置区块已启用", map[string]interface{}{
        "feature": "config_visibility",
        "status":  "enabled",
    })
}
```

### 发送简单文本消息

```go
err := client.SendText("Hello from larklogger! 🚀")
if err != nil {
    log.Printf("Failed to send text message: %v", err)
}
```

### 自定义卡片

```go
card := larklogger.NewCardBuilder().
    SetHeader("🎉 Custom Notification", "green").
    AddSection("**This is a custom card message**\n\nBuilt with larklogger SDK!").
    AddDivider().
    AddSection("**Features:**\n• Beautiful formatting\n• Multiple log levels\n• Structured data").
    Build()

err := client.SendCard(card)
```

## 📋 API 文档

### Logger 接口

```go
type Logger interface {
    Info(message string, fields map[string]interface{})
    Warn(message string, fields map[string]interface{})
    Error(message string, fields map[string]interface{})
}
```

### 客户端配置选项

```go
// 超时设置
WithTimeout(timeout time.Duration) ClientOption

// 重试配置
WithRetry(count int, delay time.Duration) ClientOption

// 用户代理
WithUserAgent(userAgent string) ClientOption

// 自定义请求头
WithHeaders(headers map[string]string) ClientOption
```

### Logger 配置选项

```go
// 服务名称（可选，默认为"larklogger"）
WithService(service string) LoggerOption

// 环境标识（可选，默认为"development"）
WithEnv(env string) LoggerOption

// 主机名（可选，默认为"localhost"）
WithHostname(hostname string) LoggerOption

// 日志卡片的自定义标题（可选，默认为"System Log"）
WithTitle(title string) LoggerOption

// 在卡片中显示配置区块（可选，默认为false）
WithShowConfig(show bool) LoggerOption
```

## 🎨 日志级别样式

| 级别 | Emoji | 颜色模板 | 用途 |
|------|-------|----------|------|
| Info | ℹ️ | blue | 一般信息 |
| Warn | ⚠️ | orange | 警告信息 |
| Error | ❌ | red | 错误信息 |

## 🧪 测试

运行所有测试：

```bash
go test -v ./...
```

运行测试并生成覆盖率报告：

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📁 项目结构

```
lark-logger/
├── .github/
│   └── workflows/          # GitHub Actions 工作流
├── cmd/                    # 示例应用程序
├── docs/                   # 文档
├── scripts/                # 工具脚本
├── src/
│   └── larklogger/         # 主包
│       ├── *.go           # 源文件
│       └── *_test.go      # 测试文件
├── go.mod                  # Go 模块文件
├── larklogger.go          # 主包导出
├── .golangci.yml          # Linter 配置
├── .goreleaser.yml         # GoReleaser 配置
├── Dockerfile              # Docker 配置
├── Makefile                # 构建自动化
├── env.example             # 环境变量示例
└── README.md               # 项目文档
```

## 🔧 开发

### 前置要求

- Go 1.19+
- Git

### 本地开发

1. 克隆仓库：
```bash
git clone https://github.com/KCNyu/lark-logger.git
cd larklogger
```

2. 安装依赖：
```bash
go mod download
```

3. 运行测试：
```bash
go test -v ./src/larklogger/...
```

4. 运行示例：
```bash
go run ./cmd/main.go
```

5. 使用Makefile运行：
```bash
make run          # 使用环境变量运行示例
make test         # 运行测试
make lint         # 运行linter
make ci           # 运行完整CI流程
```

### 代码质量

项目使用以下工具确保代码质量：

- `golangci-lint` - 代码静态分析
- `gosec` - 安全扫描
- `go test` - 单元测试
- `go vet` - 代码检查

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

## 📞 支持

如果你遇到任何问题或有建议，请：

- 提交 [Issue](https://github.com/KCNyu/lark-logger/issues)
- 发送邮件到 [KCNyu@example.com](mailto:shengyu.li.evgeny@gmail.com)

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

---

⭐ 如果这个项目对你有帮助，请给它一个星标！
