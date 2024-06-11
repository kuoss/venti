<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useDateFormat, useTimeAgo } from '@vueuse/core';
import ProgressBar from '@/components/ProgressBar.vue';

interface Alert {
  annotations: Record<string, string>;
  labels: Record<string, string>;
  state: number;
  createdAt: string;
  updatedAt: string;
}

interface Rule {
  alert: string;
  annotations: Record<string, string>;
  labels: Record<string, string>;
  expr: string;
  for: number;
}

interface AlertingRule {
  active: Record<string, Alert>;
  rule: Rule;
}

interface DatasourceSelector {
  system: string;
  type: string;
}

interface AlertingFile {
  datasourceSelector: DatasourceSelector;
  alertingRules: AlertingRule[];
  groupLabels: Record<string, string>;
}

const alertingFiles = ref([] as AlertingFile[]);
const isLoading = ref(false);
const repeat = ref(true);
const testAlertSent = ref(false);

async function fetchData() {
  isLoading.value = true;
  try {
    const resp = await fetch('/api/v1/alerts');
    const json = await resp.json();
    alertingFiles.value = json.data;
    setTimeout(() => {
      if (!repeat.value) return;
      fetchData();
    }, 3000);
  } catch (err) {
    repeat.value = false;
    console.error(err);
  }
  isLoading.value = false;
}
async function sendTestAlert() {
  try {
    const resp = await fetch('/api/v1/alerts/test');
    const json = await resp.json();
    console.log(json);
  } catch (err) {
    console.error(err);
  }
  testAlertSent.value = true;
}

function filterRecord(record: Record<string, string>, unwantedKeys: string[]) {
  let out = {} as Record<string, string>;
  for (const k in record) {
    if (unwantedKeys.includes(k)) {
      continue;
    }
    out[k] = record[k];
  }
  return out;
}

function filterGroupLabels(labels: Record<string, string>) {
  return filterRecord(labels, ['rulefile']);
}

function filterLabels(labels: Record<string, string>) {
  return filterRecord(labels, ['alertname', 'datasource', 'rulefile', 'severity', 'venti']);
}

onMounted(() => {
  fetchData();
});

onUnmounted(() => {
  repeat.value = false;
});
</script>

<template>
  <div>
    <header class="fixed right-0 w-full bg-white dark:bg-black border-b shadow z-30">
      <div class="flex items-center flex-row p-2 pb-1">
        <div class="pl-52"><i class="mdi mdi-18px mdi-database-outline" /> Alert</div>
        <div class="flex ml-auto">
          <div class="inline-flex">
            <button
              class="h-rounded-group py-2 px-4 text-gray-900 dark:text-gray-100 bg-white dark:bg-black border border-common"
              v-if="!testAlertSent"
              @click="sendTestAlert"
            >
              <i class="mdi mdi-cube-send" /> Send Test Alert
            </button>
            <button class="h-rounded-group py-2 px-4 text-gray-900 dark:text-gray-100 bg-white dark:bg-black border border-common">
              <i class="mdi mdi-refresh mdi-spin" />
            </button>
          </div>
        </div>
      </div>
      <div class="h-1">
        <Transition name="fade">
          <div v-if="isLoading">
            <ProgressBar />
          </div>
        </Transition>
      </div>
    </header>

    <main class="mt-12 w-full p-8 pb-16">
      <h1 class="py-2 font-bold">Alerting Files ({{ alertingFiles.length }} files)</h1>
      <table class="table1 w-full bg-slate-200 dark:bg-slate-800 border">
        <tr class="border-b bg-slate-50 dark:bg-slate-900">
          <th class="text-left">State</th>
          <th class="text-left">Severity</th>
          <th class="text-left">Name</th>
          <th class="text-left">Summary</th>
          <th class="text-left">Expr</th>
          <th class="text-left">For</th>
        </tr>

        <tbody v-for="f in alertingFiles">
          <tr class="border-t">
            <th class="text-left px-2 bg-slate-300 dark:bg-slate-700 p-1 pl-3" colspan="9">
              {{ f.groupLabels['rulefile'] }}
              ({{ f.datasourceSelector.type == 'prometheus' ? '🔥' : '💧' }}{{ f.datasourceSelector.system }}
              {{ f.alertingRules.length }} rules)
              <span class="bg-slate-200 dark:bg-slate-800 text-xs px-2 rounded-full" v-for="(v, k) in filterGroupLabels(f.groupLabels)">
                {{ k }}: {{ v }}
              </span>
            </th>
          </tr>
          <template v-for="(r, idx) in f.alertingRules">
            <tr class="border-t">
              <td class="text-center bg-red-400 dark:bg-red-600" v-if="r.active">
                {{ Object.keys(r.active).length }}
              </td>
              <td class="text-center bg-green-400 dark:bg-green-600" v-else>0</td>
              <td>
                {{ f.groupLabels.severity }}
              </td>
              <td>
                {{ r.rule.alert }}
              </td>
              <td>
                {{ r.rule.annotations.summary }}
              </td>
              <td>
                {{ r.rule.expr }}
              </td>
              <td>{{ r.rule.for / 1000 / 1000 / 1000 }}s</td>
            </tr>
            <tr class="bg-gray-100 dark:bg-gray-900" v-for="(alert, k) in r.active">
              <td colspan="2">&nbsp;</td>
              <td>
                {{ alert.labels['datasource'] }}
              </td>
              <td>
                {{ alert.annotations['summary'] }}
              </td>
              <td>
                <span class="bg-slate-200 dark:bg-slate-800 text-xs mr-2 px-2 rounded-full" v-for="(v, k) in filterLabels(alert.labels)">
                  {{ k }}: {{ v }}
                </span>
              </td>
              <td>
                {{ useTimeAgo(alert.createdAt).value }}
                ({{ useDateFormat(alert.createdAt, 'YYYY-MM-DD HH:mm:ss', { locales: 'ko-KR' }).value }} KST)
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </main>
  </div>
</template>

<style scoped>
.table1 th,
.table1 td {
  @apply px-2;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 1.8s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
