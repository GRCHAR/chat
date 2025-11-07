<template>
  <div class="chat-container">
    <!-- 侧边栏 -->
    <div class="sidebar">
      <div class="sidebar-header">
        <div class="user-info">
          <el-avatar :size="40" :src="authStore.user?.avatar">
            {{ authStore.user?.nickname?.charAt(0) }}
          </el-avatar>
          <div class="user-details">
            <span class="nickname">{{ authStore.user?.nickname }}</span>
            <span class="username">@{{ authStore.user?.username }}</span>
          </div>
        </div>
        <el-dropdown @command="handleCommand">
          <el-icon class="more-icon"><MoreFilled /></el-icon>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人信息</el-dropdown-item>
              <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
      
      <div class="sidebar-actions">
        <el-button type="primary" style="width: 100%" @click="showCreateRoomDialog = true">
          <el-icon><Plus /></el-icon>
          新建聊天
        </el-button>
      </div>
      
      <div class="room-list">
        <div class="room-list-header">
          <h3>聊天列表</h3>
          <el-badge :value="chatStore.totalUnreadCount" :hidden="chatStore.totalUnreadCount === 0">
            <el-icon><ChatDotRound /></el-icon>
          </el-badge>
        </div>
        
        <div class="room-items">
          <div
            v-for="room in chatStore.rooms"
            :key="room.id"
            class="room-item"
            :class="{ active: chatStore.currentRoom?.id === room.id }"
            @click="selectRoom(room)"
          >
            <div class="room-avatar">
              <el-avatar :size="40">
                {{ room.name.charAt(0).toUpperCase() }}
              </el-avatar>
            </div>
            <div class="room-info">
              <div class="room-name">{{ room.name }}</div>
              <div class="room-last-message">{{ room.last_message?.content || '暂无消息' }}</div>
            </div>
            <div class="room-meta">
              <div class="room-time">{{ formatTime(room.last_message?.created_at || '') }}</div>
              <el-badge
                :value="chatStore.unreadCounts[room.id] || 0"
                :hidden="!chatStore.unreadCounts[room.id]"
                class="unread-badge"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 聊天区域 -->
    <div class="chat-area" v-if="chatStore.currentRoom">
      <div class="chat-header">
        <div class="chat-title">
          <h3>{{ chatStore.currentRoom.name }}</h3>
          <span class="room-type">{{ chatStore.currentRoom.type === 'single' ? '私聊' : '群聊' }}</span>
        </div>
        <div class="chat-actions">
          <el-button circle @click="showRoomMembers">
            <el-icon><User /></el-icon>
          </el-button>
          <el-button circle @click="showRoomInfo">
            <el-icon><InfoFilled /></el-icon>
          </el-button>
        </div>
      </div>
      
      <div class="chat-messages" ref="messagesContainer">
        <div class="messages-container">
          <div
            v-for="message in chatStore.messages"
            :key="message.id"
            class="message-item"
            :class="{ 'message-mine': message.sender_id === authStore.user?.id }"
          >
            <div class="message-avatar">
              <el-avatar :size="32" :src="message.sender?.avatar">
                {{ message.sender?.nickname?.charAt(0) }}
              </el-avatar>
            </div>
            <div class="message-content">
              <div class="message-header">
                <span class="message-sender">{{ message.sender?.nickname }}</span>
                <span class="message-time">{{ formatTime(message.created_at) }}</span>
              </div>
              <div class="message-body">
                <div class="message-text">{{ message.content }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="chat-input">
        <div class="input-toolbar">
          <el-button text @click="selectEmoji">
            <el-icon><Avatar /></el-icon>
          </el-button>
          <el-button text @click="selectFile">
            <el-icon><Paperclip /></el-icon>
          </el-button>
        </div>
        <div class="input-area">
          <el-input
            v-model="messageInput"
            type="textarea"
            placeholder="输入消息..."
            :rows="3"
            resize="none"
            @keydown.enter.prevent="handleSendMessage"
          />
          <el-button type="primary" @click="handleSendMessage" :disabled="!messageInput.trim()">
            发送
          </el-button>
        </div>
      </div>
    </div>
    
    <!-- 空状态 -->
    <div class="empty-state" v-else>
      <el-empty description="选择一个聊天开始对话">
        <el-button type="primary" @click="showCreateRoomDialog = true">
          新建聊天
        </el-button>
      </el-empty>
    </div>
    
    <!-- 创建房间对话框 -->
    <el-dialog
      v-model="showCreateRoomDialog"
      title="新建群聊"
      width="500px"
    >
      <CreateRoomForm @success="handleCreateRoomSuccess" @cancel="showCreateRoomDialog = false" />
    </el-dialog>
    
    <!-- 用户搜索对话框 -->
    <UserSearchDialog
      v-model="showUserSearchDialog"
      @chat-created="handleChatCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { useChatStore } from '@/stores/chat'
import { formatTime } from '@/utils/format'
import CreateRoomForm from '@/components/CreateRoomForm.vue'
import UserSearchDialog from '@/components/UserSearchDialog.vue'
import type { ChatRoom } from '@/types/chat'

const router = useRouter()
const authStore = useAuthStore()
const chatStore = useChatStore()

const messagesContainer = ref<HTMLElement>()
const messageInput = ref('')
const showCreateRoomDialog = ref(false)
const showUserSearchDialog = ref(false)

// 选择房间
const selectRoom = async (room: ChatRoom) => {
  chatStore.setCurrentRoom(room)
  await chatStore.markAsRead(room.id)
  scrollToBottom()
}

// 发送消息
const handleSendMessage = async () => {
  if (!messageInput.value.trim()) return
  
  const result = await chatStore.sendChatMessage(messageInput.value.trim())
  if (result.success) {
    messageInput.value = ''
    scrollToBottom()
  } else {
    ElMessage.error(result.error || '发送失败')
  }
}

// 滚动到底部
const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

// 处理菜单命令
const handleCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      authStore.logout()
      router.push('/')
      break
  }
}

