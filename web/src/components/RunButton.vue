<template>
    <div class="relative" v-click-outside="close">
        <button
            class="border-x-group p-2 px-4 bg-blue-200 dark:bg-blue-900 hover:bg-blue-300"
            :disabled='disabled'
            :class="{'cursor-not-allowed':disabled}"
            @click="$emit('execute')"
        >
            <i class="mdi mdi-refresh" :class="{ 'mdi-spin': intervalSeconds > 0 }"></i> {{btnText}}
        </button>
        <button
            class="border-x-group p-2 bg-blue-200 dark:bg-blue-900 hover:bg-blue-300"
            @click="isOpen = !isOpen"
        >
            <span class="px-1" v-if="intervalSeconds > 0">{{ text }}</span>
            <i class="mdi" :class="[isOpen ? 'mdi-chevron-up' : 'mdi-chevron-down']"></i>
        </button>
        <div class="absolute top-[100%] right-0 bg-white dark:bg-stone-900 border" v-if="isOpen">
            <div
                class="p-1 px-2 cursor-pointer hover:bg-gray-200"
                v-for="i in intervals"
                @click="selectInterval(i)"
            >{{ i }}</div>
        </div>
    </div>
</template>

<script>
export default {
    props: ['disabled','btnText'],
    data() {
        return {
            text: '',
            isOpen: false,
            intervalSeconds: 0,
            intervals: ['Off', '5s', '10s', '30s', '1m', '5m', '15m', '30m', '1h'],
        }
    },
    methods: {
        close() {
            this.isOpen = false
        },
        selectInterval(i) {
            this.text = i
            this.intervalSeconds = this.$util.age2seconds(i)
            this.$emit('changeInterval', this.intervalSeconds)
            this.close()
        },
    }
}
</script>