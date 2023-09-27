<script setup lang="ts">
import { ref, onMounted } from 'vue';
import Util from '@/lib/util';
import LetterAvatar from '@/components/LetterAvatar.vue';
import { useDatasourceStore } from '@/stores/datasource';
import type { Datasource } from '@/types/datasource';

const datasources = ref([] as Datasource[])

async function fetchData() {
  datasources.value = await useDatasourceStore().getDatasources();
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <header class="fixed right-0 w-full bg-white border-b shadow z-30 p-2 pl-52">
    <div class="flex items-center flex-row">
      <div><i class="mdi mdi-18px mdi-database-outline"></i> Datasource</div>
      <div class="flex ml-auto">
        <div class="inline-flex">
          <button @click="fetchData()"
            class="h-rounded-group py-2 px-4 text-gray-900 bg-white border border-common hover:bg-gray-100 hover:text-blue-500 focus:text-blue-500">
            <i class="mdi mdi-refresh"></i>
          </button>
        </div>
      </div>
    </div>
  </header>

  <div class="mt-12 w-full p-8">
    <h1 class="mt-4">Datasources</h1>
    <table class="w-full bg-white border" v-if="datasources">
      <tr class="border-b bg-slate-50">
        <th>Name</th>
        <th>Type</th>
        <th>URL</th>
        <th>Main</th>
        <th>Discovered</th>
        <th>Up</th>
      </tr>
      <tr class="border-b" v-for="d in datasources">
        <td class="px-2">
          <letterAvatar :bgcolor="Util.string2color(d.name)" />
          {{ d.name }}
        </td>
        <td>{{ d.type == 'prometheus' ? '🔥' : '💧' }} {{ d.type }}</td>
        <td>{{ d.url }}</td>
        <td class="text-center">{{ d.isMain ? '✔️' : '-' }}</td>
        <td class="text-center">{{ d.isDiscovered ? '✔️' : '-' }}</td>
        <td class="text-center">
          <span :class="[d.health ? 'text-green-400' : 'text-red-400']">●</span>
        </td>
      </tr>
    </table>
    <h1 class="mt-4">Targets</h1>
    <table class="w-full bg-white border">
      <tr class="border-b bg-slate-50">
        <th>Datasource</th>
        <th>Job</th>
        <th>Address</th>
        <th>Name</th>
        <th>Last scrape</th>
        <th>Up</th>
      </tr>
      <template v-for="d of datasources">
        <tr class="border-b" v-for="t in d.targets">
          <td class="px-2" v-if="d.name">
            <LetterAvatar :bgcolor="Util.string2color(d.name)" />
            {{ d.name }}
          </td>
          <td>{{ t.discoveredLabels.job }}</td>
          <td>{{ t.discoveredLabels.__address__ }}</td>
          <td>{{ t.icon }} {{ t.name }}</td>
          <td class="text-right pr-10">{{ t.age }}s</td>
          <td class="text-center">
            <span :class="[t.health == 'up' ? 'text-green-400' : 'text-red-400']">●</span>
          </td>
        </tr>
      </template>
    </table>
  </div>
</template>
