<template>
  <div class="settings-view">
    <Header />
    <main class="main-content">
      <div class="container">
        <h1 class="page-title">系统设置</h1>
        
        <div class="settings-container">
          <!-- 左侧导航 -->
          <div class="settings-nav">
            <el-menu
              :default-active="activeTab"
              class="settings-menu"
              @select="handleTabChange"
            >
              <el-menu-item index="library">
                <el-icon><Folder /></el-icon>
                <span>媒体库</span>
              </el-menu-item>
              <el-menu-item index="scraper">
                <el-icon><DataAnalysis /></el-icon>
                <span>刮削器</span>
              </el-menu-item>
              <el-menu-item index="system">
                <el-icon><Monitor /></el-icon>
                <span>系统信息</span>
              </el-menu-item>
            </el-menu>
          </div>
          
          <!-- 右侧内容 -->
          <div class="settings-content">
            <!-- 媒体库设置 -->
            <div v-if="activeTab === 'library'" class="settings-panel">
              <h2 class="panel-title">媒体库配置</h2>
              
              <el-form :model="librarySettings" label-width="120px">
                <el-form-item label="媒体目录">
                  <el-select v-model="librarySettings.mediaPath" multiple placeholder="选择媒体目录">
                    <el-option 
                      v-for="folder in folders" 
                      :key="folder.id" 
                      :label="folder.path" 
                      :value="folder.path"
                    ></el-option>
                  </el-select>
                </el-form-item>
                
                <el-form-item label="自动扫描">
                  <el-switch v-model="librarySettings.autoScan" />
                </el-form-item>
                
                <el-form-item label="扫描间隔">
                  <el-select v-model="librarySettings.scanInterval">
                    <el-option label="每天" value="1"></el-option>
                    <el-option label="每周" value="7"></el-option>
                    <el-option label="每月" value="30"></el-option>
                  </el-select>
                </el-form-item>
                
                <el-form-item>
                  <el-button type="primary" @click="saveLibrarySettings">保存设置</el-button>
                  <el-button @click="triggerScan">立即扫描</el-button>
                </el-form-item>
              </el-form>
            </div>
            
            <!-- 刮削器设置 -->
            <div v-if="activeTab === 'scraper'" class="settings-panel">
              <h2 class="panel-title">刮削器设置</h2>
              
              <el-form :model="scraperSettings" label-width="120px">
                <el-form-item label="TMDB API Key">
                  <el-input v-model="scraperSettings.tmdbApiKey" type="password" />
                </el-form-item>
                
                <el-form-item label="语言">
                  <el-select v-model="scraperSettings.language">
                    <el-option label="中文" value="zh-CN"></el-option>
                    <el-option label="英文" value="en-US"></el-option>
                  </el-select>
                </el-form-item>
                
                <el-form-item label="自动刮削">
                  <el-switch v-model="scraperSettings.autoScrape" />
                </el-form-item>
                
                <el-form-item>
                  <el-button type="primary" @click="saveScraperSettings">保存设置</el-button>
                </el-form-item>
              </el-form>
            </div>
            
            <!-- 系统信息 -->
            <div v-if="activeTab === 'system'" class="settings-panel">
              <h2 class="panel-title">系统信息</h2>
              
              <el-descriptions :column="1" border>
                <el-descriptions-item label="系统版本">
                  {{ systemInfo.version }}
                </el-descriptions-item>
                <el-descriptions-item label="前端版本">
                  {{ systemInfo.frontendVersion }}
                </el-descriptions-item>
                <el-descriptions-item label="后端版本">
                  {{ systemInfo.backendVersion }}
                </el-descriptions-item>
                <el-descriptions-item label="媒体数量">
                  {{ systemInfo.mediaCount }}
                </el-descriptions-item>
                <el-descriptions-item label="存储使用">
                  {{ systemInfo.storageUsed }}
                </el-descriptions-item>
                <el-descriptions-item label="最后扫描时间">
                  {{ systemInfo.lastScanTime }}
                </el-descriptions-item>
              </el-descriptions>
              
              <div class="system-actions">
                <el-button @click="exportSettings">导出设置</el-button>
                <el-button @click="importSettings">导入设置</el-button>
                <el-button type="danger" @click="resetSettings">重置设置</el-button>
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
import { ref, onMounted } from 'vue'
import { Folder, DataAnalysis, Monitor } from '@element-plus/icons-vue'
import { mediaAPI } from '@/services/media'
import Header from '@/components/layout/Header.vue'
import Footer from '@/components/layout/Footer.vue'
import type { Folder as FolderType } from '@/types/media'

