<script setup>
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { useTimeStore } from '@/stores/time';
import DropdownDatasource from '@/components/DropdownDatasource.vue'
import TimeRangePicker from '@/components/TimeRangePicker.vue';
import RunButton from '@/components/RunButton.vue';
import Util from '@/lib/util';

const route = useRoute();
const timeStore = useTimeStore()

const errorResponse = ref(null)
const busy = ref(false)
const loading = ref(false)
const intervalSeconds = ref(0)
const range = ref([])
const nodes = ref([])
const namespaces = ref([])
const lastExecuted = ref(null)
const logType = ref('')
const expr = ref('pod{namespace="kube-system"}')
const result = ref([])
const resultType = ref(null)

let dsName

function detectLogType() {
  if (result.value.length < 1) {
    logType.value = '';
    return;
  }
  if (Object.prototype.hasOwnProperty.call(result.value[0], 'namespace')) {
    logType.value = 'pod';
    return;
  }
  if (Object.prototype.hasOwnProperty.call(result.value[0], 'node')) {
    logType.value = 'node';
    return;
  }
  logType.value = '';
}

function selectObject(kind, namespace, name) {
  if (kind == 'pod') {
    expr.value = `pod{namespace="${namespace}", pod="${name}"}`;
  } else {
    expr.value = `pod{namespace="${namespace}", pod=~"${name}-.*"}`;
  }
}

function changeInterval(i) {
  intervalSeconds.value = i;
  execute();
}

function updateTimeRange(r) {
  range.value = r;
}

function selectNode(name) {
  expr.value = `node{node="${name}"}`;
}

async function execute() {
  if (expr.value.length < 1) {
    console.error('emtpy expr');
    return;
  }
  const timeRange = await timeStore.toTimeRangeForQuery(range.value);
  let lastRange = timeRange.map(x => timeStore.timestamp2ymdhis(x));
  if (lastRange[0].slice(0, 10) == lastRange[1].slice(0, 10)) lastRange[1] = lastRange[1].slice(11);
  lastExecuted.value = { expr: expr.value, range: lastRange };
  loading.value = true;
  try {
    const response = await fetch(
      '/api/v1/remote/query_range?' +
      new URLSearchParams({
        logFormat: 'json',
        dsName: dsName,
        query: expr.value,
        start: timeRange[0],
        end: timeRange[1],
      }),
    );
    const jsonData = await response.json();

    loading.value = false;
    resultType.value = jsonData.data.resultType;
    result.value = jsonData.data.result;
    detectLogType();

    setTimeout(() => {
      window.scrollTo(0, 99999);
    }, 100);

    if (intervalSeconds.value > 0) {
      busy.value = true;
      setTimeout(() => timerHandler(), intervalSeconds.value * 1000);
    } else {
      busy.value = false;
    }
    errorResponse.value = null;
  } catch (error) {
    loading.value = false;
    errorResponse.value = error.response;
  }
}

function timerHandler() {
  if (timeStore.timerManager != 'LogsView' || intervalSeconds.value == 0) return;
  execute();
}

async function fetchQuery(query) {
  const now = await useTimeStore().getNow();
  const response = await fetch(
    '/api/v1/remote/query?' + new URLSearchParams({ dsType: 'prometheus', query: query, time: now }),
  );
  return await response.json();
}

async function fetchNodes() {
  const jsonData = await fetchQuery('kube_node_created');
  nodes.value = jsonData.data.result.map(x => {
    return { name: x.metric.node };
  });
}

async function fetchNamespaces() {
  const kinds = ['deployment', 'statefulset', 'daemonset', 'job', 'cronjob', 'pod'];
  const groupBy = (arr, key) =>
    arr.reduce((acc, item) => ((acc[item[key]] = [...(acc[item[key]] || []), item]), acc), {});

  let jsonData = await fetchQuery('kube_namespace_created');
  let nss = jsonData.data.result.map(x => {
    return { name: x.metric.namespace, isExpanded: false, workloads: [] };
  });

  for (const [_, kind] of kinds.entries()) {
    // console.log(`${k}: ${kind}`);
    jsonData = await fetchQuery(`kube_${kind}_created`);
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
      for (const [i, ns] of Object.entries(nss)) {
        if (ns.name == x[0])
          nss[i].workloads.push({
            kind: kind,
            isExpanded: false,
            objects: x[1],
          });
      }
    });
  }
  namespaces.value = nss;
}

