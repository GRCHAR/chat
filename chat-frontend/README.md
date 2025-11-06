# Chat Frontend

基于 Vue3 + TypeScript + Electron 的聊天应用前端项目。

## 功能特性

- 🚀 **现代化技术栈**: Vue3 + TypeScript + Vite
- 💬 **实时聊天**: WebSocket 实时通信
- 👥 **多人聊天**: 支持私聊和群聊
- 📱 **跨平台**: Electron 桌面客户端
- 🎨 **美观界面**: Element Plus UI 组件库
- 🔐 **用户认证**: JWT 认证机制
- 📱 **响应式设计**: 适配移动端和桌面端

## 技术架构

### 前端技术
- **Vue 3**: 渐进式 JavaScript 框架
- **TypeScript**: 类型安全的 JavaScript
- **Vite**: 快速的构建工具
- **Vue Router**: 路由管理
- **Pinia**: 状态管理
- **Element Plus**: UI 组件库
- **Axios**: HTTP 客户端
- **Socket.io-client**: WebSocket 客户端

### 桌面客户端
- **Electron**: 跨平台桌面应用框架
- **Electron Builder**: 应用打包工具

## 项目结构

```
chat-frontend/
├── src/                    # 源代码目录
│   ├── api/               # API 接口封装
│   ├── assets/            # 静态资源
│   ├── components/        # 公共组件
│   ├── composables/       # 组合式函数
│   ├── router/            # 路由配置
│   ├── stores/            # 状态管理
│   ├── types/             # TypeScript 类型定义
│   ├── utils/             # 工具函数
│   ├── views/             # 页面组件
│   ├── App.vue            # 根组件
│   └── main.ts            # 应用入口
├── electron/              # Electron 相关文件
│   ├── main.ts            # 主进程
│   └── preload/           # 预加载脚本
├── public/                # 公共资源
├── dist/                  # 构建输出目录
├── dist-electron/         # Electron 构建输出
└── package.json           # 项目配置
```

## 快速开始

### 环境要求
- Node.js >= 16.0.0
- npm 或 yarn

### 安装依赖
```bash
npm install
```

### 开发模式
```bash
# 启动 Web 开发服务器
npm run dev

# 启动 Electron 开发模式
npm run electron:dev
```

### 构建项目
```bash
# 构建 Web 应用
npm run build

# 构建 Electron 应用
npm run electron:build

# 构建特定平台
npm run electron:build:win    # Windows
npm run electron:build:mac    # macOS
npm run electron:build:linux  # Linux
```

## 功能模块

### 用户认证
- 用户注册
- 用户登录
- JWT Token 管理
- 用户信息管理

### 聊天功能
- 实时消息发送/接收
- 私聊和群聊
- 消息历史记录
- 未读消息提醒
- 消息已读状态

### 房间管理
- 创建聊天房间
- 加入/离开房间
- 房间成员管理
- 房间信息查看

### 桌面客户端
- 系统托盘
- 桌面通知
- 快捷键支持
- 自动更新

## API 接口

项目使用以下后端 API 接口：

### 认证相关
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录

### 用户相关
- `GET /api/v1/users/profile` - 获取用户信息
- `PUT /api/v1/users/profile` - 更新用户信息

### 聊天相关
- `GET /api/v1/rooms` - 获取房间列表
- `POST /api/v1/rooms` - 创建房间
- `GET /api/v1/rooms/:id/messages` - 获取消息历史
- `POST /api/v1/rooms/:id/read` - 标记已读
- `WebSocket /api/v1/ws` - 实时消息通信

## 配置说明

### 环境变量
- `VITE_API_BASE_URL` - API 基础地址
- `VITE_WS_URL` - WebSocket 地址

### Vite 配置
项目使用 Vite 作为构建工具，配置文件为 `vite.config.ts`。

### Electron 配置
Electron 相关配置在 `electron-builder.yml` 文件中。

## 开发指南

### 代码规范
- 使用 ESLint 进行代码检查
- 使用 Prettier 进行代码格式化
- 遵循 Vue3 组合式 API 规范

### 组件开发
- 使用 `<script setup>` 语法糖
- 合理使用 TypeScript 类型
- 遵循单一职责原则

### 状态管理
- 使用 Pinia 进行状态管理
- 按模块划分 store
- 避免过度使用全局状态

### API 封装
- 统一错误处理
- 请求/响应拦截器
- 类型安全的接口定义

## 部署说明

### Web 应用部署
1. 构建项目：`npm run build`
2. 将 `dist` 目录部署到 Web 服务器
3. 配置 API 代理

### 桌面应用发布
1. 构建应用：`npm run electron:build`
2. 在 `dist` 目录中找到安装包
3. 分发安装包给用户

## 常见问题

### Q: WebSocket 连接失败？
A: 请检查后端服务是否正常运行，以及 WebSocket 地址配置是否正确。

### Q: Electron 应用无法打包？
A: 确保已安装所有依赖，并检查 electron-builder 配置是否正确。

### Q: 跨域问题如何解决？
A: 开发模式下 Vite 已配置代理，生产环境需要在后端配置 CORS。

## 更新日志

### v1.0.0 (2023-12-01)
- 初始版本发布
- 实现基本聊天功能
- 支持私聊和群聊
- 添加 Electron 桌面客户端

## 贡献指南

欢迎提交 Issue 和 Pull Request 来改进项目。

### 提交规范
- 使用清晰的提交信息
- 遵循代码规范
- 添加必要的测试

## 许可证

MIT License

## 联系方式

如有问题或建议，请通过以下方式联系：
- 提交 Issue
- 发送邮件

---

**注意**: 本项目为演示项目，生产环境使用前请进行充分的测试和安全检查。
