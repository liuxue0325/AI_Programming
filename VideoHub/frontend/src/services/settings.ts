import apiClient from './api'
import type { Setting } from '@/types/api'
import type { ApiResponse } from '@/types/api'

export const settingsAPI = {
  // 获取系统设置
  getSettings() {
    return apiClient.get<ApiResponse<Setting[]>>('/settings')
  },
  
  // 更新系统设置
  updateSettings(settings: Setting[]) {
    return apiClient.put<ApiResponse<void>>('/settings', settings)
  },
}