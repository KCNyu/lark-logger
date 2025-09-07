# lark-logger

English | [中文文档](docs/README_zh.md)

A minimal Go SDK for sending structured logs to Lark (Feishu) via webhook using clean, mobile-friendly cards.

## Install

```bash
go get github.com/KCNyu/lark-logger
```

## Quick start (context required)

```go
import (
  "context"
  larklogger "github.com/KCNyu/lark-logger"
)

func main() {
    ctx := context.Background()
    client := larklogger.NewClient(larklogger.GetWebhookURL())
    logger := larklogger.NewLogger(ctx, client,
        larklogger.WithService("api-gateway"),
        larklogger.WithEnv("production"),
        larklogger.WithTitle("System Monitor"),
        larklogger.WithShowConfig(true),
    )

    logger.Info("Service started", map[string]interface{}{"port": 8080})
    logger.Warnf("Memory usage", "usage", "87%")
}
```

## Environment

- LARK_WEBHOOK_URL: your bot webhook
- LARK_TEST_MODE: set `true` to skip real sends in tests

## Buttons (optional)

```go
logger := larklogger.NewLogger(ctx, client,
  larklogger.WithButtons([]larklogger.Button{
    {Text: "View Logs", URL: "https://logs.example.com", Style: larklogger.ButtonStylePrimary},
    {Text: "Restart", URL: "https://ops.example.com/restart", Style: larklogger.ButtonStyleDanger, Confirm: true},
  }),
)
```

## Local testing

- `make test` sets test mode automatically and skips external sends
- Or: `LARK_TEST_MODE=true go test ./src/larklogger/...`

## Screenshots

See placeholders below; replace with real images under `docs/images/`.

### CLI output

![CLI output placeholder](docs/images/cli_output.png)

### Cards — Desktop vs Mobile

| Desktop | Mobile |
| --- | --- |
| ![Desktop card placeholder](docs/images/desktop_card.png) | ![Mobile card placeholder](docs/images/mobile_card.png) |

### Config 2×2 Grid (edge-to-edge)

| Desktop | Mobile |
| --- | --- |
| ![Desktop config placeholder](docs/images/desktop_config.png) | ![Mobile config placeholder](docs/images/mobile_config.png) |

### KV — Short pairs (table) vs Long pairs (stacked block)

| Short KV (table) | Long KV (stacked) |
| --- | --- |
| ![Short KV placeholder](docs/images/kv_short.png) | ![Long KV placeholder](docs/images/kv_long.png) |

### Buttons & confirmations

![Buttons placeholder](docs/images/buttons.png)

## License

MIT