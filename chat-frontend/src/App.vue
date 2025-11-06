<template>
  <el-config-provider :locale="locale" :size="componentSize">
    <div id="app" :class="{ 'dark-theme': isDarkTheme }">
      <!-- 加载状态 -->
      <div v-if="appLoading" class="app-loading">
        <el-loading :loading="true" text="加载中..." />
      </div>
      
      <!-- 主应用 -->
      <div v-else class="app-container">
        <!-- 侧边栏导航 -->
        <el-aside v-if="showSidebar" class="app-sidebar" :width="sidebarWidth">
          <div class="sidebar-header">
            <div class="logo">
              <el-icon size="32"><ChatDotRound /></el-icon>
              <span class="logo-text">Chat App</span>
            </div>
          </div>
          
          <el-menu
            :default-active="activeMenu"
            :collapse="isCollapse"
            :unique-opened="true"
            router
            class="sidebar-menu"
          >
            <el-menu-item
              v-for="route in menuRoutes"
              :key="route.path"
              :index="route.path"
              :route="{ path: route.path }"
            >
              <el-icon v-if="route.meta?.icon">
                <component :is="route.meta.icon" />
              </el-icon>
              <span>{{ route.meta?.title }}</span>
            </el-menu-item>
          </el-menu>
          
          <div class="sidebar-footer">
            <el-button
              text
              class="collapse-btn"
              @click="toggleSidebar"
            >
              <el-icon>
                <component :is="isCollapse ? 'Expand' : 'Fold'" />
              </el-icon>
            </el-button>
          </div>
        </el-aside>
        
        <!-- 主内容区 -->
        <el-container class="main-container">
          <!-- 顶部导航栏 -->
          <el-header v-if="showHeader" class="app-header">
            <div class="header-left">
              <el-breadcrumb separator="/">
                <el-breadcrumb-item
                  v-for="item in breadcrumbs"
                  :key="item.path"
                  :to="item.path"
                >
                  {{ item.title }}
                </el-breadcrumb-item>
              </el-breadcrumb>
            </div>
            
            <div class="header-right">
              <!-- 主题切换 -->
              <el-button
                text
                class="theme-btn"
                @click="toggleTheme"
              >
                <el-icon>
                  <component :is="isDarkTheme ? 'Sunny' : 'Moon'" />
                </el-icon>
              </el-button>
              
              <!-- 用户信息 -->
              <el-dropdown v-if="isAuthenticated" @command="handleUserAction">
                <div class="user-info">
                  <el-avatar :size="32" :src="userAvatar">
                    {{ userName?.charAt(0)?.toUpperCase() }}
                  </el-avatar>
                  <span class="user-name">{{ userName }}</span>
                  <el-icon><ArrowDown /></el-icon>
                </div>
                
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="profile">
                      <el-icon><User /></el-icon>
                      个人资料
                    </el-dropdown-item>
                    <el-dropdown-item divided command="logout">
                      <el-icon><SwitchButton /></el-icon>
                      退出登录
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </el-header>
          
          <!-- 页面内容 -->
          <el-main class="app-main">
            <router-view v-slot="{ Component }">
              <transition name="fade" mode="out-in">
                <keep-alive>
                  <component :is="Component" />
                </keep-alive>
              </transition>
            </router-view>
          </el-main>
        </el-container>
      </div>
    </div>
  </el-config-provider>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElConfigProvider } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { useChatStore } from '@/stores/chat'
import { storageUtils } from '@/utils/storage'
import { APP_CONFIG, THEME_CONFIG } from '@/config'
import { setupGlobalErrorHandler } from '@/utils/error'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import {
  ChatDotRound,
  User,
  ArrowDown,
  ArrowLeft,
  Expand,
  Fold,
  Sunny,
  Moon,
  SwitchButton
} from '@element-plus/icons-vue'

// 状态管理
const authStore = useAuthStore()
const chatStore = useChatStore()
const router = useRouter()
const route = useRoute()

// 响应式状态
const appLoading = ref(true)
const isCollapse = ref(false)
const locale = ref(zhCn)
const componentSize = ref<'default' | 'large' | 'small'>('default')

// 计算属性
const isAuthenticated = computed(() => authStore.isAuthenticated)
const userName = computed(() => authStore.userInfo?.username)
const userAvatar = computed(() => authStore.userInfo?.avatar)

const showSidebar = computed(() => {
  return isAuthenticated.value && !route.meta?.hideInMenu
})

const showHeader = computed(() => {
  return isAuthenticated.value
})

