<template>
  <table class="w-full border-collapse">
    <tr class="border-b">
      <th v-for="h in panelConfig.headers" class="px-2 py-1 bg-slate-50">
        {{ h }}
      </th>
    </tr>
    <tr v-for="(row, i) in rows" class="border-b">
      <td class="px-2 py-1">
        {{ datasources[i].host }}
      </td>
      <td v-for="cell in row" class="border-l px-2 py-1" :class="{ 'text-right': !isNaN(cell) }">
        {{ cell }}
      </td>
    </tr>
  </table>
</template>

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
      rows: [],
      datasources: [],
    };
  },
  watch: {
    count() {
      if (!this.loading) this.fetchData();
    },
  },
  mounted() {
    this.fetchDatasources();
  },
  methods: {
    async fetchDatasources() {
      try {
        const response = await fetch('/api/v1/datasources');
        const jsonData = await response.json();
        this.datasources = jsonData;
      } catch (error) {
        console.error(error);
      }
      this.fetchData();
    },
    async fetchData() {
      if (this.timeRange.length < 2) return;
      this.$emit('setIsLoading', true);
      try {
        let rows = [];
        for (let i = 0; i < this.datasources.length; i++) {
          const ds = this.datasources[i];
          if (ds.type != 'prometheus' || !ds.isDiscovered) continue;

          let row = [];
          for (const target of this.panelConfig.targets) {
            const response = await fetch(
              '/api/v1/remote/query?' + new URLSearchParams({ dsid: i, query: target.expr, time: this.timeRange[1] }),
            );
            const jsonData = await response.json();
            const result = jsonData.data.result;
            if (!target.legends) {
              row.push(result[0].value[1]);
              continue;
            }
            for (const legend of target.legends) {
              const value = result.map(x => legend.replace(/\{\{(.*?)\}\}/g, (i, m) => x.metric[m]))[0];
              row.push(value);
            }
          }
          rows.push(row);
        }
        this.rows = rows;
      } catch (error) {
        console.error(error);
      }
      this.$emit('setIsLoading', false);
    },
  },
};
</script>
