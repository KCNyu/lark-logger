#!/bin/bash

# 本地覆盖率报告脚本
# 当Codecov遇到速率限制时，可以使用此脚本生成本地覆盖率报告

echo "📊 生成本地覆盖率报告..."
echo ""

# 检查是否在项目根目录
if [ ! -f "go.mod" ]; then
    echo "❌ 请在项目根目录运行此脚本"
    exit 1
fi

# 运行测试并生成覆盖率报告
echo "🧪 运行测试..."
go test -v -race -coverprofile=coverage.out -covermode=atomic ./src/larklogger/...

if [ $? -ne 0 ]; then
    echo "❌ 测试失败"
    exit 1
fi

echo ""
echo "📈 覆盖率统计:"
go tool cover -func=coverage.out | tail -1

echo ""
echo "🌐 生成HTML报告..."
go tool cover -html=coverage.out -o coverage.html

echo ""
echo "✅ 覆盖率报告已生成:"
echo "  📄 HTML报告: coverage.html"
echo "  📊 数据文件: coverage.out"
echo ""
echo "💡 提示:"
echo "  - 在浏览器中打开 coverage.html 查看详细报告"
echo "  - 如果Codecov恢复正常，可以手动上传 coverage.out"
echo ""
echo "🔗 Codecov上传命令 (如果需要):"
echo "  curl -s https://codecov.io/bash | bash -s -- -t \$CODECOV_TOKEN"
