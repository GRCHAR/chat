<template>
  <el-dialog
    v-model="visible"
    title="搜索用户"
    width="600px"
    @close="handleClose"
  >
    <div class="search-container">
      <el-input
        v-model="searchQuery"
        placeholder="输入用户名、昵称或邮箱搜索用户"
        clearable
        size="large"
        @input="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      
      <div class="search-results" v-loading="loading">
        <div v-if="!searchQuery" class="empty-state">
          <el-icon class="empty-icon"><Search /></el-icon>
          <p>请输入关键词搜索用户</p>
        </div>
        
        <div v-else-if="searchResults.length === 0 && !loading" class="empty-state">
          <el-icon class="empty-icon"><Warning /></el-icon>
          <p>未找到匹配的用户</p>
        </div>
        
        <div v-else class="user-list">
          <div
            v-for="user in searchResults"
            :key="user.id"
            class="user-item"
            @click="handleUserClick(user)"
          >
            <el-avatar :size="40" :src="user.avatar">
              {{ user.nickname.charAt(0) }}
            </el-avatar>
            <div class="user-info">
              <div class="user-name">{{ user.nickname }}</div>
              <div class="user-username">@{{ user.username }}</div>
            </div>
            <div class="user-actions">
              <el-button type="primary" size="small" @click.stop="startPrivateChat(user)">
                私聊
              </el-button>
              <el-button size="small" @click.stop="addToGroup(user)">
                添加到群聊
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 添加成员对话框 -->
    <AddMemberDialog
      v-model="showAddMemberDialog"
      :pre-selected-user="selectedUserForGroup"
      @member-added="handleMemberAdded"
    />
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Warning } from '@element-plus/icons-vue'
import { userApi } from '@/api/user'
import { chatApi } from '@/api/chat'
import { useChatStore } from '@/stores/chat'
import AddMemberDialog from './AddMemberDialog.vue'
import type { User } from '@/types/user'

interface Props {
  modelValue: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'userSelected', user: User): void
  (e: 'chatCreated', roomId: number): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const chatStore = useChatStore()

const visible = ref(props.modelValue)
const searchQuery = ref('')
const searchResults = ref<User[]>([])
const loading = ref(false)

let searchTimer: ReturnType<typeof setTimeout> | null = null

// 监听modelValue变化
watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
})

// 监听visible变化
watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

// 搜索用户
const handleSearch = () => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  
  if (!searchQuery.value.trim()) {
    searchResults.value = []
    return
  }
  
  searchTimer = setTimeout(async () => {
    loading.value = true
    try {
      const response = await userApi.searchUsers(searchQuery.value.trim())
      searchResults.value = response.data.users || []
    } catch (error) {
      console.error('搜索用户失败:', error)
      ElMessage.error('搜索用户失败')
      searchResults.value = []
    } finally {
      loading.value = false
    }
  }, 300)
}

// 点击用户
const handleUserClick = (user: User) => {
  emit('userSelected', user)
}

// 发起私聊
const startPrivateChat = async (user: User) => {
  try {
    loading.value = true
    
    // 创建私聊房间
    const response = await chatApi.createRoom({
      name: `与${user.nickname}的聊天`,
      type: 'single',
      member_ids: [user.id]
    })
    
    if (response.data.room) {
      ElMessage.success('私聊创建成功')
      
      // 刷新房间列表
      await chatStore.fetchRooms()
      
      // 选择新创建的房间
      chatStore.setCurrentRoom(response.data.room)
      
      emit('chatCreated', response.data.room.id)
      handleClose()
    }
  } catch (error: any) {
    console.error('创建私聊失败:', error)
    ElMessage.error(error.response?.data?.error || '创建私聊失败')
  } finally {
    loading.value = false
  }
}

// 添加到群聊
const showAddMemberDialog = ref(false)
const selectedUserForGroup = ref<User>()

const addToGroup = (user: User) => {
  selectedUserForGroup.value = user
  showAddMemberDialog.value = true
}

const handleMemberAdded = () => {
  ElMessage.success('成员添加成功')
  showAddMemberDialog.value = false
}

// 关闭对话框
const handleClose = () => {
  searchQuery.value = ''
  searchResults.value = []
  visible.value = false
}
</script>

<style lang="scss" scoped>
.search-container {
  padding: 10px 0;
}

.search-results {
  margin-top: 20px;
  min-height: 300px;
  max-height: 400px;
  overflow-y: auto;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #999;
  
  .empty-icon {
    font-size: 64px;
    margin-bottom: 16px;
  }
  
  p {
    margin: 0;
    font-size: 14px;
  }
}

.user-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
  
  &:hover {
    background-color: #f5f5f5;
  }
}

.user-info {
  flex: 1;
  margin-left: 12px;
  min-width: 0;
}

.user-name {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.user-username {
  font-size: 12px;
  color: #999;
}

.user-actions {
  display: flex;
  gap: 8px;
}
</style>
