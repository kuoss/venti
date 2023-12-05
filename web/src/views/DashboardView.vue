<script setup>
import Dropdown from '@/components/Dropdown.vue';
import Panel from '@/components/Panel.vue';
import RunButton from '@/components/RunButton.vue';
import SidePanel from '@/components/SidePanel.vue';
import TimeRangePicker from '@/components/TimeRangePicker.vue';

import { ref, watch, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';

import { useDashboardStore } from '@/stores/dashboard';
import { useFilterStore } from '@/stores/filter';
import { useSidePanelStore } from '@/stores/sidePanel';
import { useTimeStore } from '@/stores/time';

const route = useRoute();
const sidePanel = useSidePanelStore();
const timeStore = useTimeStore();

const count = ref(0);
const range = ref([]);
const timeRange = ref([])
const intervalSeconds = ref(0)
const currentDropdown = ref(-1)
const dashboards = ref([])
const dashboard = ref({})
const namespaces = ref([])
const nodes = ref([])
const clientWidth = ref(100)
const root1 = ref(null)

function updateTimeRange(r) {
  range.value = r;
  execute()
}

async function execute() {
  timeRange.value = await timeStore.toTimeRangeForQuery(range.value);
  count.value++;
  if (intervalSeconds.value > 0) {
    setTimeout(() => timerHandler, intervalSeconds.value * 1000);
  }
}

function timerHandler() {
  if (timeStore.timerManager != 'DashboardView' || intervalSeconds.value == 0) return;
  execute();
}

function changeInterval(i) {
  intervalSeconds.value = i;
  execute();
}

function onDropdownOpen(uid) {
  currentDropdown.value = uid;
}

function selectNamespace(ns) {
  useFilterStore().selectedNamespace = ns;
  execute();
}

function selectNode(node) {
  useFilterStore().selectedNode = node;
  execute();
}

function onResize() {
  clientWidth.value = root1.value.clientWidth;
}

async function init() {
  namespaces.value = await useFilterStore().getNamespaces();
  nodes.value = await useFilterStore().getNodes();
  dashboards.value = await useDashboardStore().getDashboards();
  renderDashboard();
  execute();
}

function renderDashboard() {
  dashboards.value.forEach(d => {
    if (d.title == route.params.name) {
      dashboard.value = d;
      sidePanel.updateDashboardInfo(d);
    }
  });
}

watch(
  () => route.params,
  () => {
    renderDashboard();
  },
);

onMounted(() => {
  timeStore.timerManager = 'DashboardView';
  init();
  window.addEventListener('resize', onResize);
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize);
})
</script>

<template>
  <div class="w-full" ref="root1">
    <nav class="fixed left-0 w-full bg-white dark:bg-black border-b shadow p-2 pl-52 pr-4 z-10">
      <div class="flex gap-2">
        <div class="flex-none">
          <Dropdown :options="namespaces" :current-dropdown="currentDropdown" @select="selectNamespace"
            @open="onDropdownOpen" />
        </div>
        <div class="flex-none">
          <Dropdown :options="nodes" :current-dropdown="currentDropdown" @select="selectNode" @open="onDropdownOpen" />
        </div>
        <div class="grow font-bold text-center align-middle">
          {{ dashboard.title }}
          <button class="ml-2 px-2 py-2 border rounded" @click="sidePanel.toggleDashboardInfo()">
            <i class="mdi mdi-information-outline" />
          </button>
        </div>
        <div class="flex-none">
          <TimeRangePicker @updateTimeRange="updateTimeRange" />
        </div>
        <div class="flex-none">
          <RunButton :disabled="false" btn-text="Refresh" @execute="execute" @changeInterval="changeInterval" />
        </div>
      </div>
    </nav>
    <div class="block w-full mt-[56px] text-xs">
      <div class="flex">
        <div class="flex-1">
          <div class="p-4">
            <div v-for="(row, i) in dashboard.rows" class="pb-3 grid gap-3" :class="[
              'flex',
              [
                '',
                'grid-cols-1',
                'grid-cols-2',
                'grid-cols-3',
                'grid-cols-4',
                'grid-cols-5',
                'grid-cols-6',
                'grid-cols-7',
                'grid-cols-8',
                'grid-cols-9',
                'grid-cols-10',
                'grid-cols-11',
                'grid-cols-12',
              ][row.panels.length],
            ]">
              <div v-for="(panel, j) in row.panels" class="flex-1 bg-white dark:bg-black border">
                <Panel :position="`${i + 1}${j + 1}`" :count="count" :panel-config="panel"
                  :panel-width="clientWidth / row.panels.length" :time-range="timeRange" />
              </div>
            </div>
          </div>
        </div>
        <div>
          <SidePanel />
        </div>
      </div>
    </div>
  </div>
</template>
