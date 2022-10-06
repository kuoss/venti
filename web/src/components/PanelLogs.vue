<script setup>
import { useFilterStore } from "@/stores/filter";
</script>

<template>
    <div ref="logs" class="font-mono h-64 break-all overflow-y-auto">
        <template v-if="result.length > 0">
            <div v-for="line in result" v-html="colorizeLog(line)"></div>
        </template>
        <template v-else>
            <div class="text-center p-5">No data</div>
        </template>
    </div>
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
            total: "-",
            result: [],
            resultType: null,
        }
    },
    methods: {
        colorizeLog(row) {
            const classes = ['text-green-600', 'text-cyan-600', 'text-blue-600', 'text-purple-600', 'text-pink-600']
            const idx = row.indexOf(' ')
            if (idx == 20) return '<span class="text-yellow-500">' + row.substr(0, 20) + '</span> ' + row.substr(idx)
            return '<span class="text-yellow-500">' + row.substr(0, 20) + '</span>[' + row.substr(21, idx - 22).split('|').map((x, i) => `<span class="${classes[i]}">${x}</span>`).join('|') + '] ' + row.substr(idx)
        },
        async fetchData() {
            if (this.timeRange.length < 2) return
            this.$emit('setIsLoading', true)
            try {
                const response = await this.axios.get('/api/lethe/query_range', {
                    params: {
                        expr: useFilterStore().renderExpr(this.panelConfig.targets[0].expr),
                        start: this.timeRange[0],
                        end: this.timeRange[1],
                    }
                })
                this.result = response.data.data.result.slice(-100)
                this.resultType = response.data.data.resultType
                // console.log(this.result)
                // console.log(this.resultType)
                setTimeout(() => { if (this.$refs.logs) this.$refs.logs.scrollTop = 99999 }, 100)
            } catch (error) { console.error(error) }
            this.$emit('setIsLoading', false)
        },
    },
    mounted() {
        this.fetchData()
    }
}
</script>
