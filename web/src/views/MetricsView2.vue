<script setup>
import { useTimeStore } from "@/stores/time"
import TimeRangePicker from "@/components/TimeRangePicker.vue"
import RunButton from "@/components/RunButton.vue"
</script>

<template>
    <header class="fixed right-0 w-full bg-white border-b border-common shadow z-30 p-2 pl-52"
        :class="{ 'is-loading': loading }">
        <div class="flex items-center flex-row">
            <div>
                <i class="mdi mdi-18px mdi-numeric"></i> Metrics
            </div>
        </div>
    </header>
</template>

<script>
import UplotVue from 'uplot-vue'
import 'uplot/dist/uPlot.min.css'

export default {
    components: {
        RunButton,
        TimeRangePicker,
        uplotvue: UplotVue,
    },
    data() {
        return {
            loading: false,
        }
    }
}
</script>