<template>
  <footer class="footer">
    <div class="container">
      <div class="footer-content">
        <div class="footer-info">
          <h3>VideoHub</h3>
          <p>本地视频管理系统</p>
        </div>
        <div class="footer-links">
          <a href="#" class="footer-link">关于我们</a>
          <a href="#" class="footer-link">使用帮助</a>
          <a href="#" class="footer-link">隐私政策</a>
        </div>
        <div class="footer-copyright">
          <p>&copy; {{ currentYear }} VideoHub. 保留所有权利。</p>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const currentYear = computed(() => new Date().getFullYear())
</script>

<style scoped lang="scss">
.footer {
  background-color: #333;
  color: #fff;
  padding: 32px 0;
  margin-top: 40px;
  
  .container {
    width: 100%;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
  }
  
  .footer-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
    
    .footer-info {
      text-align: center;
      
      h3 {
        font-size: 20px;
        font-weight: bold;
        margin-bottom: 8px;
      }
      
      p {
        font-size: 14px;
        color: #909399;
        margin: 0;
      }
    }
    
    .footer-links {
      display: flex;
      gap: 24px;
      
      .footer-link {
        color: #909399;
        text-decoration: none;
        font-size: 14px;
        transition: color 0.3s ease;
        
        &:hover {
          color: #409EFF;
        }
      }
    }
    
    .footer-copyright {
      font-size: 12px;
      color: #909399;
      
      p {
        margin: 0;
      }
    }
  }
}

@media (max-width: 768px) {
  .footer {
    padding: 20px 0;
    
    .footer-content {
      gap: 16px;
      
      .footer-links {
        gap: 16px;
      }
    }
  }
}
</style>