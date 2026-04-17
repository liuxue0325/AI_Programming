export interface Media {
  id: number
  type: 'movie' | 'tv' | 'anime'
  title: string
  original_title?: string
  year?: number
  path: string
  poster_path?: string
  backdrop_path?: string
  overview?: string
  rating?: number
  runtime?: number
  created_at?: string
  updated_at?: string
}

export interface Episode {
  id: number
  series_id: number
  season_number: number
  episode_number: number
  title?: string
  path: string
  poster_path?: string
  overview?: string
  runtime?: number
}

export interface WatchHistory {
  id: number
  media_id: number
  episode_id?: number
  progress: number
  completed: boolean
  last_watched: string
  media?: Media
  episode?: Episode
}

export interface Favorite {
  id: number
  media_id: number
  created_at: string
  media?: Media
}

export interface MediaListParams {
  type?: 'movie' | 'tv' | 'anime'
  sort?: 'title' | 'year' | 'rating' | 'created_at'
  order?: 'asc' | 'desc'
  page?: number
  limit?: number
  search?: string
}

export interface Subtitle {
  id: number
  name: string
  url: string
}

export interface Folder {
  id: number
  name: string
  path: string
}

export interface PlayUrl {
  url: string
}