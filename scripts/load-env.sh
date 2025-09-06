#!/bin/bash

# Load environment variables from .env files if they exist
# Priority: .env.local > .env > default test config

if [ -f .env.local ]; then
    echo "Loading environment variables from .env.local file..."
    export $(cat .env.local | grep -v '^#' | xargs)
    echo "Environment variables loaded from .env.local!"
elif [ -f .env ]; then
    echo "Loading environment variables from .env file..."
    export $(cat .env | grep -v '^#' | xargs)
    echo "Environment variables loaded from .env!"
else
    echo "No .env or .env.local file found. Using default test configuration."
    echo "To use real webhook, create .env.local file with:"
    echo "LARK_WEBHOOK_URL=https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url"
    echo "LARK_TEST_MODE=false"
fi

echo "Current configuration:"
echo "LARK_WEBHOOK_URL: $LARK_WEBHOOK_URL"
echo "LARK_TEST_MODE: $LARK_TEST_MODE"
