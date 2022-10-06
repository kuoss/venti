<template>
    <section class="relative" v-click-outside="close">
        <button
            class="py-2 px-4 text-gray-900 bg-white rounded border border-gray-200 hover:bg-gray-100 hover:text-blue-700"
            @click="toggleIsOpen"
        >
            <i class="mdi mdi-clock-outline"></i>
            {{ btnText }}
            <i
                class="mdi"
                :class="isOpen ? 'mdi-chevron-up' : 'mdi-chevron-down'"
            ></i>
        </button>
        <section
            class="absolute top-[100%] right-0 w-full min-w-fit mt-1 rounded bg-white dark:bg-stone-900 border"
            v-if="isOpen"
        >
            <div class="relative">
                <div
                    class="absolute top-0 right-[100%] mr-[1px]"
                    :class="{ 'hidden': !showCalendar }"
                >
                    <Datepicker
                        ref="datepicker"
                        v-model="localDate"
                        weekStart="0"
                        range
                        inline
                        noToday
                        autoApply
                        closeOnAutoApply
                        enableSeconds
                        multiCalendars
                        :dark="dark"
                        :maxDate="new Date()"
                        :multiStatic="false"
                        :startDate="new Date(new Date().getFullYear(), new Date().getMonth() - 1)"
                        :enableTimePicker="false"
                        @update:modelValue="updateDatePicker"
                    />
                </div>
            </div>
            <div class="grid grid-cols-2 w-[24rem]">
                <div class="col-span-1">
                    <div class="font-bold p-2">Absolute time range</div>
                    <div class="p-2 px-4">
                        <div class="mt-2">From</div>
                        <input
                            type="text"
                            class="p-1 w-full border bg-white text-black"
                            v-model="startTime"
                            @click="onClickTimeInput"
                        />
                        <div class="mt-2">To</div>
                        <input
                            type="text"
                            class="p-1 w-full border bg-white text-black"
                            v-model="endTime"
                            @click="onClickTimeInput"
                        />
                        <button
                            class="p-1 mt-3 w-full bg-gray-200 border rounded"
                            @click="applyTimeRange"
                        >Apply time range</button>
                    </div>
                </div>
                <div class="col-span-1 border-l pb-3">
                    <div class="p-2 font-bold">Relative time range</div>
                    <div v-for="r in relativeTimeRanges">
                        <div
                            class="p-1 px-4 cursor-pointer hover:bg-gray-200"
                            :class="{ 'text-bold bg-gray-300': rangeToText(r) == rangeToText(range) }"
                            @click="onSelectRelativeTimeRange(r)"
                        >{{ rangeToText(r) }}</div>
                    </div>
                </div>
            </div>
            <div class="border-t p-2">
                {{ timezone }}
                <div class="float-right px-2 bg-gray-200">{{ offsetString }}</div>
            </div>
        </section>
    </section>
</template>

<script>
import Datepicker from '@vuepic/vue-datepicker'
import { useConfigStore } from "@/stores/config"
import { useTimeStore } from "@/stores/time"
import '@vuepic/vue-datepicker/dist/main.css'

export default {
    components: { Datepicker },
    data() {
        return {
            timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
            offset: null,
            offsetString: '',
            range: [],
            localDate: null,
            startTime: '',
            endTime: '',
            isOpen: false,
            showCalendar: false,
            dark: useConfigStore().dark,
            relativeTimeRanges: [
                ['now-5m', 'now'],
                ['now-15m', 'now'],
                ['now-30m', 'now'],
                ['now-1h', 'now'],
                ['now-3h', 'now'],
                ['now-6h', 'now'],
                ['now-12h', 'now'],
                ['now-1d', 'now'],
                ['now-2d', 'now'],
                ['now-7d', 'now'],
            ]
        };
    },
    watch: {
        range() {
            this.$emit('updateTimeRange', this.range)
        },
    },
    computed: {
        btnText() {
            return this.rangeToText(this.range)
        },
    },
    methods: {
        async applyTimeRange() {
            this.range = [this.startTime, this.endTime]
            this.close()
        },
        rangeToText(r) {
            if (r.length != 2) return ''
            if (r[1] == 'now' && r[0].startsWith('now') && r[0] != 'now') {
                const parts = r[0].split('-')
                const num = parts[1].slice(0, -1)
                let unit = parts[1].slice(-1)
                switch (unit) {
                    case 's': unit = 'second'; break;
                    case 'm': unit = 'minute'; break;
                    case 'h': unit = 'hour'; break;
                    case 'd': unit = 'day'; break;
                }
                if (num > 1) unit += 's'
                return `Last ${num} ${unit}`
            }
            return `${r[0]} to ${r[1]}`
        },
        async onClickTimeInput() {
            const timeRange = await useTimeStore().toTimeRangeForQuery(this.range, false)
            this.localDate = [new Date(timeRange[0]*1000-this.offset*60000), new Date(timeRange[1]*1000-this.offset*60000)]
            this.showCalendar = true
        },
        onSelectRelativeTimeRange(r) {
            this.startTime = r[0]
            this.endTime = r[1]
            this.applyTimeRange()
        },
        toggleIsOpen() {
            this.isOpen = !this.isOpen
            if (!this.isOpen) this.close()
        },
        updateDatePicker() {
            this.localDate[0].setHours(0, 0, 0, 0)
            this.localDate[1].setHours(23, 59, 59, 0)
            this.startTime = new Date(this.localDate[0] - (this.offset*60*1000)).toISOString().split('T')[0] + ' 00:00:00'
            this.endTime = new Date(this.localDate[1] - (this.offset*60*1000)).toISOString().split('T')[0] + ' 23:59:59'
        },
        close() {
            this.showCalendar = false
            this.isOpen = false
        },
    },
    mounted() {
        this.offset = useTimeStore().getOffset()
        this.offsetString = useTimeStore().getOffsetString()
        this.onSelectRelativeTimeRange(this.relativeTimeRanges[0])
        useConfigStore().$subscribe((_, state) => {
            this.dark = state.dark
        })
    },
}
</script>

<style>
.hideDP .dp__menu {
    @apply hidden;
}
.dp__menu {
    @apply text-xs;
}
.dp__arrow_top {
    @apply hidden;
}
.dp__theme_light,
.dp__theme_dark {
    --dp-hover-color: var(--dp-primary-color);
    --dp-hover-text-color: var(--dp-primary-text-color);
}
.dp__range_between {
    color: var(--dp-primary-text-color);
}
</style>