import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: () => import('@/views/DashBoard.vue'),
      meta: {
        label: 'navbar.dashboard',
      },
    },
    {
      path: '/resource',
      name: 'resource',
      component: () => import('@/views/ResourceManage.vue'),
      meta: {
        label: 'navbar.resource',
      },
    },
    {
      path: '/weather',
      name: 'weather',
      component: () => import('@/views/HistoryWeather.vue'),
      meta: {
        label: 'navbar.weather',
      },
    },
    {
      path: '/routes',
      name: 'routes',
      component: () => import('@/views/RescueRoutes.vue'),
      meta: {
        label: 'navbar.routes',
      },
    },
    {
      path: '/events',
      name: 'events',
      component: () => import('@/views/DisasterEvents.vue'),
      meta: {
        label: 'navbar.events',
      },
    },
    {
      path: '/staff',
      name: 'staff',
      component: () => import('@/views/StaffManage.vue'),
      meta: {
        label: 'navbar.staff',
      },
    },
  ],
})

export default router