// 创建房间成功
const handleCreateRoomSuccess = () => {
  showCreateRoomDialog.value = false
  chatStore.fetchRooms()
}

// 私聊创建成功
const handleChatCreated = (roomId: number) => {
  showUserSearchDialog.value = false
  // 房间已经在UserSearchDialog中被选中
}

// 显示房间成员
const showRoomMembers = () => {
  ElMessage.info('房间成员功能开发中...')
}

// 显示房间信息
const showRoomInfo = () => {
  ElMessage.info('房间信息功能开发中...')
}

// 选择表情
const selectEmoji = () => {
  ElMessage.info('表情功能开发中...')
}

// 选择文件
const selectFile = () => {
  ElMessage.info('文件发送功能开发中...')
}

// 监听消息变化，自动滚动到底部
watch(
  () => chatStore.messages.length,
  () => {
    scrollToBottom()
  }
)

onMounted(async () => {
  // 获取房间列表
  await chatStore.fetchRooms()
  
  // 连接WebSocket
  chatStore.connectWebSocket()
  
  // 如果有房间，选择第一个
  if (chatStore.rooms.length > 0) {
    selectRoom(chatStore.rooms[0])
  }
})
</script>

<style lang="scss" scoped>
.chat-container {
  display: flex;
  height: 100vh;
  background-color: #f5f5f5;
}

.sidebar {
  width: 300px;
  background: #fff;
  border-right: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-details {
  display: flex;
  flex-direction: column;
  
  .nickname {
    font-weight: 500;
    color: #333;
  }
  
  .username {
    font-size: 12px;
    color: #999;
  }
}

.more-icon {
  cursor: pointer;
  color: #999;
  font-size: 20px;
  
  &:hover {
    color: #333;
  }
}

.sidebar-actions {
  padding: 20px;
  border-bottom: 1px solid #e0e0e0;
}

.room-list {
  flex: 1;
  overflow-y: auto;
}

.room-list-header {
  padding: 15px 20px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  h3 {
    margin: 0;
    font-size: 16px;
    color: #333;
  }
}

.room-items {
  padding: 10px 0;
}

.room-item {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  cursor: pointer;
  transition: background-color 0.2s;
  
  &:hover {
    background-color: #f5f5f5;
  }
  
  &.active {
    background-color: #e3f2fd;
    border-left: 3px solid #2196f3;
  }
}

.room-avatar {
  margin-right: 12px;
}

.room-info {
  flex: 1;
  min-width: 0;
}

.room-name {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.room-last-message {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.room-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.room-time {
  font-size: 12px;
  color: #999;
}

.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #fff;
}

.chat-header {
  padding: 15px 20px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chat-title {
  display: flex;
  align-items: center;
  gap: 10px;
  
  h3 {
    margin: 0;
    font-size: 18px;
    color: #333;
  }
}

.room-type {
  padding: 2px 8px;
  background: #f0f0f0;
  border-radius: 12px;
  font-size: 12px;
  color: #666;
}

.chat-actions {
  display: flex;
  gap: 8px;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: #fafafa;
}

.messages-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message-item {
  display: flex;
  gap: 12px;
  
  &.message-mine {
    flex-direction: row-reverse;
    
    .message-content {
      align-items: flex-end;
    }
    
    .message-body {
      background: #2196f3;
      color: white;
    }
  }
}

.message-avatar {
  flex-shrink: 0;
}

.message-content {
  display: flex;
  flex-direction: column;
  max-width: 60%;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.message-sender {
  font-size: 12px;
  color: #666;
}

.message-time {
  font-size: 12px;
  color: #999;
}

.message-body {
  background: #fff;
  padding: 10px 15px;
  border-radius: 8px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.message-text {
  word-wrap: break-word;
  line-height: 1.4;
}

.chat-input {
  border-top: 1px solid #e0e0e0;
  padding: 15px;
  background: #fff;
}

.input-toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 10px;
}

.input-area {
  display: flex;
  gap: 10px;
  align-items: flex-end;
  
  .el-textarea {
    flex: 1;
  }
}

.empty-state {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  background: #fff;
}
</style>
