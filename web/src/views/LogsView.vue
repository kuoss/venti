<template>
  <header
    class="fixed right-0 w-full bg-white border-b border-common shadow z-30 p-2 pl-52"
    :class="{ 'is-loading': loading }"
  >
    <div class="flex items-center flex-row">
      <div>
        <i class="mdi mdi-18px mdi-text"></i> Logs
      </div>
      <div class="flex ml-auto">
        <span>
          <TimeRangePicker @updateTimeRange="updateTimeRange" />
        </span>
        <span class="ml-2">
          <RunButton
            :disabled="busy"
            btnText="Run query"
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
        <div class="flex">
          <input
            type="search"
            class="flex-1 relative flex-auto min-w-0 block w-full px-3 py-1.5 text-base font-normal text-gray-700 bg-white bg-clip-padding border border-solid rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none"
            placeholder="Expression"
            aria-label="Expression"
            aria-describedby="button-addon3"
            v-model="expr"
            @keyup.enter="execute()"
          />
        </div>
      </div>
      <div class="break-all">
        <table class="w-full bg-slate-300">
          <tr>
            <td
              class="px-2"
              v-if="lastExecuted"
            >{{ lastExecuted.range[0] }} - {{ lastExecuted.range[1] }}</td>
            <td class="px-2" v-if="lastExecuted">{{ lastExecuted.expr }}</td>
            <td class="px-2 font-bold text-right">{{ result.length }} rows</td>
          </tr>
        </table>
        <div v-if="errorResponse">
          <div
            class="rounded bg-slate-200 text-yellow-200 text-center py-20"
          >[{{ errorResponse.status }}] {{ errorResponse.data.error }}</div>
        </div>
        <div v-else-if="result && result.length < 1">
          <div class="rounded bg-slate-200 text-center py-20">empty query result</div>
        </div>
        <div v-else-if="result">
          <div ref="logs" class="font-mono overflow-y-auto border bg-white max-h-[80vh]">
            <div v-if="resultType == 'vector'">
              <div class="border-b text-lg p-2" v-for="row in result">{{ row }}</div>
            </div>
            <div v-else>
              <table>
                <tr v-for="row in rows">
                  <td class="border-b border-slate-100 hover:bg-slate-200">
                    <span
                      v-for="column in row.columns"
                      :class="column.class"
                      @click="onClickColumn(column)"
                    >{{ column.text }}</span>
                  </td>
                </tr>
              </table>
            </div>
          </div>
        </div>
        <table class="w-full bg-slate-300">
          <tr>
            <td
              class="px-2"
              v-if="lastExecuted"
            >{{ lastExecuted.range[0] }} - {{ lastExecuted.range[1] }}</td>
            <td class="px-2" v-if="lastExecuted">{{ lastExecuted.expr }}</td>
            <td class="px-2 font-bold text-right">{{ result.length }} rows</td>
          </tr>
        </table>
      </div>
    </div>
    <div class="w-80">
      <div class="fixed right-0 bottom-0 bg-slate-300 text-xs pt-4 w-80">
        <div>
          <div class="w-full py-2 text-center">Targets</div>
        </div>
        <div
          class="overflow-y-auto border-l border-2 w-full bg-slate-200 border-b"
          style="height: calc(100vh - 90px);"
        >
          <div class="bg-gray-100 text-center cursor-pointer" @click="expr = 'audit'">audit</div>
          <div class="bg-gray-100 text-center cursor-pointer" @click="expr = 'pod{namespace=\'kube-system\',pod=\'eventrouter-.*\'}'">event</div>
          <div v-for="ns in namespaces">
            <div class="cursor-pointer" @click="selectEventNamespace(ns.name)">{{ ns.name }}</div>
          </div>
          <div class="bg-gray-100 text-center cursor-pointer" @click="expr = 'node'">node</div>
          <div
            class="cursor-pointer"
            v-for="node in nodes"
            @click="selectNode(node.name)"
          >{{ node.name }}</div>
          <div class="bg-gray-100 text-center cursor-pointer" @click="expr = 'pod'">pod</div>
          <div v-for="ns in namespaces">
            <div class="cursor-pointer" @click="ns.isExpanded = !ns.isExpanded">
              <i class="mdi" :class="[ns.isExpanded ? 'mdi-chevron-down' : 'mdi-chevron-right']"></i>
              {{ ns.name }}
            </div>
            <div v-if="ns.isExpanded" v-for="w in ns.workloads">
              <div class="pl-2 cursor-pointer" @click="w.isExpanded = !w.isExpanded">
                <i class="mdi" :class="[w.isExpanded ? 'mdi-chevron-down' : 'mdi-chevron-right']"></i>
                {{ w.kind }} ({{ w.objects.length }})
              </div>
              <div
                class="pl-5 cursor-pointer"
                v-if="w.isExpanded"
                v-for="object in w.objects"
                @click="selectObject(w.kind, ns.name, object.name)"
              >{{ object.name }}</div>
            </div>
          </div>
          <div class="pb-16">
            <hr />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useTimeStore } from "@/stores/time"
import TimeRangePicker from "@/components/TimeRangePicker.vue"
import RunButton from "@/components/RunButton.vue"

