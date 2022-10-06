<template>
    <table class="w-full border-collapse">
        <tr>
            <template v-for="t in panelConfig.targets">
                <th class="border px-2 py-1 bg-slate-50" v-for="h in t.headers">{{ h }}</th>
            </template>
        </tr>
        <tr v-for="row in rows">
            <template v-for="(t, i) in panelConfig.targets">
                <td
                    class="border px-2 py-1"
                    v-for="c in t.columns"
                >{{ c == 'VALUE' ? row['VALUE' + i] : row[c] }}</td>
            </template>
        </tr>
    </table>
</template>

<script>
import { useTimeStore } from "@/stores/time"
export default {
    props: {
        count: Number,
        isLoading: Boolean,
        panelConfig: Object,
        panelWidth: Number,
        timeRange: Object,
    },
    watch: {
        count() {
            if (!this.isLoading) this.fetchData()
        },
    },
    data() {
        return {
            loading: false,
            rows: [],
        }
    },
    methods: {
        async fetchData() {
            if (this.timeRange.length < 2) return
            const targets = this.panelConfig.targets
            let merged = {}
            this.$emit('setIsLoading', true)
            for (const i in targets) {
                try {
                    const target = targets[i]
                    const response = await this.axios.get('/api/prometheus/query', {
                        params: {
                            expr: target.expr,
                            time: this.timeRange[1],
                        }
                    })
                    let rows = {}
                    response.data.data.result.forEach(x => {
                        const key = x.metric[target.key]
                        let row = {}
                        target.columns.forEach(c => {
                            if (c == 'VALUE') row['VALUE' + i] = x.value[1]
                            else row[c] = x.metric[c]
                        })
                        rows[key] = row
                    })
                    Object.entries(rows).forEach(r => {
                        const k = r[0]
                        const v = r[1]
                        merged[k] = { ...merged[k], ...rows[k] }
                    })
                } catch (error) { console.error(error) }
            }
            this.$emit('setIsLoading', false)
            this.rows = merged
        },
    },
    timerHandler() {
        if (useTimeStore().timerManager != this.timerManager || this.intervalSeconds == 0) return
        this.execute()
    },
    mounted() {
        this.fetchData()
    }
}
</script>
