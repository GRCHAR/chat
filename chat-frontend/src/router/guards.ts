import { ROUTE_CONFIG } from '@/config'
import { useAuthStore } from '@/stores/auth'
import { showError } from '@/utils/error'
import type { Router } from 'vue-router'

/**
 * 设置路由守卫
 */
export function setupRouterGuards(router: Router): void {
  // 全局前置守卫
  router.beforeEach(async (to, _from, next) => {
    const authStore = useAuthStore()
    
    // 检查是否需要认证
    const requiresAuth = ROUTE_CONFIG.authRequired.some(path => 
      to.path.startsWith(path)
    )
    
    // 检查是否是公开路由
    const isPublicRoute = ROUTE_CONFIG.publicRoutes.includes(to.path)
    
    try {
      // 如果路由需要认证
      if (requiresAuth) {
        // 检查用户是否已登录
        if (!authStore.isAuthenticated) {
          // 尝试自动登录（检查本地存储的token）
          const token = authStore.getStoredToken()
          if (token) {
            try {
              await authStore.verifyToken()
              next()
              return
            } catch (error) {
              // Token验证失败，清除token并跳转到登录页
              authStore.clearAuth()
              next('/login')
              return
            }
          }
          
          // 未登录，跳转到登录页
          next('/login')
          return
        }
        
        // 已登录，检查用户权限
        if (to.meta?.requiredRole) {
          const requiredRole = to.meta.requiredRole as string
          if (!authStore.hasRole(requiredRole)) {
            showError('您没有访问该页面的权限')
            next('/')
            return
          }
        }
      }
      
      // 如果用户已登录，尝试访问登录/注册页面，跳转到主页
      if (authStore.isAuthenticated && isPublicRoute) {
        next(ROUTE_CONFIG.defaultRedirect)
        return
      }
      
      // 正常访问
      next()
      
    } catch (error) {
      console.error('路由守卫错误:', error)
      showError('页面加载失败，请重试')
      next('/')
    }
  })
  
  // 全局后置守卫
  router.afterEach((to, _from) => {
    // 设置页面标题
    const title = to.meta?.title as string
    if (title) {
      document.title = `${title} - Chat Application`
    } else {
      document.title = 'Chat Application'
    }
    
    // 滚动到页面顶部
    window.scrollTo(0, 0)
  })
  
  // 错误处理
  router.onError((error) => {
    console.error('路由错误:', error)
    showError('页面加载失败，请刷新重试')
  })
}

/**
 * 路由元信息类型定义
 */
export interface RouteMeta extends Record<PropertyKey, unknown> {
  // 页面标题
  title?: string
  
  // 是否需要认证
  requiresAuth?: boolean
  
  // 需要的角色
  requiredRole?: string
  
  // 是否缓存页面
  keepAlive?: boolean
  
  // 页面图标
  icon?: string
  
  // 是否在菜单中隐藏
  hideInMenu?: boolean
  
  // 排序
  order?: number
}
