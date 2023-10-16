import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { useTimeStore } from '@/stores/time';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: () => import('@/components/Layout.vue'),
      children: [
        {
          path: '',
          name: 'home',
          component: () => import('@/views/HomeView.vue'),
          meta: { requireAuth: true },
        },
        {
          path: 'metrics',
          name: 'metrics',
          component: () => import('@/views/MetricsView.vue'),
          meta: { requireAuth: true },
        },
        {
          path: 'logs',
          name: 'logs',
          component: () => import('@/views/LogsView.vue'),
          meta: { requireAuth: true },
        },
        {
          path: 'dashboard/:name',
          name: 'dashboard',
          component: () => import('@/views/DashboardView.vue'),
          meta: { requireAuth: true },
        },
        {
          path: 'alert',
          name: 'alert',
          component: () => import('@/views/AlertView.vue'),
          meta: { requireAuth: true },
        },
        {
          path: 'datasource',
          name: 'datasource',
          component: () => import('@/views/DatasourceView.vue'),
          meta: { requireAuth: true },
        },
        {
          path: 'status',
          name: 'status',
          component: () => import('@/views/StatusView.vue'),
          meta: { requireAuth: true },
        },
      ],
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      beforeEnter: (to, from, next) => {
        if (useAuthStore().loggedIn) {
          next({ name: 'home' });
          return;
        }
        next();
      },
    },
    {
      path: '/logout',
      name: 'logout',
      component: () => import('@/views/LogoutView.vue'),
    },
  ],
});

router.beforeEach((to, _, next) => {
  useTimeStore().timerManager = '';
  if (to.meta.requireAuth && !useAuthStore().loggedIn) {
    next({ name: 'login' });
    return;
  }
  next();
});

export default router;