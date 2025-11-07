import { chatApi } from '@/api/chat'
import { useWebSocket } from '@/composables/useWebSocket'
import type { ChatRoom, Message } from '@/types/chat'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useChatStore = defineStore('chat', () => {
  const rooms = ref<ChatRoom[]>([])
  const currentRoom = ref<ChatRoom | null>(null)
  const messages = ref<Message[]>([])
  const unreadCounts = ref<Record<number, number>>({})
  const isLoading = ref(false)
  
  const { connect, disconnect, onMessage } = useWebSocket()
  
  const totalUnreadCount = computed(() => {
    return Object.values(unreadCounts.value).reduce((sum, count) => sum + count, 0)
  })
  
  const setRooms = (roomList: ChatRoom[]) => {
    rooms.value = roomList
  }
  
  const setCurrentRoom = (room: ChatRoom | null) => {
    currentRoom.value = room
    if (room) {
      // 清除当前房间的未读数
      unreadCounts.value[room.id] = 0
      // 获取房间消息
      fetchRoomMessages(room.id)
    }
  }
  
  const addMessage = (message: Message) => {
    messages.value.push(message)
    // 保持消息列表长度，避免内存溢出
    if (messages.value.length > 1000) {
      messages.value = messages.value.slice(-500)
    }
  }
  
  const updateUnreadCount = (roomId: number, count: number) => {
    unreadCounts.value[roomId] = count
  }
  
  const fetchRooms = async () => {
    try {
      isLoading.value = true
      const response = await chatApi.getRooms()
      setRooms(response.data.rooms)
      
      // 获取未读数
      const unreadResponse = await chatApi.getRoomsWithUnread()
      unreadResponse.data.rooms.forEach((room: any) => {
        updateUnreadCount(room.id, room.unread_count)
      })
    } catch (error) {
      console.error('获取房间列表失败:', error)
    } finally {
      isLoading.value = false
    }
  }
  
  const fetchRoomMessages = async (roomId: number, page = 1, pageSize = 20) => {
    try {
      const response = await chatApi.getMessages(roomId, page, pageSize)
      if (page === 1) {
        messages.value = response.data.messages
      } else {
        messages.value = [...response.data.messages, ...messages.value]
      }
      return response.data
    } catch (error) {
      console.error('获取消息失败:', error)
      return { messages: [], page: 1, page_size: pageSize }
    }
  }
  
  const createRoom = async (roomData: {
    name: string
    description: string
    type: 'single' | 'group'
    member_ids: number[]
  }) => {
    try {
      const response = await chatApi.createRoom(roomData)
      const newRoom = response.data.room
      rooms.value.push(newRoom)
      return { success: true, room: newRoom }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.error || '创建房间失败'
      }
    }
  }
  
  const joinRoom = async (roomId: number) => {
    try {
      await chatApi.joinRoom(roomId)
      await fetchRooms()
      return { success: true }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.error || '加入房间失败'
      }
    }
  }
  
  const leaveRoom = async (roomId: number) => {
    try {
      await chatApi.leaveRoom(roomId)
      rooms.value = rooms.value.filter(room => room.id !== roomId)
      if (currentRoom.value?.id === roomId) {
        setCurrentRoom(null)
      }
      return { success: true }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.error || '离开房间失败'
      }
    }
  }
  
  const sendChatMessage = async (content: string, type: 'text' | 'image' | 'file' = 'text') => {
    if (!currentRoom.value) return { success: false, error: '未选择聊天房间' }
    
    try {
      const response = await chatApi.sendMessage({
        room_id: currentRoom.value.id,
        content,
        type
      })
      
      const message = response.data.message
      addMessage(message)
      
      return { success: true, message }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.error || '发送消息失败'
      }
    }
  }
  
  const markAsRead = async (roomId: number) => {
    try {
      await chatApi.markAsRead(roomId)
      updateUnreadCount(roomId, 0)
    } catch (error) {
      console.error('标记已读失败:', error)
    }
  }
  
  const getUnreadCount = async (roomId: number) => {
    try {
      const response = await chatApi.getUnreadCount(roomId)
      updateUnreadCount(roomId, response.data.unread_count)
    } catch (error) {
      console.error('获取未读数失败:', error)
    }
  }
  
  // WebSocket相关
  const connectWebSocket = () => {
    connect()
    
    onMessage((data: any) => {
      if (data.type === 'message') {
        const message = data.message
        addMessage(message)
        
        // 如果不是当前房间的消息，增加未读数
        if (currentRoom.value?.id !== message.room_id) {
          const currentCount = unreadCounts.value[message.room_id] || 0
          updateUnreadCount(message.room_id, currentCount + 1)
        }
      }
    })
  }
  
  const disconnectWebSocket = () => {
    disconnect()
  }
  
  return {
    rooms,
    currentRoom,
    messages,
    unreadCounts,
    isLoading,
    totalUnreadCount,
    setRooms,
    setCurrentRoom,
    addMessage,
    updateUnreadCount,
    fetchRooms,
    fetchRoomMessages,
    createRoom,
    joinRoom,
    leaveRoom,
    sendChatMessage,
    markAsRead,
    getUnreadCount,
    connectWebSocket,
    disconnectWebSocket
  }
})
