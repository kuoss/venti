<script setup>
import UplotVue from 'uplot-vue'
import 'uplot/dist/uPlot.min.css'

import { useFilterStore } from "@/stores/filter";
import { useTimeStore } from "@/stores/time";
import { useSidePanelStore } from "@/stores/sidePanel";
</script>

<template>
    <div v-if="!isLoading && isNoData" class="h-[150px] grid grid-cols-1 content-center">
        <div class="text-center text-lg">No data</div>
    </div>
    <div v-else>
        <UplotVue @click="onChartClick" :data="data" :options="options" ref="main" />
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
        panelWidth() { this.resize() },
    },
    data() {
        return {
            isNoData: true,
            expr: this.panelConfig.targets[0].expr,
            data: [],
            options: {
                axes: [
                    {
                        stroke: "#888",
                        grid: { stroke: "#8885", width: 1 },
                        ticks: { stroke: "#8885", width: 1 },
                        values: [
                            [3600 * 24 * 365, "{YYYY}", null, null, null, null, null, null, 1],
                            [3600 * 24 * 28, "{MM}", "\n{YYYY}", null, null, null, null, null, 1],
                            [3600 * 24, "{MM}-{DD}", "\n{YYYY}", null, null, null, null, null, 1],
                            [3600, "{HH}:00", "\n{YYYY}-{MM}-{DD}", null, "\n{MM}-{DD}", null, null, null, 1],
                            [80, "{HH}:{mm}", "\n{YYYY}-{MM}-{DD}", null, "\n{MM}-{DD}", null, null, null, 1],
                            [1, "{HH}:{mm}:{ss}", "\n{YYYY}-{MM}-{DD}", null, "\n{MM}-{DD}", null, null, null, 1],
                        ],
                    },
                    {
                        stroke: "#888",
                        grid: { stroke: "#8885", width: 1 },
                        ticks: { stroke: "#8885", width: 1 },
                    },
                ],
                width: 100,
                height: 180,
                legend: { show: false },
                cursor: { points: false },
                scales: { x: { time: true }, y: { auto: true } },
                select: { show: false },
                series: [],
                plugins: [this.tooltipPlugin()],
            }
        }
    },
    methods: {
        async fetchData() {
            if (this.timeRange.length < 2) return
            const target = this.panelConfig.targets[0]
            this.$emit('setIsLoading', true)
            try {
                const response = await this.axios.get('/api/prometheus/query_range', {
                    params: {
                        expr: useFilterStore().renderExpr(this.expr),
                        start: this.timeRange[0],
                        end: this.timeRange[1],
                        step: (this.timeRange[1] - this.timeRange[0]) / 120,
                    }
                })
                const result = response.data.data.result
                if (result.length < 1) {
                    this.isNoData = true
                    return
                }
                this.isNoData = false

                const temp = result.map(x => x.values)
                const timestamps = Array.from(new Set(temp.map(a => a.map(b => b[0])).flat())).sort()
                const seriesData = temp.map(a => {
                    let newA = []
                    timestamps.forEach(t => {
                        const newPoint = a.filter(b => t == b[0])
                        if (newPoint.length != 1 || isNaN(parseFloat(newPoint[0][1]))) {
                            newA.push(null)
                            return
                        }
                        newA.push(parseFloat(newPoint[0][1]))
                    })
                    return newA
                })
                // labels
                const labels = result.map(x => target.legend.replace(/\{\{(.*?)\}\}/g, (i, m) => x.metric[m]))
                let newSeries = []
                newSeries.push({})
                labels.forEach((x, i) => newSeries.push({ label: x, stroke: this.$util.string2color(x), width: 1, points: { size: 0 } }))
                let options = { ...this.options, series: newSeries, scales: useTimeStore().scales }
                const yMax = this.panelConfig.chartOptions?.yMax ?? 0
                options.scales = { y: { range: (self, fromMin, fromMax) => [0, Math.max(fromMax, yMax)] } }
                this.options = options
                this.data = [timestamps, ...seriesData]
            } catch (error) { console.error(error) }
            this.$emit('setIsLoading', false)
        },
        resize() {
            this.options = { ...this.options, width: this.panelWidth - 15 }
        },
        onChartClick() {
            useSidePanelStore().toggleShow('DataTable')
        },
        tooltipPlugin() {
            return {
                hooks: {
                    setCursor: u => {
                        let columnData = u.data.map(x => x[u.cursor.idx])
                        const time = columnData.shift()
                        if (!time) return
                        const labels = this.options.series.map(x => x['label']).slice(1)
                        const rows = labels.map((x, i) => [x, columnData[i]]).filter(x => x[1] != undefined)
                        useSidePanelStore().updatetDataTable({
                            title: this.panelConfig.title,
                            time: time,
                            rows: rows,
                        })
                    },
                }
            }
        },
        async init() {
            await this.fetchData()
            this.resize()
        },
    },
    mounted() {
        this.init()
    },
}
</script>