import { defineStore } from 'pinia'
import { mediaAPI } from '@/services/media'
import type { Favorite } from '@/types/media'

export const useFavoritesStore = defineStore('favorites', {
  state: () => ({
    favoritesList: [] as Favorite[],
    loading: false,
    error: null as string | null,
  }),
  
  getters: {
    isFavorite: (state) => (mediaId: number) => {
      return state.favoritesList.some(fav => fav.media_id === mediaId)
    },
  },
  
  actions: {
    async fetchFavorites() {
      this.loading = true
      this.error = null
      try {
        const response = await mediaAPI.getFavorites()
        this.favoritesList = response.data.data
      } catch (error) {
        this.error = '获取收藏列表失败'
        console.error(error)
      } finally {
        this.loading = false
      }
    },
    
    async toggleFavorite(mediaId: number) {
      try {
        await mediaAPI.toggleFavorite(mediaId)
        // 重新获取收藏列表
        await this.fetchFavorites()
      } catch (error) {
        console.error('切换收藏状态失败', error)
      }
    },
  },
})