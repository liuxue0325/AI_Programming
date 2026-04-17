<template>
  <div class="home-view">
    <Header />
    <main class="main-content">
      <div class="container">
        <!-- 轮播图 -->
        <section class="carousel-section">
          <Carousel 
            :items="recommendedMedia" 
            @click="handleMediaClick"
          />
        </section>
        
        <!-- 分类导航 -->
        <section class="category-section">
          <div class="category-nav">
            <button 
              v-for="category in categories" 
              :key="category.type"
              class="category-btn"
              :class="{ active: currentCategory === category.type }"
              @click="currentCategory = category.type"
            >
              {{ category.name }}
            </button>
          </div>
        </section>
        
        <!-- 最近添加 -->
        <section class="recent-section">
          <h2 class="section-title">最近添加</h2>
          <div class="media-grid">
            <MediaCard 
              v-for="media in recentMedia" 
              :key="media.id"
              :media="media"
              @click="handleMediaClick"
            />
          </div>
        </section>
        
        <!-- 推荐电影 -->
        <section class="recommended-section" v-if="currentCategory === 'movie'">
          <h2 class="section-title">推荐电影</h2>
          <div class="media-grid">
            <MediaCard 
              v-for="media in recommendedMovies" 
              :key="media.id"
              :media="media"
              @click="handleMediaClick"
            />
          </div>
        </section>
        
        <!-- 推荐电视剧 -->
        <section class="recommended-section" v-if="currentCategory === 'tv'">
          <h2 class="section-title">推荐电视剧</h2>
          <div class="media-grid">
            <MediaCard 
              v-for="media in recommendedTV" 
              :key="media.id"
              :media="media"
              @click="handleMediaClick"
            />
          </div>
        </section>
        
        <!-- 推荐动画 -->
        <section class="recommended-section" v-if="currentCategory === 'anime'">
          <h2 class="section-title">推荐动画</h2>
          <div class="media-grid">
            <MediaCard 
              v-for="media in recommendedAnime" 
              :key="media.id"
              :media="media"
              @click="handleMediaClick"
            />
          </div>
        </section>
      </div>
    </main>
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import Header from '@/components/layout/Header.vue'
import Footer from '@/components/layout/Footer.vue'
import Carousel from '@/components/media/Carousel.vue'
import MediaCard from '@/components/media/MediaCard.vue'

const router = useRouter()

const currentCategory = ref('movie')
const categories = [
  { type: 'movie', name: '电影' },
  { type: 'tv', name: '电视剧' },
  { type: 'anime', name: '动画' },
]

// 模拟数据
const recommendedMedia = ref([
  {
    id: 1,
    type: 'movie' as const,
    title: '复仇者联盟4：终局之战',
    year: 2019,
    backdrop_path: 'https://via.placeholder.com/1200x600?text=Avengers+Endgame',
    overview: '复仇者联盟与他们的盟友必须愿意牺牲一切，以在灭霸造成的破坏后扭转局势。',
    rating: 8.4,
    path: '/path/to/movie1.mp4'
  },
  {
    id: 2,
    type: 'tv' as const,
    title: '权力的游戏',
    year: 2011,
    backdrop_path: 'https://via.placeholder.com/1200x600?text=Game+of+Thrones',
    overview: '九大家族为夺取维斯特洛大陆的控制权而相互征战。',
    rating: 9.2,
    path: '/path/to/tv1.mp4'
  },
  {
    id: 3,
    type: 'anime' as const,
    title: '进击的巨人',
    year: 2013,
    backdrop_path: 'https://via.placeholder.com/1200x600?text=Attack+on+Titan',
    overview: '人类在巨大围墙内生存，对抗吃人的巨人。',
    rating: 9.0,
    path: '/path/to/anime1.mp4'
  }
])

const recentMedia = ref([
  {
    id: 4,
    type: 'movie' as const,
    title: '蜘蛛侠：英雄远征',
    year: 2019,
    poster_path: 'https://via.placeholder.com/200x300?text=Spider-Man',
    rating: 7.8,
    path: '/path/to/movie2.mp4'
  },
  {
    id: 5,
    type: 'tv' as const,
    title: '怪奇物语',
    year: 2016,
    poster_path: 'https://via.placeholder.com/200x300?text=Stranger+Things',
    rating: 8.7,
    path: '/path/to/tv2.mp4'
  },
  {
    id: 6,
    type: 'anime' as const,
    title: '我的英雄学院',
    year: 2016,
    poster_path: 'https://via.placeholder.com/200x300?text=My+Hero+Academia',
    rating: 8.5,
    path: '/path/to/anime2.mp4'
  },
  {
    id: 7,
    type: 'movie' as const,
    title: '银河护卫队3',
    year: 2023,
    poster_path: 'https://via.placeholder.com/200x300?text=Guardians+of+the+Galaxy',
    rating: 8.2,
    path: '/path/to/movie3.mp4'
  }
])

