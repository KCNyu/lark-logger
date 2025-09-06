#!/bin/bash

# Load environment variables from .env file if it exists
if [ -f .env ]; then
    echo "Loading environment variables from .env file..."
    export $(cat .env | grep -v '^#' | xargs)
    echo "Environment variables loaded successfully!"
    echo "LARK_WEBHOOK_URL: $LARK_WEBHOOK_URL"
    echo "LARK_TEST_MODE: $LARK_TEST_MODE"
else
    echo "No .env file found. Using default test configuration."
    echo "To use real webhook, create .env file with:"
    echo "LARK_WEBHOOK_URL=https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url"
    echo "LARK_TEST_MODE=false"
fi
