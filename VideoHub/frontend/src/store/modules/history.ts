import { defineStore } from 'pinia'
import { mediaAPI } from '@/services/media'
import type { WatchHistory } from '@/types/media'

export const useHistoryStore = defineStore('history', {
  state: () => ({
    historyList: [] as WatchHistory[],
    loading: false,
    error: null as string | null,
  }),
  
  actions: {
    async fetchHistory() {
      this.loading = true
      this.error = null
      try {
        const response = await mediaAPI.getHistory()
        this.historyList = response.data.data
      } catch (error) {
        this.error = '获取观看历史失败'
        console.error(error)
      } finally {
        this.loading = false
      }
    },
    
    async addHistory(data: { media_id: number; episode_id?: number; progress: number; completed: boolean }) {
      try {
        await mediaAPI.addHistory(data)
        // 重新获取历史列表
        await this.fetchHistory()
      } catch (error) {
        console.error('添加观看历史失败', error)
      }
    },
  },
})