<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import UplotVue from 'uplot-vue';
import 'uplot/dist/uPlot.min.css';
import { vElementSize } from '@vueuse/components'

import { useFilterStore } from '@/stores/filter';
import { useTimeStore } from '@/stores/time';
import { useSidePanelStore } from '@/stores/sidePanel';

import Util from '@/lib/util';

const props = defineProps({
  count: Number,
  isLoading: Boolean,
  panelConfig: { type: Object, required: true },
  panelWidth: { type: Number, required: true },
  timeRange: { type: Object, required: true },
})

const emit = defineEmits<{
  (e: 'setIsLoading', value: Boolean): void
}>()

const timeStore = useTimeStore()
const filter = useFilterStore()
const sidePanel = useSidePanelStore()

const myWidth = ref(0)

const isNoData = ref(true)
const data = ref([] as any[])
const options = ref({
  axes: [
    {
      stroke: '#888',
      grid: { stroke: '#8885', width: 1 },
      ticks: { stroke: '#8885', width: 1 },
      values: [
        [3600 * 24 * 365, '{YYYY}', null, null, null, null, null, null, 1],
        [3600 * 24 * 28, '{MM}', '\n{YYYY}', null, null, null, null, null, 1],
        [3600 * 24, '{MM}-{DD}', '\n{YYYY}', null, null, null, null, null, 1],
        [3600, '{HH}:00', '\n{YYYY}-{MM}-{DD}', null, '\n{MM}-{DD}', null, null, null, 1],
        [80, '{HH}:{mm}', '\n{YYYY}-{MM}-{DD}', null, '\n{MM}-{DD}', null, null, null, 1],
        [1, '{HH}:{mm}:{ss}', '\n{YYYY}-{MM}-{DD}', null, '\n{MM}-{DD}', null, null, null, 1],
      ],
    },
    {
      stroke: '#888',
      grid: { stroke: '#8885', width: 1 },
      ticks: { stroke: '#8885', width: 1 },
    },
  ],
  width: 100,
  height: 180,
  legend: { show: false },
  cursor: { points: false },
  scales: { x: { time: true }, y: { auto: true } },
  select: { show: false },
  series: [],
  plugins: [tooltipPlugin()],
} as any)

async function fetchData() {
  if (props.timeRange.length < 2) return;
  const target = props.panelConfig.targets[0];
  emit('setIsLoading', true);
  try {
    const response = await fetch(
      '/api/v1/remote/query_range?dsType=prometheus&' +
      new URLSearchParams({
        query: filter.renderExpr(props.panelConfig.targets[0].expr),
        start: props.timeRange[0],
        end: props.timeRange[1],
        step: String((props.timeRange[1] - props.timeRange[0]) / 120),
      }),
    );
    const jsonData = await response.json();

    const result = jsonData.data.result;
    if (result.length < 1) {
      isNoData.value = true;
      return;
    }
    isNoData.value = false;

    const temp = result.map((x: any) => x.values);
    const timestamps = Array.from(new Set(temp.map((a: any) => a.map((b: any) => b[0])).flat())).sort();
    const seriesData = temp.map((a: any) => {
      let newA: any = [];
      timestamps.forEach(t => {
        const newPoint = a.filter((b: any) => t == b[0]);
        if (newPoint.length != 1 || isNaN(parseFloat(newPoint[0][1]))) {
          newA.push(null);
          return;
        }
        newA.push(parseFloat(newPoint[0][1]));
      });
      return newA;
    });
    // labels
    const labels = result.map((x: any) => target.legend.replace(/\{\{(.*?)\}\}/g, (_: any, m: any) => x.metric[m]));
    let newSeries = [];
    newSeries.push({});
    labels.forEach((x: any) =>
      newSeries.push({
        label: x,
        stroke: Util.string2color(x),
        width: 1,
        points: { size: 0 },
      }),
    );
    let newOptions = { ...options.value }
    newOptions.series = newSeries
    newOptions.scales = {
      y: {
        range: (_a: any, _b: any, fromMax: any) => [0, Math.max(fromMax, props.panelConfig.chartOptions?.yMax ?? 0)],
      },
    };
    data.value = [timestamps, ...seriesData];
    options.value = newOptions;
  } catch (error) {
    console.error(error);
  }
  emit('setIsLoading', false);
}

function onResize({ width, height }: { width: number; height: number }) {
  myWidth.value = width
  const newOptions = {
    ...options.value,
    width: width,
  }
  options.value = newOptions
}

function onChartClick() {
  sidePanel.toggleShow('DataTable');
}

function tooltipPlugin() {
  return {
    hooks: {
      setCursor: (u: any) => {
        let columnData = u.data.map((x: any) => x[u.cursor.idx]);
        const time = columnData.shift();
        if (!time) return;
        const labels = options.value.series.map((x: any) => x['label']).slice(1);
        const rows = labels.map((x: any, i: any) => [x, columnData[i]]).filter((x: any) => x[1] != undefined);
        sidePanel.updatetDataTable({
          title: props.panelConfig.title,
          time: time,
          rows: rows,
        });
      },
    },
  };
}

async function init() {
  await fetchData();
  // onResize();
}


watch(() => props.count, () => {
  if (!props.isLoading) fetchData();
})

onMounted(() => {
  init();
})

</script>

<template>
  <div v-if="!isLoading && isNoData" class="h-[150px] grid grid-cols-1 content-center">
    <div class="text-center text-lg">No data</div>
  </div>
  <div v-else>
    <div v-element-size="onResize">
      <UplotVue ref="main" :data="data" :options="options" @click="onChartClick" />
    </div>
  </div>
</template>
