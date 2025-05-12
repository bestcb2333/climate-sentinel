import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: () => import('@/views/DashBoard.vue'),
      meta: {
        label: 'dashboard',
      },
    },
    {
      path: '/resource',
      name: 'resource',
      component: () => import('@/views/ResourceManage.vue'),
      meta: {
        label: 'resource',
      },
    },
    {
      path: '/weather',
      name: 'weather',
      component: () => import('@/views/HistoryWeather.vue'),
      meta: {
        label: 'weather',
      },
    },
    {
      path: '/routes',
      name: 'routes',
      component: () => import('@/views/RescueRoutes.vue'),
      meta: {
        label: 'routes',
      },
    },
    {
      path: '/events',
      name: 'events',
      component: () => import('@/views/DisasterEvents.vue'),
      meta: {
        label: 'events',
      },
    },
    {
      path: '/users',
      name: 'users',
      component: () => import('@/views/UserManage.vue'),
      meta: {
        label: 'staff',
      },
    },
  ],
})

export default router