function onChangeDatasource(value) {
  dsName = value
}

onMounted(() => {
  timeStore.timerManager = 'LogsView';
  fetchNodes();
  fetchNamespaces();
  if (route.query?.query) {
    expr.value = route.query.query;
    setTimeout(execute, 500);
  }
})
</script>

<template>
  <header class="fixed right-0 w-full bg-white dark:bg-black border-b border-common shadow z-30 p-2 pl-52"
    :class="{ 'is-loading': loading }">
    <div class="flex items-center flex-row">
      <div>
        <i class="mdi mdi-18px mdi-text" /> Logs
        <DropdownDatasource dsType="lethe" @change="onChangeDatasource" />
      </div>
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
          <input v-model="expr" type="search"
            class="flex-auto relative min-w-0 block w-full px-3 py-1.5 text-base font-normal text-gray-700 dark:text-gray-300 bg-white dark:bg-black bg-clip-padding border border-solid rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none"
            placeholder="Expression" aria-label="Expression" aria-describedby="button-addon3" @keyup.enter="execute()" />
        </div>
      </div>
      <div class="break-all">
        <table class="w-full bg-slate-300 dark:bg-slate-700">
          <tr>
            <td v-if="lastExecuted" class="px-2">{{ lastExecuted.range[0] }} - {{ lastExecuted.range[1] }}</td>
            <td v-if="lastExecuted" class="px-2">
              {{ lastExecuted.expr }}
            </td>
            <td class="px-2 font-bold text-right">{{ result.length }} rows</td>
          </tr>
        </table>
        <div v-if="errorResponse">
          <div class="rounded bg-slate-200 dark:bg-slate-800 text-yellow-200 dark:text-yellow-800 text-center py-20">
            [{{ errorResponse.status }}] {{ errorResponse.data.error }}
          </div>
        </div>
        <div v-else-if="result.length > 0">
          <div class="text-xs font-mono bg-white dark:bg-black">
            <div v-for="row in result" class="border-b">
              <span class="bg-slate-100 dark:bg-slate-900">
                <span class="mr-1 text-yellow-500">{{ Util.utc2local(row.time) }}</span>
                <template v-if="logType == 'pod'">
                  <span class="mr-1 text-green-500">{{ row.namespace }}</span>
                  <span class="mr-1 text-teal-500">{{ row.pod }}</span>
                  <span class="mr-1 text-sky-500">{{ row.container }}</span>
                </template>
                <template v-if="logType == 'node'">
                  <span class="mr-1 text-green-500">{{ row.node }}</span>
                  <span class="mr-1 text-teal-500">{{ row.process }}</span>
                </template>
              </span>
              {{ row.log }}
            </div>
          </div>
        </div>
        <table class="w-full bg-slate-300 dark:bg-slate-700">
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
      <div class="fixed right-0 bottom-0 bg-slate-300 dark:bg-slate-700 text-xs pt-4 w-80">
        <div>
          <div class="w-full py-2 text-center">Targets</div>
        </div>
        <div class="overflow-y-auto border-l border-2 w-full bg-slate-200 dark:bg-slate-800 border-b" style="height: calc(100vh - 90px)">
          <div class="bg-gray-100 dark:bg-gray-900 text-center cursor-pointer">node</div>
          <div v-for="node in nodes" class="cursor-pointer" @click="selectNode(node.name)">
            {{ node.name }}
          </div>
          <div class="bg-gray-100 dark:bg-gray-900 text-center cursor-pointer">pod</div>
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
                  <div v-for="object in w.objects" class="pl-5 cursor-pointer"
                    @click="selectObject(w.kind, ns.name, object.name)">
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
