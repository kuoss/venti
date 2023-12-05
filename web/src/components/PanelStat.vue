<script setup>
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
      value: '',
      thresholdClass: '',
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
      this.$emit('setIsLoading', true);
      try {
        const target = this.panelConfig.targets[0];
        const response = await fetch(
          '/api/v1/remote/query?' +
            new URLSearchParams({ dsType: 'prometheus', query: target.expr, time: this.timeRange[1] }),
        );
        const jsonData = await response.json();

        const result = jsonData.data.result;
        const resultType = jsonData.data.resultType;

        let value;
        if (resultType == 'scalar') value = result[1];
        else {
          // vector
          if (target.legend) value = result.map(x => target.legend.replace(/\{\{(.*?)\}\}/g, (i, m) => x.metric[m]))[0];
          else value = result[0].value[1];
        }
        // unit
        if (target.unit == 'dateTimeAsLocal') value = Util.dateTimeAsLocal(value);
        this.value = value;

        // thresholds
        let level = 0;
        if (target.thresholds) {
          level = 1;
          target.thresholds.forEach(threshold => {
            if (threshold.values) {
              // console.log('threshold.values=', threshold.values)
              threshold.values.forEach(v => {
                if ((threshold.invert && value < v) || (!threshold.invert && value > v)) {
                  level += threshold.values.length - 1;
                }
              });
            }
          });
        }
        // console.log('level=', level)
        this.thresholdClass = ['', 'bg-green-100 dark:bg-green-800', 'bg-orange-100 dark:bg-orange-800', 'bg-red-100 dark:bg-red-800'][level];
      } catch (error) {
        console.error(error);
      }
      this.$emit('setIsLoading', false);
    },
  },
};
</script>

<template>
  <div class="text-center py-1" :class="thresholdClass">
    {{ value }}
  </div>
</template>
