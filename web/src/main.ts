import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

import authPlugin from './plugins/authPlugin';
import axiosPlugin from './plugins/axiosPlugin';
import clickOutsidePlugin from './plugins/clickOutsidePlugin';

import { useThemeStore } from './stores/theme';

import '@/assets/base.scss';
import '@mdi/font/css/materialdesignicons.css';

const app = createApp(App)
app.use(createPinia())
app.use(router);
app.use(authPlugin);
app.use(axiosPlugin);
app.use(clickOutsidePlugin);
app.mount('#app')

useThemeStore().init();
