<script setup>
import { useTimeStore } from "@/stores/time"
import TimeRangePicker from "@/components/TimeRangePicker.vue"
import RunButton from "@/components/RunButton.vue"
</script>

<template>
  <header
    class="fixed right-0 w-full bg-white border-b border-common shadow z-30 p-2 pl-52"
    :class="{ 'is-loading': loading }"
  >
    <div class="flex items-center flex-row">
      <div>
        <i class="mdi mdi-18px mdi-numeric"></i> Metrics
      </div>
      <div class="flex ml-auto">
        <span>
          <TimeRangePicker @updateTimeRange="updateTimeRange" />
        </span>
        <span class="ml-2">
          <RunButton
            btnText="Run query"
            :disabled="busy"
            @execute="execute"
            @changeInterval="changeInterval"
          />
        </span>
      </div>
    </div>
  </header>
  <div class="w-full flex mt-[3.6rem]">
    <div class="flex-1 py-4 px-4">
      <div class="pb-4">
        <div class="relative w-full">
          <input
            type="search"
            class="flex-1 relative flex-auto min-w-0 block w-full px-3 py-1.5 text-base font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none"
            placeholder="Expression"
            aria-label="Expression"
            aria-describedby="button-addon3"
            v-model="expr"
            @keyup="searchKeyUp"
          />
          <ul
            class="absolute bg-white border max-h-[70vh] overflow-y-auto z-20"
            v-if="searchMode && expr"
          >
            <li
              class="flex gap-3 hover:bg-gray-200 cursor-pointer"
              v-for="item in items"
              @click="expr = item[0]; searchMode = false"
            >
              <div class="text-gray-600" v-html="item[2]"></div>
              <div class="flex-auto text-right text-gray-500">{{ item[1][0].type }}</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="break-all">
        <div v-if="result.length < 1">
          <div class="rounded bg-slate-200 text-center p-8">Empty query result</div>
        </div>
        <div v-else>
          <div class="border">
            <uplotvue :data="chartData" :options="chartOptions" />
          </div>

          <div class="mt-4 py-1 font-bold">
            <span v-if="result && result.length > 0">{{ result.length }} rows</span>
            <div
              class="float-right"
              v-if="cursorTime"
            >{{ useTimeStore().timestamp2ymdhis(cursorTime) }}</div>
          </div>
          <div
            class="overflow-x-auto overflow-y-auto margin-l-[5em] max-h-[50vh]"
            :style="{ width: tableWitdh + 'px' }"
          >
            <table class="whitespace-nowrap border-separate w-full" style="border-spacing:0">
              <tr class="sticky z-10 top-0 border-y bg-slate-200 text-left">
                <th
                  class="font-normal max-w-[100px] px-2 border border-r-0 text-ellipsis overflow-hidden hover:whitespace-normal hover:min-w-[200px]"
                  v-for="key in keys"
                >{{ key }}</th>
                <th
                  class="min-w-[120px] sticky top-0 right-0 font-normal border bg-slate-200 text-center"
                >VALUE</th>
              </tr>
              <tr v-if="result && result.length>0" v-for="row in result" class="border-b hover:bg-gray-200">
                <td
                  class="max-w-[250px] px-2 border border-r-0 text-ellipsis overflow-hidden hover:whitespace-normal hover:min-w-[200px]"
                  v-for="key in keys"
                  @mouseover="row.hover = row.hover || {}; row.hover[key] = true"
                  @mouseleave="row.hover[key] = false"
                >
                  {{ row.metric[key] }}
                  <span class="inline-flex" v-if="row.hover && row.hover[key]">
                    <button
                      class="rounded px-1 border bg-slate-50 ml-1"
                      @click="addLabel('', key, row.metric[key])"
                    >
                      <i class="mdi mdi-plus-circle-outline"></i>
                    </button>
                    <button
                      class="rounded px-1 border bg-slate-50"
                      @click="addLabel('!', key, row.metric[key])"
                    >
                      <i class="mdi mdi-minus-circle-outline"></i>
                    </button>
                  </span>
                </td>
                <td class="sticky right-0 top-auto px-4 border bg-slate-50 text-right">
                  <span v-if="cursorIdx">{{ row.values[cursorIdx][1] }}</span>
                </td>
              </tr>
            </table>
          </div>
        </div>
      </div>
    </div>
    <div class="w-80">
      <div class="fixed right-0 bottom-0 bg-slate-300 text-xs pt-4 w-80">
        <div>
          <ul class="w-full flex list-none">
            <li
              class="py-3 basis-1/2 text-center hover:bg-slate-50 cursor-pointer border-b-2 border-transparent"
              :class="tab == 0 ? 'active' : ''"
              @click="tab = 0"
            >Metrics ({{ Object.keys(metadata).length }})</li>
            <li
              class="py-3 basis-1/2 text-center hover:bg-slate-50 cursor-pointer border-b-2 border-transparent"
              :class="tab == 1 ? 'active' : ''"
              @click="tab = 1"
            >Labels ({{ keys.length }})</li>
          </ul>
        </div>
        <div
          class="overflow-y-auto scrollbar-thin scrollbar-thumb-rounded scrollbar-track-rounded scrollbar-thumb-slate-300 scrollbar-track-transparnt dark:scrollbar-thumb-slate-900 border-l border-2 border-slate-300 w-full bg-slate-200 border-b"
          style="height: calc(100vh - 100px);"
        >
          <div v-if="tab == 0">
            <div class v-for="(d, k) in metaDict">
              <div
                class="pl-1 cursor-pointer text-stone-600"
                @click="d.showMetrics = !d.showMetrics"
              >{{ k }} ({{ d.metrics.length }})</div>
              <div
                class="pl-4 overflow-hidden text-ellipsis hover:bg-white cursor-pointer"
                v-if="d.showMetrics"
                v-for="m in d.metrics"
                @click="applyMetric(m)"
                @mouseover="selectMetric(m)"
              >{{ m.name }}</div>
            </div>
          </div>
          <div v-else>
            <div v-for="(d, key) in keyDict">
              <div
                class="pl-1 overflow-hidden text-ellipsis hover:bg-white cursor-pointer"
                @click="d.show = !d.show"
              >{{ key }} ({{ d.values.length }})</div>
              <div class="pl-4" v-if="d.show" v-for="v in d.values">
                {{ v }}
                <button
                  class="rounded px-1 border bg-slate-50 ml-1"
                  @click="addLabel('', key, v)"
                >
                  <i class="mdi mdi-plus-circle-outline"></i>
                </button>
                <button class="rounded px-1 border bg-slate-50" @click="addLabel('!', key, v)">
                  <i class="mdi mdi-minus-circle-outline"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div
    v-if="metricInfo"
    v-click-outside="clickOutside"
    class="fixed z-50 top-[9rem] right-[20.5rem] w-80 bg-white border border-slate-300 rounded opacity-[.9] hover:opacity-100"
  >
    <div class="border-b border-slate-300 p-2 break-all font-bold">{{ metricInfo.name }}</div>
    <div class="px-2 py-1 word-break" v-for="(v, k) in metricInfo.data[0]">{{ v }}</div>
  </div>
