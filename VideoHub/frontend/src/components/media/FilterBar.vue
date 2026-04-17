<template>
  <div class="filter-bar">
    <div class="filter-item">
      <label>媒体类型</label>
      <el-select v-model="localFilters.type" placeholder="全部" @change="handleFilterChange">
        <el-option label="全部" value=""></el-option>
        <el-option label="电影" value="movie"></el-option>
        <el-option label="电视剧" value="tv"></el-option>
        <el-option label="动画" value="anime"></el-option>
      </el-select>
    </div>
    <div class="filter-item">
      <label>年份</label>
      <el-select v-model="localFilters.year" placeholder="全部" @change="handleFilterChange">
        <el-option label="全部" value=""></el-option>
        <el-option 
          v-for="year in years" 
          :key="year" 
          :label="year.toString()" 
          :value="year.toString()"
        ></el-option>
      </el-select>
    </div>
    <div class="filter-item">
      <label>搜索</label>
      <el-input 
        v-model="localFilters.search" 
        placeholder="搜索媒体" 
        clearable
        @input="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Search } from '@element-plus/icons-vue'

const props = defineProps<{
  filters: {
    type?: string
    year?: string
    search?: string
  }
}>()

const emit = defineEmits<{
  (e: 'filterChange', filters: typeof props.filters): void
}>()

const localFilters = ref({ ...props.filters })

const years = computed(() => {
  const currentYear = new Date().getFullYear()
  const yearList = []
  for (let i = currentYear; i >= 1990; i--) {
    yearList.push(i)
  }
  return yearList
})

const handleFilterChange = () => {
  emit('filterChange', localFilters.value)
}

const handleSearch = () => {
  emit('filterChange', localFilters.value)
}

watch(
  () => props.filters,
  (newFilters) => {
    localFilters.value = { ...newFilters }
  },
  { deep: true }
)
</script>

<style scoped lang="scss">
.filter-bar {
  display: flex;
  gap: 16px;
  padding: 16px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
  
  .filter-item {
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-width: 120px;
    
    label {
      font-size: 12px;
      color: #606266;
      font-weight: 500;
    }
    
    :deep(.el-select),
    :deep(.el-input) {
      width: 100%;
    }
  }
  
  .filter-item:nth-child(3) {
    flex: 1;
  }
}

@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
    
    .filter-item {
      min-width: auto;
    }
  }
}
</style>