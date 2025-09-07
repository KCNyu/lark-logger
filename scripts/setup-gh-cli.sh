#!/bin/bash

# GitHub CLI 认证设置脚本

echo "🔐 GitHub CLI 认证设置"
echo "========================"
echo ""

# 检查是否已经认证
if gh auth status &>/dev/null; then
    echo "✅ GitHub CLI已经认证成功！"
    echo ""
    gh auth status
    echo ""
    echo "📝 如需重新认证，请运行: gh auth logout && gh auth login"
    exit 0
fi

echo "❌ GitHub CLI尚未认证"
echo ""
echo "请选择认证方式："
echo "  1) 通过浏览器认证（推荐）"
echo "  2) 使用Personal Access Token"
echo "  3) 稍后手动认证"
echo ""
read -p "请选择 (1/2/3): " choice

case $choice in
    1)
        echo ""
        echo "🌐 即将打开浏览器进行认证..."
        echo ""
        echo "认证步骤："
        echo "  1. 按Enter继续"
        echo "  2. 选择 GitHub.com"
        echo "  3. 选择 HTTPS 协议"
        echo "  4. 使用浏览器认证"
        echo "  5. 在浏览器中完成认证"
        echo ""
        read -p "按Enter开始认证..."
        gh auth login
        ;;
    2)
        echo ""
        echo "📝 使用Personal Access Token认证"
        echo ""
        echo "请先在GitHub创建Personal Access Token："
        echo "  1. 访问: https://github.com/settings/tokens/new"
        echo "  2. 设置token名称（如: gh-cli-access）"
        echo "  3. 选择权限："
        echo "     - repo (完整仓库访问)"
        echo "     - workflow (如果需要管理Actions)"
        echo "     - admin:public_key (如果需要管理SSH keys)"
        echo "  4. 生成token并复制"
        echo ""
        read -p "请粘贴您的Personal Access Token: " token
        echo "$token" | gh auth login --with-token
        ;;
    3)
        echo ""
        echo "📋 手动认证指南："
        echo ""
        echo "稍后运行以下命令进行认证："
        echo "  gh auth login"
        echo ""
        echo "或使用token认证："
        echo "  echo 'YOUR_TOKEN' | gh auth login --with-token"
        echo ""
        exit 0
        ;;
    *)
        echo "❌ 无效选择"
        exit 1
        ;;
esac

echo ""
echo "🔍 检查认证状态..."
echo ""

if gh auth status &>/dev/null; then
    echo "✅ 认证成功！"
    echo ""
    gh auth status
    echo ""
    echo "🎉 您现在可以使用GitHub CLI了！"
    echo ""
    echo "下一步："
    echo "  运行 ./scripts/cleanup-releases.sh 来清理旧的releases"
else
    echo "❌ 认证失败，请重试"
    exit 1
fi