</template>

<script>
import UplotVue from 'uplot-vue'
import 'uplot/dist/uPlot.min.css'

export default {
  components: {
    RunButton,
    TimeRangePicker,
    uplotvue: UplotVue,
  },
  computed: {
    items() {
      const keyword = this.expr
      if (!keyword || keyword.length < 1) return []
      return Object.entries(this.metadata).filter(x => x[0].indexOf(keyword) >= 0).map(x => {
        x.push(x[0].replaceAll(keyword, `<span class="text-blue-600 font-bold">${keyword}</span>`))
        return x
      })
    }
  },
  data() {
    return {
      searchMode: false,
      cursorIdx: null,
      cursorTime: null,
      errorResponse: null,
      busy: false,
      loading: false,
      intervalSeconds: 0,
      range: [],
      lastExecuted: null,
      metadata: {},
      metaDict: {},
      metricInfo: null,
      queryType: 'raw',
      expr: `container_memory_working_set_bytes{namespace="kube-system"}`,
      keys: [],
      keyDict: {},
      result: [],
      tab: 0,
      time: null,
      chartData: [],
      chartOptions: {
        axes: [
          {
            stroke: "#888",
            grid: { stroke: "#8885", width: 1 },
            ticks: { stroke: "#8885", width: 1 },
            values: [
              [3600 * 24 * 365, "{YYYY}", null, null, null, null, null, null, 1],
              [3600 * 24 * 28, "{MM}", "\n{YYYY}", null, null, null, null, null, 1],
              [3600 * 24, "{MM}-{DD}", "\n{YYYY}", null, null, null, null, null, 1],
              [3600, "{HH}:00", "\n{YYYY}-{MM}-{DD}", null, "\n{MM}-{DD}", null, null, null, 1],
              [60, "{HH}:{mm}", "\n{YYYY}-{MM}-{DD}", null, "\n{MM}-{DD}", null, null, null, 1],
              [1, "{HH}:{mm}:{ss}", "\n{YYYY}-{MM}-{DD}", null, "\n{MM}-{DD}", null, null, null, 1],
            ],
          },
          {
            stroke: "#888",
            grid: { stroke: "#8885", width: 1 },
            ticks: { stroke: "#8885", width: 1 },
            size(self, values, axisIdx, cycleNum) {
              const axis = self.axes[axisIdx]
              if (cycleNum > 1) return axis._size
              let axisSize = axis.ticks.size + axis.gap;
              let longestVal = (values ?? []).reduce((acc, val) => (
                val.length > acc.length ? val : acc
              ), "");
              if (longestVal != "") {
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
        plugins: [this.tooltipPlugin()],
      },
    }
  },
  methods: {
    searchKeyUp(e) {
      if (e.keyCode == 13) {
        this.searchMode = false
        this.execute()
        return
      }
      this.searchMode = true
    },
    addLabel(not, key, value) {
      const where = `${key}${not}="${value}"`
      const idx = this.expr.indexOf('}')
      if (idx < 0) {
        this.expr += `{${where}}`
        return
      }
      this.expr = this.expr.slice(0, -1) + `,${where}` + this.expr.slice(-1)
    },
    changeInterval(i) {
      this.intervalSeconds = i
      this.execute()
    },
    updateTimeRange(r) {
      this.range = r
    },
    async execute() {
      if (this.expr.length < 1) {
        console.error('emtpy expr')
        return
      }
      const timeRange = await useTimeStore().toTimeRangeForQuery(this.range)
      let lastRange = timeRange.map(x => useTimeStore().timestamp2ymdhis(x))
      if (lastRange[0].slice(0, 10) == lastRange[1].slice(0, 10)) lastRange[1] = lastRange[1].slice(11)
      this.lastExecuted = { expr: this.expr, range: lastRange }
      this.loading = true
      try {
        const response = await this.axios.get('/api/prometheus/query_range', {
          params: {
            expr: this.expr,
            start: timeRange[0],
            end: timeRange[1],
            step: (timeRange[1] - timeRange[0]) / 120,
          }
        })
        this.loading = false
        this.result = response.data.data.result
        this.keys = this.result.map(x => Object.keys(x.metric)).flat().filter((v, i, s) => s.indexOf(v) === i).sort().slice(1, 99)
        this.renderChart()
        if (this.intervalSeconds > 0) {
          this.busy = true
          setTimeout(() => this.timerHandler(), this.intervalSeconds * 1000)
        }
        else {
          this.busy = false
        }
        this.errorResponse = null
      } catch (error) {
        this.loading = false
        this.errorResponse = error.response
      }
    },
    timerHandler() {
      if (useTimeStore().timerManager != 'MetricsView' || this.intervalSeconds == 0) return
      this.execute()
    },
    renderChart() {
      const temp = this.result.map(x => x.values)
      const timestamps = Array.from(new Set(temp.map(a => a.map(b => b[0])).flat())).sort()
      let seriesData = temp.map(a => {
        let newA = []
        timestamps.forEach(t => {
          const newPoint = a.filter(b => t == b[0])
          if (newPoint.length != 1 || isNaN(parseFloat(newPoint[0][1]))) {
            newA.push(null)
            return
          }
          newA.push(parseFloat(newPoint[0][1]))
        })
        return newA
      })
      let m = Math.max(...seriesData.flat())
      let c = 0
      while (m > 1000) {
        m /= 1000
        c++
      }
      const metrics = this.result.map(x => x.metric)
      let newSeries = []
      newSeries.push({})
      this.keyDict = {}
      metrics.forEach(x => {
        delete x.__name__
        const entries = Object.entries(x)
        entries.forEach((a) => {
          this.keyDict[a[0]] = this.keyDict[a[0]] || { show: false, values: [] }
          this.keyDict[a[0]].values.push(a[1])
          this.keyDict[a[0]].values = this.keyDict[a[0]].values.filter((v, i, s) => s.indexOf(v) === i)
        })
        x = '{' + entries.map(v => `${v[0]}="${v[1]}"`).join(',') + '}'

        newSeries.push({
          label: x,
          stroke: this.$util.string2color(x),
          points: { size: 1 },
        })
      })
      // this.chartOptions.axes[1].values = (self, ticks) => ticks.map(rawValue => rawValue / Math.pow(1000, c) + ['', 'k', 'M', 'G', 'T', 'P', 'E', 'Z', 'Y'][c])
      this.chartOptions = {
        ...this.chartOptions,
        series: newSeries,
        scales: useTimeStore().scales,
      }
      this.chartData = [timestamps, ...seriesData]
      this.chartResize()
    },
    selectMetric(m) {
      this.metricInfo = m
    },
    applyMetric(m) {
      this.metricInfo.selected = null
      this.expr = m.name
    },
    clickOutside() {
      this.selectMetric(null)
    },
    async fetchMetadata() {
      try {
        const response = await this.axios.get('/api/prometheus/metadata')
        this.metadata = response.data.data
        this.metaDict = Object.keys(this.metadata).reduce((a, k) => {
          const p = k.slice(0, k.indexOf('_'))
          a[p] = a[p] || { showMetrics: false }
          a[p].metrics = a[p].metrics || []
          a[p].metrics.push({ name: k, data: this.metadata[k] })
          return a
        }, {})
      } catch (error) { console.error(error) }
    },
    chartResize() {
      const width = document.body.clientWidth - 545
      this.chartOptions = { ...this.chartOptions, width: width }
      this.tableWitdh = width
    },
    tooltipPlugin() {
      return {
        hooks: {
          setCursor: u => {
            if (!u.cursor.idx) return
            this.cursorIdx = u.cursor.idx
            this.cursorTime = u.data[0][u.cursor.idx]
          }
        }
      }
    },
  },
  mounted() {
    useTimeStore().timerManager = 'MetricsView'
    this.fetchMetadata()
    if(this.$route.query?.query) {
      this.expr = this.$route.query.query
      setTimeout(this.execute, 500)
    }
    window.addEventListener("resize", this.chartResize)
  },
  unmounted() {
    window.removeEventListener("resize", this.chartResize)
  },
}
</script>

<style scope>
.headcol-before:before {
  content: "Row ";
}
.u-inline.u-live th::after {
  content: "";
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
  @apply border border-slate-200 table-cell;
}

.u-legend th {
  padding-left: 1.4rem;
  text-indent: -1rem;
  @apply font-medium text-left;
}

.u-legend th > .u-marker {
  @apply w-3 h-3;
}

.u-legend th > .u-label {
  @apply inline;
}

.u-value {
  @apply w-28 text-right;
}
</style>