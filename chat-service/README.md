# 聊天服务 (Chat Service)

一个高性能的实时聊天服务后端系统，支持单聊、群聊、历史消息等功能。

## 技术栈

- **Web框架**: Gin
- **数据库**: MySQL + GORM
- **缓存**: Redis
- **消息队列**: RabbitMQ
- **实时通信**: WebSocket
- **认证**: JWT

## 功能特性

- ✅ 用户注册/登录
- ✅ 单人聊天
- ✅ 群组聊天
- ✅ 历史聊天记录
- ✅ 实时消息推送
- ✅ 未读消息计数
- ✅ 在线状态管理
- ✅ 消息持久化
- ✅ 高并发支持 (10,000+ 同时在线)

## 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │    │  Mobile Client  │    │   Other Apps    │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │      API Gateway          │
                    │      (Gin Framework)      │
                    └─────────────┬─────────────┘
                                 │
          ┌──────────────────────┼──────────────────────┐
          │                      │                      │
    ┌─────┴─────┐        ┌───────┴───────┐      ┌───────┴───────┐
    │WebSocket  │        │    HTTP API   │      │ Message Queue │
    │   Hub     │        │   Controller  │      │  (RabbitMQ)   │
    └─────┬─────┘        └───────┬───────┘      └───────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │     Service Layer         │
                    └─────────────┬─────────────┘
                                 │
          ┌──────────────────────┼──────────────────────┐
          │                      │                      │
    ┌─────┴─────┐        ┌───────┴───────┐      ┌───────┴───────┐
    │   MySQL   │        │    Redis      │      │  Online Cache  │
    │ (GORM)    │        │   Cache       │      │   (Redis)      │
    └───────────┘        └───────────────┘      └────────────────┘
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- RabbitMQ 3.8+

### 安装依赖

```bash
make deps
```

### 配置文件

复制并编辑配置文件 `config.yaml`:

```yaml
server:
  host: "0.0.0.0"
  port: "8080"
  mode: "debug"

database:
  host: "localhost"
  port: "3306"
  user: "root"
  password: "your_password"
  dbname: "chat_service"

redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 0

rabbitmq:
  url: "amqp://guest:guest@localhost:5672/"
  exchange: "chat.exchange"
  queue: "chat.queue"

jwt:
  secret: "your-secret-key"
  expire_hour: 24
```

### 运行服务

```bash
# 开发环境
make run

# 或者
make dev

# 构建后运行
make build
./bin/chat-service
```

## API 文档

### 认证相关

#### 用户注册
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "nickname": "测试用户",
  "email": "test@example.com",
  "password": "123456"
}
```

#### 用户登录
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456"
}
```

### 聊天相关

#### 获取聊天室列表
```
GET /api/v1/rooms
Authorization: Bearer <token>
```

#### 创建聊天室
```
POST /api/v1/rooms
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "群聊名称",
  "description": "群聊描述",
  "type": "group",
  "member_ids": [2, 3, 4]
}
```

#### 获取聊天记录
```
GET /api/v1/rooms/{id}/messages?page=1&page_size=20
Authorization: Bearer <token>
```

### WebSocket 连接

```
ws://localhost:8080/api/v1/ws?token=<jwt_token>
```

WebSocket 消息格式:

```json
{
  "type": "message",
  "room_id": 1,
  "content": "Hello, World!"
}
```

支持的消息类型:
- `join_room`: 加入房间
- `leave_room`: 离开房间
- `message`: 发送消息

## 性能优化

### 数据库优化

- 合理的索引设计
- 读写分离
- 连接池优化

### 缓存策略

- Redis 缓存热点数据
- 最近消息缓存
- 用户在线状态缓存

### 消息队列

- 异步处理消息
- 消息持久化
- 消费者组

### 并发处理

- WebSocket 连接池
- 协程池管理
- 内存复用

## 部署

### Docker 部署

```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

### 生产环境部署

```bash
# 构建生产版本
make build-prod

# 运行
./bin/chat-service-linux-amd64
```

## 监控指标

- 并发连接数
- 消息吞吐量
- 系统资源使用
- 错误率统计

## 开发指南

### 项目结构

```
chat-service/
├── cmd/server/          # 应用入口
├── internal/
│   ├── api/            # API控制器和路由
│   ├── config/         # 配置管理
│   ├── database/       # 数据库连接
│   ├── middleware/     # 中间件
│   ├── models/         # 数据模型
│   ├── service/        # 业务逻辑
│   └── websocket/      # WebSocket处理
├── pkg/
│   ├── cache/          # Redis缓存
│   ├── queue/          # 消息队列
│   └── utils/          # 工具函数
├── config.yaml         # 配置文件
├── Makefile           # 构建脚本
├── Dockerfile         # Docker配置
└── README.md          # 项目文档
```

### 开发命令

```bash
# 运行测试
make test

# 代码格式化
make fmt

# 代码检查
make lint

# 测试覆盖率
make test-coverage
```

## 许可证

MIT License

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request
