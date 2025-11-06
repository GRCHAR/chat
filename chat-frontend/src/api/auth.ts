import type { AuthResponse, LoginRequest, RegisterRequest, User } from '@/types/user'
import request from './request'

export const authApi = {
  // 用户登录
  login: (data: LoginRequest) => {
    return request.post<AuthResponse>('/auth/login', data)
  },
  
  // 用户注册
  register: (data: RegisterRequest) => {
    return request.post<AuthResponse>('/auth/register', data)
  },
  
  // 获取用户信息
  getProfile: () => {
    return request.get<{ user: User }>('/users/profile')
  },
  
  // 更新用户信息
  updateProfile: (data: Partial<User>) => {
    return request.put<{ user: User }>('/users/profile', data)
  }
}

export default authApi
