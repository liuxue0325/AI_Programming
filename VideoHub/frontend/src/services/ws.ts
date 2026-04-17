import type { WebSocketMessage } from '@/types/api'

class WebSocketService {
  private ws: WebSocket | null = null
  private listeners: Map<string, Function[]> = new Map()
  
  connect() {
    if (this.ws) {
      return
    }
    
    this.ws = new WebSocket('ws://localhost:8080/api/ws')
    
    this.ws.onopen = () => {
      console.log('WebSocket connected')
    }
    
    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data) as WebSocketMessage
        const { type, data } = message
        
        if (this.listeners.has(type)) {
          const callbacks = this.listeners.get(type)
          callbacks?.forEach(callback => callback(data))
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error)
      }
    }
    
    this.ws.onclose = () => {
      console.log('WebSocket disconnected')
      this.ws = null
      // 尝试重连
      setTimeout(() => this.connect(), 5000)
    }
    
    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
  }
  
  on(type: string, callback: Function) {
    if (!this.listeners.has(type)) {
      this.listeners.set(type, [])
    }
    this.listeners.get(type)?.push(callback)
  }
  
  off(type: string, callback: Function) {
    if (this.listeners.has(type)) {
      const callbacks = this.listeners.get(type)
      if (callbacks) {
        this.listeners.set(type, callbacks.filter(cb => cb !== callback))
      }
    }
  }
  
  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.listeners.clear()
  }
}

export const wsService = new WebSocketService()