export default {
  components: {
    TimeRangePicker,
    RunButton,
  },
  data() {
    return {
      errorResponse: null,
      busy: false,
      loading: false,
      intervalSeconds: 0,
      range: [],
      nodes: [],
      namespaces: [],
      workloads: [],
      lastExecuted: null,
      metadata: {},
      queryType: 'raw',
      expr: 'pod{namespace="kube-system"}',
      showKeys: [],
      keys: [],
      topKeys: [],
      result: [],
      resultType: null,
      tab: 0,
      metricInfo: {
        selected: null,
        timer: 0,
      },
      chartData: [],
      chartOptions: {
        width: 500,
        height: 300,
      },
      timerID: null,
    }
  },
  computed: {
    rows() {
      const labels = this.getLogLabels()
      const classes = ['text-green-600', 'text-cyan-600', 'text-blue-600', 'text-purple-600', 'text-pink-600']
      return this.result.map(x => {
        const idx = x.indexOf(' ')

        let columns = [{ text: x.substr(0, 20), class: 'text-yellow-500' }]
        if (idx == 20) {
          return { columns: [...columns, { text: ' ' }, { text: x.substr(idx) }] }
        }
        columns.push({ text: '[' },)
        const parts = x.substr(21, idx - 22).split('|')
        parts.forEach((text, i) => {
          columns.push({ text: text, class: classes[i] + ' cursor-pointer hover:underline', label: labels[i] })
          columns.push({ text: '|' })
        })
        columns.pop()
        return { columns: [...columns, { text: ']' }, { text: x.substr(idx) }] }
      })
    }
  },
  methods: {
    selectObject(kind, namespace, name) {
      if (kind != 'pod') name += '-.*'
      this.expr = `pod{namespace="${namespace}", pod="${name}"}`
    },
    onClickColumn(c) {
      if (!c.label) return
      const idx = this.expr.indexOf('}')
      if (idx < 0) {
        this.expr += `{${c.label}="${c.text}"}`
        return
      }
      this.expr = `${this.expr.slice(0, idx)}, ${c.label}="${c.text}"${this.expr.slice(idx)}`
    },
    getLogLabels() {
      let podIndex = this.expr.indexOf('pod')
      if (podIndex < 0) podIndex = 9999
      let nodeIndex = this.expr.indexOf('node')
      if (nodeIndex < 0) nodeIndex = 9999
      if (podIndex < nodeIndex ) return ['namespace', 'pod', 'container']
      if (nodeIndex < podIndex ) return ['node', 'process']
      return []
    },
    changeInterval(i) {
      this.intervalSeconds = i
      this.execute()
    },
    updateTimeRange(r) {
      this.range = r
    },
    selectEventNamespace(ns_name) {
      this.expr = `pod{namespace="kube-system",pod="eventrouter-.*"} | "namespace":"${ns_name}"`
    },
    selectNode(name) {
      this.expr = `node{node="${name}"}`
    },
    selectNamespace(name, idx) {
      this.expr = `pod{namespace="${name}"}`
      this.namespaces[idx].showPods = !this.namespaces[idx].showPods
    },
    selectPod(ns_name, ns_idx, pod_name, pod_idx) {
      this.expr = `pod{namespace="${ns_name}", pod="${pod_name}"}`
      this.namespaces[ns_idx].pods[pod_idx].showContainers = !this.namespaces[ns_idx].pods[pod_idx].showContainers
    },
    selectContainer(namespace, pod, container) {
      this.expr = `pod{namespace="${namespace}", pod="${pod}", container="${container}"}`
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
        const response = await this.axios.get('/api/lethe/query_range', {
          params: {
            expr: this.expr,
            start: timeRange[0],
            end: timeRange[1],
          }
        })
        this.loading = false
        this.resultType = response.data.data.resultType
        this.result = response.data.data.result
        setTimeout(() => { if (this.$refs.logs) this.$refs.logs.scrollTop = 99999 }, 100)
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
      if (useTimeStore().timerManager != 'LogsView' || this.intervalSeconds == 0) return
      this.execute()
    },
    selectMetric(k) {
      this.metricInfo.selected = k ? [k, this.metadata[k]] : null
    },
    applyTarget(k) {
      this.expr = k
    },
    clickOutside() {
      this.selectMetric(null)
    },
    async fetchMetadata() {
      try {
        const now = await useTimeStore().getNow()
        const kinds = ['deployment', 'statefulset', 'daemonset', 'job', 'cronjob', 'pod']
        const groupBy = (arr, key) => arr.reduce((acc, item) => ((acc[item[key]] = [...(acc[item[key]] || []), item]), acc), {})

        let response, namespaces
        response = await this.axios.get('/api/prometheus/query', { params: { expr: 'kube_node_created', time: now } })
        this.nodes = response.data.data.result.map(x => { return { name: x.metric.node } })

        response = await this.axios.get('/api/prometheus/query', { params: { expr: 'kube_namespace_created', time: now } })
        namespaces = response.data.data.result.map(x => { return { name: x.metric.namespace, isExpanded: false, workloads: [] } })

        for (const [i, kind] of kinds.entries()) {
          response = await this.axios.get('/api/prometheus/query', { params: { expr: `kube_${kind}_created`, time: now } })
          const objects = groupBy(response.data.data.result.map(x => { return { namespace: x.metric.namespace, isExpanded: false, pods: [], name: x.metric[kind == 'job' ? 'job_name' : kind] } }), 'namespace')
          Object.entries(objects).forEach(x => {
            for (const [i, ns] of Object.entries(namespaces)) {
              if (ns.name == x[0]) namespaces[i].workloads.push({ kind: kind, isExpanded: false, objects: x[1] })
            }
          })
        }
        this.namespaces = namespaces
      } catch (error) { console.error(error) }
    },
  },
  mounted() {
    useTimeStore().timerManager = 'LogsView'
    this.fetchMetadata()
    if (this.$route.query?.query) {
      this.expr = this.$route.query.query
      setTimeout(this.execute, 500)
    }
    window.addEventListener("resize", this.chartResize)
  },
  unmounted() {
    window.removeEventListener("resize", this.chartResize)
  }
}
</script>
