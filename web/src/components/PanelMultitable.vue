<template>
    <table class="w-full border-collapse">
        <tr class="border-b">
            <th class="px-2 py-1 bg-slate-50" v-for="h in panelConfig.headers">{{ h }}</th>
        </tr>
        <tr class="border-b" v-for="(row,i) in rows">
            <td class="px-2 py-1">{{ datasources[i].host }}</td>
            <td
                class="border-l px-2 py-1"
                v-for="cell in row"
                :class="{ 'text-right': !isNaN(cell) }"
            >{{ cell }}</td>
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
    watch: {
        count() { if (!this.loading) this.fetchData() },
    },
    data() {
        return {
            rows: [],
            datasources: [],
        }
    },
    methods: {
        async init() {
            try {
                const response = await this.axios.get('/api/datasources')
                this.datasources = response.data.filter(x => x.type == 'Prometheus' && x.is_discovered)
            } catch (error) { console.error(error) }
            this.fetchData()
        },
        async fetchData() {
            if (this.timeRange.length < 2) return
            this.$emit('setIsLoading', true)
            try {
                let rows = []
                for (const datasource of this.datasources) {
                    let row = []
                    for (const target of this.panelConfig.targets) {
                        const response = await this.axios.get('/api/prometheus/query', {
                            params: {
                                host: datasource.host,
                                expr: target.expr,
                                time: this.timeRange[1],
                            }
                        })
                        const result = response.data.data.result
                        if (!target.legends) {
                            row.push(result[0].value[1])
                            continue
                        }
                        for (const legend of target.legends) {
                            const value = result.map(x => legend.replace(/\{\{(.*?)\}\}/g, (i, m) => x.metric[m]))[0]
                            row.push(value)
                        }
                    }
                    rows.push(row)
                }
                this.rows = rows
            } catch (error) { console.error(error) }
            this.$emit('setIsLoading', false)
        },
    },
    mounted() {
        this.init()
    }
}
</script>
