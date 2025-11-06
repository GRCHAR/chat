#!/bin/bash

echo "🚀 Chat Frontend - 快速启动脚本"
echo "================================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查Node.js
 check_node() {
    if command -v node &> /dev/null; then
        NODE_VERSION=$(node --version)
        echo -e "${GREEN}✓${NC} Node.js 已安装: $NODE_VERSION"
        return 0
    else
        echo -e "${RED}✗${NC} Node.js 未安装，请先安装 Node.js (>= 16.0.0)"
        return 1
    fi
}

# 检查npm
 check_npm() {
    if command -v npm &> /dev/null; then
        NPM_VERSION=$(npm --version)
        echo -e "${GREEN}✓${NC} npm 已安装: $NPM_VERSION"
        return 0
    else
        echo -e "${RED}✗${NC} npm 未安装，请先安装 npm"
        return 1
    fi
}

# 检查依赖
 check_dependencies() {
    if [ -d "node_modules" ]; then
        echo -e "${GREEN}✓${NC} 项目依赖已安装"
        return 0
    else
        echo -e "${YELLOW}⚠${NC} 项目依赖未安装，需要运行 npm install"
        return 1
    fi
}

# 检查后端服务
 check_backend() {
    echo -n "检查后端服务状态... "
    if curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/health 2>/dev/null | grep -q "200\|404"; then
        echo -e "${GREEN}✓${NC} 后端服务运行正常"
        return 0
    else
        echo -e "${YELLOW}⚠${NC} 后端服务未运行 (http://localhost:8080)"
        echo "  请确保后端服务已启动，或修改 .env 文件中的 API 地址"
        return 1
    fi
}

# 显示菜单
 show_menu() {
    echo ""
    echo "请选择操作:"
    echo "1) 🌐 启动 Web 开发服务器"
    echo "2) 🖥️  启动 Electron 客户端"
    echo "3) 📦 安装项目依赖"
    echo "4) 🔧 检查项目环境"
    echo "5) 📊 验证项目完整性"
    echo "6) 🏗️  构建项目"
    echo "7) 📖 查看项目文档"
    echo "8) ❌ 退出"
    echo ""
    echo -n "请输入选项 (1-8): "
}

# 启动Web开发服务器
 start_web_dev() {
    echo -e "${BLUE}🌐 启动 Web 开发服务器...${NC}"
    if npm run dev; then
        echo -e "${GREEN}✓ Web 开发服务器启动成功！${NC}"
        echo "请访问: http://localhost:5173"
    else
        echo -e "${RED}✗ Web 开发服务器启动失败${NC}"
        echo "请检查错误信息并确保依赖已安装"
    fi
}

# 启动Electron客户端
 start_electron() {
    echo -e "${BLUE}🖥️  启动 Electron 客户端...${NC}"
    if npm run electron:dev; then
        echo -e "${GREEN}✓ Electron 客户端启动成功！${NC}"
    else
        echo -e "${RED}✗ Electron 客户端启动失败${NC}"
        echo "请检查错误信息并确保依赖已安装"
    fi
}

# 安装依赖
 install_dependencies() {
    echo -e "${BLUE}📦 安装项目依赖...${NC}"
    echo "这可能需要几分钟时间，请耐心等待..."
    
    if npm install; then
        echo -e "${GREEN}✓ 依赖安装成功！${NC}"
    else
        echo -e "${RED}✗ 依赖安装失败${NC}"
        echo "请检查网络连接和npm配置"
    fi
}

# 检查项目环境
 check_environment() {
    echo -e "${BLUE}🔧 检查项目环境...${NC}"
    echo ""
    
    check_node
    check_npm
    check_dependencies
    check_backend
    
    echo ""
    echo -e "${BLUE}项目信息:${NC}"
    echo "- 项目名称: Chat Frontend"
    echo "- 技术栈: Vue3 + TypeScript + Electron"
    echo "- 默认端口: 5173 (Web)"
    echo "- API地址: http://localhost:8080/api"
    echo "- WebSocket: ws://localhost:8080/ws"
}

# 验证项目完整性
 validate_project() {
    echo -e "${BLUE}📊 验证项目完整性...${NC}"
    if [ -f "validate-project.sh" ]; then
        ./validate-project.sh
    else
        echo -e "${RED}✗ 验证脚本不存在${NC}"
    fi
}

