<script setup lang="ts">
import { ref } from 'vue';
import { useThemeStore } from '@/stores/theme';
import { useDashboardStore } from '@/stores/dashboard';
import type { Dashboard } from '@/types/dashboard';

const dashboards = ref([] as Dashboard[]);
const version = ref('');

async function fetchData() {
  dashboards.value = await useDashboardStore().getDashboards();
  try {
    const resp = await fetch('/api/v1/status/buildinfo');
    const json = await resp.json();
    version.value = json.data.version;
  } catch (error) {
    console.error(error);
  }
}

fetchData();
</script>

<template>
  <div class="flex flex-row" v-if="$auth.loggedIn">
    <div class="flex-none w-48 h-screen"></div>
    <aside class="fixed w-48 shadow bg-slate-700 h-screen dark:bg-slate-700 text-slate-100 dark:text-slate-100 z-40">
      <div class="py-4 text-center mx-auto">
        <a href="#" class="mx-auto">
          <img class="inline" src="@/assets/venti-logo.svg" width="20" height="20" />
          <span class="ml-2 text-lg">venti</span>
        </a>
        <div class="text-xs text-center text-gray-300 dark:text-gray-300">{{ version }}</div>
      </div>
      <div class="py-5">
        <RouterLink class="block hover:bg-slate-600 dark:hover:bg-slate-600 px-8 py-2" to="/metrics">Metrics</RouterLink>
        <RouterLink class="block hover:bg-slate-600 dark:hover:bg-slate-600 px-8 py-2" to="/logs">Logs</RouterLink>

        <div class="block px-8 text-slate-400 dark:text-slate-400 py-2">Dashboards ({{ dashboards.length }})</div>
        <div v-for="dashboard in dashboards">
          <RouterLink class="block hover:bg-slate-600 dark:hover:bg-slate-600 px-12 py-2"
            :to="'/dashboard/' + dashboard.title">
            {{ dashboard.title }}
          </RouterLink>
        </div>

        <RouterLink class="block hover:bg-slate-600 dark:hover:bg-slate-600 px-8 py-2" to="/alert"> Alert </RouterLink>
        <RouterLink class="block hover:bg-slate-600 dark:hover:bg-slate-600 px-8 py-2" to="/datasource">
          Datasource
        </RouterLink>
        <RouterLink class="block hover:bg-slate-600 dark:hover:bg-slate-600 px-8 py-2" to="/status">
          Status
        </RouterLink>
      </div>
      <div class="text-center text-slate-100 dark:text-slate-100 py-1">
        <button @click="useThemeStore().setDark(false)"
          class="rounded py-2 px-4 bg-slate-600 dark:bg-slate-600 border border-slate-700 dark:border-slate-700 hover:bg-slate-800 dark:hover:bg-slate-800">
          light
        </button>
        <button @click="useThemeStore().setDark(true)"
          class="rounded py-2 px-4 bg-slate-600 dark:bg-slate-600 border border-slate-700 dark:border-slate-700 hover:bg-slate-800 dark:hover:bg-slate-800">
          dark
        </button>
      </div>
      <RouterLink class="block hover:bg-slate-600 dark:hover:bg-slate-600 px-8 py-4" to="/logout">Logout</RouterLink>
    </aside>
    <main class="main flex flex-col flex-grow transition-all duration-150 ease-in">
      <div class="w-full flex flex-grow">
        <!-- @vue-ignore -->
        <RouterView :key="$route.fullPath" />
      </div>
    </main>
  </div>
</template>
