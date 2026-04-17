import { defineStore } from 'pinia'

export const useUIStore = defineStore('ui', {
  state: () => ({
    sidebarCollapsed: false,
    currentTheme: 'light',
    isMobile: false,
  }),
  
  actions: {
    toggleSidebar() {
      this.sidebarCollapsed = !this.sidebarCollapsed
    },
    
    setTheme(theme: 'light' | 'dark') {
      this.currentTheme = theme
    },
    
    setIsMobile(isMobile: boolean) {
      this.isMobile = isMobile
    },
  },
})