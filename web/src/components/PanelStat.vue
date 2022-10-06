<template>
    <div class="text-center py-1" :class="thresholdClass">{{ value }}</div>
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
    watch: {
        count() { if (!this.isLoading) this.fetchData() },
    },
    data() {
        return {
            value: '',
            thresholdClass: '',
        }
    },
    methods: {
        async fetchData() {
            if (this.timeRange.length < 2) return
            this.$emit('setIsLoading', true)
            try {
                const target = this.panelConfig.targets[0]
                const response = await this.axios.get('/api/prometheus/query', {
                    params: {
                        expr: this.panelConfig.targets[0].expr,
                        time: this.timeRange[1],
                    }
                })
                const result = response.data.data.result
                const resultType = response.data.data.resultType

                let value
                if (resultType == 'scalar') value = result[1]
                else {
                    // vector
                    if (target.legend) value = result.map(x => target.legend.replace(/\{\{(.*?)\}\}/g, (i, m) => x.metric[m]))[0]
                    else value = result[0].value[1]
                }
                // unit
                if (target.unit == "dateTimeAsLocal") value = this.$util.dateTimeAsLocal(value)
                this.value = value

                // thresholds
                let level = 0
                if (target.thresholds) {
                    level = 1
                    target.thresholds.forEach((threshold, i) => {
                        if (threshold.values) {
                            // console.log('threshold.values=', threshold.values)
                            threshold.values.forEach((v, i) => {
                                if ((threshold.invert && value < v) || (!threshold.invert && value > v)) {
                                    level += threshold.values.length - 1
                                }
                            })
                        }
                    })
                }
                // console.log('level=', level)
                this.thresholdClass = ['', 'bg-green-100', 'bg-orange-100', 'bg-red-100'][level]
            } catch (error) { console.error(error) }
            this.$emit('setIsLoading', false)
        },
    },
    mounted() {
        this.fetchData()
    }
}
</script>
