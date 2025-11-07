import type { User } from './user'

export type { User } from './user'

export interface ChatRoom {
  id: number
  name: string
  description?: string
  type: 'single' | 'group'
  owner_id: number
  max_members: number
  created_at: string
  updated_at: string
  members?: User[]
  last_message?: Message
  unread_count?: number
}

export interface Message {
  id: number
  room_id: number
  sender_id: number
  sender?: User
  content: string
  type: 'text' | 'image' | 'file'
  created_at: string
  updated_at: string
  is_read: boolean
}

export interface CreateRoomRequest {
  name: string
  description?: string
  type: 'single' | 'group'
  member_ids: number[]
}

export interface SendMessageRequest {
  room_id: number
  content: string
  type: 'text' | 'image' | 'file'
}

export interface MessageResponse {
  message: Message
}

export interface RoomResponse {
  room: ChatRoom
}

export interface RoomsResponse {
  rooms: ChatRoom[]
}

export interface MessagesResponse {
  messages: Message[]
  page: number
  page_size: number
}
