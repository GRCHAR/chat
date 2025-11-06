import type {
    CreateRoomRequest,
    MessageResponse,
    MessagesResponse,
    RoomResponse,
    RoomsResponse,
    SendMessageRequest,
    User
} from '@/types/chat'
import request from './request'

export const chatApi = {
  // 获取房间列表
  getRooms: () => {
    return request.get<RoomsResponse>('/rooms')
  },
  
  // 获取带未读数的房间列表
  getRoomsWithUnread: () => {
    return request.get<RoomsResponse>('/rooms/unread')
  },
  
  // 创建房间
  createRoom: (data: CreateRoomRequest) => {
    return request.post<RoomResponse>('/rooms', data)
  },
  
  // 获取房间详情
  getRoom: (roomId: number) => {
    return request.get<RoomResponse>(`/rooms/${roomId}`)
  },
  
  // 加入房间
  joinRoom: (roomId: number) => {
    return request.post(`/rooms/${roomId}/join`)
  },
  
  // 离开房间
  leaveRoom: (roomId: number) => {
    return request.post(`/rooms/${roomId}/leave`)
  },
  
  // 获取房间消息
  getMessages: (roomId: number, page = 1, pageSize = 20) => {
    return request.get<MessagesResponse>(`/rooms/${roomId}/messages`, {
      params: { page, page_size: pageSize }
    })
  },
  
  // 发送消息
  sendMessage: (data: SendMessageRequest) => {
    return request.post<MessageResponse>('/messages', data)
  },
  
  // 标记消息已读
  markAsRead: (roomId: number) => {
    return request.post(`/rooms/${roomId}/read`)
  },
  
  // 获取未读消息数
  getUnreadCount: (roomId: number) => {
    return request.get<{ unread_count: number }>(`/rooms/${roomId}/unread`)
  },
  
  // 获取房间成员
  getRoomMembers: (roomId: number) => {
    return request.get<{ members: User[] }>(`/rooms/${roomId}/members`)
  }
}

export default chatApi
