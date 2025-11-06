# Chat Frontend - 完整运行指南

## 项目概述

这是一个基于 Vue3 + TypeScript + Electron 的现代化聊天应用前端，支持实时聊天、用户管理、主题切换等功能。

## 技术栈

- **前端框架**: Vue 3.3+
- **编程语言**: TypeScript 5.0+
- **构建工具**: Vite 4.3+
- **状态管理**: Pinia 2.1+
- **路由管理**: Vue Router 4.2+
- **UI框架**: Element Plus 2.3+
- **HTTP客户端**: Axios 1.4+
- **桌面客户端**: Electron 25+
- **代码规范**: ESLint + Prettier

## 项目结构

```
chat-frontend/
├── electron/                 # Electron相关文件
│   ├── main.ts              # 主进程文件
│   └── preload/             # 预加载脚本
├── src/
│   ├── api/                 # API接口封装
│   ├── assets/              # 静态资源
│   ├── components/          # 公共组件
│   ├── composables/         # 组合式函数
│   ├── config/              # 配置文件
│   ├── router/              # 路由配置
│   ├── stores/              # 状态管理
│   ├── types/               # 类型定义
│   ├── utils/               # 工具函数
│   └── views/               # 页面组件
├── build/                   # 构建资源
├── .env                     # 环境变量
├── vite.config.ts           # Vite配置
└── electron-builder.yml     # Electron构建配置
```

## 快速开始

### 1. 环境要求

- Node.js >= 16.0.0
- npm >= 8.0.0 或 yarn >= 1.22.0

### 2. 安装依赖

```bash
# 进入项目目录
cd chat-frontend

# 安装依赖
npm install

# 或者使用yarn
yarn install
```

### 3. 环境配置

编辑 `.env` 文件配置API端点：

```env
# API基础URL
VITE_API_BASE_URL=http://localhost:8080/api

# WebSocket URL
VITE_WS_BASE_URL=ws://localhost:8080/ws

# 应用名称
VITE_APP_NAME=Chat Application

# 运行模式
VITE_NODE_ENV=development
```

### 4. 开发模式运行

```bash
# Web版本开发
npm run dev

# Electron版本开发
npm run electron:dev
```

### 5. 构建项目

```bash
# 构建Web版本
npm run build

# 构建Electron版本
npm run electron:build

# 预览构建结果
npm run preview
```

## 功能特性

### 1. 用户认证
- 用户注册/登录
- JWT Token认证
- 自动登录
- 权限控制

### 2. 实时聊天
- WebSocket连接
- 多聊天室支持
- 消息实时推送
- 在线用户列表

### 3. 用户管理
- 个人资料管理
- 头像上传
- 用户状态显示

### 4. 界面特性
- 响应式设计
- 主题切换（亮色/暗色/自动）
- 国际化支持
- 加载状态管理

### 5. 桌面客户端
- Electron打包
- 系统托盘
- 桌面通知
- 自动更新

## API接口对接

项目已配置好与后端chat-service的接口对接：

### 认证接口
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/verify` - Token验证
- `POST /api/auth/logout` - 用户登出

### 用户接口
- `GET /api/users/profile` - 获取用户信息
- `PUT /api/users/profile` - 更新用户信息
- `POST /api/users/avatar` - 上传头像

### 聊天接口
- `GET /api/chat/rooms` - 获取聊天室列表
- `POST /api/chat/rooms` - 创建聊天室
- `GET /api/chat/rooms/:id/messages` - 获取消息历史
- `POST /api/chat/rooms/:id/join` - 加入聊天室
- `POST /api/chat/rooms/:id/leave` - 离开聊天室

### WebSocket接口
- `ws://localhost:8080/ws` - WebSocket连接

## 开发指南

### 1. 组件开发

```vue
<template>
  <div class="my-component">
    <h1>{{ title }}</h1>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const title = ref('My Component')
</script>

<style scoped lang="scss">
.my-component {
  padding: 20px;
}
</style>
```

### 2. API调用

```typescript
import { getUserProfile } from '@/api/user'

const loadUserProfile = async () => {
  try {
    const response = await getUserProfile()
    console.log('用户信息:', response.data)
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
}
```

### 3. 状态管理

```typescript
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// 登录
await authStore.login({ username: 'user', password: 'pass' })

// 获取用户信息
const userInfo = authStore.userInfo
```

### 4. 路由跳转

```typescript
import { useRouter } from 'vue-router'

const router = useRouter()

// 跳转到聊天页面
router.push('/chat')

// 带参数跳转
router.push({ name: 'Profile', params: { id: 1 } })
```

## 构建配置

### Vite配置

项目使用Vite作为构建工具，配置文件 `vite.config.ts` 包含：
- Vue插件配置
- TypeScript支持
- 路径别名配置
- 开发服务器配置
- 构建优化配置

### Electron配置

Electron配置文件 `electron-builder.yml` 包含：
- 应用信息配置
- 构建目标平台
- 图标和启动配置
- 自动更新配置

## 部署指南

### Web版本部署

1. 构建项目：
```bash
npm run build
```

2. 部署dist目录到Web服务器
3. 配置Nginx/Apache反向代理

### Electron版本发布

1. 构建应用：
```bash
npm run electron:build
```

2. 发布到GitHub Releases
3. 配置自动更新服务器

## 常见问题

### 1. 依赖安装失败

- 清除npm缓存：`npm cache clean --force`
- 使用淘宝镜像：`npm config set registry https://registry.npmmirror.com`
- 删除node_modules重新安装

### 2. 构建失败

- 检查Node.js版本是否符合要求
- 检查内存是否足够（Electron构建需要较多内存）
- 尝试删除dist目录重新构建

### 3. WebSocket连接失败

- 检查后端服务是否启动
- 检查WebSocket URL配置是否正确
- 检查防火墙和网络设置

### 4. Electron打包失败

- 安装必要的构建工具：
```bash
# Ubuntu/Debian
sudo apt-get install build-essential

# macOS
xcode-select --install

# Windows
npm install --global windows-build-tools
```

## 开发脚本

```bash
# 开发模式
npm run dev

# 代码检查
npm run lint

# 代码格式化
npm run format

# 类型检查
npm run type-check

# 构建Web版本
npm run build

# 构建Electron版本
npm run electron:build

# 预览构建结果
npm run preview

# 清理构建文件
npm run clean
```

## 联系支持

如遇到问题，请检查：
1. 后端服务是否正常运行
2. API端点配置是否正确
3. 网络连接是否正常
4. 浏览器控制台错误信息

项目已完整配置，可以直接运行和构建。祝您开发愉快！
