import { authApi } from '@/api/auth'
import type { User } from '@/types/user'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('token'))
  
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  
  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }
  
  const clearToken = () => {
    token.value = null
    localStorage.removeItem('token')
  }
  
  const setUser = (userData: User) => {
    user.value = userData
  }
  
  const login = async (username: string, password: string) => {
    try {
      const response = await authApi.login(username, password)
      const { user: userData, token: userToken } = response.data
      
      setToken(userToken)
      setUser(userData)
      
      return { success: true }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.error || '登录失败'
      }
    }
  }
  
  const register = async (userData: {
    username: string
    nickname: string
    email: string
    password: string
  }) => {
    try {
      const response = await authApi.register(userData)
      const { user: newUser, token: userToken } = response.data
      
      setToken(userToken)
      setUser(newUser)
      
      return { success: true }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.error || '注册失败'
      }
    }
  }
  
  const logout = () => {
    user.value = null
    clearToken()
  }
  
  const fetchUserInfo = async () => {
    try {
      const response = await authApi.getProfile()
      setUser(response.data.user)
    } catch (error) {
      // Token可能过期，清除token
      clearToken()
    }
  }
  
  const updateProfile = async (userData: Partial<User>) => {
    try {
      const response = await authApi.updateProfile(userData)
      setUser(response.data.user)
      return { success: true }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.error || '更新失败'
      }
    }
  }
  
  return {
    user,
    token,
    isAuthenticated,
    setToken,
    setUser,
    login,
    register,
    logout,
    fetchUserInfo,
    updateProfile
  }
})
