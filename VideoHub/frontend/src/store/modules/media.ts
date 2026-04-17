import { defineStore } from 'pinia'
import { mediaAPI } from '@/services/media'
import type { Media, MediaListParams } from '@/types/media'

export const useMediaStore = defineStore('media', {
  state: () => ({
    mediaList: [] as Media[],
    currentMedia: null as Media | null,
    loading: false,
    error: null as string | null,
  }),
  
  getters: {
    getMediaById: (state) => (id: number) => {
      return state.mediaList.find(media => media.id === id) || state.currentMedia
    },
  },
  
  actions: {
    async fetchMediaList(params?: MediaListParams) {
      this.loading = true
      this.error = null
      try {
        const response = await mediaAPI.getMediaList(params)
        this.mediaList = response.data.data
      } catch (error) {
        this.error = '获取媒体列表失败'
        console.error(error)
      } finally {
        this.loading = false
      }
    },
    
    async fetchMediaDetail(id: number) {
      this.loading = true
      this.error = null
      try {
        const response = await mediaAPI.getMediaDetail(id)
        this.currentMedia = response.data.data
      } catch (error) {
        this.error = '获取媒体详情失败'
        console.error(error)
      } finally {
        this.loading = false
      }
    },
    
    async searchMedia(keyword: string) {
      this.loading = true
      this.error = null
      try {
        const response = await mediaAPI.getMediaList({ search: keyword })
        this.mediaList = response.data.data
      } catch (error) {
        this.error = '搜索媒体失败'
        console.error(error)
      } finally {
        this.loading = false
      }
    },
  },
})