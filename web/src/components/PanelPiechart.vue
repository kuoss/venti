<script setup>
import { useFilterStore } from "@/stores/filter";
import { Pie } from 'vue-chartjs'
import { Chart as ChartJS, ArcElement } from 'chart.js'
ChartJS.register(ArcElement)
</script>

<template>
    <div v-if="chartData">
        <div class="piechart-wrapper p-2">
            <pie :chart-data="chartData" :chart-options="chartOptions"></pie>
        </div>
        <table class="border-t border-common w-full">
            <tr class="border-b border-common" v-for="(label, idx) in chartData.labels">
                <td>
                    <span
                        class="px-2"
                        :style="'color:' + chartData.datasets[0].backgroundColor[idx]"
                    >●</span>
                    {{ label }}
                </td>
                <td
                    :class="{ 'opacity-50': chartData.datasets[0].data[idx] < 1 }"
                    class="px-2 text-right"
                >{{ chartData.datasets[0].data[idx] }}</td>
            </tr>
        </table>
    </div>
</template>

<style>
</style>

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
            chartData: null,
            chartOptions: { maintainAspectRatio: false, elements: { arc: { borderWidth: 0.5 } } },
        }
    },
    methods: {
        getColorsForLabels(labels) {
            const labelColors = {
                'active': '#2b8',
                'bound': '#2b8',
                'available': '#176',
                'failed': '#dc2626',
                'pending': '#ffad20',
                'ready': '#2b8',
                'running': '#2b8',
                'succeeded': '#176',
            }
            let colors = []
            labels.forEach((a) => {
                const x = Object.entries(labelColors).filter(b => b[0] == a.toLowerCase())
                if (x.length > 0) colors.push(x[0][1])
                else colors.push("#aaa")
            })
            return colors;
        },
        async fetchData() {
            if (this.timeRange.length < 2) return
            this.$emit('setIsLoading', true)
            try {
                let data = []
                for (const target of this.panelConfig.targets) {
                    const response = await this.axios.get('/api/prometheus/query', {
                        params: {
                            expr: useFilterStore().renderExpr(target.expr),
                            time: this.timeRange[1],
                        }
                    })
                    data = [...data, ...response.data.data.result.map(x => ({
                        label: target.legend.replace(/\{\{(.*?)\}\}/g, (i, m) => x.metric[m]),
                        value: 1 * x.value[1]
                    }))].sort((a, b) => 1 * a.value > 1 * b.value ? -1 : 1)
                }
                const labels = data.map(x => x.label)
                const values = data.map(x => x.value)
                this.total = values.reduce((a, b) => a + b)
                this.chartData = {
                    labels: labels,
                    datasets: [
                        {
                            data: values,
                            backgroundColor: this.getColorsForLabels(labels),
                        }
                    ]
                }
            } catch (error) { console.error(error) }
            this.$emit('setIsLoading', false)
        },
    },
    mounted() {
        this.fetchData()
    }
}
</script>
