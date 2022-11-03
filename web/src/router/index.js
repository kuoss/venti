// import { createRouter, createWebHistory } from 'vue-router'
import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from "@/stores/auth"
import { useTimeStore } from "@/stores/time"

const router = createRouter({
  // history: createWebHistory(import.meta.env.BASE_URL),
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
      meta: { requireAuth: true, layout: 'sidebar' },
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      beforeEnter: (to, from, next) => {
        if (useAuthStore().loggedIn) next({ name: "home" })
        else next()
      },
    },
    {
      path: '/logout',
      name: 'logout',
      component: () => import('@/views/LogoutView.vue'),
      meta: { requireAuth: true, layout: 'nosidebar' },
    },
    {
      path: '/metrics',
      name: 'metrics',
      component: () => import('@/views/MetricsView.vue'),
      meta: { requireAuth: true, layout: 'sidebar' },
    },
    {
      path: '/logs',
      name: 'logs',
      component: () => import('@/views/LogsView.vue'),
      meta: { requireAuth: true, layout: 'sidebar' },
    },
    {
      path: '/dashboard/:name',
      name: 'dashboard',
      component: () => import('@/views/DashboardView.vue'),
      meta: { requireAuth: true, layout: 'sidebar' },
    },
    {
      path: '/alert',
      name: 'alert',
      component: () => import('@/views/AlertView.vue'),
      meta: { requireAuth: true, layout: 'sidebar' },
    },
    {
      path: '/datasource',
      name: 'datasource',
      component: () => import('@/views/DatasourceView.vue'),
      meta: { requireAuth: true, layout: 'sidebar' },
    },
    {
      path: '/config',
      name: 'config',
      component: () => import('@/views/ConfigView.vue'),
      meta: { requireAuth: true, layout: 'sidebar' },
    },
  ],
})

router.beforeEach((to, from, next) => {
  useTimeStore().timerManager = ''
  if (to.meta.requireAuth && !useAuthStore().loggedIn) next({ name: "login" })
  else next()
})

export default router
