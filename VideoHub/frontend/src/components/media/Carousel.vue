<template>
  <div class="carousel" ref="carouselRef">
    <div class="carousel-container" :style="carouselStyle">
      <div 
        v-for="item in items" 
        :key="item.id"
        class="carousel-item"
        @click="$emit('click', item.id)"
      >
        <img 
          :src="item.backdrop_path || defaultBackdrop" 
          :alt="item.title" 
          class="backdrop"
        />
        <div class="overlay">
          <h2 class="title">{{ item.title }}</h2>
          <p class="overview" v-if="item.overview">{{ truncateOverview(item.overview) }}</p>
        </div>
      </div>
    </div>
    <div class="carousel-indicators">
      <span 
        v-for="(item, idx) in items" 
        :key="item.id"
        class="indicator"
        :class="{ active: currentIndex === idx }"
        @click="goTo(idx)"
      ></span>
    </div>
    <button class="carousel-control prev" @click="prev">
      <span class="control-icon">‹</span>
    </button>
    <button class="carousel-control next" @click="next">
      <span class="control-icon">›</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import type { Media } from '@/types/media'

const props = defineProps<{
  items: Media[]
  interval?: number
  autoplay?: boolean
}>()

defineEmits<{
  (e: 'click', id: number): void
}>()

const carouselRef = ref<HTMLElement | null>(null)
const currentIndex = ref(0)
const timer = ref<number | null>(null)

const intervalValue = props.interval || 5000
const autoplay = props.autoplay !== false

const carouselStyle = computed(() => {
  return {
    transform: `translateX(-${currentIndex.value * 100}%)`,
    transition: 'transform 0.5s ease'
  }
})

const truncateOverview = (overview: string) => {
  return overview.length > 100 ? overview.substring(0, 100) + '...' : overview
}

const goTo = (idx: number) => {
  currentIndex.value = idx
  resetTimer()
}

const next = () => {
  currentIndex.value = (currentIndex.value + 1) % props.items.length
  resetTimer()
}

const prev = () => {
  currentIndex.value = (currentIndex.value - 1 + props.items.length) % props.items.length
  resetTimer()
}

const startAutoplay = () => {
  if (autoplay && props.items.length > 1) {
    timer.value = window.setInterval(next, intervalValue)
  }
}

const stopAutoplay = () => {
  if (timer.value) {
    clearInterval(timer.value)
    timer.value = null
  }
}

const resetTimer = () => {
  stopAutoplay()
  startAutoplay()
}

onMounted(() => {
  startAutoplay()
})

onUnmounted(() => {
  stopAutoplay()
})

const defaultBackdrop = 'https://via.placeholder.com/1200x600?text=No+Backdrop'
</script>

<style scoped lang="scss">
.carousel {
  position: relative;
  width: 100%;
  height: 400px;
  overflow: hidden;
  border-radius: 8px;
  
  .carousel-container {
    display: flex;
    height: 100%;
    
    .carousel-item {
      position: relative;
      flex: 0 0 100%;
      height: 100%;
      cursor: pointer;
      
      .backdrop {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
      
      .overlay {
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        padding: 32px;
        background: linear-gradient(to top, rgba(0, 0, 0, 0.8), transparent);
        color: #fff;
        
        .title {
          font-size: 24px;
          font-weight: bold;
          margin-bottom: 8px;
        }
        
        .overview {
          font-size: 14px;
          line-height: 1.5;
          max-width: 80%;
        }
      }
    }
  }
  
  .carousel-indicators {
    position: absolute;
    bottom: 16px;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    gap: 8px;
    
    .indicator {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background-color: rgba(255, 255, 255, 0.5);
      cursor: pointer;
      transition: all 0.3s ease;
      
      &.active {
        background-color: #fff;
        width: 24px;
        border-radius: 4px;
      }
    }
  }
  
  .carousel-control {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background-color: rgba(0, 0, 0, 0.5);
    color: #fff;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s ease;
    
    &:hover {
      background-color: rgba(0, 0, 0, 0.8);
    }
    
    &.prev {
      left: 16px;
    }
    
    &.next {
      right: 16px;
    }
    
    .control-icon {
      font-size: 24px;
      font-weight: bold;
    }
  }
}

@media (max-width: 768px) {
  .carousel {
    height: 200px;
    
    .carousel-item {
      .overlay {
        padding: 16px;
        
        .title {
          font-size: 18px;
        }
        
        .overview {
          display: none;
        }
      }
    }
  }
}
</style>