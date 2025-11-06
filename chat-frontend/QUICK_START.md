# 快速开始指南

本指南将帮助您快速启动和运行聊天应用前端项目。

## 环境准备

### 必需环境
- Node.js >= 16.0.0
- npm >= 7.0.0 或 yarn >= 1.22.0

### 可选环境
- Git (用于版本控制)
- VS Code (推荐编辑器)

## 安装步骤

### 1. 克隆项目
```bash
git clone <repository-url>
cd chat-frontend
```

### 2. 安装依赖
```bash
npm install
```

### 3. 启动开发服务器
```bash
# 启动 Web 开发服务器
npm run dev

# 启动 Electron 桌面应用开发模式
npm run electron:dev
```

### 4. 访问应用
- Web 应用: http://localhost:3000
- Electron 应用: 自动打开桌面窗口

## 开发流程

### 1. 项目结构
```
src/
├── api/          # API 接口封装
├── assets/       # 静态资源
├── components/   # 公共组件
├── composables/  # 组合式函数
├── router/       # 路由配置
├── stores/       # 状态管理
├── types/        # TypeScript 类型
├── utils/        # 工具函数
└── views/        # 页面组件
```

### 2. 添加新页面
1. 在 `src/views/` 目录下创建新的 Vue 组件
2. 在 `src/router/index.ts` 中添加路由配置
3. 在需要的地方添加导航链接

### 3. 添加新组件
1. 在 `src/components/` 目录下创建组件文件
2. 遵循组件命名规范 (PascalCase)
3. 使用 TypeScript 和组合式 API

### 4. 状态管理
1. 在 `src/stores/` 目录下创建新的 store
2. 使用 Pinia 进行状态管理
3. 遵循模块化设计原则

### 5. API 接口
1. 在 `src/api/` 目录下创建 API 模块
2. 使用统一的请求封装
3. 添加 TypeScript 类型定义

## 常用命令

### 开发命令
```bash
npm run dev              # 启动开发服务器
npm run electron:dev     # 启动 Electron 开发模式
```

### 构建命令
```bash
npm run build            # 构建 Web 应用
npm run electron:build   # 构建 Electron 应用
```

### 代码质量
```bash
npm run lint             # 运行 ESLint 检查
npm run format           # 格式化代码
npm run type-check       # TypeScript 类型检查
```

## 配置说明

### 环境变量
在 `.env` 文件中配置：
```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_WS_URL=ws://localhost:8080/api/ws
```

### Vite 配置
`vite.config.ts` - 构建工具配置
`vite.config.electron.ts` - Electron 构建配置

### Electron 配置
`electron-builder.yml` - 打包配置
`electron/main.ts` - 主进程代码
`electron/preload/index.ts` - 预加载脚本

## 开发规范

### 代码风格
- 使用 TypeScript 严格模式
- 遵循 Vue3 组合式 API 规范
- 使用 Element Plus 组件库
- 遵循 ESLint 和 Prettier 配置

### 文件命名
- 组件文件: PascalCase (如: `ChatView.vue`)
- 工具函数: camelCase (如: `formatTime.ts`)
- 类型定义: PascalCase (如: `User.ts`)

### 组件开发
```vue
<template>
  <div class="component-name">
    <!-- 模板内容 -->
  </div>
</template>

<script setup lang="ts">
// TypeScript 代码
import { ref } from 'vue'

const props = defineProps<{
  // 属性定义
}>()

const emit = defineEmits<{
  // 事件定义
}>()
</script>

<style lang="scss" scoped>
// 样式代码
</style>
```

## 常见问题

### Q: 安装依赖失败？
A: 尝试清除缓存后重新安装：
```bash
npm cache clean --force
rm -rf node_modules
npm install
```

### Q: 开发服务器无法启动？
A: 检查端口是否被占用，或查看控制台错误信息。

### Q: Electron 应用无法打开？
```bash
# 重新构建 Electron
npm run electron:build

# 或尝试清除缓存
rm -rf dist-electron
npm run electron:dev
```

### Q: TypeScript 类型错误？
A: 运行类型检查命令：
```bash
npm run type-check
```

### Q: 代码格式问题？
A: 运行格式化命令：
```bash
npm run format
```

## 调试技巧

### 1. Vue DevTools
安装 Vue DevTools 浏览器扩展进行调试。

### 2. Electron DevTools
在 Electron 应用中按 `F12` 打开开发者工具。

### 3. 日志调试
```typescript
console.log('调试信息')
console.error('错误信息')
```

### 4. 断点调试
在 VS Code 中设置断点进行调试。

## 性能优化

### 1. 组件懒加载
```typescript
const ChatView = () => import('@/views/ChatView.vue')
```

### 2. 虚拟滚动
对于长列表使用虚拟滚动优化性能。

### 3. 图片优化
使用适当的图片格式和大小。

### 4. 代码分割
利用 Vite 的代码分割功能。

## 部署指南

### Web 应用部署
1. 构建项目：`npm run build`
2. 将 `dist` 目录部署到 Web 服务器
3. 配置反向代理到后端 API

### Electron 应用发布
1. 构建应用：`npm run electron:build`
2. 在 `dist` 目录中找到安装包
3. 分发安装包给用户

## 获取更多帮助

- 查看项目文档
- 提交 Issue
- 查看示例代码
- 参考官方文档

---

祝您开发愉快！🚀
