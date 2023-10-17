<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import Util from '@/lib/util'

interface Alert {
  annotations: Record<string, string>
  labels: Record<string, string>
  state: number
  createdAt: string
  updatedAt: string
}

interface Rule {
  alert: string
  annotations: Record<string, string>
  labels: Record<string, string>
  expr: string
  for: number
}

interface AlertingRule {
  active: Record<string, Alert>
  rule: Rule
}

interface DatasourceSelector {
  system: string
  type: string
}

interface AlertingFile {
  datasourceSelector: DatasourceSelector
  alertingRules: AlertingRule[]
  groupLabels: Record<string, string>
}

const alertingFiles = ref([] as AlertingFile[])
const isLoading = ref(false)
const repeat = ref(true)
const testAlertSent = ref(false)

async function fetchData() {
  isLoading.value = true
  try {
    const resp = await fetch('/api/v1/alerts')
    const json = await resp.json()
    alertingFiles.value = json.data
    setTimeout(() => {
      if (!repeat.value) return
      fetchData()
    }, 3000);
  } catch (err) {
    repeat.value = false
    console.error(err)
  }
  isLoading.value = false
}
async function sendTestAlert() {
  try {
    const resp = await fetch('/api/v1/alerts/test')
    const json = await resp.json()
    console.log(json)
  } catch (err) {
    console.error(err)
  }
  testAlertSent.value = true
}

onMounted(() => {
  fetchData()
})

onUnmounted(() => {
  repeat.value = false
})

</script>

<template>
  <div>
    <header class="fixed right-0 w-full bg-white border-b shadow z-30 p-2 pl-52" :class="{ 'is-loading': isLoading }">
      <div class="flex items-center flex-row">
        <div><i class="mdi mdi-18px mdi-database-outline" /> Alert</div>
        <div class="flex ml-auto">
          <div class="inline-flex">
            <button class="h-rounded-group py-2 px-4 text-gray-900 bg-white border border-common" v-if="!testAlertSent"
              @click="sendTestAlert">
              <i class="mdi mdi-cube-send" /> Send Test Alert
            </button>
            <button class="h-rounded-group py-2 px-4 text-gray-900 bg-white border border-common">
              <i class="mdi mdi-refresh mdi-spin" />
            </button>
          </div>
        </div>
      </div>
    </header>

    <main class="mt-12 w-full p-8 pb-16">
      <h1 class="py-2 font-bold">Alerting Files ({{ alertingFiles.length }} files)</h1>
      <table class="w-full bg-white border">
        <tr class="border-b bg-slate-50">
          <th class="text-left px-2">State</th>
          <th class="text-left px-2">Severity</th>
          <th class="text-left px-2">Name</th>
          <th class="text-left px-2">Summary</th>
          <th class="text-left px-2">Expr</th>
          <th class="text-left px-2">For</th>
        </tr>

        <tbody v-for="f in alertingFiles">
          <tr class="border-b">
            <th class="text-left px-2 bg-slate-300 p-1 pl-3" colspan="9">
              {{ f.datasourceSelector.type == 'prometheus' ? '🔥' : '💧' }} file ({{ f.alertingRules.length }} rules)
            </th>
          </tr>
          <template v-for="r in f.alertingRules">
            <tr class="border-b">
              <td class="text-center">
                <span v-if="r.active" class="px-2 rounded-full bg-red-400">
                  {{ Object.keys(r.active).length }}
                </span>
                <span v-else>
                  ·
                </span>
              </td>
              <td class="px-2">
                {{ f.groupLabels.severity }}
              </td>
              <td class="px-2">
                {{ r.rule.alert }}
              </td>
              <td class="px-2">
                {{ r.rule.annotations.summary }}
              </td>
              <td class="px-2">
                {{ r.rule.expr }}
              </td>
              <td class="px-2">
                {{ r.rule.for/100000000 }}s
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </main>
  </div>
</template>
