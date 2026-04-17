export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface Setting {
  key: string
  value: string
}

export interface WebSocketMessage {
  type: string
  data: any
}