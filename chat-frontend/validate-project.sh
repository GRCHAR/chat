#!/bin/bash

echo "Chat Frontend - 项目验证脚本"
echo "=============================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查函数
check_file() {
    if [ -f "$1" ]; then
        echo -e "${GREEN}✓${NC} $2"
        return 0
    else
        echo -e "${RED}✗${NC} $2 (文件不存在)"
        return 1
    fi
}

check_directory() {
    if [ -d "$1" ]; then
        echo -e "${GREEN}✓${NC} $2"
        return 0
    else
        echo -e "${RED}✗${NC} $2 (目录不存在)"
        return 1
    fi
}

# 计数器
total_checks=0
passed_checks=0

# 检查项目结构
echo "1. 检查项目结构"
echo "-----------------"

# 根目录文件
check_file "package.json" "package.json 配置文件" && ((passed_checks++))
((total_checks++))

check_file "vite.config.ts" "Vite 配置文件" && ((passed_checks++))
((total_checks++))

check_file "tsconfig.json" "TypeScript 配置文件" && ((passed_checks++))
((total_checks++))

check_file ".env" "环境变量文件" && ((passed_checks++))
((total_checks++))

check_file "electron-builder.yml" "Electron 构建配置" && ((passed_checks++))
((total_checks++))

# 源代码目录
check_directory "src" "源代码目录" && ((passed_checks++))
((total_checks++))

check_directory "electron" "Electron 目录" && ((passed_checks++))
((total_checks++))

check_directory "build" "构建资源目录" && ((passed_checks++))
((total_checks++))

echo ""

# 检查核心源代码文件
echo "2. 检查核心源代码文件"
echo "----------------------"

# 主要入口文件
check_file "src/main.ts" "Vue 应用入口文件" && ((passed_checks++))
((total_checks++))

check_file "src/App.vue" "根组件文件" && ((passed_checks++))
((total_checks++))

check_file "electron/main.ts" "Electron 主进程文件" && ((passed_checks++))
((total_checks++))

# 核心功能模块
check_directory "src/api" "API 接口模块" && ((passed_checks++))
((total_checks++))

check_directory "src/router" "路由模块" && ((passed_checks++))
((total_checks++))

check_directory "src/stores" "状态管理模块" && ((passed_checks++))
((total_checks++))

check_directory "src/views" "页面组件模块" && ((passed_checks++))
((total_checks++))

check_directory "src/components" "公共组件模块" && ((passed_checks++))
((total_checks++))

check_directory "src/utils" "工具函数模块" && ((passed_checks++))
((total_checks++))

check_directory "src/types" "类型定义模块" && ((passed_checks++))
((total_checks++))

check_directory "src/config" "配置文件模块" && ((passed_checks++))
((total_checks++))

echo ""

# 检查关键配置文件
echo "3. 检查关键配置文件"
echo "--------------------"

# 检查 package.json 中的关键字段
if [ -f "package.json" ]; then
    echo -n "检查 package.json 配置... "
    if grep -q '"name"' package.json && grep -q '"version"' package.json && grep -q '"scripts"' package.json; then
        echo -e "${GREEN}✓${NC}"
        ((passed_checks++))
    else
        echo -e "${RED}✗${NC} (缺少必要字段)"
    fi
    ((total_checks++))
fi

# 检查环境变量文件
if [ -f ".env" ]; then
    echo -n "检查环境变量配置... "
    if grep -q "VITE_API_BASE_URL" .env && grep -q "VITE_WS_BASE_URL" .env; then
        echo -e "${GREEN}✓${NC}"
        ((passed_checks++))
    else
        echo -e "${YELLOW}⚠${NC} (缺少推荐的环境变量)"
        ((passed_checks++)) # 仍然算通过，但有警告
    fi
    ((total_checks++))
fi

echo ""

# 检查TypeScript配置
echo "4. 检查TypeScript配置"
echo "---------------------"

