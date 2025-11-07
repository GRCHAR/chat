export interface User {
  id: number
  username: string
  nickname: string
  email: string
  avatar?: string
  status: 'active' | 'inactive' | 'offline'
  role?: string
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  nickname: string
  email: string
  password: string
}

export interface AuthResponse {
  user: User
  token: string
}
