<script setup>
import { useConfigStore } from "@/stores/config";
import { useDashboardStore } from "@/stores/dashboard";
</script>

<template>
  <div class="flex flex-row" v-if="$auth.loggedIn">
    <div class="flex-none w-48 h-screen"></div>
    <aside
      class="fixed w-48 shadow bg-slate-700 h-screen dark:bg-slate-700 text-slate-100 dark:text-slate-100 z-40"
    >
      <div class="py-4 text-center mx-auto">
        <a href="#" class="mx-auto">
          <img class="inline" src="@/assets/venti-logo.svg" width="20" height="20" />
          <span class="ml-2 text-lg">venti</span>
        </a>
        <div class="text-xs text-center opacity-25">{{version}}</div>
      </div>
      <div class="py-5">
        <RouterLink
          class="block hover:bg-slate-600 dark:hover:bg-slate-600 px-8 py-2"
          to="/metrics"
        >Metrics</RouterLink>
        <RouterLink
          class="block hover:bg-slate-600 dark:hover:bg-slate-600 px-8 py-2"
          to="/logs"
        >Logs</RouterLink>
        <div
          class="block px-8 text-slate-400 dark:text-slate-400 py-2"
        >Dashboards ({{ dashboards.length }})</div>
        <div v-for="dashboard in dashboards">
          <RouterLink
            class="block hover:bg-slate-600 dark:hover:bg-slate-600 px-12 py-2"
            :to="'/dashboard/' + dashboard.title"
          >{{ dashboard.title }}</RouterLink>
        </div>
        <RouterLink
          class="block hover:bg-slate-600 dark:hover:bg-slate-600 px-8 py-2"
          to="/alert"
        >Alert</RouterLink>
        <RouterLink
          class="block hover:bg-slate-600 dark:hover:bg-slate-600 px-8 py-2"
          to="/datasource"
        >Datasource</RouterLink>
      </div>
      <div class="text-center text-slate-100 dark:text-slate-100 py-1">
        <button
          @click="useConfigStore().setDark(false)"
          class="h-rounded-group py-2 px-4 bg-slate-600 dark:bg-slate-600 border border-slate-700 dark:border-slate-700 hover:bg-slate-800 dark:hover:bg-slate-800"
        >light</button>
        <button
          @click="useConfigStore().setDark(true)"
          class="h-rounded-group py-2 px-4 bg-slate-600 dark:bg-slate-600 border border-slate-700 dark:border-slate-700 hover:bg-slate-800 dark:hover:bg-slate-800"
        >dark</button>
      </div>
      <RouterLink
        class="block hover:bg-slate-600 dark:hover:bg-slate-600 px-8 py-4"
        to="/logout"
      >Logout</RouterLink>
    </aside>
    <main class="main flex flex-col flex-grow transition-all duration-150 ease-in">
      <div class="w-full flex flex-grow">
        <RouterView :key="$route.fullPath" />
      </div>
    </main>
  </div>
</template>

<script>
export default {
  data() {
    return {
      dashboards: [],
      version: '',
    }
  },
  methods: {
    async init() {
      this.dashboards = await useDashboardStore().getDashboards()
      try {
        const response = await this.axios.get('api/config/version')
        this.version = response.data
      } catch (error) {
        console.error(error)
      }
    },
  },
  mounted() {
    this.init()
  },
}
</script>
