import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/home/HomeView.vue'),
    meta: { title: '首页' },
  },
  {
    path: '/library',
    name: 'Library',
    component: () => import('@/views/library/LibraryView.vue'),
    meta: { title: '媒体库' },
  },
  {
    path: '/detail/:id',
    name: 'Detail',
    component: () => import('@/views/detail/DetailView.vue'),
    props: true,
    meta: { title: '媒体详情' },
  },
  {
    path: '/player/:id',
    name: 'Player',
    component: () => import('@/views/player/PlayerView.vue'),
    props: true,
    meta: { title: '播放' },
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/settings/SettingsView.vue'),
    meta: { title: '设置' },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: '页面不存在' },
  },
]

export default routes