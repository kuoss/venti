<script setup>
import { useTimeStore } from '@/stores/time';
</script>
<script>
export default {
  props: {
    count: Number,
    isLoading: Boolean,
    panelConfig: Object,
    panelWidth: Number,
    timeRange: Object,
  },
  data() {
    return {
      loading: false,
      rows: [],
    };
  },
  watch: {
    count() {
      if (!this.isLoading) this.fetchData();
    },
  },
  mounted() {
    this.fetchData();
  },
  methods: {
    async fetchData() {
      if (this.timeRange.length < 2) return;
      const targets = this.panelConfig.targets;
      let merged = {};
      this.$emit('setIsLoading', true);
      for (const i in targets) {
        try {
          const target = targets[i];

          const response = await fetch(
            '/api/v1/remote/query?' +
              new URLSearchParams({ dsType: 'prometheus', query: target.expr, time: this.timeRange[1] }),
          );
          const jsonData = await response.json();
          let rows = {};
          jsonData.data.result.forEach(x => {
            const key = x.metric[target.key];
            let row = {};
            target.columns.forEach(c => {
              if (c == 'VALUE') row['VALUE' + i] = x.value[1];
              else row[c] = x.metric[c];
            });
            rows[key] = row;
          });
          Object.entries(rows).forEach(r => {
            const k = r[0];
            merged[k] = { ...merged[k], ...rows[k] };
          });
        } catch (error) {
          console.error(error);
        }
      }
      this.$emit('setIsLoading', false);
      this.rows = merged;
    },
  },
  timerHandler() {
    if (useTimeStore().timerManager != this.timerManager || this.intervalSeconds == 0) return;
    this.execute();
  },
};
</script>

<template>
  <table class="w-full border-collapse">
    <tr>
      <template v-for="t in panelConfig.targets">
        <th v-for="h in t.headers" class="border px-2 py-1 bg-slate-50">
          {{ h }}
        </th>
      </template>
    </tr>
    <tr v-for="row in rows">
      <template v-for="(t, i) in panelConfig.targets">
        <td v-for="c in t.columns" class="border px-2 py-1">
          {{ c == 'VALUE' ? row['VALUE' + i] : row[c] }}
        </td>
      </template>
    </tr>
  </table>
</template>
