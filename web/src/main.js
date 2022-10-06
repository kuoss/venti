import { createPinia } from 'pinia'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import VueAxios from 'vue-axios'
import vClickOutside from "click-outside-vue3"

import auth from './plugins/auth'
import axios from './plugins/axios'
import util from './plugins/util'

import sidebarLayout from './layouts/sidebar.vue'
import nosidebarLayout from './layouts/nosidebar.vue'

try {
    if (localStorage.theme === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
        document.documentElement.classList.add('dark')
        document.documentElement.style.setProperty('color-scheme', 'dark');
    } else {
        document.documentElement.classList.remove('dark')
        document.documentElement.style.setProperty('color-scheme', 'normal');
    }
} catch (_) { }

const app = createApp(App)
app.use(createPinia())
app.use(auth)
app.use(VueAxios, axios)
app.use(util)
app.use(router)
app.use(vClickOutside)
app.component('sidebar', sidebarLayout)
app.component('nosidebar', nosidebarLayout)
app.mount('#app')
