import { defineStore } from 'pinia'
import { settingsAPI } from '@/services/settings'
import type { Setting } from '@/types/api'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    settings: [] as Setting[],
    loading: false,
    error: null as string | null,
  }),
  
  getters: {
    getSetting: (state) => (key: string) => {
      return state.settings.find(setting => setting.key === key)?.value
    },
  },
  
  actions: {
    async fetchSettings() {
      this.loading = true
      this.error = null
      try {
        const response = await settingsAPI.getSettings()
        this.settings = response.data.data
      } catch (error) {
        this.error = '获取系统设置失败'
        console.error(error)
      } finally {
        this.loading = false
      }
    },
    
    async updateSettings(settings: Setting[]) {
      this.loading = true
      this.error = null
      try {
        await settingsAPI.updateSettings(settings)
        // 重新获取设置
        await this.fetchSettings()
      } catch (error) {
        this.error = '更新系统设置失败'
        console.error(error)
      } finally {
        this.loading = false
      }
    },
  },
})