const recommendedMovies = computed(() => {
  return [
    {
      id: 8,
      type: 'movie' as const,
      title: '盗梦空间',
      year: 2010,
      poster_path: 'https://via.placeholder.com/200x300?text=Inception',
      rating: 8.8,
      path: '/path/to/movie4.mp4'
    },
    {
      id: 9,
      type: 'movie' as const,
      title: '星际穿越',
      year: 2014,
      poster_path: 'https://via.placeholder.com/200x300?text=Interstellar',
      rating: 8.6,
      path: '/path/to/movie5.mp4'
    },
    {
      id: 10,
      type: 'movie' as const,
      title: '蝙蝠侠：黑暗骑士',
      year: 2008,
      poster_path: 'https://via.placeholder.com/200x300?text=The+Dark+Knight',
      rating: 9.0,
      path: '/path/to/movie6.mp4'
    },
    {
      id: 11,
      type: 'movie' as const,
      title: '黑客帝国',
      year: 1999,
      poster_path: 'https://via.placeholder.com/200x300?text=The+Matrix',
      rating: 8.7,
      path: '/path/to/movie7.mp4'
    }
  ]
})

const recommendedTV = computed(() => {
  return [
    {
      id: 12,
      type: 'tv' as const,
      title: '绝命毒师',
      year: 2008,
      poster_path: 'https://via.placeholder.com/200x300?text=Breaking+Bad',
      rating: 9.5,
      path: '/path/to/tv3.mp4'
    },
    {
      id: 13,
      type: 'tv' as const,
      title: '老友记',
      year: 1994,
      poster_path: 'https://via.placeholder.com/200x300?text=Friends',
      rating: 8.9,
      path: '/path/to/tv4.mp4'
    },
    {
      id: 14,
      type: 'tv' as const,
      title: '生活大爆炸',
      year: 2007,
      poster_path: 'https://via.placeholder.com/200x300?text=The+Big+Bang+Theory',
      rating: 8.2,
      path: '/path/to/tv5.mp4'
    },
    {
      id: 15,
      type: 'tv' as const,
      title: '西部世界',
      year: 2016,
      poster_path: 'https://via.placeholder.com/200x300?text=Westworld',
      rating: 8.7,
      path: '/path/to/tv6.mp4'
    }
  ]
})

const recommendedAnime = computed(() => {
  return [
    {
      id: 16,
      type: 'anime' as const,
      title: '钢之炼金术师FA',
      year: 2009,
      poster_path: 'https://via.placeholder.com/200x300?text=Fullmetal+Alchemist',
      rating: 9.2,
      path: '/path/to/anime3.mp4'
    },
    {
      id: 17,
      type: 'anime' as const,
      title: '死亡笔记',
      year: 2006,
      poster_path: 'https://via.placeholder.com/200x300?text=Death+Note',
      rating: 9.0,
      path: '/path/to/anime4.mp4'
    },
    {
      id: 18,
      type: 'anime' as const,
      title: '海贼王',
      year: 1999,
      poster_path: 'https://via.placeholder.com/200x300?text=One+Piece',
      rating: 8.9,
      path: '/path/to/anime5.mp4'
    },
    {
      id: 19,
      type: 'anime' as const,
      title: '火影忍者',
      year: 2002,
      poster_path: 'https://via.placeholder.com/200x300?text=Naruto',
      rating: 8.6,
      path: '/path/to/anime6.mp4'
    }
  ]
})

const handleMediaClick = (id: number) => {
  router.push(`/detail/${id}`)
}

onMounted(async () => {
  // 实际项目中，这里会从API获取数据
  // await mediaStore.fetchMediaList()
})
</script>

<style scoped lang="scss">
.home-view {
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
    
    .carousel-section {
      margin-bottom: 32px;
    }
    
    .category-section {
      margin-bottom: 32px;
      
      .category-nav {
        display: flex;
        gap: 12px;
        
        .category-btn {
          padding: 8px 16px;
          border: 1px solid #409EFF;
          border-radius: 20px;
          background-color: #fff;
          color: #409EFF;
          cursor: pointer;
          transition: all 0.3s ease;
          
          &:hover {
            background-color: #409EFF;
            color: #fff;
          }
          
          &.active {
            background-color: #409EFF;
            color: #fff;
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
    
    .media-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
      gap: 20px;
      margin-bottom: 32px;
    }
  }
}

@media (max-width: 768px) {
  .home-view {
    .main-content {
      .media-grid {
        grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
        gap: 15px;
      }
    }
  }
}
</style>