# 📚 Scripts Documentation

This directory contains utility scripts for the lark-logger project.

## 🔧 Available Scripts

### `load-env.sh`
**Purpose**: Load environment variables from `.env` file  
**Usage**: `source scripts/load-env.sh`

This script:
- Automatically loads environment variables from `.env` file
- Falls back to `.env.local` if `.env` doesn't exist
- Used for local development and testing
- Required for running tests and examples

## 🚀 Quick Start

1. **Setup environment**:
   ```bash
   # Copy environment template
   cp env.example .env
   
   # Edit .env with your webhook URL
   vim .env
   
   # Load environment variables
   source scripts/load-env.sh
   ```

2. **Run tests with environment**:
   ```bash
   # Load env and run tests
   source scripts/load-env.sh && go test ./...
   
   # Or use Makefile which handles this automatically
   make test
   ```

## 📝 Notes

- The `load-env.sh` script is automatically called by the Makefile during testing
- Environment variables are essential for webhook configuration
- Never commit `.env` files with real webhook URLs to version control