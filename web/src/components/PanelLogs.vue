<script setup>
import { useFilterStore } from '@/stores/filter';
import Util from '@/lib/util';
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
      logType: '',
      total: '-',
      result: [],
      resultType: null,
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
    async fetchData() {
      if (this.timeRange.length < 2) return;
      this.$emit('setIsLoading', true);
      try {
        const query = useFilterStore().renderExpr(this.panelConfig.targets[0].expr);
        const response = await fetch(
          '/api/v1/remote/query_range?' +
            new URLSearchParams({
              logFormat: 'json',
              dsType: 'lethe',
              query: query,
              start: this.timeRange[0],
              end: this.timeRange[1],
            }),
        );
        const jsonData = await response.json();

        this.result = jsonData.data.result.slice(-100);
        this.resultType = jsonData.data.resultType;
        this.detectLogType();
        // console.log('jsonData=', jsonData);
        setTimeout(() => {
          if (this.$refs.logs) this.$refs.logs.scrollTop = 99999;
        }, 100);
      } catch (error) {
        console.error(error);
      }
      this.$emit('setIsLoading', false);
    },
  },
};
</script>

<template>
  <div ref="logs" class="text-xs font-mono h-64 break-all overflow-y-auto">
    <template v-if="result.length > 0">
      <template v-if="logType == 'pod'">
        <div v-for="row in result" class="border-b">
          <span class="bg-slate-100">
            <span class="mr-1 text-yellow-400">{{ Util.utc2local(row.time) }}</span>
            <span class="mr-1 text-green-400">{{ row.namespace }}</span>
            <span class="mr-1 text-teal-400">{{ row.pod }}</span>
            <span class="mr-1 text-sky-400">{{ row.container }}</span>
          </span>
          {{ row.log }}
        </div>
      </template>
      <template v-if="logType == 'node'">
        <div v-for="row in result" class="border-b">
          <span class="bg-slate-100">
            <span class="mr-1 text-yellow-400">{{ Util.utc2local(row.time) }}</span>
            <span class="mr-1 text-green-400">{{ row.node }}</span>
            <span class="mr-1 text-teal-400">{{ row.process }}</span>
          </span>
          {{ row.log }}
        </div>
      </template>
    </template>
    <template v-else>
      <div class="text-center p-5">No data</div>
    </template>
  </div>
</template>
