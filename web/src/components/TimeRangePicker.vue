<script setup>
import VueDatePicker from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';

import { useConfigStore } from '@/stores/config';
import { useTimeStore } from '@/stores/time';
</script>
<script>
export default {
  components: { VueDatePicker },
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
      ],
    };
  },
  computed: {
    btnText() {
      return this.rangeToText(this.range);
    },
  },
  watch: {
    range() {
      this.$emit('updateTimeRange', this.range);
    },
  },
  mounted() {
    this.offset = useTimeStore().getOffset();
    this.offsetString = useTimeStore().getOffsetString();
    this.onSelectRelativeTimeRange(this.relativeTimeRanges[0]);
    useConfigStore().$subscribe((_, state) => {
      this.dark = state.dark;
    });
  },
  methods: {
    async applyTimeRange() {
      this.range = [this.startTime, this.endTime];
      this.close();
    },
    rangeToText(r) {
      if (r.length != 2) return '';
      if (r[1] == 'now' && r[0].startsWith('now') && r[0] != 'now') {
        const parts = r[0].split('-');
        const num = parts[1].slice(0, -1);
        let unit = parts[1].slice(-1);
        switch (unit) {
          case 's':
            unit = 'second';
            break;
          case 'm':
            unit = 'minute';
            break;
          case 'h':
            unit = 'hour';
            break;
          case 'd':
            unit = 'day';
            break;
        }
        if (num > 1) unit += 's';
        return `Last ${num} ${unit}`;
      }
      return `${r[0]} to ${r[1]}`;
    },
    async onClickTimeInput() {
      const timeRange = await useTimeStore().toTimeRangeForQuery(this.range, false);
      this.localDate = [
        new Date(timeRange[0] * 1000 - this.offset * 60000),
        new Date(timeRange[1] * 1000 - this.offset * 60000),
      ];
      this.showCalendar = true;
    },
    onSelectRelativeTimeRange(r) {
      this.startTime = r[0];
      this.endTime = r[1];
      this.applyTimeRange();
    },
    toggleIsOpen() {
      this.isOpen = !this.isOpen;
      if (!this.isOpen) this.close();
    },
    updateDatePicker() {
      this.localDate[0].setHours(0, 0, 0, 0);
      this.localDate[1].setHours(23, 59, 59, 0);
      this.startTime = new Date(this.localDate[0] - this.offset * 60 * 1000).toISOString().split('T')[0] + ' 00:00:00';
      this.endTime = new Date(this.localDate[1] - this.offset * 60 * 1000).toISOString().split('T')[0] + ' 23:59:59';
    },
    close() {
      this.showCalendar = false;
      this.isOpen = false;
    },
  },
};
</script>

<template>
  <section v-click-outside="close" class="relative">
    <button
      class="py-2 px-4 text-gray-900 dark:text-gray-100 bg-white dark:bg-black rounded border border-gray-200 dark:border-gray-800 hover:bg-gray-100 dark:hover:bg-gray-900 hover:text-blue-700 dark:hover:text-blue-300"
      @click="toggleIsOpen"
    >
      <i class="mdi mdi-clock-outline" />
      {{ btnText }}
      <i class="mdi" :class="isOpen ? 'mdi-chevron-up' : 'mdi-chevron-down'" />
    </button>
    <section
      v-if="isOpen"
      class="absolute top-[100%] right-0 w-full min-w-fit mt-1 rounded bg-white dark:bg-stone-900 border"
    >
      <div class="relative">
        <div class="absolute top-0 right-[100%] mr-[1px]" :class="{ hidden: !showCalendar }">
          <VueDatePicker
            ref="datepicker"
            v-model="localDate"
            week-start="0"
            range
            inline
            no-today
            auto-apply
            close-on-auto-apply
            enable-seconds
            multi-calendars
            :dark="dark"
            :max-date="new Date()"
            :multi-static="false"
            :start-date="new Date(new Date().getFullYear(), new Date().getMonth() - 1)"
            :enable-time-picker="false"
            @update:modelValue="updateDatePicker"
            timezone="UTC"
          />
        </div>
      </div>
      <div class="grid grid-cols-2 w-[24rem]">
        <div class="col-span-1">
          <div class="font-bold p-2">Absolute time range</div>
          <div class="p-2 px-4">
            <div class="mt-2">From</div>
            <input
              v-model="startTime"
              type="text"
              class="p-1 w-full border bg-white dark:bg-black text-black dark:text-white"
              @click="onClickTimeInput"
            />
            <div class="mt-2">To</div>
            <input
              v-model="endTime"
              type="text"
              class="p-1 w-full border bg-white dark:bg-black text-black dark:text-white"
              @click="onClickTimeInput"
            />
            <button class="p-1 mt-3 w-full bg-gray-200 dark:bg-gray-800 border rounded" @click="applyTimeRange">Apply time range</button>
          </div>
        </div>
        <div class="col-span-1 border-l pb-3">
          <div class="p-2 font-bold">Relative time range</div>
          <div v-for="r in relativeTimeRanges">
            <div
              class="p-1 px-4 cursor-pointer hover:bg-gray-200 dark:hover:bg-gray-800"
              :class="{
                'text-bold bg-gray-300 dark:bg-gray-700': rangeToText(r) == rangeToText(range),
              }"
              @click="onSelectRelativeTimeRange(r)"
            >
              {{ rangeToText(r) }}
            </div>
          </div>
        </div>
      </div>
      <div class="border-t p-2">
        {{ timezone }}
        <div class="float-right px-2 bg-gray-200 dark:bg-gray-800">
          {{ offsetString }}
        </div>
      </div>
    </section>
  </section>
</template>

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
