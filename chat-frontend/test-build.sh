#!/bin/bash

echo "🚀 开始测试聊天应用前端项目构建..."

# 检查Node.js版本
echo "📋 检查环境..."
node --version
npm --version

# 安装依赖
echo "📦 安装依赖..."
npm install

# TypeScript类型检查
echo "🔍 TypeScript类型检查..."
npm run type-check

# 构建Web应用
echo "🏗️ 构建Web应用..."
npm run build

# 检查构建结果
if [ -d "dist" ]; then
    echo "✅ Web应用构建成功！"
    ls -la dist/
else
    echo "❌ Web应用构建失败！"
    exit 1
fi

echo "🎉 项目构建测试完成！"
echo ""
echo "📖 下一步："
echo "1. 启动开发服务器: npm run dev"
echo "2. 启动Electron应用: npm run electron:dev"
echo "3. 查看文档: README.md 和 QUICK_START.md"
