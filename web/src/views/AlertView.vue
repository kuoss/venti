<script setup>
import { useTimeStore } from "@/stores/time";
</script>

<script>
export default {
  data() {
    return {
      alertRuleFiles: [],
      isLoading: false,
      repeat: true,
    };
  },
  methods: {
    async fetchData() {
      this.isLoading = true;
      try {
        const response = await fetch("/api/alerts");
        const data = await response.json();
        this.alertRuleFiles = data;
        setTimeout(() => {
          if (!this.repeat) return;
          this.fetchData()
        }, 3000);
      } catch (error) {
        this.repeat = false;
        console.error(error);
      }
      this.isLoading = false;
    },
  },
  mounted() {
    this.fetchData();
  },
  beforeUnmount() {
    this.repeat = false;
  }
};
</script>

<template>
  <div>
    <header class="fixed right-0 w-full bg-white border-b shadow z-30 p-2 pl-52" :class="{ 'is-loading': isLoading }">
      <div class="flex items-center flex-row">
        <div><i class="mdi mdi-18px mdi-database-outline"></i> Alert</div>
        <div class="flex ml-auto">
          <div class="inline-flex">
            <button class="h-rounded-group py-2 px-4 text-gray-900 bg-white border border-common">
              <i class="mdi mdi-refresh mdi-spin"></i>
            </button>
          </div>
        </div>
      </div>
    </header>

    <main class="mt-12 w-full p-8 pb-16">
      <h1 class="py-2 font-bold">
        Alert Rule Files ({{ alertRuleFiles.length }} files)
      </h1>
      <table class="w-full bg-white border">
        <tr class="border-b bg-slate-50">
          <th>State</th>
          <th>Severity</th>
          <th>Name</th>
          <th>Summary</th>
          <th>Expr</th>
          <th>For</th>
        </tr>
        <tbody v-for="(f, i) in alertRuleFiles">
          <tr class="border-b">
            <th class="bg-slate-300 p-1 pl-3" colspan="6">
              {{ f.datasourceSelector.type == "prometheus" ? "🔥" : "💧" }} file ({{ f.groups.length }} groups)
            </th>
          </tr>
          <template v-for="g in f.groups">
            <tr class="border-b">
              <th class="bg-slate-200 p-1 pl-6" colspan="6">
                group: {{ g.name }} ({{ g.rules.length }} rules)
              </th>
            </tr>
            <tr class="border-b" v-for="r in g.rules">
              <td class="p-1 pl-9">
                <div class="rounded w-24 text-center bg-red-300" v-if="r.state == 'firing'">
                  firing
                </div>
                <div class="rounded w-24 text-center bg-yellow-300" v-else-if="r.state == 'pending'">
                  pending
                </div>
                <div class="rounded w-24 text-center bg-green-300" v-else>
                  normal
                </div>
              </td>
              <td class="px-2">
                {{ f.commonLabels.severity }}
              </td>
              <td class="px-2">{{ r.alert }}</td>
              <td class="px-2 max-w-[30vw] truncate hover:whitespace-normal">
                {{ r.annotations.summary }}
              </td>
              <td class="px-2 max-w-[20vw] truncate hover:whitespace-normal text-cyan-500">
                <a class="hover:underline" :href="`/${f.datasourceSelector.type == 'prometheus' ? 'metrics' : 'logs'
                  }?query=${encodeURIComponent(r.expr)}`">{{ r.expr }}</a>
              </td>
              <td class="px-2">{{ $util.nanoseconds2human(r.for) }}</td>
            </tr>
          </template>
        </tbody>
      </table>
    </main>
  </div>
</template>

<style>
th {
  @apply text-left px-2;
}
</style>