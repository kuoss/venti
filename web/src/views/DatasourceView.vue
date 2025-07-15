<script setup lang="ts">
import { ref } from 'vue';
import Util from '@/lib/util';
import { formatTimeAgo } from '@vueuse/core';
import LetterAvatar from '@/components/LetterAvatar.vue';
import { useDatasourceStore } from '@/stores/datasource';
import type { Datasource, Target } from '@/types/datasource';

const dsStore = useDatasourceStore();

const datasources = ref([] as Datasource[]);
const datasource = ref({} as Datasource);
const targets = ref([] as Target[]);
const healthClasses = ref(['text-gray-500', 'text-green-500', 'text-red-500']);

async function fetchData() {
  const dss = await dsStore.getDatasources();
  datasources.value = dss;
  for (const i in dss) {
    // @ts-ignore
    dss[i].health = await dsStore.getDatasourceHealthy(dss[i]);
  }
}

async function showTargets(ds: Datasource) {
  datasource.value = ds;
  // @ts-ignore
  targets.value = await dsStore.getTargets(ds);
}
fetchData();
</script>

<template>
  <header class="fixed right-0 w-full border-b shadow z-30 p-2 pl-52 bg-white dark:bg-black">
    <div class="flex items-center flex-row">
      <div><i class="mdi mdi-18px mdi-database-outline"></i> Datasource</div>
      <div class="flex ml-auto">
        <div class="inline-flex">
          <button @click="fetchData()"
            class="h-rounded-group py-2 px-4 border border-common text-gray-900 dark:text-gray-100 bg-white dark:bg-black hover:bg-gray-100 dark:hover:bg-gray-900 hover:text-blue-500 focus:text-blue-500">
            <i class="mdi mdi-refresh"></i>
          </button>
        </div>
      </div>
    </div>
  </header>

  <div class="mt-12 w-full">
    <div class="p-8">
      <h2 class="text-lg font-bold">Datasources</h2>
      <table class="w-full border bg-white dark:bg-black" v-if="datasources">
        <tr class="border-b bg-slate-50 dark:bg-slate-900">
          <th>Name</th>
          <th>Type</th>
          <th>URL</th>
          <th>Main</th>
          <th>Discovered</th>
          <th>Up</th>
          <th>Actions</th>
        </tr>
        <tr class="border-b" v-for="d in datasources" :class="{ 'bg-blue-50 dark:bg-blue-900': d.name == datasource.name }">
          <td class="px-2">
            <LetterAvatar :letters="d.name.charAt(0)" :bgcolor="Util.string2color(d.name)" />
            {{ d.name }}
          </td>
          <td>{{ d.type == 'prometheus' ? '🔥' : '💧' }} {{ d.type }}</td>
          <td>{{ d.url }}</td>
          <td class="text-center">{{ d.isMain ? '✔️' : '-' }}</td>
          <td class="text-center">{{ d.isDiscovered ? '✔️' : '-' }}</td>
          <td class="text-center">
            <span :class="d.health ? healthClasses[d.health] : 'text-gray-500'">●</span>
          </td>
          <td class="text-center">
            <button class="btn" @click="showTargets(d)">Show Targets</button>
          </td>
        </tr>
      </table>
    </div>
    <div class="p-8" v-if="Object.keys(datasource).length">
      <div class="overflow-auto">
        <span class="float-left text-lg font-bold">Targets</span>
        <span class="float-left p-1 px-2">
          <LetterAvatar :letters="datasource.name.charAt(0)" :bgcolor="Util.string2color(datasource.name)" />
          {{ datasource.name }}
        </span>
      </div>
      <table class="w-full bg-white dark:bg-black border">
        <tr class="border-b bg-slate-50 dark:bg-slate-900">
          <th>Job</th>
          <th>Address</th>
          <th>Info</th>
          <th>Last scrape</th>
          <th>Up</th>
        </tr>
        <tr class="border-b" v-for="t in targets">
          <td class="px-2">{{ t.discoveredLabels.job }}</td>
          <td>{{ t.discoveredLabels.__address__ }}</td>
          <td v-if="t.discoveredLabels.__meta_kubernetes_namespace">
            🖼️ {{ t.discoveredLabels.__meta_kubernetes_namespace }}
            <span v-if="t.discoveredLabels.__meta_kubernetes_service_name">
              / 🐕‍🦺 {{ t.discoveredLabels.__meta_kubernetes_service_name }}
            </span>
            <span v-if="t.discoveredLabels.__meta_kubernetes_pod_name">
              / 🐱 {{ t.discoveredLabels.__meta_kubernetes_pod_name }}
            </span>
          </td>
          <td v-else-if="t.discoveredLabels.__meta_kubernetes_node_name">
            🏝 {{ t.discoveredLabels.__meta_kubernetes_node_name }}
          </td>
          <td v-else>🔥 prometheus</td>
          <td class="text-right pr-10">{{ formatTimeAgo(new Date(t.lastScrape)) }}</td>
          <td class="text-center px-2">
            <span :class="[t.health == 'up' ? 'text-green-500' : 'text-red-500']">●</span>
          </td>
        </tr>
      </table>
    </div>
  </div>
</template>
