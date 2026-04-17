<template>
  <div class="detail-view">
    <Header />
    <main class="main-content">
      <div class="container">
        <!-- 背景图 -->
        <div class="backdrop-container" v-if="media?.backdrop_path">
          <img :src="media.backdrop_path" :alt="media.title" class="backdrop" />
          <div class="backdrop-overlay"></div>
        </div>
        
        <!-- 媒体信息 -->
        <div class="media-info">
          <div class="poster-container">
            <img :src="media?.poster_path || defaultPoster" :alt="media?.title" class="poster" />
          </div>
          <div class="info-panel">
            <h1 class="title">{{ media?.title }}</h1>
            <div class="meta-info">
              <span class="year">{{ media?.year }}</span>
              <span class="type">{{ getMediaType(media?.type) }}</span>
              <span class="rating" v-if="media?.rating">
                <el-rate v-model="media.rating" :max="10" disabled />
                <span class="rating-text">{{ media.rating.toFixed(1) }}</span>
              </span>
            </div>
            <p class="overview" v-if="media?.overview">{{ media.overview }}</p>
            <div class="actions">
              <el-button type="primary" size="large" @click="handlePlay" class="play-btn">
                <el-icon><VideoPlay /></el-icon> 播放
              </el-button>
              <el-button size="large" @click="handleToggleFavorite" :type="isFavorite ? 'danger' : 'default'">
                <el-icon><Star /></el-icon> {{ isFavorite ? '取消收藏' : '收藏' }}
              </el-button>
            </div>
          </div>
        </div>
        
        <!-- 剧集列表（电视剧） -->
        <div class="episodes-section" v-if="media?.type === 'tv'">
          <h2 class="section-title">剧集</h2>
          <div class="episodes-list">
            <div 
              v-for="episode in episodes" 
              :key="episode.id"
              class="episode-item"
              @click="handlePlayEpisode(episode.id)"
            >
              <div class="episode-info">
                <h3 class="episode-title">S{{ episode.season_number.toString().padStart(2, '0') }}E{{ episode.episode_number.toString().padStart(2, '0') }} {{ episode.title }}</h3>
                <p class="episode-overview" v-if="episode.overview">{{ episode.overview }}</p>
              </div>
              <div class="episode-duration" v-if="episode.runtime">{{ formatDuration(episode.runtime) }}</div>
            </div>
          </div>
        </div>
        
        <!-- 相关推荐 -->
        <div class="related-section">
          <h2 class="section-title">相关推荐</h2>
          <div class="media-grid">
            <MediaCard 
              v-for="item in relatedMedia" 
              :key="item.id"
              :media="item"
              @click="handleMediaClick"
            />
          </div>
        </div>
      </div>
    </main>
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Star, VideoPlay } from '@element-plus/icons-vue'
import { useFavoritesStore } from '@/store'
import Header from '@/components/layout/Header.vue'
import Footer from '@/components/layout/Footer.vue'
import MediaCard from '@/components/media/MediaCard.vue'
import type { Media, Episode } from '@/types/media'

const route = useRoute()
const router = useRouter()
const favoritesStore = useFavoritesStore()

const mediaId = computed(() => Number(route.params.id))
const media = ref<Media | null>(null)
const episodes = ref<Episode[]>([])
const relatedMedia = ref<Media[]>([])

// 模拟数据
const mockMedia: Media = {
  id: 1,
  type: 'movie',
  title: '复仇者联盟4：终局之战',
  original_title: 'Avengers: Endgame',
  year: 2019,
  path: '/path/to/movie1.mp4',
  poster_path: 'https://via.placeholder.com/300x450?text=Avengers+Endgame',
  backdrop_path: 'https://via.placeholder.com/1200x600?text=Avengers+Endgame',
  overview: '复仇者联盟与他们的盟友必须愿意牺牲一切，以在灭霸造成的破坏后扭转局势。在《复仇者联盟3：无限战争》的毁灭性事件之后，宇宙处于混乱状态。剩下的复仇者们必须面对他们最大的挑战，因为他们试图扭转灭霸的行为并恢复宇宙的平衡。',
  rating: 8.4,
  runtime: 181
}

const mockEpisodes: Episode[] = [
  {
    id: 1,
    series_id: 2,
    season_number: 1,
    episode_number: 1,
    title: '凛冬将至',
    path: '/path/to/episode1.mp4',
    overview: '史塔克家族的领主艾德·史塔克被国王之手琼恩·艾林的死讯召回君临城。',
    runtime: 60
  },
  {
    id: 2,
    series_id: 2,
    season_number: 1,
    episode_number: 2,
    title: '王者之路',
    path: '/path/to/episode2.mp4',
    overview: '艾德·史塔克成为新的国王之手，而琼恩·雪诺决定加入守夜人。',
    runtime: 55
  }
]

const mockRelatedMedia: Media[] = [
  {
    id: 4,
    type: 'movie',
    title: '复仇者联盟3：无限战争',
    year: 2018,
    poster_path: 'https://via.placeholder.com/200x300?text=Avengers+Infinity+War',
    rating: 8.5,
    path: '/path/to/movie4.mp4'
  },
  {
    id: 7,
    type: 'movie',
    title: '银河护卫队3',
    year: 2023,
    poster_path: 'https://via.placeholder.com/200x300?text=Guardians+of+the+Galaxy',
    rating: 8.2,
    path: '/path/to/movie7.mp4'
  },
  {
    id: 8,
    type: 'movie',
    title: '美国队长3：内战',
    year: 2016,
    poster_path: 'https://via.placeholder.com/200x300?text=Captain+America+Civil+War',
    rating: 8.0,
    path: '/path/to/movie8.mp4'
  },
  {
    id: 9,
    type: 'movie',
    title: '雷神3：诸神黄昏',
    year: 2017,
    poster_path: 'https://via.placeholder.com/200x300?text=Thor+Ragnarok',
    rating: 7.9,
    path: '/path/to/movie9.mp4'
  }
]

