<template>
  <div class="not-found">
    <Header />
    <main class="main-content">
      <div class="container">
        <div class="error-container">
          <h1 class="error-code">404</h1>
          <h2 class="error-message">页面不存在</h2>
          <p class="error-description">您访问的页面不存在或已被删除</p>
          <el-button type="primary" @click="goHome">返回首页</el-button>
        </div>
      </div>
    </main>
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import Header from '@/components/layout/Header.vue'
import Footer from '@/components/layout/Footer.vue'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<style scoped lang="scss">
.not-found {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  
  .main-content {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    
    .container {
      width: 100%;
      max-width: 1200px;
      margin: 0 auto;
      padding: 0 20px;
    }
    
    .error-container {
      text-align: center;
      padding: 40px;
      background-color: #fff;
      border-radius: 8px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      
      .error-code {
        font-size: 120px;
        font-weight: bold;
        color: #409EFF;
        margin-bottom: 16px;
      }
      
      .error-message {
        font-size: 24px;
        font-weight: bold;
        color: #303133;
        margin-bottom: 8px;
      }
      
      .error-description {
        font-size: 16px;
        color: #606266;
        margin-bottom: 24px;
      }
    }
  }
}
</style>