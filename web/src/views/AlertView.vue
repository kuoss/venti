<script setup>
import { useTimeStore } from "@/stores/time"
</script>

<template>
    <header
        class="fixed right-0 w-full bg-white border-b shadow z-30 p-2 pl-52"
        :class="{ 'is-loading': isLoading }"
    >
        <div class="flex items-center flex-row">
            <div>
                <i class="mdi mdi-18px mdi-database-outline"></i> Alert
            </div>
            <div class="flex ml-auto">
                <div class="inline-flex">
                    <button
                        class="h-rounded-group py-2 px-4 text-gray-900 bg-white border border-common"
                    >
                        <i class="mdi mdi-refresh mdi-spin"></i>
                    </button>
                </div>
            </div>
        </div>
    </header>

    <main class="mt-12 w-full p-8 pb-16">
        <h1 class="py-2 font-bold">Alert Groups ({{ alertGroups.length }} groups)</h1>
        <table class="w-full bg-white border">
            <tr class="border-b bg-slate-50">
                <th>State</th>
                <th>Severity</th>
                <th>Name</th>
                <th>Summary</th>
                <th>Expr</th>
                <th>For</th>
            </tr>
            <tbody v-for="g in alertGroups">
                <tr class="border-b">
                    <th
                        class="text-left bg-slate-200 p-3"
                        colspan="6"
                    >{{ g.datasource == 'Prometheus' ? '🔥' : '💧' }} {{ g.name }} ({{ g.rules.length }} rules)</th>
                </tr>
                <tr class="border-b" v-for="r in g.rules">
                    <td class="px-2">
                        <div
                            class="rounded w-24 text-center bg-red-300"
                            v-if="r.state == 'firing'"
                        >Firing</div>
                        <div
                            class="rounded w-24 text-center bg-yellow-300"
                            v-else-if="r.state == 'pending'"
                        >Pending</div>
                        <div class="rounded w-24 text-center bg-green-300" v-else>Normal</div>
                    </td>
                    <td class="px-2">{{ r.labels.severity }}</td>
                    <td class="px-2">{{ r.alert }}</td>
                    <td
                        class="px-2 max-w-[30vw] truncate hover:whitespace-normal"
                    >{{ r.annotations.summary }}</td>
                    <td class="px-2 max-w-[20vw] truncate hover:whitespace-normal text-cyan-500">
                        <a
                            class="hover:underline"
                            :href="`/#/${(g.datasource == 'Prometheus') ? 'metrics' : 'logs'}?query=${encodeURIComponent(r.expr)}`"
                        >{{ r.expr }}</a>
                    </td>
                    <td class="px-2">{{ $util.nanoseconds2human(r.for) }}</td>
                </tr>
            </tbody>
        </table>
    </main>
</template>

<script>
export default {
    data() {
        return {
            alertGroups: [],
            isLoading: false,
        }
    },
    methods: {
        async fetchData() {
            this.isLoading = true
            try {
                const response = await this.axios.get('/api/alerts')
                this.alertGroups = response.data
                setTimeout(() => this.timerHandler(), 3000)
            } catch (error) { console.error(error) }
            this.isLoading = false
        },
        timerHandler() {
            if (useTimeStore().timerManager != 'AlertView' || this.intervalSeconds == 0) return
            this.fetchData()
        },
    },
    mounted() {
        useTimeStore().timerManager = 'AlertView'
        this.fetchData()
    }
}
</script>