const isFavorite = computed(() => {
  return favoritesStore.isFavorite(mediaId.value)
})

const getMediaType = (type?: string) => {
  const typeMap: Record<string, string> = {
    movie: '电影',
    tv: '电视剧',
    anime: '动画'
  }
  return typeMap[type || ''] || ''
}

const formatDuration = (minutes: number) => {
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return `${hours}小时${mins}分钟`
}

const handlePlay = () => {
  router.push(`/player/${mediaId.value}`)
}

const handlePlayEpisode = (episodeId: number) => {
  router.push(`/player/${mediaId.value}?episode=${episodeId}`)
}

const handleToggleFavorite = async () => {
  await favoritesStore.toggleFavorite(mediaId.value)
}

const handleMediaClick = (id: number) => {
  router.push(`/detail/${id}`)
}

onMounted(async () => {
  // 实际项目中，这里会从API获取数据
  // await mediaStore.fetchMediaDetail(mediaId.value)
  // media.value = mediaStore.currentMedia
  
  // 模拟数据
  media.value = mockMedia
  if (media.value.type === 'tv') {
    episodes.value = mockEpisodes
  }
  relatedMedia.value = mockRelatedMedia
  
  // 获取收藏状态
  await favoritesStore.fetchFavorites()
})

const defaultPoster = 'https://via.placeholder.com/300x450?text=No+Poster'
</script>

<style scoped lang="scss">
.detail-view {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  
  .main-content {
    flex: 1;
    padding: 20px 0;
    
    .container {
      width: 100%;
      max-width: 1200px;
      margin: 0 auto;
      padding: 0 20px;
    }
    
    .backdrop-container {
      position: relative;
      height: 400px;
      margin-bottom: -100px;
      z-index: 1;
      
      .backdrop {
        width: 100%;
        height: 100%;
        object-fit: cover;
        border-radius: 8px;
      }
      
      .backdrop-overlay {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(to bottom, transparent, rgba(0, 0, 0, 0.8));
        border-radius: 8px;
      }
    }
    
    .media-info {
      display: flex;
      gap: 32px;
      background-color: #fff;
      border-radius: 8px;
      padding: 32px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      margin-bottom: 32px;
      z-index: 2;
      position: relative;
      
      .poster-container {
        flex-shrink: 0;
        
        .poster {
          width: 200px;
          border-radius: 8px;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
        }
      }
      
      .info-panel {
        flex: 1;
        
        .title {
          font-size: 28px;
          font-weight: bold;
          margin-bottom: 16px;
          color: #303133;
        }
        
        .meta-info {
          display: flex;
          gap: 16px;
          margin-bottom: 16px;
          
          .year,
          .type {
            font-size: 14px;
            color: #606266;
            padding: 4px 8px;
            background-color: #f5f7fa;
            border-radius: 4px;
          }
          
          .rating {
            display: flex;
            align-items: center;
            gap: 8px;
            
            .rating-text {
              font-size: 14px;
              color: #606266;
            }
          }
        }
        
        .overview {
          font-size: 16px;
          line-height: 1.5;
          color: #303133;
          margin-bottom: 24px;
        }
        
        .actions {
          display: flex;
          gap: 16px;
          
          .play-btn {
            min-width: 120px;
          }
        }
      }
    }
    
    .section-title {
      font-size: 20px;
      font-weight: bold;
      margin-bottom: 16px;
      color: #303133;
    }
    
    .episodes-section {
      background-color: #fff;
      border-radius: 8px;
      padding: 24px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      margin-bottom: 32px;
      
      .episodes-list {
        display: flex;
        flex-direction: column;
        gap: 16px;
        
        .episode-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 16px;
          border-radius: 8px;
          background-color: #f5f7fa;
          cursor: pointer;
          transition: all 0.3s ease;
          
          &:hover {
            background-color: #e4e7ed;
          }
          
          .episode-info {
            flex: 1;
            
            .episode-title {
              font-size: 16px;
              font-weight: 500;
              margin-bottom: 8px;
              color: #303133;
            }
            
            .episode-overview {
              font-size: 14px;
              color: #606266;
              line-height: 1.4;
            }
          }
          
          .episode-duration {
            font-size: 14px;
            color: #909399;
            margin-left: 16px;
          }
        }
      }
    }
    
    .related-section {
      
      .media-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
        gap: 20px;
      }
    }
  }
}

@media (max-width: 768px) {
  .detail-view {
    .main-content {
      .backdrop-container {
        height: 200px;
        margin-bottom: -50px;
      }
      
      .media-info {
        flex-direction: column;
        padding: 20px;
        
        .poster-container {
          align-self: center;
          
          .poster {
            width: 150px;
          }
        }
        
        .info-panel {
          .title {
            font-size: 24px;
          }
          
          .actions {
            flex-direction: column;
            
            .play-btn {
              width: 100%;
            }
          }
        }
      }
      
      .episodes-section {
        padding: 16px;
      }
      
      .related-section {
        .media-grid {
          grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
          gap: 15px;
        }
      }
    }
  }
}
</style>