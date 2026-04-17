<template>
  <div class="sort-options">
    <label>排序方式</label>
    <el-select v-model="localSort" @change="handleSortChange">
      <el-option label="按标题" value="title"></el-option>
      <el-option label="按年份" value="year"></el-option>
      <el-option label="按评分" value="rating"></el-option>
      <el-option label="按添加时间" value="created_at"></el-option>
    </el-select>
    <el-select v-model="localOrder" @change="handleSortChange">
      <el-option label="升序" value="asc"></el-option>
      <el-option label="降序" value="desc"></el-option>
    </el-select>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  sort: string
  order: string
}>()

const emit = defineEmits<{
  (e: 'sortChange', sort: string, order: string): void
}>()

const localSort = ref(props.sort)
const localOrder = ref(props.order)

const handleSortChange = () => {
  emit('sortChange', localSort.value, localOrder.value)
}

watch(
  () => [props.sort, props.order],
  ([newSort, newOrder]) => {
    localSort.value = newSort
    localOrder.value = newOrder
  }
)
</script>

<style scoped lang="scss">
.sort-options {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 16px;
  background-color: #f5f7fa;
  border-radius: 8px;
  
  label {
    font-size: 14px;
    color: #606266;
    font-weight: 500;
  }
  
  :deep(.el-select) {
    width: 120px;
  }
}

@media (max-width: 768px) {
  .sort-options {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
    
    :deep(.el-select) {
      width: 100%;
    }
  }
}
</style>