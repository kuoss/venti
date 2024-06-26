<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import Util from '@/lib/util';
import { useTimeStore } from '@/stores/time';
import DropdownDatasource from '@/components/DropdownDatasource.vue'

import TimeRangePicker from '@/components/TimeRangePicker.vue';
import RunButton from '@/components/RunButton.vue';
import UplotVue from 'uplot-vue';
import 'uplot/dist/uPlot.min.css';

const route = useRoute();
const timeStore = useTimeStore();

const tableWidth = ref(0);
const searchMode = ref(false);
const cursorIdx = ref(null);
const cursorTime = ref(null);
const busy = ref(false);
const loading = ref(false);

const metadata = ref({});
const metaDict = ref({});
const metricInfo = ref(null);

const expr = ref(`container_memory_working_set_bytes{namespace="kube-system"}`);
const keys = ref([]);
const keyDict = ref({});
const result = ref([]);
const tab = ref(0);
// const time = ref(null)

let dsName = '';

const chartData = ref([]);
const chartOptions = ref({
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
        [60, '{HH}:{mm}', '\n{YYYY}-{MM}-{DD}', null, '\n{MM}-{DD}', null, null, null, 1],
        [1, '{HH}:{mm}:{ss}', '\n{YYYY}-{MM}-{DD}', null, '\n{MM}-{DD}', null, null, null, 1],
      ],
    },
    {
      stroke: '#888',
      grid: { stroke: '#8885', width: 1 },
      ticks: { stroke: '#8885', width: 1 },
      size(self, values, axisIdx, cycleNum) {
        const axis = self.axes[axisIdx];
        if (cycleNum > 1) return axis._size;
        let axisSize = axis.ticks.size + axis.gap;
        let longestVal = (values ?? []).reduce((acc, val) => (val.length > acc.length ? val : acc), '');
        if (longestVal != '') {
          self.ctx.font = axis.font[0];
          axisSize += self.ctx.measureText(longestVal).width / devicePixelRatio;
        }
        return Math.ceil(axisSize);
      },
    },
  ],
  width: 400,
  height: 280,
  legend: { show: false },
  cursor: { points: false },
  scales: { x: { time: true }, y: { auto: true } },
  select: { show: false },
  series: [],
  plugins: [tooltipPlugin()],
});

let errorResponse;
let intervalSeconds = 0;
let range = [];
let lastExecuted = null;

const items = computed(() => {
  const keyword = expr.value;
  if (!keyword || keyword.length < 1) return [];
  return Object.entries(metadata.value)
    .filter(x => x[0].indexOf(keyword) >= 0)
    .map(x => {
      x.push(x[0].replaceAll(keyword, `<span class="text-blue-600 font-bold">${keyword}</span>`));
      return x;
    });
});

