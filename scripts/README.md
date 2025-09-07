# 📚 Scripts Documentation

This directory contains utility scripts for managing the lark-logger project.

## 🔧 Available Scripts

### 1. `setup-gh-cli.sh`
**Purpose**: Setup GitHub CLI authentication  
**Usage**: `./scripts/setup-gh-cli.sh`

This script helps you:
- Check GitHub CLI authentication status
- Guide through browser-based authentication (recommended)
- Setup authentication using Personal Access Token
- Verify successful authentication

Run this script first if you haven't authenticated GitHub CLI yet.

### 2. `cleanup-releases.sh`
**Purpose**: Clean up old GitHub releases, keeping only v1.2.1  
**Usage**: `./scripts/cleanup-releases.sh`

This script will:
- Check GitHub CLI authentication
- List releases to be deleted (v1.2.0, v1.1.0, v1.0.3, v1.0.2)
- Ask for confirmation before deletion
- Delete old releases
- Show remaining releases

⚠️ **Note**: Requires GitHub CLI authentication. Run `setup-gh-cli.sh` first if not authenticated.

### 3. `load-env.sh`
**Purpose**: Load environment variables from `.env` file  
**Usage**: `source scripts/load-env.sh`

This script:
- Loads environment variables from `.env` file
- Used for local development and testing

## 🚀 Quick Start

1. **First-time setup**:
   ```bash
   # Install GitHub CLI (if not installed)
   brew install gh
   
   # Setup authentication
   ./scripts/setup-gh-cli.sh
   ```

2. **Clean up old releases**:
   ```bash
   ./scripts/cleanup-releases.sh
   ```

3. **For development**:
   ```bash
   # Copy environment template
   cp env.example .env
   
   # Edit .env with your webhook URL
   
   # Load environment
   source scripts/load-env.sh
   ```

## 📝 Notes

- All scripts have executable permissions
- Scripts are designed to be safe with confirmation prompts
- GitHub CLI is required for release management scripts
- Environment scripts are for local development only
