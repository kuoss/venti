<script setup>
import { useTimeStore } from '@/stores/time';
import TimeRangePicker from '@/components/TimeRangePicker.vue';
import RunButton from '@/components/RunButton.vue';
import Util from '@/lib/util';
</script>
<script>
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
      logType: '',
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
    };
  },
  mounted() {
    useTimeStore().timerManager = 'LogsView';
    this.fetchMetadata();
    if (this.$route.query?.query) {
      this.expr = this.$route.query.query;
      setTimeout(this.execute, 500);
    }
    window.addEventListener('resize', this.chartResize);
  },
  unmounted() {
    window.removeEventListener('resize', this.chartResize);
  },
  methods: {
    detectLogType() {
      if (this.result.length < 1) {
        this.logType = '';
        return;
      }
      if (Object.prototype.hasOwnProperty.call(this.result[0], 'namespace')) {
        this.logType = 'pod';
        return;
      }
      if (Object.prototype.hasOwnProperty.call(this.result[0], 'node')) {
        this.logType = 'node';
        return;
      }
      this.logType = '';
    },
    selectObject(kind, namespace, name) {
      if (kind == 'pod') this.expr = `pod{namespace="${namespace}", pod="${name}"}`;
      else this.expr = `pod{namespace="${namespace}", pod=~"${name}-.*"}`;
    },
    onClickColumn(c) {
      if (!c.label) return;
      const idx = this.expr.indexOf('}');
      if (idx < 0) {
        this.expr += `{${c.label}="${c.text}"}`;
        return;
      }
      this.expr = `${this.expr.slice(0, idx)}, ${c.label}="${c.text}"${this.expr.slice(idx)}`;
    },
    getLogLabels() {
      let podIndex = this.expr.indexOf('pod');
      if (podIndex < 0) podIndex = 9999;
      let nodeIndex = this.expr.indexOf('node');
      if (nodeIndex < 0) nodeIndex = 9999;
      if (podIndex < nodeIndex) return ['namespace', 'pod', 'container'];
      if (nodeIndex < podIndex) return ['node', 'process'];
      return [];
    },
    changeInterval(i) {
      this.intervalSeconds = i;
      this.execute();
    },
    updateTimeRange(r) {
      this.range = r;
    },
    selectEvent() {
      this.expr = `pod{namespace="kube-system",container="eventrouter"}`;
    },
    selectNode(name) {
      this.expr = `node{node="${name}"}`;
    },
    async execute() {
      if (this.expr.length < 1) {
        console.error('emtpy expr');
        return;
      }
      const timeRange = await useTimeStore().toTimeRangeForQuery(this.range);
      let lastRange = timeRange.map(x => useTimeStore().timestamp2ymdhis(x));
      if (lastRange[0].slice(0, 10) == lastRange[1].slice(0, 10)) lastRange[1] = lastRange[1].slice(11);
      this.lastExecuted = { expr: this.expr, range: lastRange };
      this.loading = true;
      try {
        const response = await fetch(
          '/api/v1/remote/query_range?' +
            new URLSearchParams({
              logFormat: 'json',
              dstype: 'lethe',
              query: this.expr,
              start: timeRange[0],
              end: timeRange[1],
            }),
        );
        const jsonData = await response.json();
        // console.log('jsonData.data.result=', jsonData.data.result)
        this.loading = false;
        this.resultType = jsonData.data.resultType;
        this.result = jsonData.data.result;
        this.detectLogType();

        setTimeout(() => {
          if (this.$refs.logs) this.$refs.logs.scrollTop = 99999;
        }, 100);
        if (this.intervalSeconds > 0) {
          this.busy = true;
          setTimeout(() => this.timerHandler(), this.intervalSeconds * 1000);
        } else {
          this.busy = false;
        }
        this.errorResponse = null;
      } catch (error) {
        this.loading = false;
        this.errorResponse = error.response;
      }
    },
    timerHandler() {
      if (useTimeStore().timerManager != 'LogsView' || this.intervalSeconds == 0) return;
      this.execute();
    },
    selectMetric(k) {
      this.metricInfo.selected = k ? [k, this.metadata[k]] : null;
    },
    applyTarget(k) {
      this.expr = k;
    },
    clickOutside() {
      this.selectMetric(null);
    },
    async fetchQuery(query) {
      const now = await useTimeStore().getNow();
      const response = await fetch(
        '/api/v1/remote/query?' + new URLSearchParams({ dstype: 'prometheus', query: query, time: now }),
      );
      return await response.json();
    },
    async fetchNode() {
      const jsonData = await this.fetchQuery('kube_node_created');
      this.nodes = jsonData.data.result.map(x => {
        return { name: x.metric.node };
      });
    },
    async fetchNamespace() {
      const kinds = ['deployment', 'statefulset', 'daemonset', 'job', 'cronjob', 'pod'];
      const groupBy = (arr, key) =>
        arr.reduce((acc, item) => ((acc[item[key]] = [...(acc[item[key]] || []), item]), acc), {});

      let jsonData = await this.fetchQuery('kube_namespace_created');
      let namespaces = jsonData.data.result.map(x => {
        return { name: x.metric.namespace, isExpanded: false, workloads: [] };
      });

      for (const [_, kind] of kinds.entries()) {
        // console.log(`${k}: ${kind}`);
        jsonData = await this.fetchQuery(`kube_${kind}_created`);
        const objects = groupBy(
          jsonData.data.result.map(x => {
            return {
              namespace: x.metric.namespace,
              isExpanded: false,
              pods: [],
              name: x.metric[kind == 'job' ? 'job_name' : kind],
            };
          }),
          'namespace',
        );
        Object.entries(objects).forEach(x => {
          for (const [i, ns] of Object.entries(namespaces)) {
            if (ns.name == x[0])
              namespaces[i].workloads.push({
                kind: kind,
                isExpanded: false,
                objects: x[1],
              });
          }
        });
      }
      this.namespaces = namespaces;
    },
    async fetchMetadata() {
      this.fetchNode();
      this.fetchNamespace();
    },
  },
};
</script>

