<template>
  <div class="library-view">
    <Header />
    <main class="main-content">
      <div class="container">
        <h1 class="page-title">媒体库</h1>
        
        <!-- 筛选和排序 -->
        <div class="filter-sort-container">
          <FilterBar 
            :filters="filters" 
            @filterChange="handleFilterChange"
          />
          <SortOptions 
            :sort="sort" 
            :order="order"
            @sortChange="handleSortChange"
          />
        </div>
        
        <!-- 媒体列表 -->
        <div class="media-grid">
          <MediaCard 
            v-for="media in mediaList" 
            :key="media.id"
            :media="media"
            @click="handleMediaClick"
          />
        </div>
        
        <!-- 分页 -->
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[12, 24, 36]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </main>
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Header from '@/components/layout/Header.vue'
import Footer from '@/components/layout/Footer.vue'
import FilterBar from '@/components/media/FilterBar.vue'
import SortOptions from '@/components/media/SortOptions.vue'
import MediaCard from '@/components/media/MediaCard.vue'
import type { Media } from '@/types/media'

const router = useRouter()

const filters = ref({
  type: '',
  year: '',
  search: ''
})

const sort = ref('title')
const order = ref('asc')
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(100)

// 模拟数据
const mediaList = ref<Media[]>([
  {
    id: 1,
    type: 'movie',
    title: '复仇者联盟4：终局之战',
    year: 2019,
    poster_path: 'https://via.placeholder.com/200x300?text=Avengers+Endgame',
    rating: 8.4,
    path: '/path/to/movie1.mp4'
  },
  {
    id: 2,
    type: 'tv',
    title: '权力的游戏',
    year: 2011,
    poster_path: 'https://via.placeholder.com/200x300?text=Game+of+Thrones',
    rating: 9.2,
    path: '/path/to/tv1.mp4'
  },
  {
    id: 3,
    type: 'anime',
    title: '进击的巨人',
    year: 2013,
    poster_path: 'https://via.placeholder.com/200x300?text=Attack+on+Titan',
    rating: 9.0,
    path: '/path/to/anime1.mp4'
  },
  {
    id: 4,
    type: 'movie',
    title: '蜘蛛侠：英雄远征',
    year: 2019,
    poster_path: 'https://via.placeholder.com/200x300?text=Spider-Man',
    rating: 7.8,
    path: '/path/to/movie2.mp4'
  },
  {
    id: 5,
    type: 'tv',
    title: '怪奇物语',
    year: 2016,
    poster_path: 'https://via.placeholder.com/200x300?text=Stranger+Things',
    rating: 8.7,
    path: '/path/to/tv2.mp4'
  },
  {
    id: 6,
    type: 'anime',
    title: '我的英雄学院',
    year: 2016,
    poster_path: 'https://via.placeholder.com/200x300?text=My+Hero+Academia',
    rating: 8.5,
    path: '/path/to/anime2.mp4'
  },
  {
    id: 7,
    type: 'movie',
    title: '银河护卫队3',
    year: 2023,
    poster_path: 'https://via.placeholder.com/200x300?text=Guardians+of+the+Galaxy',
    rating: 8.2,
    path: '/path/to/movie3.mp4'
  },
  {
    id: 8,
    type: 'movie',
    title: '盗梦空间',
    year: 2010,
    poster_path: 'https://via.placeholder.com/200x300?text=Inception',
    rating: 8.8,
    path: '/path/to/movie4.mp4'
  },
  {
    id: 9,
    type: 'movie',
    title: '星际穿越',
    year: 2014,
    poster_path: 'https://via.placeholder.com/200x300?text=Interstellar',
    rating: 8.6,
    path: '/path/to/movie5.mp4'
  },
  {
    id: 10,
    type: 'tv',
    title: '绝命毒师',
    year: 2008,
    poster_path: 'https://via.placeholder.com/200x300?text=Breaking+Bad',
    rating: 9.5,
    path: '/path/to/tv3.mp4'
  },
  {
    id: 11,
    type: 'anime',
    title: '钢之炼金术师FA',
    year: 2009,
    poster_path: 'https://via.placeholder.com/200x300?text=Fullmetal+Alchemist',
    rating: 9.2,
    path: '/path/to/anime3.mp4'
  },
  {
    id: 12,
    type: 'movie',
    title: '蝙蝠侠：黑暗骑士',
    year: 2008,
    poster_path: 'https://via.placeholder.com/200x300?text=The+Dark+Knight',
    rating: 9.0,
    path: '/path/to/movie6.mp4'
  }
])

const handleFilterChange = (newFilters: { type?: string; year?: string; search?: string }) => {
  filters.value = newFilters as typeof filters.value
  // 实际项目中，这里会根据筛选条件重新获取数据
  // await mediaStore.fetchMediaList({ ...newFilters, sort, order, page: currentPage.value, limit: pageSize.value })
}

const handleSortChange = (newSort: string, newOrder: string) => {
  sort.value = newSort
  order.value = newOrder
  // 实际项目中，这里会根据排序条件重新获取数据
  // await mediaStore.fetchMediaList({ ...filters.value, sort: newSort, order: newOrder, page: currentPage.value, limit: pageSize.value })
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  // 实际项目中，这里会根据新的页面大小重新获取数据
  // await mediaStore.fetchMediaList({ ...filters.value, sort, order, page: 1, limit: size })
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  // 实际项目中，这里会根据新的页码重新获取数据
  // await mediaStore.fetchMediaList({ ...filters.value, sort, order, page, limit: pageSize.value })
}

const handleMediaClick = (id: number) => {
  router.push(`/detail/${id}`)
}

onMounted(async () => {
  // 实际项目中，这里会从API获取数据
  // await mediaStore.fetchMediaList({ sort, order, page: currentPage.value, limit: pageSize.value })
})
</script>

<style scoped lang="scss">
.library-view {
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
    
    .page-title {
      font-size: 24px;
      font-weight: bold;
      margin-bottom: 20px;
      color: #303133;
    }
    
    .filter-sort-container {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;
      gap: 20px;
    }
    
    .media-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
      gap: 20px;
      margin-bottom: 32px;
    }
    
    .pagination-container {
      display: flex;
      justify-content: center;
      margin-top: 20px;
    }
  }
}

@media (max-width: 768px) {
  .library-view {
    .main-content {
      .filter-sort-container {
        flex-direction: column;
        align-items: stretch;
      }
      
      .media-grid {
        grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
        gap: 15px;
      }
    }
  }
}
</style>