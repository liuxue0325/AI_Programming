<template>
  <div class="media-card" @click="$emit('click', media.id)">
    <div class="poster-container">
      <img 
        :src="media.poster_path || defaultPoster" 
        :alt="media.title" 
        class="poster"
      />
      <div v-if="showRating && media.rating" class="rating">
        {{ media.rating.toFixed(1) }}
      </div>
    </div>
    <div class="info">
      <h3 class="title">{{ media.title }}</h3>
      <p class="year" v-if="media.year">{{ media.year }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'
import type { Media } from '@/types/media'

defineProps<{
  media: Media
  type?: string
  showRating?: boolean
}>()

defineEmits<{
  (e: 'click', id: number): void
}>()

const defaultPoster = 'https://via.placeholder.com/200x300?text=No+Poster'
</script>

<style scoped lang="scss">
.media-card {
  background-color: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
  
  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 20px rgba(0, 0, 0, 0.15);
  }
  
  .poster-container {
    position: relative;
    width: 100%;
    padding-top: 150%; // 2:3 aspect ratio
    overflow: hidden;
    
    .poster {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform 0.3s ease;
    }
    
    &:hover .poster {
      transform: scale(1.05);
    }
    
    .rating {
      position: absolute;
      top: 8px;
      right: 8px;
      background-color: rgba(0, 0, 0, 0.8);
      color: #fff;
      padding: 4px 8px;
      border-radius: 4px;
      font-size: 12px;
      font-weight: bold;
    }
  }
  
  .info {
    padding: 12px;
    
    .title {
      font-size: 14px;
      font-weight: 500;
      margin-bottom: 4px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    
    .year {
      font-size: 12px;
      color: #909399;
    }
  }
}
</style>