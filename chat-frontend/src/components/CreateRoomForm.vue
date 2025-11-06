<template>
  <el-form
    ref="formRef"
    :model="formData"
    :rules="formRules"
    label-width="100px"
  >
    <el-form-item label="房间类型" prop="type">
      <el-radio-group v-model="formData.type">
        <el-radio value="single">私聊</el-radio>
        <el-radio value="group">群聊</el-radio>
      </el-radio-group>
    </el-form-item>
    
    <el-form-item label="房间名称" prop="name" v-if="formData.type === 'group'">
      <el-input
        v-model="formData.name"
        placeholder="请输入房间名称"
        clearable
      />
    </el-form-item>
    
    <el-form-item label="房间描述" prop="description" v-if="formData.type === 'group'">
      <el-input
        v-model="formData.description"
        type="textarea"
        placeholder="请输入房间描述（可选）"
        :rows="3"
        clearable
      />
    </el-form-item>
    
    <el-form-item label="选择成员" prop="member_ids">
      <div class="member-selector">
        <el-input
          v-model="searchQuery"
          placeholder="搜索用户"
          clearable
          style="margin-bottom: 10px"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        
        <div class="member-list">
          <el-checkbox-group v-model="formData.member_ids">
            <el-checkbox
              v-for="user in filteredUsers"
              :key="user.id"
              :value="user.id"
              class="member-item"
            >
              <div class="member-info">
                <el-avatar :size="24" :src="user.avatar">
                  {{ user.nickname.charAt(0) }}
                </el-avatar>
                <span class="member-name">{{ user.nickname }}</span>
                <span class="member-username">@{{ user.username }}</span>
              </div>
            </el-checkbox>
          </el-checkbox-group>
        </div>
        
        <div class="selected-count">
          已选择 {{ formData.member_ids.length }} 个成员
        </div>
      </div>
    </el-form-item>
  </el-form>
  
  <div class="form-actions">
    <el-button @click="handleCancel">取消</el-button>
    <el-button type="primary" @click="handleSubmit" :loading="loading">
      创建
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { chatApi } from '@/api/chat'
import type { User } from '@/types/user'

const emit = defineEmits<{
  success: []
  cancel: []
}>()

const authStore = useAuthStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const searchQuery = ref('')
const allUsers = ref<User[]>([])

const formData = reactive({
  name: '',
  description: '',
  type: 'single' as 'single' | 'group',
  member_ids: [] as number[]
})

const formRules: FormRules = {
  type: [
    { required: true, message: '请选择房间类型', trigger: 'change' }
  ],
  name: [
    { required: true, message: '请输入房间名称', trigger: 'blur' },
    { min: 2, max: 100, message: '房间名称长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 500, message: '房间描述不能超过 500 个字符', trigger: 'blur' }
  ],
  member_ids: [
    { 
      required: true, 
      message: '请至少选择一个成员', 
      trigger: 'change',
      validator: (rule, value, callback) => {
        if (value.length === 0) {
          callback(new Error('请至少选择一个成员'))
        } else {
          callback()
        }
      }
    }
  ]
}

const filteredUsers = computed(() => {
  const query = searchQuery.value.toLowerCase()
  if (!query) {
    return allUsers.value.filter(user => user.id !== authStore.user?.id)
  }
  
  return allUsers.value.filter(user => 
    user.id !== authStore.user?.id && (
      user.nickname.toLowerCase().includes(query) ||
      user.username.toLowerCase().includes(query)
    )
  )
})

const fetchUsers = async () => {
  try {
    // 这里应该调用获取用户列表的API，现在模拟数据
    // const response = await userApi.getUsers()
    // allUsers.value = response.data.users
    
    // 模拟数据
    allUsers.value = [
      { id: 2, username: 'user2', nickname: '用户2', email: 'user2@example.com', status: 'active', created_at: '2023-01-01', updated_at: '2023-01-01' },
      { id: 3, username: 'user3', nickname: '用户3', email: 'user3@example.com', status: 'active', created_at: '2023-01-01', updated_at: '2023-01-01' },
      { id: 4, username: 'user4', nickname: '用户4', email: 'user4@example.com', status: 'offline', created_at: '2023-01-01', updated_at: '2023-01-01' }
    ]
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const requestData = {
          ...formData,
          name: formData.type === 'single' 
            ? `与${getSelectedUserNames()}的聊天`
            : formData.name
        }
        
        const response = await chatApi.createRoom(requestData)
        
        if (response.data.room) {
          ElMessage.success('房间创建成功')
          emit('success')
        }
      } catch (error: any) {
        ElMessage.error(error.response?.data?.error || '创建房间失败')
      } finally {
        loading.value = false
      }
    }
  })
}

const handleCancel = () => {
  emit('cancel')
}

const getSelectedUserNames = () => {
  const selectedUsers = allUsers.value.filter(user => 
    formData.member_ids.includes(user.id)
  )
  return selectedUsers.map(user => user.nickname).join('、')
}

// 监听房间类型变化
const handleTypeChange = () => {
  formData.member_ids = []
  if (formData.type === 'single') {
    formData.name = ''
    formData.description = ''
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style lang="scss" scoped>
.member-selector {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 10px;
  max-height: 300px;
  overflow-y: auto;
}

.member-list {
  margin-bottom: 10px;
}

.member-item {
  display: block;
  margin-bottom: 8px;
  width: 100%;
  
  &:last-child {
    margin-bottom: 0;
  }
}

.member-info {
  display: flex;
  align-items: center;
  gap: 8px;
  
  .member-name {
    font-weight: 500;
    color: #333;
  }
  
  .member-username {
    font-size: 12px;
    color: #999;
  }
}

.selected-count {
  text-align: center;
  color: #666;
  font-size: 14px;
  padding-top: 10px;
  border-top: 1px solid #e0e0e0;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}
</style>
