import { ElMessage } from 'element-plus'

// 错误类型定义
export interface AppError {
  code: number | string
  message: string
  details?: any
}

// HTTP状态码映射
const HTTP_STATUS_MESSAGES: Record<number, string> = {
  400: '请求参数错误',
  401: '未授权，请重新登录',
  403: '没有权限访问',
  404: '请求的资源不存在',
  408: '请求超时',
  409: '资源冲突',
  422: '请求参数验证失败',
  500: '服务器内部错误',
  502: '网关错误',
  503: '服务不可用',
  504: '网关超时'
}

// 错误码映射
const ERROR_CODE_MESSAGES: Record<string, string> = {
  // 认证相关
  'AUTH_001': '用户名或密码错误',
  'AUTH_002': '用户不存在',
  'AUTH_003': '用户已被禁用',
  'AUTH_004': 'Token已过期',
  'AUTH_005': 'Token无效',
  'AUTH_006': '权限不足',
  
  // 用户相关
  'USER_001': '用户已存在',
  'USER_002': '用户注册失败',
  'USER_003': '用户信息更新失败',
  'USER_004': '用户不存在',
  
  // 聊天相关
  'CHAT_001': '聊天室不存在',
  'CHAT_002': '聊天室创建失败',
  'CHAT_003': '聊天室已满',
  'CHAT_004': '消息发送失败',
  'CHAT_005': '消息格式错误',
  'CHAT_006': '用户不在聊天室中',
  
  // WebSocket相关
  'WS_001': 'WebSocket连接失败',
  'WS_002': 'WebSocket认证失败',
  'WS_003': 'WebSocket消息格式错误'
}

/**
 * 处理HTTP错误
 */
export function handleHttpError(error: any): AppError {
  let appError: AppError
  
  if (error.response) {
    // 服务器响应错误
    const { status, data } = error.response
    const message = data?.message || HTTP_STATUS_MESSAGES[status] || '请求失败'
    
    appError = {
      code: status,
      message,
      details: data
    }
  } else if (error.request) {
    // 请求发送失败
    appError = {
      code: 'NETWORK_ERROR',
      message: '网络连接失败，请检查网络设置',
      details: error.request
    }
  } else {
    // 请求配置错误
    appError = {
      code: 'REQUEST_ERROR',
      message: error.message || '请求配置错误',
      details: error
    }
  }
  
  return appError
}

/**
 * 处理业务错误
 */
export function handleBusinessError(errorCode: string, message?: string): AppError {
  const errorMessage = message || ERROR_CODE_MESSAGES[errorCode] || '操作失败'
  
  return {
    code: errorCode,
    message: errorMessage
  }
}

/**
 * 显示错误消息
 */
export function showError(error: AppError | string, duration: number = 3000): void {
  const message = typeof error === 'string' ? error : error.message
  
  ElMessage({
    message,
    type: 'error',
    duration,
    showClose: true
  })
}

/**
 * 显示成功消息
 */
export function showSuccess(message: string, duration: number = 2000): void {
  ElMessage({
    message,
    type: 'success',
    duration,
    showClose: true
  })
}

/**
 * 显示警告消息
 */
export function showWarning(message: string, duration: number = 3000): void {
  ElMessage({
    message,
    type: 'warning',
    duration,
    showClose: true
  })
}

/**
 * 显示信息消息
 */
export function showInfo(message: string, duration: number = 3000): void {
  ElMessage({
    message,
    type: 'info',
    duration,
    showClose: true
  })
}

/**
 * 全局错误处理器
 */
export function setupGlobalErrorHandler(): void {
  // 处理未捕获的Promise错误
  window.addEventListener('unhandledrejection', (event) => {
    console.error('未处理的Promise错误:', event.reason)
    const error = handleHttpError(event.reason)
    showError(error)
  })
  
  // 处理未捕获的JavaScript错误
  window.addEventListener('error', (event) => {
    console.error('未捕获的JavaScript错误:', event.error)
    showError('页面出现错误，请刷新重试')
  })
}
