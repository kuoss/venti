<script setup>
import PanelLogs from '@/components/PanelLogs.vue'
import PanelMultitable from '@/components/PanelMultitable.vue'
import PanelPiechart from '@/components/PanelPiechart.vue'
import PanelStat from '@/components/PanelStat.vue'
import PanelTable from '@/components/PanelTable.vue'
import PanelTimeSeries from '@/components/PanelTimeSeries.vue'

import { useSidePanelStore } from '@/stores/sidePanel'
</script>

<template>
    <div class="flex border-b">
        <button class="p-1 bg-cyan-100" v-if="showPanelPosition" @click="useSidePanelStore().goToPanelConfig(position)">{{ position }}</button>
        <div class="flex-1 py-1 text-center font-bold" :class="{ 'is-loading': isLoading }">
            <button @click="togglePanelInfo" class="hover:underline">
                {{ panelConfig.title }}
                <span v-for="v in usingVariables">
                    <span class="w-2 rounded-full" :class="v.class">n</span>
                </span>
            </button>
        </div>
    </div>
    <component
        :is="componentName"
        :count="count"
        :isLoading="isLoading"
        @setIsLoading="setIsLoading"
        :panelConfig="panelConfig"
        :panelWidth="panelWidth"
        :timeRange="timeRange"
    />
</template>

<script>
export default {
    components: {
        PanelLogs,
        PanelMultitable,
        PanelPiechart,
        PanelStat,
        PanelTable,
        PanelTimeSeries,
    },
    props: {
        position: String,
        count: Number,
        panelConfig: Object,
        panelWidth: Number,
        timeRange: Object,
    },
    computed: {
        componentName() {
            switch (this.panelConfig.type) {
                case 'logs': return 'PanelLogs'
                case 'multitable': return 'PanelMultitable'
                case 'piechart': return 'PanelPiechart'
                case 'stat': return 'PanelStat'
                case 'table': return 'PanelTable'
                case 'time_series': return 'PanelTimeSeries'
            }
        },
    },
    data() {
        return {
            isLoading: false,
            showPanelPosition: false,
            usingVariables: [],
        }
    },
    methods: {
        setIsLoading(b) {
            this.isLoading = b
        },
        togglePanelInfo() {
            useSidePanelStore().goToPanelConfig(this.position)
        }
    },
    mounted() {
        const variables = [{ name: '$namespace', class: 'namespace' }, { name: '$node', class: 'node' }]
        variables.forEach((v, i) => {
            if (this.panelConfig.targets[0].expr.indexOf(v.name) > 0) {
                this.usingVariables.push(v)
            }
        })
        useSidePanelStore().$subscribe((mutation, state) => {
            this.showPanelPosition = state.show && state.type == 'DashboardInfo'
        })
    }
}
</script>

<style scoped>
.node {
    @apply text-yellow-500;
}
.namespace {
    @apply text-green-500;
}
</style>