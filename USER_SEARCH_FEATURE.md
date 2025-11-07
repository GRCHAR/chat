# 用户搜索和私聊/群聊功能

## 功能概述

本次更新为聊天系统添加了以下功能：

1. **用户搜索功能**：根据用户名、昵称或邮箱搜索用户
2. **发起私聊**：搜索到用户后可以直接发起私聊
3. **添加成员到群聊**：将搜索到的用户添加到现有群聊中
4. **创建群聊时搜索成员**：在创建群聊时可以搜索并选择成员

## 后端更新

### 新增API接口

#### 1. 搜索用户
- **路径**: `GET /api/v1/users/search`
- **参数**: `q` (查询关键词)
- **返回**: 用户列表
- **说明**: 支持按用户名、昵称、邮箱模糊搜索

```bash
curl -X GET "http://localhost:8080/api/v1/users/search?q=test" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 2. 获取用户详情
- **路径**: `GET /api/v1/users/:id`
- **参数**: `id` (用户ID)
- **返回**: 用户详细信息

```bash
curl -X GET "http://localhost:8080/api/v1/users/1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 3. 添加成员到房间
- **路径**: `POST /api/v1/rooms/:id/members`
- **参数**: `user_id` (要添加的用户ID)
- **返回**: 成功消息
- **说明**: 只能向群聊添加成员

```bash
curl -X POST "http://localhost:8080/api/v1/rooms/1/members" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 2}'
```

### 后端代码修改

1. **controllers.go**
   - 添加 `SearchUsers` 方法：搜索用户
   - 添加 `GetUserByID` 方法：获取用户详情
   - 添加 `AddMember` 方法：添加成员到房间

2. **chat_service.go**
   - 添加 `SearchUsers` 方法：实现用户搜索逻辑

3. **routes.go**
   - 添加用户搜索路由
   - 添加获取用户详情路由
   - 添加添加成员路由

## 前端更新

### 新增组件

#### 1. UserSearchDialog.vue
用户搜索对话框组件，提供以下功能：
- 实时搜索用户（防抖处理）
- 显示搜索结果列表
- 发起私聊
- 添加用户到群聊

**使用方式**：
```vue
<UserSearchDialog
  v-model="showDialog"
  @chat-created="handleChatCreated"
/>
```

#### 2. AddMemberDialog.vue
添加成员到群聊对话框组件，提供以下功能：
- 选择目标群聊
- 搜索用户
- 将用户添加到选定的群聊

**使用方式**：
```vue
<AddMemberDialog
  v-model="showDialog"
  :pre-selected-user="selectedUser"
  @member-added="handleMemberAdded"
/>
```

### 修改的组件

#### 1. ChatView.vue
- 添加"搜索用户"按钮
- 集成 UserSearchDialog 组件
- 将"新建聊天"改为"新建群聊"以区分功能

#### 2. CreateRoomForm.vue
- 实现实时用户搜索
- 添加防抖处理
- 优化用户选择体验
- 添加空状态提示

### 新增API文件

#### user.ts
提供用户相关的API接口：
- `searchUsers`: 搜索用户
- `getUserById`: 获取用户详情
- `getUsers`: 获取用户列表

## 使用说明

### 1. 搜索用户并发起私聊

1. 在聊天界面点击"搜索用户"按钮
2. 在搜索框中输入用户名、昵称或邮箱
3. 从搜索结果中选择用户
4. 点击"私聊"按钮即可创建私聊房间

### 2. 添加用户到群聊

1. 在聊天界面点击"搜索用户"按钮
2. 搜索要添加的用户
3. 点击"添加到群聊"按钮
4. 选择目标群聊
5. 点击"添加"按钮完成操作

### 3. 创建群聊时选择成员

1. 点击"新建群聊"按钮
2. 选择房间类型为"群聊"
3. 输入房间名称和描述
4. 在成员选择区域搜索用户
5. 勾选要添加的成员
6. 点击"创建"按钮

## 技术特点

### 1. 实时搜索
- 使用防抖技术，避免频繁请求
- 默认延迟300ms后发起搜索请求

### 2. 用户体验优化
- 加载状态提示
- 空状态提示
- 错误处理和提示
- 搜索结果实时更新

### 3. 权限控制
- 所有接口都需要JWT认证
- 只能向群聊添加成员
- 自动过滤当前用户

## 数据库查询

用户搜索使用模糊查询：
```sql
SELECT * FROM users 
WHERE username LIKE '%keyword%' 
   OR nickname LIKE '%keyword%' 
   OR email LIKE '%keyword%'
LIMIT 20
```

## 注意事项

1. **搜索限制**：每次搜索最多返回20个用户
2. **私聊创建**：如果与某用户的私聊已存在，会创建新的私聊房间
3. **群聊成员**：添加成员时会检查用户是否已在群聊中
4. **权限验证**：所有操作都需要用户登录认证

## 未来改进

1. 添加用户在线状态显示
2. 支持批量添加成员
3. 添加最近联系人功能
4. 优化搜索算法（如支持拼音搜索）
5. 添加用户标签和分组功能
6. 检查私聊是否已存在，避免重复创建

## 测试建议

1. 测试用户搜索功能的准确性
2. 测试私聊创建流程
3. 测试添加成员到群聊的权限控制
4. 测试搜索防抖功能
5. 测试各种边界情况（空搜索、无结果等）
