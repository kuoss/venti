import type { App } from 'vue';
import type { AuthStore } from '@/stores/auth';
import { useAuthStore } from '@/stores/auth';

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $auth: AuthStore,
  }
}

export default {
  install: (app: App) => {
    app.config.globalProperties.$auth = useAuthStore();
  },
};