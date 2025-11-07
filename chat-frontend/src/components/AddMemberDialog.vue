<template>
  <el-dialog
    v-model="visible"
    title="添加成员到群聊"
    width="600px"
    @close="handleClose"
  >
    <div class="add-member-container">
      <!-- 选择群聊 -->
      <el-form-item label="选择群聊">
        <el-select
          v-model="selectedRoomId"
          placeholder="请选择要添加成员的群聊"
          style="width: 100%"
          filterable
        >
          <el-option
            v-for="room in groupRooms"
            :key="room.id"
            :label="room.name"
            :value="room.id"
          >
            <div class="room-option">
              <el-avatar :size="24">{{ room.name.charAt(0) }}</el-avatar>
              <span class="room-name">{{ room.name }}</span>
            </div>
          </el-option>
        </el-select>
      </el-form-item>
      
      <!-- 搜索用户 -->
      <el-form-item label="搜索用户" style="margin-top: 20px">
        <el-input
          v-model="searchQuery"
          placeholder="输入用户名、昵称或邮箱搜索用户"
          clearable
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </el-form-item>
      
      <!-- 搜索结果 -->
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
          >
            <el-avatar :size="32" :src="user.avatar">
              {{ user.nickname.charAt(0) }}
            </el-avatar>
            <div class="user-info">
              <div class="user-name">{{ user.nickname }}</div>
              <div class="user-username">@{{ user.username }}</div>
            </div>
            <el-button
              type="primary"
              size="small"
              :disabled="!selectedRoomId"
              @click="addMember(user)"
            >
              添加
            </el-button>
          </div>
        </div>
      </div>
    </div>
    
    <template #footer>
      <el-button @click="handleClose">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Warning } from '@element-plus/icons-vue'
import { userApi } from '@/api/user'
import { chatApi } from '@/api/chat'
import { useChatStore } from '@/stores/chat'
import type { User } from '@/types/user'
import type { ChatRoom } from '@/types/chat'

interface Props {
  modelValue: boolean
  preSelectedUser?: User
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'memberAdded', roomId: number, user: User): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const chatStore = useChatStore()

const visible = ref(props.modelValue)
const selectedRoomId = ref<number>()
const searchQuery = ref('')
const searchResults = ref<User[]>([])
const loading = ref(false)

let searchTimer: ReturnType<typeof setTimeout> | null = null

// 获取群聊列表
const groupRooms = computed(() => {
  return chatStore.rooms.filter(room => room.type === 'group')
})

// 监听modelValue变化
watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal && props.preSelectedUser) {
    searchQuery.value = props.preSelectedUser.username
    searchResults.value = [props.preSelectedUser]
  }
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

// 添加成员
const addMember = async (user: User) => {
  if (!selectedRoomId.value) {
    ElMessage.warning('请先选择群聊')
    return
  }
  
  try {
    loading.value = true
    
    // 调用添加成员API
    await chatApi.addMember(selectedRoomId.value, user.id)
    
    ElMessage.success(`已将 ${user.nickname} 添加到群聊`)
    emit('memberAdded', selectedRoomId.value, user)
    
    // 刷新房间信息
    await chatStore.fetchRooms()
  } catch (error: any) {
    console.error('添加成员失败:', error)
    ElMessage.error(error.response?.data?.error || '添加成员失败')
  } finally {
    loading.value = false
  }
}

// 关闭对话框
const handleClose = () => {
  searchQuery.value = ''
  searchResults.value = []
  selectedRoomId.value = undefined
  visible.value = false
}
</script>

<style lang="scss" scoped>
.add-member-container {
  padding: 10px 0;
}

.room-option {
  display: flex;
  align-items: center;
  gap: 8px;
  
  .room-name {
    font-size: 14px;
  }
}

.search-results {
  margin-top: 10px;
  min-height: 200px;
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  padding: 10px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: #999;
  
  .empty-icon {
    font-size: 48px;
    margin-bottom: 12px;
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
  padding: 8px;
  border-radius: 6px;
  transition: background-color 0.2s;
  
  &:hover {
    background-color: #f5f5f5;
  }
}

.user-info {
  flex: 1;
  margin-left: 10px;
  min-width: 0;
}

.user-name {
  font-weight: 500;
  color: #333;
  margin-bottom: 2px;
}

.user-username {
  font-size: 12px;
  color: #999;
}
</style>