<template>
  <header
    class="fixed right-0 w-full bg-white border-b border-common shadow z-30 p-2 pl-52"
    :class="{ 'is-loading': loading }"
  >
    <div class="flex items-center flex-row">
      <div><i class="mdi mdi-18px mdi-text" /> Logs</div>
      <div class="flex ml-auto">
        <span>
          <TimeRangePicker @updateTimeRange="updateTimeRange" />
        </span>
        <span class="ml-2">
          <RunButton :disabled="busy" btn-text="Run query" @execute="execute" @changeInterval="changeInterval" />
        </span>
      </div>
    </div>
  </header>
  <div class="w-full flex mt-[3.6rem]">
    <div class="flex-1 py-4 px-4">
      <div class="pb-4">
        <div class="flex">
          <input
            v-model="expr"
            type="search"
            class="flex-1 relative flex-auto min-w-0 block w-full px-3 py-1.5 text-base font-normal text-gray-700 bg-white bg-clip-padding border border-solid rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none"
            placeholder="Expression"
            aria-label="Expression"
            aria-describedby="button-addon3"
            @keyup.enter="execute()"
          />
        </div>
      </div>
      <div class="break-all">
        <table class="w-full bg-slate-300">
          <tr>
            <td v-if="lastExecuted" class="px-2">{{ lastExecuted.range[0] }} - {{ lastExecuted.range[1] }}</td>
            <td v-if="lastExecuted" class="px-2">
              {{ lastExecuted.expr }}
            </td>
            <td class="px-2 font-bold text-right">{{ result.length }} rows</td>
          </tr>
        </table>
        <div v-if="errorResponse">
          <div class="rounded bg-slate-200 text-yellow-200 text-center py-20">
            [{{ errorResponse.status }}] {{ errorResponse.data.error }}
          </div>
        </div>
        <div v-else-if="result.length > 0">
          <div class="text-xs font-mono bg-white">
            <div v-for="row in result" class="border-b">
              <span class="bg-slate-100">
                <span class="mr-1 text-yellow-400">{{ Util.utc2local(row.time) }}</span>
                <template v-if="logType == 'pod'">
                  <span class="mr-1 text-green-400">{{ row.namespace }}</span>
                  <span class="mr-1 text-teal-400">{{ row.pod }}</span>
                  <span class="mr-1 text-sky-400">{{ row.container }}</span>
                </template>
                <template v-if="logType == 'node'">
                  <span class="mr-1 text-green-400">{{ row.node }}</span>
                  <span class="mr-1 text-teal-400">{{ row.process }}</span>
                </template>
              </span>
              {{ row.log }}
            </div>
          </div>
        </div>
        <table class="w-full bg-slate-300">
          <tr>
            <td v-if="lastExecuted" class="px-2">{{ lastExecuted.range[0] }} - {{ lastExecuted.range[1] }}</td>
            <td v-if="lastExecuted" class="px-2">
              {{ lastExecuted.expr }}
            </td>
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
        <div class="overflow-y-auto border-l border-2 w-full bg-slate-200 border-b" style="height: calc(100vh - 90px)">
          <div class="bg-gray-100 text-center cursor-pointer">event</div>
          <div class="cursor-pointer" @click="selectEvent">event</div>
          <div class="bg-gray-100 text-center cursor-pointer">node</div>
          <div v-for="node in nodes" class="cursor-pointer" @click="selectNode(node.name)">
            {{ node.name }}
          </div>
          <div class="bg-gray-100 text-center cursor-pointer">pod</div>
          <div v-for="ns in namespaces">
            <div class="cursor-pointer" @click="ns.isExpanded = !ns.isExpanded">
              <i class="mdi" :class="[ns.isExpanded ? 'mdi-chevron-down' : 'mdi-chevron-right']" />
              {{ ns.name }}
            </div>
            <template v-if="ns.isExpanded">
              <div v-for="w in ns.workloads">
                <div class="pl-2 cursor-pointer" @click="w.isExpanded = !w.isExpanded">
                  <i class="mdi" :class="[w.isExpanded ? 'mdi-chevron-down' : 'mdi-chevron-right']" />
                  {{ w.kind }} ({{ w.objects.length }})
                </div>
                <template v-if="w.isExpanded">
                  <div
                    v-for="object in w.objects"
                    class="pl-5 cursor-pointer"
                    @click="selectObject(w.kind, ns.name, object.name)"
                  >
                    {{ object.name }}
                  </div>
                </template>
              </div>
            </template>
          </div>
          <div class="pb-16">
            <hr />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
