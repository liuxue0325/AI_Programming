<template>
  <div class="player-view">
    <Header />
    <main class="main-content">
      <div class="container">
        <!-- 播放器容器 -->
        <div class="player-container">
          <div id="dplayer" ref="playerRef"></div>
        </div>
        
        <!-- 播放控制 -->
        <div class="player-controls">
          <div class="control-section">
            <h3>字幕选择</h3>
            <el-select v-model="selectedSubtitle" @change="changeSubtitle">
              <el-option label="无字幕" value=""></el-option>
              <el-option 
                v-for="subtitle in subtitles" 
                :key="subtitle.id" 
                :label="subtitle.name" 
                :value="subtitle.url"
              ></el-option>
            </el-select>
          </div>
          <div class="control-section">
            <h3>画质选择</h3>
            <el-select v-model="selectedQuality" @change="changeQuality">
              <el-option label="自动" value="auto"></el-option>
              <el-option label="1080p" value="1080p"></el-option>
              <el-option label="720p" value="720p"></el-option>
              <el-option label="480p" value="480p"></el-option>
            </el-select>
          </div>
          <div class="control-section">
            <h3>播放控制</h3>
            <el-button-group>
              <el-button @click="playPause">{{ isPlaying ? '暂停' : '播放' }}</el-button>
              <el-button @click="stop">停止</el-button>
              <el-button @click="fullscreen">全屏</el-button>
            </el-button-group>
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
              :class="{ active: currentEpisodeId === episode.id }"
              @click="switchEpisode(episode.id)"
            >
              <div class="episode-info">
                <h3 class="episode-title">S{{ episode.season_number.toString().padStart(2, '0') }}E{{ episode.episode_number.toString().padStart(2, '0') }} {{ episode.title }}</h3>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useHistoryStore } from '@/store'
import Header from '@/components/layout/Header.vue'
import Footer from '@/components/layout/Footer.vue'
import type { Media, Episode, Subtitle } from '@/types/media'

// 导入DPlayer
import DPlayer from 'dplayer'

const route = useRoute()
const historyStore = useHistoryStore()

const mediaId = computed(() => Number(route.params.id))
const episodeId = computed(() => route.query.episode ? Number(route.query.episode) : undefined)

const playerRef = ref<HTMLElement | null>(null)
const player = ref<DPlayer | null>(null)
const isPlaying = ref(false)
const selectedSubtitle = ref('')
const selectedQuality = ref('auto')
const currentEpisodeId = ref(episodeId.value || 0)

// 模拟数据
const media = ref<Media>({
  id: 1,
  type: 'movie',
  title: '复仇者联盟4：终局之战',
  path: '/path/to/movie1.mp4',
  runtime: 181
})

const episodes = ref<Episode[]>([
  {
    id: 1,
    series_id: 2,
    season_number: 1,
    episode_number: 1,
    title: '凛冬将至',
    path: '/path/to/episode1.mp4'
  },
  {
    id: 2,
    series_id: 2,
    season_number: 1,
    episode_number: 2,
    title: '王者之路',
    path: '/path/to/episode2.mp4'
  }
])

const subtitles = ref<Subtitle[]>([
  { id: 1, name: '中文', url: '/subtitles/chinese.srt' },
  { id: 2, name: '英文', url: '/subtitles/english.srt' }
])

const playPause = () => {
  if (player.value) {
    if (isPlaying.value) {
      player.value.pause()
    } else {
      player.value.play()
    }
  }
}

const stop = () => {
  if (player.value) {
    player.value.seek(0)
    player.value.pause()
  }
}

const fullscreen = () => {
  if (player.value) {
    const video = player.value.video
    if (video.requestFullscreen) {
      video.requestFullscreen()
    } else if ((video as any).webkitRequestFullscreen) {
      (video as any).webkitRequestFullscreen()
    } else if ((video as any).msRequestFullscreen) {
      (video as any).msRequestFullscreen()
    }
  }
}

const changeSubtitle = () => {
  // 实际项目中，这里会切换字幕
  console.log('切换字幕到:', selectedSubtitle.value)
}

const changeQuality = () => {
  // 实际项目中，这里会切换视频画质
  console.log('切换画质到:', selectedQuality.value)
}

const switchEpisode = (id: number) => {
  currentEpisodeId.value = id
  // 实际项目中，这里会切换到对应剧集
  console.log('切换到剧集:', id)
}

const updatePlayHistory = () => {
  if (player.value) {
    const progress = player.value.video.currentTime / player.value.video.duration
    const completed = progress >= 0.95
    
    historyStore.addHistory({
      media_id: mediaId.value,
      episode_id: currentEpisodeId.value || undefined,
      progress,
      completed
    })
  }
}

onMounted(() => {
    if (playerRef.value) {
      // 初始化DPlayer
      player.value = new DPlayer({
        container: playerRef.value,
        video: {
          url: 'https://example.com/video.mp4', // 实际项目中，这里会从API获取播放地址
          pic: 'https://via.placeholder.com/1200x600?text=Video+Poster',
          thumbnails: 'https://example.com/thumbnails.jpg',
          type: 'auto'
        },
        danmaku: {
          id: 'demo',
          api: 'https://api.prprpr.me/dplayer/',
          token: 'demo'
        },
        autoplay: true,
        mutex: true
      })

      // 监听播放状态
      player.value.on('play' as any, () => {
        isPlaying.value = true
      })

      player.value.on('pause' as any, () => {
        isPlaying.value = false
      })

      player.value.on('ended' as any, () => {
        isPlaying.value = false
        updatePlayHistory()
      })

      // 定期更新播放历史
      setInterval(updatePlayHistory, 30000) // 每30秒更新一次
    }
  })

onUnmounted(() => {
  if (player.value) {
    player.value.destroy()
  }
})
</script>

<style scoped lang="scss">
.player-view {
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
    
    .player-container {
      width: 100%;
      aspect-ratio: 16/9;
      background-color: #000;
      border-radius: 8px;
      overflow: hidden;
      margin-bottom: 24px;
      
      #dplayer {
        width: 100%;
        height: 100%;
      }
    }
    
    .player-controls {
      display: flex;
      gap: 24px;
      background-color: #fff;
      border-radius: 8px;
      padding: 20px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      margin-bottom: 32px;
      
      .control-section {
        flex: 1;
        
        h3 {
          font-size: 14px;
          font-weight: 500;
          margin-bottom: 8px;
          color: #606266;
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
      
      .episodes-list {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
        gap: 12px;
        
        .episode-item {
          padding: 12px;
          border-radius: 8px;
          background-color: #f5f7fa;
          cursor: pointer;
          transition: all 0.3s ease;
          text-align: center;
          
          &:hover {
            background-color: #e4e7ed;
          }
          
          &.active {
            background-color: #409EFF;
            color: #fff;
          }
          
          .episode-title {
            font-size: 14px;
            font-weight: 500;
            margin: 0;
          }
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .player-view {
    .main-content {
      .player-controls {
        flex-direction: column;
        gap: 16px;
      }
      
      .episodes-section {
        .episodes-list {
          grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
        }
      }
    }
  }
}
</style>