const sidebarWidth = computed(() => {
  return isCollapse.value ? '64px' : '240px'
})

const isDarkTheme = computed(() => {
  const theme = storageUtils.getTheme()
  if (theme === 'auto') {
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  return theme === 'dark'
})

const activeMenu = computed(() => {
  return route.path
})

const menuRoutes = computed(() => {
  return router.getRoutes()
    .filter(route => 
      route.meta?.title && 
      !route.meta?.hideInMenu && 
      route.path !== '/'
    )
    .sort((a, b) => (a.meta?.order || 0) - (b.meta?.order || 0))
})

const breadcrumbs = computed(() => {
  const matched = route.matched.filter(item => item.meta?.title)
  return matched.map(item => ({
    path: item.path,
    title: item.meta?.title
  }))
})

// 方法
const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
}

const toggleTheme = () => {
  const currentTheme = storageUtils.getTheme()
  const themes = THEME_CONFIG.availableThemes
  const currentIndex = themes.indexOf(currentTheme)
  const nextIndex = (currentIndex + 1) % themes.length
  const newTheme = themes[nextIndex]
  
  storageUtils.setTheme(newTheme)
  applyTheme(newTheme)
}

const applyTheme = (theme: string) => {
  const root = document.documentElement
  if (theme === 'dark' || (theme === 'auto' && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    root.classList.add('dark-theme')
  } else {
    root.classList.remove('dark-theme')
  }
}

const handleUserAction = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      handleLogout()
      break
  }
}

const handleLogout = async () => {
  try {
    await authStore.logout()
    router.push('/login')
  } catch (error) {
    console.error('退出登录失败:', error)
  }
}

// 初始化
const initializeApp = async () => {
  try {
    // 设置全局错误处理
    setupGlobalErrorHandler()
    
    // 应用主题
    const theme = storageUtils.getTheme()
    applyTheme(theme)
    
    // 检查认证状态
    const token = storageUtils.getToken()
    if (token) {
      await authStore.verifyToken()
    }
    
    // 初始化WebSocket连接（如果已登录）
    if (authStore.isAuthenticated) {
      chatStore.connectWebSocket()
    }
    
  } catch (error) {
    console.error('应用初始化失败:', error)
  } finally {
    appLoading.value = false
  }
}

// 生命周期
onMounted(() => {
  initializeApp()
})
</script>

<style lang="scss">
#app {
  height: 100vh;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
}

.app-container {
  display: flex;
  height: 100vh;
}

.app-loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.app-sidebar {
  background-color: #fff;
  border-right: 1px solid #e4e7ed;
  transition: width 0.3s;
  display: flex;
  flex-direction: column;
  
  .sidebar-header {
    padding: 20px;
    border-bottom: 1px solid #e4e7ed;
    
    .logo {
      display: flex;
      align-items: center;
      gap: 12px;
      
      .logo-text {
        font-size: 18px;
        font-weight: 600;
        color: #303133;
      }
    }
  }
  
  .sidebar-menu {
    flex: 1;
    border-right: none;
  }
  
  .sidebar-footer {
    padding: 16px;
    border-top: 1px solid #e4e7ed;
    
    .collapse-btn {
      width: 100%;
      justify-content: center;
    }
  }
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.app-header {
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 60px;
  
  .header-left {
    flex: 1;
  }
  
  .header-right {
    display: flex;
    align-items: center;
    gap: 16px;
    
    .theme-btn {
      padding: 8px;
    }
    
    .user-info {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
      padding: 8px;
      border-radius: 6px;
      transition: background-color 0.3s;
      
      &:hover {
        background-color: #f5f7fa;
      }
      
      .user-name {
        font-size: 14px;
        color: #606266;
      }
    }
  }
}

.app-main {
  flex: 1;
  padding: 20px;
  background-color: #f5f7fa;
  overflow-y: auto;
}

// 过渡动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

// 暗色主题
.dark-theme {
  background-color: #141414;
  color: #fff;
  
  .app-sidebar {
    background-color: #1f1f1f;
    border-right-color: #303030;
    
    .sidebar-header {
      border-bottom-color: #303030;
      
      .logo-text {
        color: #fff;
      }
    }
    
    .sidebar-footer {
      border-top-color: #303030;
    }
  }
  
  .app-header {
    background-color: #1f1f1f;
    border-bottom-color: #303030;
    
    .user-info:hover {
      background-color: #303030;
    }
  }
  
  .app-main {
    background-color: #141414;
  }
}
</style>
