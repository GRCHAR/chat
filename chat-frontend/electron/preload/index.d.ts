import { ElectronAPI } from '@electron-toolkit/preload'

declare global {
  interface Window {
    electron: ElectronAPI
    electronAPI: {
      // 应用信息
      getAppVersion: () => Promise<string>
      getPlatform: () => Promise<string>
      
      // 窗口控制
      minimizeWindow: () => void
      maximizeWindow: () => void
      closeWindow: () => void
      
      // 菜单事件监听
      onMenuEvent: (callback: (event: string) => void) => () => void
      
      // 系统通知
      showNotification: (title: string, body: string) => void
      
      // 剪贴板
      writeText: (text: string) => Promise<void>
      readText: () => Promise<string>
      
      // 本地存储增强
      localStorage: {
        setItem: (key: string, value: string) => void
        getItem: (key: string) => string | null
        removeItem: (key: string) => void
        clear: () => void
      }
    }
  }
}
