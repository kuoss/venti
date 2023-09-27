<script setup>
import Util from '@/lib/util';
</script>

<script>
export default {
  data() {
    return {
      alertFiles: [],
      isLoading: false,
      repeat: true,
      testAlertSent: false,
    };
  },
  mounted() {
    this.fetchData();
  },
  beforeUnmount() {
    this.repeat = false;
  },
  methods: {
    async fetchData() {
      this.isLoading = true;
      try {
        const response = await fetch('/api/v1/alerts');
        const jsonData = await response.json();
        this.alertFiles = jsonData;
        setTimeout(() => {
          if (!this.repeat) return;
          this.fetchData();
        }, 3000);
      } catch (error) {
        this.repeat = false;
        console.error(error);
      }
      this.isLoading = false;
    },
    async sendTestAlert() {
      try {
        const response = await fetch('/api/v1/alerts/test');
        const jsonData = await response.json();
        console.log(jsonData)
      } catch (error) {
        console.error(error);
      }
      this.testAlertSent = true;
    }
  },
};
</script>

<template>
  <div>
    <header class="fixed right-0 w-full bg-white border-b shadow z-30 p-2 pl-52" :class="{ 'is-loading': isLoading }">
      <div class="flex items-center flex-row">
        <div><i class="mdi mdi-18px mdi-database-outline" /> Alert</div>
        <div class="flex ml-auto">
          <div class="inline-flex">
            <button class="h-rounded-group py-2 px-4 text-gray-900 bg-white border border-common" v-if="!testAlertSent" @click="sendTestAlert">
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
      <h1 class="py-2 font-bold">Alert Rule Files ({{ alertFiles.length }} files)</h1>
      <table class="w-full bg-white border">
        <tr class="border-b bg-slate-50">
          <th class="text-left px-2">State</th>
          <th class="text-left px-2">Severity</th>
          <th class="text-left px-2">Name</th>
          <th class="text-left px-2">Summary</th>
          <th class="text-left px-2">Expr</th>
          <th class="text-left px-2">For</th>
        </tr>
        <tbody v-for="f in alertFiles">
          <tr class="border-b">
            <th class="text-left px-2 bg-slate-300 p-1 pl-3" colspan="6">
              {{ f.datasourceSelector.type == 'prometheus' ? '🔥' : '💧' }} file ({{ f.groups.length }} groups)
            </th>
          </tr>
          <template v-for="g in f.groups">
            <tr class="border-b">
              <th class="text-left px-2 bg-slate-200 p-1 pl-6" colspan="6">group: {{ g.name }} ({{ g.ruleAlerts.length }} rules)</th>
            </tr>
            <tr v-for="ra in g.ruleAlerts" class="border-b">
              <td class="p-1 pl-9">
                <span v-for="a in ra.alerts">
                  <div v-if="a.state == 'firing'" class="rounded w-24 text-center bg-red-300">firing</div>
                  <div v-else-if="a.state == 'pending'" class="rounded w-24 text-center bg-yellow-300">pending</div>
                  <div v-else class="rounded w-24 text-center bg-green-300">normal</div>
                </span>
              </td>
              <td class="px-2">
                {{ f.commonLabels?.severity }}
              </td>
              <td class="px-2">
                {{ ra.rule.alert }}
              </td>
              <td class="px-2 max-w-[30vw] truncate hover:whitespace-normal">
                {{ ra.rule.annotations?.summary }}
              </td>
              <td class="px-2 max-w-[20vw] truncate hover:whitespace-normal text-cyan-500">
                <a
                  class="hover:underline"
                  :href="`/${f.datasourceSelector.type == 'prometheus' ? 'metrics' : 'logs'}?query=${encodeURIComponent(
                    ra.rule.expr,
                  )}`"
                  >{{ ra.rule.expr }}</a
                >
              </td>
              <td class="px-2">
                {{ Util.nanoseconds2human(ra.rule.for) }}
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </main>
  </div>
</template>
