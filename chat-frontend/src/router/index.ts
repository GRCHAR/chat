import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import type { RouteMeta } from './guards'
import { setupRouterGuards } from './guards'

// 路由配置
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    redirect: '/chat',
    meta: {
      title: '首页',
      requiresAuth: true
    } as RouteMeta
  },
  
  // 认证相关
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
    meta: {
      title: '登录',
      hideInMenu: true
    } as RouteMeta
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/RegisterView.vue'),
    meta: {
      title: '注册',
      hideInMenu: true
    } as RouteMeta
  },
  
  // 聊天相关
  {
    path: '/chat',
    name: 'Chat',
    component: () => import('@/views/ChatView.vue'),
    meta: {
      title: '聊天室',
      requiresAuth: true,
      icon: 'ChatDotRound',
      order: 1
    } as RouteMeta
  },
  
  // 用户相关
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/ProfileView.vue'),
    meta: {
      title: '个人资料',
      requiresAuth: true,
      icon: 'User',
      order: 2
    } as RouteMeta
  },
  
  // 404页面
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFoundView.vue'),
    meta: {
      title: '页面未找到',
      hideInMenu: true
    } as RouteMeta
  }
]

// 创建路由器
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0 }
  }
})

// 设置路由守卫
setupRouterGuards(router)

export default router
