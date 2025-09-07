# lark-logger

English | [中文文档](docs/README_zh.md)

A minimal Go SDK for sending structured logs to Lark (Feishu) via webhook using clean, mobile-friendly cards.

## Install

```bash
go get github.com/KCNyu/lark-logger
```

## Quick start

```go
import larklogger "github.com/KCNyu/lark-logger"

func main() {
    logger := larklogger.New(
        larklogger.GetWebhookURL(),
        larklogger.WithService("api-gateway"),
        larklogger.WithEnv("production"),
        larklogger.WithTitle("System Monitor"),
        larklogger.WithShowConfig(true),
    )

    // Map-style structured logging
    logger.Info("Service started", map[string]interface{}{
        "port": 8080,
        "version": "1.2.0",
        "features": []string{"auth", "rate_limit", "metrics"},
    })

    // Pair-style logging (Infof/Warnf/Errorf)
    logger.Warnf("Memory usage", "usage", "87%", "threshold", "85%")
    logger.Errorf("Database error", "error", "connection timeout", "retry_count", 3)
}
```

## Environment

- LARK_WEBHOOK_URL: your bot webhook
- LARK_TEST_MODE: set `true` to skip real sends in tests

Create `.env.local` (optional):

```bash
LARK_WEBHOOK_URL=https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook
LARK_TEST_MODE=false
```

## Optional actions (buttons)

```go
logger := larklogger.New(
  larklogger.GetWebhookURL(),
  larklogger.WithTitle("Action Required"),
  larklogger.WithButtons([]larklogger.Button{
    {Text: "View Logs", URL: "https://logs.example.com/entry", Style: larklogger.ButtonStylePrimary},
    {Text: "Restart", URL: "https://ops.example.com/restart", Style: larklogger.ButtonStyleDanger, Confirm: true},
  }),
)
logger.Error("Critical incident", map[string]interface{}{"error_code": "SYS_001"})
```

## Local testing

- `make test` sets test mode automatically and skips external sends
- You can also run: `LARK_TEST_MODE=true go test ./src/larklogger/...`

## License

MIT