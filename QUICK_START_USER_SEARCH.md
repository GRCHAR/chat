# 用户搜索和私聊功能 - 快速启动指南

## 前置条件

- Go 1.16+ 已安装
- Node.js 14+ 已安装
- MySQL 5.7+ 已安装并运行
- Git 已安装

## 快速启动步骤

### 1. 启动后端服务

```bash
# 进入后端目录
cd /data/workspace/chat/chat-service

# 确保数据库已初始化
mysql -u root -p < init.sql

# 启动服务
make run
# 或者
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动

### 2. 启动前端服务

```bash
# 进入前端目录
cd /data/workspace/chat/chat-frontend

# 安装依赖（首次运行）
npm install

# 启动开发服务器
npm run dev
```

前端服务将在 `http://localhost:5173` 启动

### 3. 测试新功能

#### 方式一：使用自动化测试脚本

```bash
cd /data/workspace/chat
./test-user-search.sh
```

#### 方式二：手动测试

1. 打开浏览器访问 `http://localhost:5173`
2. 注册两个测试账号
3. 使用第一个账号登录
4. 点击"搜索用户"按钮
5. 搜索第二个用户
6. 测试以下功能：
   - 发起私聊
   - 创建群聊
   - 添加成员到群聊

## 功能使用说明

### 搜索用户并发起私聊

1. 在聊天界面，点击左侧边栏的 **"搜索用户"** 按钮
2. 在弹出的对话框中输入用户名、昵称或邮箱
3. 等待搜索结果显示（自动防抖，300ms延迟）
4. 在搜索结果中找到目标用户
5. 点击 **"私聊"** 按钮
6. 系统自动创建私聊房间并跳转

### 添加用户到群聊

1. 在聊天界面，点击 **"搜索用户"** 按钮
2. 搜索要添加的用户
3. 点击用户旁边的 **"添加到群聊"** 按钮
4. 在弹出的对话框中选择目标群聊
5. 点击 **"添加"** 按钮完成

### 创建群聊并添加成员

1. 点击 **"新建群聊"** 按钮
2. 选择房间类型为 **"群聊"**
3. 输入房间名称（必填）
4. 输入房间描述（可选）
5. 在成员选择区域输入关键词搜索用户
6. 勾选要添加的成员
7. 点击 **"创建"** 按钮

## API测试示例

### 使用 curl 测试

```bash
# 1. 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "nickname": "测试用户",
    "email": "test@example.com",
    "password": "password123"
  }'

# 2. 登录获取token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }' | jq -r '.token')

# 3. 搜索用户
curl -X GET "http://localhost:8080/api/v1/users/search?q=test" \
  -H "Authorization: Bearer $TOKEN"

# 4. 创建私聊
curl -X POST http://localhost:8080/api/v1/rooms \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "私聊",
    "type": "single",
    "member_ids": [2]
  }'

# 5. 添加成员到群聊
curl -X POST http://localhost:8080/api/v1/rooms/1/members \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 3
  }'
```

### 使用 Postman 测试

1. 导入以下环境变量：
   - `base_url`: `http://localhost:8080/api/v1`
   - `token`: 登录后获取的JWT token

2. 测试接口：
   - POST `/auth/register` - 注册
   - POST `/auth/login` - 登录
   - GET `/users/search?q=keyword` - 搜索用户
   - GET `/users/:id` - 获取用户详情
   - POST `/rooms` - 创建房间
   - POST `/rooms/:id/members` - 添加成员

## 常见问题

### Q1: 搜索不到用户？
**A:** 确保：
- 用户已注册
- 搜索关键词正确
- 后端服务正常运行
- 数据库连接正常

### Q2: 创建私聊失败？
**A:** 检查：
- 目标用户ID是否正确
- 是否有权限创建房间
- 网络连接是否正常

### Q3: 添加成员失败？
**A:** 确认：
- 目标房间是群聊类型
- 用户未在房间中
- 有足够的权限

### Q4: 前端无法连接后端？
**A:** 检查：
- 后端服务是否启动
- 端口是否被占用
- CORS配置是否正确
- 前端配置的API地址是否正确

## 配置说明

### 后端配置 (config.yaml)

```yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 3306
  user: root
  password: your_password
  dbname: chat_db

jwt:
  secret: your_jwt_secret
  expire: 24h
```

### 前端配置 (.env)

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_WS_URL=ws://localhost:8080/api/v1/ws
```

## 开发调试

### 后端调试

```bash
# 查看日志
tail -f logs/app.log

# 使用调试模式
go run -race cmd/server/main.go
```

### 前端调试

```bash
# 开启详细日志
npm run dev -- --debug

# 构建生产版本
npm run build

# 预览生产版本
npm run preview
```

## 性能优化建议

1. **搜索优化**
   - 添加索引：`CREATE INDEX idx_username ON users(username)`
   - 使用全文搜索
   - 实现搜索结果缓存

2. **前端优化**
   - 使用虚拟滚动处理大量搜索结果
   - 实现搜索历史记录
   - 添加搜索结果分页

3. **后端优化**
   - 使用Redis缓存热门搜索
   - 实现搜索结果预加载
   - 优化数据库查询

## 安全注意事项

1. 确保JWT密钥安全
2. 使用HTTPS传输敏感数据
3. 实现请求频率限制
4. 验证用户输入
5. 防止SQL注入

## 下一步

- 查看 [USER_SEARCH_FEATURE.md](./USER_SEARCH_FEATURE.md) 了解详细功能说明
- 查看 [IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md) 了解实现细节
- 运行 `./test-user-search.sh` 进行自动化测试

## 获取帮助

如有问题，请：
1. 查看日志文件
2. 检查配置文件
3. 运行测试脚本
4. 查看文档

祝使用愉快！🎉
