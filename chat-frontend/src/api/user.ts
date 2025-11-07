import type { User } from '@/types/user'
import request from './request'

export const userApi = {
  // 搜索用户
  searchUsers: (query: string) => {
    return request.get<{ users: User[] }>('/users/search', {
      params: { q: query }
    })
  },
  
  // 获取用户详情
  getUserById: (userId: number) => {
    return request.get<{ user: User }>(`/users/${userId}`)
  },
  
  // 获取用户列表（用于创建房间时选择成员）
  getUsers: () => {
    return request.get<{ users: User[] }>('/users/search', {
      params: { q: '' }
    })
  }
}

export default userApi
