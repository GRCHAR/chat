import { useAuthStore } from '@/stores/auth'
import { onUnmounted, ref } from 'vue'

export function useWebSocket() {
  const socket = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const reconnectAttempts = ref(0)
  const maxReconnectAttempts = 5
  const reconnectDelay = 3000
  
  const messageHandlers: Array<(data: any) => void> = []
  
  const connect = () => {
    const authStore = useAuthStore()
    if (!authStore.token) return
    
    const wsUrl = `ws://localhost:8080/api/v1/ws?token=${authStore.token}`
    
    try {
      socket.value = new WebSocket(wsUrl)
      
      socket.value.onopen = () => {
        console.log('WebSocket连接成功')
        isConnected.value = true
        reconnectAttempts.value = 0
      }
      
      socket.value.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          console.log('收到WebSocket消息:', data)
          
          // 调用所有消息处理器
          messageHandlers.forEach(handler => {
            try {
              handler(data)
            } catch (error) {
              console.error('消息处理器错误:', error)
            }
          })
        } catch (error) {
          console.error('解析WebSocket消息失败:', error)
        }
      }
      
      socket.value.onclose = () => {
        console.log('WebSocket连接关闭')
        isConnected.value = false
        socket.value = null
        
        // 尝试重连
        if (reconnectAttempts.value < maxReconnectAttempts) {
          setTimeout(() => {
            reconnectAttempts.value++
            console.log(`尝试重连 (${reconnectAttempts.value}/${maxReconnectAttempts})`)
            connect()
          }, reconnectDelay)
        }
      }
      
      socket.value.onerror = (error) => {
        console.error('WebSocket错误:', error)
        isConnected.value = false
      }
      
    } catch (error) {
      console.error('创建WebSocket连接失败:', error)
    }
  }
  
  const disconnect = () => {
    if (socket.value) {
      socket.value.close()
      socket.value = null
    }
    isConnected.value = false
    reconnectAttempts.value = maxReconnectAttempts // 停止重连
  }
  
  const sendMessage = (data: any) => {
    if (socket.value && isConnected.value) {
      try {
        socket.value.send(JSON.stringify(data))
        return true
      } catch (error) {
        console.error('发送WebSocket消息失败:', error)
        return false
      }
    } else {
      console.warn('WebSocket未连接，无法发送消息')
      return false
    }
  }
  
  const onMessage = (handler: (data: any) => void) => {
    messageHandlers.push(handler)
    
    // 返回取消订阅函数
    return () => {
      const index = messageHandlers.indexOf(handler)
      if (index > -1) {
        messageHandlers.splice(index, 1)
      }
    }
  }
  
  // 组件卸载时断开连接
  onUnmounted(() => {
    disconnect()
  })
  
  return {
    socket,
    isConnected,
    connect,
    disconnect,
    sendMessage,
    onMessage
  }
}