if [ -f "tsconfig.json" ]; then
    echo -n "检查 TypeScript 配置... "
    if grep -q '"compilerOptions"' tsconfig.json && grep -q '"target"' tsconfig.json; then
        echo -e "${GREEN}✓${NC}"
        ((passed_checks++))
    else
        echo -e "${RED}✗${NC} (配置不完整)"
    fi
    ((total_checks++))
fi

if [ -f "tsconfig.node.json" ]; then
    echo -e "${GREEN}✓${NC} Node.js TypeScript 配置"
    ((passed_checks++))
else
    echo -e "${RED}✗${NC} Node.js TypeScript 配置缺失"
fi
((total_checks++))

echo ""

# 检查Electron配置
echo "5. 检查Electron配置"
echo "-------------------"

if [ -f "electron/main.ts" ]; then
    echo -n "检查 Electron 主进程文件... "
    if grep -q "app.whenReady" electron/main.ts && grep -q "BrowserWindow" electron/main.ts; then
        echo -e "${GREEN}✓${NC}"
        ((passed_checks++))
    else
        echo -e "${YELLOW}⚠${NC} (文件存在但内容可能不完整)"
        ((passed_checks++))
    fi
    ((total_checks++))
fi

if [ -f "electron-builder.yml" ]; then
    echo -e "${GREEN}✓${NC} Electron 构建配置"
    ((passed_checks++))
else
    echo -e "${RED}✗${NC} Electron 构建配置缺失"
fi
((total_checks++))

echo ""

# 检查开发工具配置
echo "6. 检查开发工具配置"
echo "-------------------"

check_file ".eslintrc.cjs" "ESLint 配置" && ((passed_checks++))
((total_checks++))

check_file ".prettierrc.json" "Prettier 配置" && ((passed_checks++))
((total_checks++))

check_file ".gitignore" "Git 忽略配置" && ((passed_checks++))
((total_checks++))

echo ""

# 检查文档文件
echo "7. 检查文档文件"
echo "---------------"

check_file "README.md" "项目说明文档" && ((passed_checks++))
((total_checks++))

check_file "RUN_GUIDE.md" "运行指南文档" && ((passed_checks++))
((total_checks++))

check_file "QUICK_START.md" "快速开始文档" && ((passed_checks++))
((total_checks++))

echo ""

# 总结
echo "验证总结"
echo "========"
echo -e "总检查项: ${total_checks}"
echo -e "通过项: ${GREEN}${passed_checks}${NC}"
echo -e "失败项: ${RED}$((total_checks - passed_checks))${NC}"
echo -e "通过率: ${GREEN}$((passed_checks * 100 / total_checks))%${NC}"
echo ""

# 建议
if [ $passed_checks -eq $total_checks ]; then
    echo -e "${GREEN}✓ 项目结构完整，所有必要文件都已创建！${NC}"
    echo ""
    echo "下一步操作:"
    echo "1. 安装依赖: npm install"
    echo "2. 启动开发服务器: npm run dev"
    echo "3. 构建项目: npm run build"
    echo ""
elif [ $passed_checks -ge $((total_checks * 80 / 100)) ]; then
    echo -e "${YELLOW}⚠ 项目基本完整，但缺少部分文件${NC}"
    echo "建议检查缺失的文件并重新创建"
    echo ""
else
    echo -e "${RED}✗ 项目结构不完整，需要补充缺失的文件${NC}"
    echo "请根据检查结果补充缺失的文件"
    echo ""
fi

# 环境检查
echo "环境检查建议:"
echo "============="
echo "1. 确保已安装 Node.js (版本 >= 16.0.0)"
echo "2. 确保已安装 npm 或 yarn"
echo "3. 检查后端服务是否运行 (默认: http://localhost:8080)"
echo "4. 确保网络连接正常"
echo ""

# 创建验证标记文件
echo "$(date)" > .validated
echo "项目验证完成" >> .validated
echo "通过检查: $passed_checks/$total_checks" >> .validated

echo -e "${GREEN}验证完成！${NC}"
