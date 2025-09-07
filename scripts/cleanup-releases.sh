#!/bin/bash

# 脚本：清理旧的GitHub releases，只保留v1.2.1
# 使用前请确保已经通过 gh auth login 完成GitHub CLI认证

echo "🧹 开始清理旧的releases..."
echo ""

# 检查GitHub CLI认证状态
if ! gh auth status &>/dev/null; then
    echo "❌ 错误：GitHub CLI未认证"
    echo "请先运行以下命令进行认证："
    echo "  gh auth login"
    echo ""
    echo "认证步骤："
    echo "  1. 选择 GitHub.com"
    echo "  2. 选择 HTTPS"
    echo "  3. 选择通过浏览器认证"
    echo "  4. 完成浏览器中的认证流程"
    exit 1
fi

echo "✅ GitHub CLI已认证"
echo ""

# 要删除的release IDs（除了v1.2.1）
declare -a releases_to_delete=(
    "245328927"  # v1.2.0
    "245291255"  # v1.1.0
    "245290049"  # v1.0.3
    "245288491"  # v1.0.2
)

declare -a release_names=(
    "v1.2.0"
    "v1.1.0"
    "v1.0.3"
    "v1.0.2"
)

echo "📋 将要删除以下releases："
for i in "${!release_names[@]}"; do
    echo "  - ${release_names[$i]} (ID: ${releases_to_delete[$i]})"
done
echo ""

read -p "确认删除这些releases吗？(y/N): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ 操作已取消"
    exit 0
fi

echo ""
echo "🗑️  开始删除releases..."

# 删除每个release
for i in "${!releases_to_delete[@]}"; do
    release_id="${releases_to_delete[$i]}"
    release_name="${release_names[$i]}"
    
    echo -n "  删除 ${release_name}... "
    
    if gh release delete "${release_name}" --yes --repo KCNyu/lark-logger 2>/dev/null; then
        echo "✅"
    else
        echo "❌ 失败（可能已经被删除）"
    fi
done

echo ""
echo "🎉 清理完成！"
echo ""
echo "📊 当前releases状态："
gh release list --repo KCNyu/lark-logger

echo ""
echo "✨ 只保留了最新的v1.2.1版本"
