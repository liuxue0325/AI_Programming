import apiClient from './api'
import type { Media, MediaListParams, WatchHistory, Favorite, Subtitle, Folder, PlayUrl } from '@/types/media'
import type { ApiResponse } from '@/types/api'

export const mediaAPI = {
  // 获取媒体列表
  getMediaList(params?: MediaListParams) {
    return apiClient.get<ApiResponse<Media[]>>('/media', { params })
  },
  
  // 获取媒体详情
  getMediaDetail(id: number) {
    return apiClient.get<ApiResponse<Media>>(`/media/${id}`)
  },
  
  // 获取播放地址
  getPlayUrl(id: number) {
    return apiClient.get<ApiResponse<PlayUrl>>(`/media/${id}/play`)
  },
  
  // 获取字幕列表
  getSubtitles(id: number) {
    return apiClient.get<ApiResponse<Subtitle[]>>(`/media/${id}/subtitles`)
  },
  
  // 获取文件夹列表
  getFolders() {
    return apiClient.get<ApiResponse<Folder[]>>('/folders')
  },
  
  // 触发媒体扫描
  triggerScan() {
    return apiClient.post<ApiResponse<void>>('/scan')
  },
  
  // 获取观看历史
  getHistory() {
    return apiClient.get<ApiResponse<WatchHistory[]>>('/history')
  },
  
  // 记录观看历史
  addHistory(data: { media_id: number; episode_id?: number; progress: number; completed: boolean }) {
    return apiClient.post<ApiResponse<void>>('/history', data)
  },
  
  // 获取收藏列表
  getFavorites() {
    return apiClient.get<ApiResponse<Favorite[]>>('/favorites')
  },
  
  // 添加/删除收藏
  toggleFavorite(mediaId: number) {
    return apiClient.post<ApiResponse<void>>('/favorites', { media_id: mediaId })
  },
}