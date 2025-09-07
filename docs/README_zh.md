# lark-logger

[English](../README.md) | 中文文档

一个简洁的 Go SDK，用于通过 Lark（飞书）Webhook 发送结构化日志，卡片样式清晰、移动端友好。

## 安装

```bash
go get github.com/KCNyu/lark-logger
```

## 快速开始

```go
import larklogger "github.com/KCNyu/lark-logger"

func main() {
    logger := larklogger.New(
        larklogger.GetWebhookURL(),
        larklogger.WithService("api-gateway"),
        larklogger.WithEnv("production"),
        larklogger.WithTitle("系统监控"),
        larklogger.WithShowConfig(true),
    )

    // Map 形式结构化日志
    logger.Info("服务启动", map[string]interface{}{
        "port": 8080,
        "version": "1.2.0",
        "features": []string{"auth", "rate_limit", "metrics"},
    })

    // 键值对形式（Infof/Warnf/Errorf）
    logger.Warnf("内存使用率", "usage", "87%", "threshold", "85%")
    logger.Errorf("数据库错误", "error", "connection timeout", "retry_count", 3)
}
```

## 环境变量

- LARK_WEBHOOK_URL：你的机器人 webhook
- LARK_TEST_MODE：测试模式（`true` 可跳过真实发送）

可选创建 `.env.local`：

```bash
LARK_WEBHOOK_URL=https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook
LARK_TEST_MODE=false
```

## 可选操作按钮

```go
logger := larklogger.New(
  larklogger.GetWebhookURL(),
  larklogger.WithTitle("需要操作"),
  larklogger.WithButtons([]larklogger.Button{
    {Text: "查看日志", URL: "https://logs.example.com/entry", Style: larklogger.ButtonStylePrimary},
    {Text: "重启服务", URL: "https://ops.example.com/restart", Style: larklogger.ButtonStyleDanger, Confirm: true},
  }),
)
logger.Error("关键事件", map[string]interface{}{"error_code": "SYS_001"})
```

## 本地测试

- `make test` 会自动启用测试模式并跳过外部发送
- 也可直接运行：`LARK_TEST_MODE=true go test ./src/larklogger/...`

## 许可证

MIT