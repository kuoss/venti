<script setup>
import { computed, ref, onMounted } from 'vue';

import PanelLogs from '@/components/PanelLogs.vue';
import PanelMultitable from '@/components/PanelMultitable.vue';
import PanelPiechart from '@/components/PanelPiechart.vue';
import PanelStat from '@/components/PanelStat.vue';
import PanelTable from '@/components/PanelTable.vue';
import PanelTimeSeries from '@/components/PanelTimeSeries.vue';

import { useSidePanelStore } from '@/stores/sidePanel';

const sidePanel = useSidePanelStore()

const props = defineProps({
  position: String,
  count: Number,
  panelConfig: Object,
  panelWidth: Number,
  timeRange: Object,
})

const isLoading = ref(false)
const showPanelPosition = ref(false)
const usingVariables = ref([])

const componentName = computed(() => {
  switch (props.panelConfig.type) {
    case 'logs':
      return PanelLogs;
    case 'multitable':
      return PanelMultitable;
    case 'piechart':
      return PanelPiechart;
    case 'stat':
      return PanelStat;
    case 'table':
      return PanelTable;
  }
  // case 'time_series'
  return PanelTimeSeries;
})

onMounted(() => {
  const variables = [
    { name: '$namespace', class: 'namespace' },
    { name: '$node', class: 'node' },
  ];
  variables.forEach(v => {
    if (props.panelConfig.targets[0].expr.indexOf(v.name) > 0) {
      usingVariables.value.push(v);
    }
  });
  sidePanel.$subscribe((_, state) => {
    showPanelPosition.value = state.show && state.type == 'DashboardInfo';
  });
})

function setIsLoading(b) {
  isLoading.value = b;
}

function togglePanelInfo() {
  sidePanel.goToPanelConfig(props.position);
}
</script>
<template>
  <div class="flex border-b">
    <button v-if="showPanelPosition" class="p-1 bg-cyan-100 dark:bg-cyan-900" @click="sidePanel.goToPanelConfig(position)">
      {{ position }}
    </button>
    <div class="flex-1 py-1 text-center font-bold" :class="{ 'is-loading': isLoading }">
      <button class="hover:underline" @click="togglePanelInfo">
        {{ panelConfig.title }}
        <span v-for="v in usingVariables">
          <span class="w-2 rounded-full" :class="v.class">n</span>
        </span>
      </button>
    </div>
  </div>
  <component :is="componentName" :count="count" :is-loading="isLoading" :panel-config="panelConfig"
    :panel-width="panelWidth" :time-range="timeRange" @setIsLoading="setIsLoading" />
</template>

<style scoped>
@reference "tailwindcss";

.node {
  @apply text-yellow-500;
}

.namespace {
  @apply text-green-500;
}
</style>
