// 应用配置
export const APP_CONFIG = {
  // 应用名称
  name: 'Chat Application',
  
  // 应用版本
  version: '1.0.0',
  
  // API基础URL
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  
  // WebSocket URL
  wsBaseUrl: import.meta.env.VITE_WS_BASE_URL || 'ws://localhost:8080/ws',
  
  // 本地存储键前缀
  storagePrefix: 'chat_app_',
  
  // 分页配置
  pagination: {
    defaultPageSize: 20,
    pageSizeOptions: [10, 20, 50, 100]
  },
  
  // 消息配置
  message: {
    maxLength: 1000,
    showTimestamp: true,
    autoScroll: true
  },
  
  // 文件上传配置
  upload: {
    maxSize: 10 * 1024 * 1024, // 10MB
    allowedTypes: ['image/jpeg', 'image/png', 'image/gif', 'application/pdf']
  }
}

// 路由配置
export const ROUTE_CONFIG = {
  // 需要认证的路由
  authRequired: ['/chat', '/profile'],
  
  // 游客可访问的路由
  publicRoutes: ['/login', '/register'],
  
  // 默认重定向路径
  defaultRedirect: '/chat'
}

// WebSocket配置
export const WS_CONFIG = {
  // 重连配置
  reconnect: {
    enabled: true,
    maxAttempts: 5,
    delay: 3000 // 3秒
  },
  
  // 心跳配置
  heartbeat: {
    enabled: true,
    interval: 30000 // 30秒
  }
}

// 主题配置
export const THEME_CONFIG = {
  // 默认主题
  defaultTheme: 'light',
  
  // 可用主题
  availableThemes: ['light', 'dark', 'auto'],
  
  // 主题色
  primaryColor: '#409EFF'
}