onMounted(() => {
  timeStore.timerManager = 'MetricsView';
  fetchMetadata();
  if (route.query?.query) {
    expr.value = '' + route.query.query;
    setTimeout(execute, 500);
  }
  window.addEventListener('resize', chartResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', chartResize);
});

function searchKeyUp(e) {
  if (e.keyCode == 13) {
    searchMode.value = false;
    execute();
    return;
  }
  searchMode.value = true;
}

function addLabel(not, key, value) {
  const where = `${key}${not}="${value}"`;
  const idx = expr.value.indexOf('}');
  if (idx < 0) {
    expr.value += `{${where}}`;
    return;
  }
  expr.value = expr.value.slice(0, -1) + `,${where}` + expr.value.slice(-1);
}

function changeInterval(i) {
  intervalSeconds = i;
  execute();
}

function updateTimeRange(r) {
  range = r;
}

async function execute() {
  if (expr.value.length < 1) {
    console.error('emtpy expr');
    return;
  }
  const timeRange = await timeStore.toTimeRangeForQuery(range);
  let lastRange = timeRange.map(x => timeStore.timestamp2ymdhis(x));
  if (lastRange[0].slice(0, 10) == lastRange[1].slice(0, 10)) lastRange[1] = lastRange[1].slice(11);
  lastExecuted = { expr: expr.value, range: lastRange };
  loading.value = true;
  try {
    const response = await fetch(
      '/api/v1/remote/query_range?' +
      new URLSearchParams({
        dsName: dsName,
        query: expr.value,
        start: timeRange[0],
        end: timeRange[1],
        step: (timeRange[1] - timeRange[0]) / 120,
      }),
    );
    const jsonData = await response.json();
    loading.value = false;
    result.value = jsonData.data.result;

    let keySet = new Set()
    for (const v of result.value) {
      keySet = new Set([...keySet, ...Object.keys(v.metric)])
    }
    keySet.delete('__name__')
    keys.value = Array.from(keySet)
    renderChart();
    if (intervalSeconds > 0) {
      busy.value = true;
      setTimeout(() => timerHandler(), intervalSeconds * 1000);
    } else {
      busy.value = false;
    }
    errorResponse = null;
  } catch (error) {
    loading.value = false;
    errorResponse = error.response;
  }
}

function timerHandler() {
  if (timeStore.timerManager != 'MetricsView' || intervalSeconds == 0) return;
  execute();
}

function renderChart() {
  const temp = result.value.map(x => x.values);
  const timestamps = Array.from(new Set(temp.map(a => a.map(b => b[0])).flat())).sort();
  let seriesData = temp.map(a => {
    let newA = [];
    timestamps.forEach(t => {
      const newPoint = a.filter(b => t == b[0]);
      if (newPoint.length != 1 || isNaN(parseFloat(newPoint[0][1]))) {
        newA.push(null);
        return;
      }
      newA.push(parseFloat(newPoint[0][1]));
    });
    return newA;
  });
  const metrics = result.value.map(x => x.metric);
  let newSeries = [];
  newSeries.push({});
  keyDict.value = {};

  metrics.forEach(x => {
    delete x.__name__;
    const entries = Object.entries(x);

    entries.forEach(a => {
      keyDict.value[a[0]] = keyDict.value[a[0]] || {
        show: false,
        values: [],
      };
      keyDict.value[a[0]].values.push(a[1]);
      keyDict.value[a[0]].values = keyDict.value[a[0]].values.filter((v, i, s) => s.indexOf(v) === i);
    });
    x = '{' + entries.map(v => `${v[0]}="${v[1]}"`).join(',') + '}';
    newSeries.push({
      label: x,
      stroke: Util.string2color(x),
      points: { size: 1 },
    });
  });
  // this.chartOptions.axes[1].values = (self, ticks) => ticks.map(rawValue => rawValue / Math.pow(1000, c) + ['', 'k', 'M', 'G', 'T', 'P', 'E', 'Z', 'Y'][c])
  chartOptions.value = {
    ...chartOptions.value,
    series: newSeries,
    scales: timeStore.scales,
  };
  chartData.value = [timestamps, ...seriesData];
  chartResize();
}

function selectMetric(m) {
  metricInfo.value = m;
}

function applyMetric(m) {
  metricInfo.value.selected = null;
  expr.value = m.name;
}

function clickOutside() {
  selectMetric(null);
}
async function fetchMetadata() {
  try {
    const resp = await fetch('/api/v1/remote/metadata?dsType=prometheus');
    const jsonData = await resp.json();
    metadata.value = jsonData.data;
    metaDict.value = Object.keys(metadata.value).reduce((a, k) => {
      const p = k.slice(0, k.indexOf('_'));
      a[p] = a[p] || { showMetrics: false };
      a[p].metrics = a[p].metrics || [];
      a[p].metrics.push({ name: k, data: metadata.value[k] });
      return a;
    }, {});
  } catch (err) {
    console.error(err);
  }
}

const chartResize = () => {
  const width = document.body.clientWidth - 545;
  chartOptions.value = { ...chartOptions.value, width: width };
  tableWidth.value = width;
};

function tooltipPlugin() {
  return {
    hooks: {
      setCursor: u => {
        if (!u.cursor.idx) return;
        cursorIdx.value = u.cursor.idx;
        cursorTime.value = u.data[0][u.cursor.idx];
      },
    },
  };
}

function clickItem(item) {
  expr.value = item[0];
  searchMode.value = false;
}

function onChangeDatasource(value) {
  dsName = value
}
</script>

<template>
  <header class="fixed right-0 w-full border-b border-common shadow z-30 p-2 pl-52 bg-white dark:bg-black"
    :class="{ 'is-loading': loading }">
    <div class="flex items-center flex-row">
      <div>
        <i class="mdi mdi-18px mdi-numeric" /> Metrics
        <DropdownDatasource dsType="prometheus" @change="onChangeDatasource" />
      </div>
      <div class="flex ml-auto">
        <span>
          <TimeRangePicker @updateTimeRange="updateTimeRange" />
        </span>
        <span class="ml-2">
          <RunButton btn-text="Run query" :disabled="busy" @execute="execute" @changeInterval="changeInterval" />
        </span>
      </div>
    </div>
  </header>
  <div class="w-full flex mt-[3.6rem]">
    <div class="flex-1 py-4 px-4">
      <div class="pb-4">
        <div class="relative w-full">
          <input v-model="expr" type="search"
            class="flex-auto relative min-w-0 block w-full px-3 py-1.5 text-base font-normal text-gray-700 dark:text-gray-300 bg-white dark:bg-black bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none"
            placeholder="Expression" aria-label="Expression" aria-describedby="button-addon3" @keyup="searchKeyUp" />
          <ul v-if="searchMode && expr" class="absolute bg-white dark:bg-black border max-h-[70vh] overflow-y-auto z-20">
            <li v-for="item in items" class="flex gap-3 hover:bg-gray-200 dark:hover:bg-gray-800 selection:cursor-pointer" @click="clickItem(item)">
              <div class="text-gray-600 dark:text-gray-300" v-html="item[2]" />
              <div class="flex-auto text-right text-gray-300 dark:text-gray-500">
                {{ item[1][0].type }}
              </div>
            </li>
          </ul>
        </div>
      </div>
      <div class="break-all">
        <div v-if="result.length < 1">
          <div class="rounded bg-slate-200 dark:bg-slate-800 text-center p-8">Empty query result</div>
        </div>
        <div v-else>
          <div class="border">
            <UplotVue :data="chartData" :options="chartOptions" />
          </div>

          <div class="mt-4 py-1 font-bold">
            <span v-if="result && result.length > 0">{{ result.length }} rows</span>
            <div v-if="cursorTime" class="float-right">
              {{ timeStore.timestamp2ymdhis(cursorTime) }}
            </div>
          </div>
          <div class="overflow-x-auto overflow-y-auto margin-l-[5em] max-h-[50vh]" :style="{ width: tableWidth + 'px' }">
            <table class="whitespace-nowrap border-separate w-full" style="border-spacing: 0">
              <tr class="sticky z-10 top-0 border-y bg-slate-200 dark:bg-slate-800 text-left">
                <th v-for="key in keys"
                  class="font-normal max-w-[100px] px-2 border border-r-0 text-ellipsis overflow-hidden hover:whitespace-normal hover:min-w-[200px]">
                  {{ key }}
                </th>
                <th class="min-w-[120px] sticky top-0 right-0 font-normal border bg-slate-200 dark:bg-slate-800 text-center">VALUE</th>
              </tr>
              <template v-if="result && result.length > 0">
                <tr v-for="row in result" class="border-b hover:bg-gray-200 dark:hover:bg-gray-800">
                  <td v-for="key in keys"
                    class="max-w-[250px] px-2 border border-r-0 text-ellipsis overflow-hidden hover:whitespace-normal hover:min-w-[200px]"
                    @mouseover="
                      row.hover = row.hover || {};
                    row.hover[key] = true;
                    " @mouseleave="row.hover[key] = false">
                    {{ row.metric[key] }}
                    <span v-if="row.hover && row.hover[key]" class="inline-flex">
                      <button class="rounded px-1 border bg-slate-50 dark:bg-slate-900 ml-1" @click="addLabel('', key, row.metric[key])">
                        <i class="mdi mdi-plus-circle-outline" />
                      </button>
                      <button class="rounded px-1 border bg-slate-50 dark:bg-slate-900" @click="addLabel('!', key, row.metric[key])">
                        <i class="mdi mdi-minus-circle-outline" />
                      </button>
                    </span>
                  </td>
                  <td class="sticky right-0 top-auto px-4 border bg-slate-50 dark:bg-slate-900 text-right">
                    <span v-if="cursorIdx">{{ row.values[cursorIdx][1] }}</span>
                  </td>
                </tr>
              </template>
            </table>
          </div>
        </div>
      </div>
    </div>
    <div class="w-80">
      <div class="fixed right-0 bottom-0 bg-slate-300 dark:bg-slate-700 text-xs pt-4 w-80">
        <div>
          <ul class="w-full flex list-none">
            <li class="py-3 basis-1/2 text-center hover:bg-slate-50 dark:bg-slate-900 cursor-pointer border-b-2 border-transparent"
              :class="tab == 0 ? 'active' : ''" @click="tab = 0">
              Metrics ({{ Object.keys(metadata).length }})
            </li>
            <li class="py-3 basis-1/2 text-center hover:bg-slate-50 dark:bg-slate-900 cursor-pointer border-b-2 border-transparent"
              :class="tab == 1 ? 'active' : ''" @click="tab = 1">
              Labels ({{ keys.length }})
            </li>
          </ul>
        </div>
        <div
          class="overflow-y-auto bg-slate-200 dark:bg-slate-800 scrollbar-thin scrollbar-thumb-rounded scrollbar-track-rounded scrollbar-thumb-slate-300 scrollbar-track-transparnt dark:scrollbar-thumb-slate-500 border-l border-2 border-slate-300 dark:border-slate-700 w-full border-b"
          style="height: calc(100vh - 100px)">
          <div v-if="tab == 0">
            <div v-for="(d, k) in metaDict" class>
              <div class="pl-1 cursor-pointer text-stone-600 dark:text-stone-400" @click="d.showMetrics = !d.showMetrics">
                {{ k }} ({{ d.metrics.length }})
              </div>
              <template v-if="d.showMetrics">
                <div v-for="m in d.metrics" class="pl-4 overflow-hidden text-ellipsis hover:bg-white cursor-pointer"
                  @click="applyMetric(m)" @mouseover="selectMetric(m)">
                  {{ m.name }}
                </div>
              </template>
            </div>
          </div>
          <div v-else>
            <div v-for="(d, key) in keyDict">
              <div class="pl-1 overflow-hidden text-ellipsis hover:bg-white dark:bg-black cursor-pointer" @click="d.show = !d.show">
                {{ key }} ({{ d.values.length }})
              </div>
              <template v-if="d.show">
                <div v-for="v in d.values" class="pl-4">
                  {{ v }}
                  <button class="rounded px-1 border bg-slate-50 dark:bg-slate-900 ml-1" @click="addLabel('', key, v)">
                    <i class="mdi mdi-plus-circle-outline" />
                  </button>
                  <button class="rounded px-1 border bg-slate-50 dark:bg-slate-900" @click="addLabel('!', key, v)">
                    <i class="mdi mdi-minus-circle-outline" />
                  </button>
                </div>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div v-if="metricInfo" v-click-outside="clickOutside"
    class="fixed z-50 top-[9rem] right-[20.5rem] w-80 bg-white dark:bg-black border border-slate-300 dark:border-slate-700 rounded opacity-[.9] hover:opacity-100">
    <div class="border-b border-slate-300 dark:border-slate-700 p-2 break-all font-bold">
      {{ metricInfo.name }}
    </div>
    <div v-for="v in metricInfo.data[0]" class="px-2 py-1 word-break">
      {{ v }}
    </div>
  </div>
</template>

<style scope>
.headcol-before:before {
  content: 'Row ';
}

.u-inline.u-live th::after {
  content: '';
}

.u-series:first-child {
  display: contents;
}

.u-series:first-child th {
  display: none;
}

.u-series:first-child td.u-value {
  display: block;
  width: 100%;
  text-align: right;
}

.u-series {
  @apply table w-full text-xs;
}

.u-legend th,
.u-legend td {
  @apply border border-slate-200 dark:border-slate-800 table-cell;
}

.u-legend th {
  padding-left: 1.4rem;
  text-indent: -1rem;
  @apply font-medium text-left;
}

.u-legend th>.u-marker {
  @apply w-3 h-3;
}

.u-legend th>.u-label {
  @apply inline;
}

.u-value {
  @apply w-28 text-right;
}
</style>
