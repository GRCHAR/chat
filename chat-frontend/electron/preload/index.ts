import { contextBridge, ipcRenderer } from 'electron'

// Custom APIs for renderer
const api = {
  // 应用信息
  getAppVersion: () => ipcRenderer.invoke('app-version'),
  getPlatform: () => ipcRenderer.invoke('platform'),
  
  // 窗口控制
  minimizeWindow: () => ipcRenderer.send('minimize-window'),
  maximizeWindow: () => ipcRenderer.send('maximize-window'),
  closeWindow: () => ipcRenderer.send('close-window'),
  
  // 菜单事件监听
  onMenuEvent: (callback: (event: string) => void) => {
    const menuEvents = ['menu-new-chat', 'menu-about']
    menuEvents.forEach(event => {
      ipcRenderer.on(event, () => callback(event))
    })
    
    // 返回取消监听函数
    return () => {
      menuEvents.forEach(event => {
        ipcRenderer.removeAllListeners(event)
      })
    }
  },
  
  // 系统通知
  showNotification: (title: string, body: string) => {
    if ('Notification' in window) {
      if (Notification.permission === 'granted') {
        new Notification(title, { body })
      } else if (Notification.permission !== 'denied') {
        Notification.requestPermission().then(permission => {
          if (permission === 'granted') {
            new Notification(title, { body })
          }
        })
      }
    }
  },
  
  // 剪贴板
  writeText: (text: string) => navigator.clipboard.writeText(text),
  readText: () => navigator.clipboard.readText(),
  
  // 本地存储增强
  localStorage: {
    setItem: (key: string, value: string) => localStorage.setItem(key, value),
    getItem: (key: string) => localStorage.getItem(key),
    removeItem: (key: string) => localStorage.removeItem(key),
    clear: () => localStorage.clear()
  }
}

// Use `contextBridge` APIs to expose Electron APIs to
// renderer only if context isolation is enabled, otherwise
// just add to the DOM global.
if (process.contextIsolated) {
  try {
    contextBridge.exposeInMainWorld('electronAPI', api)
  } catch (error) {
    console.error(error)
  }
} else {
  // @ts-ignore (define in dts for window)
  window.electronAPI = api
}