# 构建项目
 build_project() {
    echo -e "${BLUE}🏗️  构建项目...${NC}"
    echo "请选择构建类型:"
    echo "1) 🌐 构建 Web 版本"
    echo "2) 🖥️  构建 Electron 版本"
    echo "3) 🔙 返回主菜单"
    echo ""
    echo -n "请输入选项 (1-3): "
    
    read build_choice
    case $build_choice in
        1)
            echo -e "${BLUE}构建 Web 版本...${NC}"
            if npm run build; then
                echo -e "${GREEN}✓ Web 版本构建成功！${NC}"
                echo "构建文件位于: dist/"
            else
                echo -e "${RED}✗ Web 版本构建失败${NC}"
            fi
            ;;
        2)
            echo -e "${BLUE}构建 Electron 版本...${NC}"
            if npm run electron:build; then
                echo -e "${GREEN}✓ Electron 版本构建成功！${NC}"
                echo "构建文件位于: dist-electron/"
            else
                echo -e "${RED}✗ Electron 版本构建失败${NC}"
            fi
            ;;
        3)
            return
            ;;
        *)
            echo -e "${RED}无效选项${NC}"
            ;;
    esac
}

# 查看文档
 show_documentation() {
    echo -e "${BLUE}📖 项目文档${NC}"
    echo ""
    echo "可用文档:"
    echo "1) 📋 项目总结 (PROJECT_SUMMARY.md)"
    echo "2) 📖 运行指南 (RUN_GUIDE.md)"
    echo "3) 🚀 快速开始 (QUICK_START.md)"
    echo "4) 📄 README.md"
    echo "5) 🔙 返回主菜单"
    echo ""
    echo -n "请输入选项 (1-5): "
    
    read doc_choice
    case $doc_choice in
        1)
            less PROJECT_SUMMARY.md 2>/dev/null || cat PROJECT_SUMMARY.md
            ;;
        2)
            less RUN_GUIDE.md 2>/dev/null || cat RUN_GUIDE.md
            ;;
        3)
            less QUICK_START.md 2>/dev/null || cat QUICK_START.md
            ;;
        4)
            less README.md 2>/dev/null || cat README.md
            ;;
        5)
            return
            ;;
        *)
            echo -e "${RED}无效选项${NC}"
            ;;
    esac
}

# 主程序
 main() {
    echo -e "${BLUE}环境检查:${NC}"
    check_node
    check_npm
    echo ""
    
    # 如果依赖未安装，提示安装
    if ! check_dependencies; then
        echo ""
        echo -e "${YELLOW}项目依赖未安装，建议先安装依赖${NC}"
        echo -n "是否现在安装依赖? (y/n): "
        read install_choice
        if [ "$install_choice" = "y" ] || [ "$install_choice" = "Y" ]; then
            install_dependencies
        fi
    fi
    
    while true; do
        show_menu
        read choice
        
        case $choice in
            1)
                start_web_dev
                ;;
            2)
                start_electron
                ;;
            3)
                install_dependencies
                ;;
            4)
                check_environment
                ;;
            5)
                validate_project
                ;;
            6)
                build_project
                ;;
            7)
                show_documentation
                ;;
            8)
                echo -e "${GREEN}感谢使用 Chat Frontend！再见！${NC}"
                exit 0
                ;;
            *)
                echo -e "${RED}无效选项，请输入 1-8${NC}"
                ;;
        esac
        
        echo ""
        echo -n "按回车键继续..."
        read
        clear
        echo "🚀 Chat Frontend - 快速启动脚本"
        echo "================================="
        echo ""
    done
}

# 脚本入口
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "Chat Frontend 快速启动脚本"
    echo ""
    echo "用法: ./start.sh [选项]"
    echo ""
    echo "选项:"
    echo "  --help, -h     显示帮助信息"
    echo "  --web          直接启动 Web 开发服务器"
    echo "  --electron     直接启动 Electron 客户端"
    echo "  --install      直接安装依赖"
    echo "  --check        直接检查项目环境"
    echo "  --validate     直接验证项目完整性"
    echo "  --build        直接构建项目"
    echo ""
    echo "无参数运行时进入交互式菜单"
    exit 0
fi

# 处理命令行参数
case "$1" in
    --web)
        check_dependencies && start_web_dev
        ;;
    --electron)
        check_dependencies && start_electron
        ;;
    --install)
        install_dependencies
        ;;
    --check)
        check_environment
        ;;
    --validate)
        validate_project
        ;;
    --build)
        build_project
        ;;
    "")
        main
        ;;
    *)
        echo -e "${RED}未知选项: $1${NC}"
        echo "使用 --help 查看可用选项"
        exit 1
        ;;
esac
