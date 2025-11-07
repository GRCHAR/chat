import { APP_CONFIG, THEME_CONFIG } from '@/config'

// 存储键名定义
export const STORAGE_KEYS = {
  // 用户认证信息
  TOKEN: 'token',
  USER_INFO: 'user_info',
  
  // 应用设置
  THEME: 'theme',
  LANGUAGE: 'language',
  
  // 聊天设置
  CHAT_SETTINGS: 'chat_settings',
  RECENT_ROOMS: 'recent_rooms',
  
  // 系统设置
  SYSTEM_SETTINGS: 'system_settings'
} as const

// 存储工具类
class StorageManager {
  private prefix: string
  
  constructor(prefix: string = APP_CONFIG.storagePrefix) {
    this.prefix = prefix
  }
  
  /**
   * 获取完整的存储键名
   */
  private getKey(key: string): string {
    return `${this.prefix}${key}`
  }
  
  /**
   * 设置本地存储
   */
  set<T>(key: string, value: T, expireTime?: number): void {
    try {
      const data = {
        value,
        expireTime: expireTime ? Date.now() + expireTime : null,
        timestamp: Date.now()
      }
      
      localStorage.setItem(this.getKey(key), JSON.stringify(data))
    } catch (error) {
      console.error('本地存储设置失败:', error)
      throw new Error('本地存储设置失败')
    }
  }
  
  /**
   * 获取本地存储
   */
  get<T>(key: string, defaultValue?: T): T | null {
    try {
      const item = localStorage.getItem(this.getKey(key))
      
      if (!item) {
        return defaultValue !== undefined ? defaultValue : null
      }
      
      const data = JSON.parse(item)
      
      // 检查是否过期
      if (data.expireTime && Date.now() > data.expireTime) {
        this.remove(key)
        return defaultValue !== undefined ? defaultValue : null
      }
      
      return data.value
    } catch (error) {
      console.error('本地存储读取失败:', error)
      return defaultValue !== undefined ? defaultValue : null
    }
  }
  
  /**
   * 删除本地存储
   */
  remove(key: string): void {
    try {
      localStorage.removeItem(this.getKey(key))
    } catch (error) {
      console.error('本地存储删除失败:', error)
    }
  }
  
  /**
   * 清空所有本地存储（只清空当前应用的前缀）
   */
  clear(): void {
    try {
      const keys = []
      for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i)
        if (key && key.startsWith(this.prefix)) {
          keys.push(key)
        }
      }
      
      keys.forEach(key => localStorage.removeItem(key))
    } catch (error) {
      console.error('本地存储清空失败:', error)
    }
  }
  
  /**
   * 检查是否存在
   */
  has(key: string): boolean {
    return this.get(key) !== null
  }
  
  /**
   * 获取所有键名
   */
  keys(): string[] {
    const keys = []
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key && key.startsWith(this.prefix)) {
        keys.push(key.substring(this.prefix.length))
      }
    }
    return keys
  }
}

// 创建存储管理器实例
export const storage = new StorageManager()

// 会话存储工具类
class SessionStorageManager {
  private prefix: string
  
  constructor(prefix: string = APP_CONFIG.storagePrefix) {
    this.prefix = prefix
  }
  
  /**
   * 获取完整的存储键名
   */
  private getKey(key: string): string {
    return `${this.prefix}${key}`
  }
  
  /**
   * 设置会话存储
   */
  set<T>(key: string, value: T): void {
    try {
      sessionStorage.setItem(this.getKey(key), JSON.stringify({
        value,
        timestamp: Date.now()
      }))
    } catch (error) {
      console.error('会话存储设置失败:', error)
      throw new Error('会话存储设置失败')
    }
  }
  
  /**
   * 获取会话存储
   */
  get<T>(key: string, defaultValue?: T): T | null {
    try {
      const item = sessionStorage.getItem(this.getKey(key))
      
      if (!item) {
        return defaultValue !== undefined ? defaultValue : null
      }
      
      const data = JSON.parse(item)
      return data.value
    } catch (error) {
      console.error('会话存储读取失败:', error)
      return defaultValue !== undefined ? defaultValue : null
    }
  }
  
  /**
   * 删除会话存储
   */
  remove(key: string): void {
    try {
      sessionStorage.removeItem(this.getKey(key))
    } catch (error) {
      console.error('会话存储删除失败:', error)
    }
  }
  
  /**
   * 清空所有会话存储（只清空当前应用的前缀）
   */
  clear(): void {
    try {
      const keys = []
      for (let i = 0; i < sessionStorage.length; i++) {
        const key = sessionStorage.key(i)
        if (key && key.startsWith(this.prefix)) {
          keys.push(key)
        }
      }
      
      keys.forEach(key => sessionStorage.removeItem(key))
    } catch (error) {
      console.error('会话存储清空失败:', error)
    }
  }
  
  /**
   * 检查是否存在
   */
  has(key: string): boolean {
    return this.get(key) !== null
  }
}

// 创建会话存储管理器实例
export const sessionStorageManager = new SessionStorageManager()

// 快捷方法
export const storageUtils = {
  // 认证相关
  setToken(token: string, expireTime?: number): void {
    storage.set(STORAGE_KEYS.TOKEN, token, expireTime)
  },
  
  getToken(): string | null {
    return storage.get<string>(STORAGE_KEYS.TOKEN)
  },
  
  removeToken(): void {
    storage.remove(STORAGE_KEYS.TOKEN)
  },
  
  // 用户信息
  setUserInfo(userInfo: any): void {
    storage.set(STORAGE_KEYS.USER_INFO, userInfo)
  },
  
  getUserInfo(): any {
    return storage.get(STORAGE_KEYS.USER_INFO)
  },
  
  removeUserInfo(): void {
    storage.remove(STORAGE_KEYS.USER_INFO)
  },
  
  // 主题
  setTheme(theme: string): void {
    storage.set(STORAGE_KEYS.THEME, theme)
  },
  
  getTheme(): string {
    return storage.get<string>(STORAGE_KEYS.THEME) || THEME_CONFIG.defaultTheme
  },
  
  // 清除所有用户数据
  clearUserData(): void {
    storage.remove(STORAGE_KEYS.TOKEN)
    storage.remove(STORAGE_KEYS.USER_INFO)
    sessionStorageManager.clear()
  }
}
