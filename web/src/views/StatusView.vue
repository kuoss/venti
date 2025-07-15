<script setup lang="ts">
import { ref } from 'vue';

export interface RuntimeInfo {
  startTime: string;
  CWD: string;
  reloadConfigSuccess: boolean;
  lastConfigTime: string;
  goroutineCount: number;
  GOMAXPROCS: number;
  GOMEMLIMIT: number;
  GOGC: string;
  GODEBUG: string;
}

export interface BuildInfo {
  version: string;
  revision: string;
  branch: string;
  buildUser: string;
  buildDate: string;
  goVersion: string;
}

interface Alertmanager {
  url: string;
}

export interface Alertmanagers {
  activeAlertmanagers: Alertmanager[];
  droppedAlertmanagers: Alertmanager[];
}

const runtimeInfo = ref({} as RuntimeInfo);
const buildInfo = ref({} as BuildInfo);
const alertmanagers = ref({} as Alertmanagers);

async function fetchRuntimeInfo() {
  const response = await fetch('/api/v1/status/runtimeinfo');
  const respJSON = await response.json();
  if (respJSON.status != 'success') {
    console.error('fetchRuntimeInfo', 'unsuccessful', respJSON);
    return;
  }
  runtimeInfo.value = respJSON.data;
}

async function fetchBuildInfo() {
  const response = await fetch('/api/v1/status/buildinfo');
  const respJSON = await response.json();
  if (respJSON.status != 'success') {
    console.error('fetchBuildInfo', 'unsuccessful', respJSON);
    return;
  }
  buildInfo.value = respJSON.data;
}

async function fetchAlertmanagers() {
  const response = await fetch('/api/v1/alertmanagers');
  const respJSON = await response.json();
  if (respJSON.status != 'success') {
    console.error('fetchAlertmanagers', 'unsuccessful', respJSON);
    return;
  }
  alertmanagers.value = respJSON.data;
}

fetchRuntimeInfo();
fetchBuildInfo();
fetchAlertmanagers();
</script>

<template>
  <header class="fixed right-0 w-full bg-white dark:bg-black border-b shadow z-30 p-2 pl-52">
    <div><i class="mdi mdi-18px mdi-database-outline"></i> Status</div>
  </header>
  <div class="w-full py-8">
    <div class="w-full p-8">
      <h2 class="text-lg font-bold">Runtime Information</h2>
      <table class="w-full">
        <tr>
          <th>Start time</th>
          <td>{{ runtimeInfo.startTime }}</td>
        </tr>
        <tr>
          <th>Working directory</th>
          <td>{{ runtimeInfo.CWD }}</td>
        </tr>
        <tr>
          <th>Configuration reload</th>
          <td>{{ runtimeInfo.reloadConfigSuccess ? 'Successful' : 'Unsuccessful' }}</td>
        </tr>
        <tr>
          <th>Last successful configuration reload</th>
          <td>{{ runtimeInfo.lastConfigTime }}</td>
        </tr>
        <tr>
          <th>Goroutines</th>
          <td>{{ runtimeInfo.goroutineCount }}</td>
        </tr>
        <tr>
          <th>GOMAXPROCS</th>
          <td>{{ runtimeInfo.GOMAXPROCS }}</td>
        </tr>
        <tr>
          <th>GOMEMLIMIT</th>
          <td>{{ runtimeInfo.GOMEMLIMIT }}</td>
        </tr>
        <tr>
          <th>GOGC</th>
          <td>{{ runtimeInfo.GOGC }}</td>
        </tr>
        <tr>
          <th>GODEBUG</th>
          <td>{{ runtimeInfo.GODEBUG }}</td>
        </tr>
      </table>
    </div>
    <div class="w-full p-8">
      <h2 class="text-lg font-bold">Build Information</h2>
      <table class="w-full">
        <tr>
          <th>Version</th>
          <td>{{ buildInfo.version }}</td>
        </tr>
        <tr>
          <th>Revision</th>
          <td>{{ buildInfo.revision }}</td>
        </tr>
        <tr>
          <th>Branch</th>
          <td>{{ buildInfo.branch }}</td>
        </tr>
        <tr>
          <th>BuildUser</th>
          <td>{{ buildInfo.buildUser }}</td>
        </tr>
        <tr>
          <th>BuildDate</th>
          <td>{{ buildInfo.buildDate }}</td>
        </tr>
        <tr>
          <th>GoVersion</th>
          <td>{{ buildInfo.goVersion }}</td>
        </tr>
      </table>
    </div>
    <div class="w-full p-8">
      <h2 class="text-lg font-bold">Alertmanagers</h2>
      <table class="w-full">
        <tr>
          <th>Endpoints</th>
        </tr>
        <tr v-for="el in alertmanagers.activeAlertmanagers">
          <td>
            <a class="text-blue-500" :href="`${el.url}/api/v2/alerts`">{{ el.url }}</a
            >/api/v2/alerts
          </td>
        </tr>
      </table>
    </div>
  </div>
</template>

<style scoped>
@import "tailwindcss";

th {
  width: 30%;
}

td {
  width: 70%;
}

th,
td {
  @apply p-2 border border-slate-300 dark:border-slate-700;
}
</style>