const activeTab = ref('library')
const folders = ref<FolderType[]>([])

// 模拟数据
const librarySettings = ref({
  mediaPath: ['/path/to/media'],
  autoScan: true,
  scanInterval: '7'
})

const scraperSettings = ref({
  tmdbApiKey: 'your_api_key_here',
  language: 'zh-CN',
  autoScrape: true
})

const systemInfo = ref({
  version: '1.0.0',
  frontendVersion: '1.0.0',
  backendVersion: '1.0.0',
  mediaCount: 100,
  storageUsed: '100GB / 500GB',
  lastScanTime: '2026-04-17 14:00:00'
})

const handleTabChange = (tab: string) => {
  activeTab.value = tab
}

const saveLibrarySettings = async () => {
  // 实际项目中，这里会保存设置到API
  console.log('保存媒体库设置:', librarySettings.value)
  // await settingsStore.updateSettings([
  //   { key: 'mediaPath', value: librarySettings.value.mediaPath.join(',') },
  //   { key: 'autoScan', value: librarySettings.value.autoScan.toString() },
  //   { key: 'scanInterval', value: librarySettings.value.scanInterval }
  // ])
  
  // 显示成功消息
  ElMessage.success('媒体库设置保存成功')
}

const saveScraperSettings = async () => {
  // 实际项目中，这里会保存设置到API
  console.log('保存刮削器设置:', scraperSettings.value)
  // await settingsStore.updateSettings([
  //   { key: 'tmdbApiKey', value: scraperSettings.value.tmdbApiKey },
  //   { key: 'language', value: scraperSettings.value.language },
  //   { key: 'autoScrape', value: scraperSettings.value.autoScrape.toString() }
  // ])
  
  // 显示成功消息
  ElMessage.success('刮削器设置保存成功')
}

const triggerScan = async () => {
  try {
    await mediaAPI.triggerScan()
    ElMessage.success('扫描任务已启动')
  } catch (error) {
    ElMessage.error('扫描任务启动失败')
  }
}

const exportSettings = () => {
  // 实际项目中，这里会导出设置
  console.log('导出设置')
  ElMessage.success('设置导出成功')
}

const importSettings = () => {
  // 实际项目中，这里会导入设置
  console.log('导入设置')
  ElMessage.success('设置导入成功')
}

const resetSettings = () => {
  // 实际项目中，这里会重置设置
  console.log('重置设置')
  ElMessage.success('设置已重置')
}

onMounted(async () => {
  // 实际项目中，这里会从API获取设置和文件夹列表
  // await settingsStore.fetchSettings()
  // const foldersResponse = await mediaAPI.getFolders()
  // folders.value = foldersResponse.data.data
  
  // 模拟数据
  folders.value = [
    { id: 1, name: 'Movies', path: '/path/to/movies' },
    { id: 2, name: 'TV Shows', path: '/path/to/tvshows' },
    { id: 3, name: 'Anime', path: '/path/to/anime' }
  ]
})

// 导入ElMessage
import { ElMessage } from 'element-plus'
</script>

<style scoped lang="scss">
.settings-view {
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
    
    .settings-container {
      display: flex;
      gap: 24px;
      
      .settings-nav {
        flex: 0 0 200px;
        
        .settings-menu {
          background-color: #fff;
          border-radius: 8px;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }
      }
      
      .settings-content {
        flex: 1;
        
        .settings-panel {
          background-color: #fff;
          border-radius: 8px;
          padding: 24px;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
          
          .panel-title {
            font-size: 18px;
            font-weight: bold;
            margin-bottom: 20px;
            color: #303133;
          }
        }
        
        .system-actions {
          margin-top: 24px;
          display: flex;
          gap: 12px;
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .settings-view {
    .main-content {
      .settings-container {
        flex-direction: column;
        
        .settings-nav {
          flex: none;
          
          .settings-menu {
            display: flex;
            overflow-x: auto;
            
            .el-menu-item {
              white-space: nowrap;
            }
          }
        }
      }
    }
  }
}
</style>