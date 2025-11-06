<template>
  <div class="profile-container">
    <div class="profile-header">
      <el-button :icon="ArrowLeft" @click="$router.back()">返回</el-button>
      <h2>个人信息</h2>
    </div>
    
    <div class="profile-content">
      <el-card class="profile-card">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
            <el-button
              v-if="!isEditing"
              type="primary"
              @click="startEdit"
            >
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <div v-else class="edit-actions">
              <el-button @click="cancelEdit">取消</el-button>
              <el-button type="primary" @click="saveProfile" :loading="saving">
                保存
              </el-button>
            </div>
          </div>
        </template>
        
        <el-form
          ref="profileFormRef"
          :model="profileForm"
          :rules="profileRules"
          label-width="80px"
          :disabled="!isEditing"
        >
          <div class="avatar-section">
            <el-avatar :size="80" :src="profileForm.avatar" class="avatar">
              {{ profileForm.nickname?.charAt(0) || 'U' }}
            </el-avatar>
            <el-button
              v-if="isEditing"
              type="primary"
              size="small"
              class="avatar-upload-btn"
              @click="uploadAvatar"
            >
              <el-icon><Camera /></el-icon>
              更换头像
            </el-button>
            <input
              ref="avatarInputRef"
              type="file"
              accept="image/*"
              style="display: none"
              @change="handleAvatarChange"
            />
          </div>
          
          <el-form-item label="用户名">
            <el-input v-model="profileForm.username" disabled />
          </el-form-item>
          
          <el-form-item label="昵称" prop="nickname">
            <el-input v-model="profileForm.nickname" placeholder="请输入昵称" />
          </el-form-item>
          
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="profileForm.email" placeholder="请输入邮箱地址" />
          </el-form-item>
          
          <el-form-item label="状态">
            <el-tag :type="getStatusType(profileForm.status)">
              {{ getStatusText(profileForm.status) }}
            </el-tag>
          </el-form-item>
          
          <el-form-item label="注册时间">
            <el-input :value="formatDate(profileForm.created_at)" disabled />
          </el-form-item>
        </el-form>
      </el-card>
      
      <el-card class="security-card">
        <template #header>
          <span>账户安全</span>
        </template>
        
        <div class="security-items">
          <div class="security-item">
            <div class="security-info">
              <h4>修改密码</h4>
              <p>定期更换密码可以提高账户安全性</p>
            </div>
            <el-button @click="showPasswordDialog = true">修改</el-button>
          </div>
        </div>
      </el-card>
    </div>
    
    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="showPasswordDialog"
      title="修改密码"
      width="400px"
    >
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="100px"
      >
        <el-form-item label="当前密码" prop="currentPassword">
          <el-input
            v-model="passwordForm.currentPassword"
            type="password"
            placeholder="请输入当前密码"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showPasswordDialog = false">取消</el-button>
        <el-button type="primary" @click="changePassword" :loading="changingPassword">
          确认修改
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { ArrowLeft, Edit, Camera } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { formatDate } from '@/utils/format'

const authStore = useAuthStore()

const profileFormRef = ref<FormInstance>()
const avatarInputRef = ref<HTMLInputElement>()
const passwordFormRef = ref<FormInstance>()

const isEditing = ref(false)
const saving = ref(false)
const showPasswordDialog = ref(false)
const changingPassword = ref(false)

const originalProfile = ref({
  nickname: '',
  email: '',
  avatar: ''
})

const profileForm = reactive({
  username: '',
  nickname: '',
  email: '',
  avatar: '',
  status: 'active',
  created_at: ''
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule: any, value: string, callback: any) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入密码不一致!'))
  } else {
    callback()
  }
}

const profileRules: FormRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 100, message: '昵称长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const passwordRules: FormRules = {
  currentPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const loadUserProfile = () => {
  if (authStore.user) {
    const user = authStore.user
    profileForm.username = user.username
    profileForm.nickname = user.nickname
    profileForm.email = user.email
    profileForm.avatar = user.avatar || ''
    profileForm.status = user.status
    profileForm.created_at = user.created_at
    
    // 保存原始数据
    originalProfile.value = {
      nickname: user.nickname,
      email: user.email,
      avatar: user.avatar || ''
    }
  }
}

const startEdit = () => {
  isEditing.value = true
}

const cancelEdit = () => {
  isEditing.value = false
  // 恢复原始数据
  profileForm.nickname = originalProfile.value.nickname
  profileForm.email = originalProfile.value.email
  profileForm.avatar = originalProfile.value.avatar
}

const saveProfile = async () => {
  if (!profileFormRef.value) return
  
  await profileFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        const result = await authStore.updateProfile({
          nickname: profileForm.nickname,
          email: profileForm.email,
          avatar: profileForm.avatar
        })
        
        if (result.success) {
          ElMessage.success('个人信息更新成功')
          isEditing.value = false
          // 更新原始数据
          originalProfile.value = {
            nickname: profileForm.nickname,
            email: profileForm.email,
            avatar: profileForm.avatar
          }
        } else {
          ElMessage.error(result.error || '更新失败')
        }
      } catch (error) {
        ElMessage.error('更新失败，请重试')
      } finally {
        saving.value = false
      }
    }
  })
}

const uploadAvatar = () => {
  avatarInputRef.value?.click()
}

const handleAvatarChange = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    // 这里应该上传文件到服务器，现在只是模拟
    const reader = new FileReader()
    reader.onload = (e) => {
      profileForm.avatar = e.target?.result as string
    }
    reader.readAsDataURL(file)
  }
}

const changePassword = async () => {
  if (!passwordFormRef.value) return
  
  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      changingPassword.value = true
      try {
        // 这里应该调用修改密码的API
        ElMessage.success('密码修改成功')
        showPasswordDialog.value = false
        passwordFormRef.value?.resetFields()
      } catch (error) {
        ElMessage.error('密码修改失败')
      } finally {
        changingPassword.value = false
      }
    }
  })
}

const getStatusType = (status: string) => {
  switch (status) {
    case 'active': return 'success'
    case 'inactive': return 'warning'
    case 'offline': return 'info'
    default: return 'info'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'active': return '活跃'
    case 'inactive': return '未激活'
    case 'offline': return '离线'
    default: return '未知'
  }
}

onMounted(() => {
  loadUserProfile()
})
</script>

<style lang="scss" scoped>
.profile-container {
  height: 100vh;
  background-color: #f5f5f5;
  overflow-y: auto;
}

.profile-header {
  background: #fff;
  padding: 20px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  gap: 20px;
  
  h2 {
    margin: 0;
    color: #333;
  }
}

.profile-content {
  max-width: 800px;
  margin: 20px auto;
  padding: 0 20px;
}

.profile-card,
.security-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.edit-actions {
  display: flex;
  gap: 10px;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 30px;
  position: relative;
}

.avatar {
  margin-bottom: 10px;
}

.avatar-upload-btn {
  position: absolute;
  bottom: 0;
  right: calc(50% - 40px);
}

.security-items {
  padding: 20px 0;
}

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 0;
  border-bottom: 1px solid #f0f0f0;
  
  &:last-child {
    border-bottom: none;
  }
}

.security-info {
  h4 {
    margin: 0 0 8px 0;
    color: #333;
  }
  
  p {
    margin: 0;
    color: #666;
    font-size: 14px;
  }
}
</style>
