<template>
  <header class="header">
    <div class="container">
      <div class="header-content">
        <div class="logo" @click="navigateTo('/')">
          <h1>VideoHub</h1>
        </div>
        <nav class="nav" :class="{ 'collapsed': isMobile && !menuOpen }">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/library" class="nav-link">媒体库</router-link>
          <router-link to="/settings" class="nav-link">设置</router-link>
        </nav>
        <div class="header-actions">
          <button class="nav-toggle" @click="toggleMenu" v-if="isMobile">
            <span class="toggle-icon"></span>
            <span class="toggle-icon"></span>
            <span class="toggle-icon"></span>
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUIStore } from '@/store'

const router = useRouter()
const uiStore = useUIStore()

const menuOpen = ref(false)
const isMobile = computed(() => uiStore.isMobile)

const toggleMenu = () => {
  menuOpen.value = !menuOpen.value
}

const navigateTo = (path: string) => {
  router.push(path)
  if (isMobile.value) {
    menuOpen.value = false
  }
}

const checkMobile = () => {
  uiStore.setIsMobile(window.innerWidth < 768)
  if (window.innerWidth >= 768) {
    menuOpen.value = false
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped lang="scss">
.header {
  background-color: #333;
  color: #fff;
  padding: 16px 0;
  position: sticky;
  top: 0;
  z-index: 100;
  
  .container {
    width: 100%;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
  }
  
  .header-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    
    .logo {
      cursor: pointer;
      
      h1 {
        font-size: 24px;
        font-weight: bold;
        margin: 0;
      }
    }
    
    .nav {
      display: flex;
      gap: 24px;
      
      .nav-link {
        color: #fff;
        text-decoration: none;
        font-size: 16px;
        font-weight: 500;
        transition: color 0.3s ease;
        
        &:hover {
          color: #409EFF;
        }
        
        &.router-link-active {
          color: #409EFF;
        }
      }
      
      &.collapsed {
        position: fixed;
        top: 64px;
        left: 0;
        right: 0;
        background-color: #333;
        flex-direction: column;
        padding: 20px;
        gap: 16px;
        transform: translateY(-100%);
        transition: transform 0.3s ease;
        
        &.open {
          transform: translateY(0);
        }
      }
    }
    
    .header-actions {
      display: flex;
      gap: 12px;
      
      .nav-toggle {
        display: none;
        flex-direction: column;
        gap: 4px;
        background: none;
        border: none;
        cursor: pointer;
        padding: 8px;
        
        .toggle-icon {
          width: 24px;
          height: 2px;
          background-color: #fff;
          transition: all 0.3s ease;
        }
        
        &.open .toggle-icon:nth-child(1) {
          transform: rotate(45deg) translate(5px, 5px);
        }
        
        &.open .toggle-icon:nth-child(2) {
          opacity: 0;
        }
        
        &.open .toggle-icon:nth-child(3) {
          transform: rotate(-45deg) translate(5px, -5px);
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .header {
    .header-content {
      .nav {
        display: none;
        
        &.collapsed {
          display: flex;
        }
      }
      
      .header-actions {
        .nav-toggle {
          display: flex;
        }
      }
    }
  }
}
</style>