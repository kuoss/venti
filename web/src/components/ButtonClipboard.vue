<script setup>
import { ClipboardIcon, CheckIcon } from '@heroicons/vue/outline'
</script>

<template>
    <div class="relative">
        <button
            :class="`p-1 rounded border ${buttonClass}`"
            @mouseenter="hover = true"
            @mouseleave="hover = false"
            @click="copy"
        >
            <component
                class="w-4 h-4 mt-[-1px] inline"
                :class="{ 'text-green-400': copied }"
                :is="copied ? 'CheckIcon' : 'ClipboardIcon'"
            />
            <span v-if="text" class="px-2">{{ text }}</span>
        </button>
        <div
            v-if="hover || copied"
            class="absolute w-32 px-3 py-1 z-50 bg-white rounded text-center m-[2px] mx-[3px]"
            :class="tooltipClass"
        >{{ copied ? 'Copied!' : 'Copy to clipboard' }}</div>
    </div>
</template>

<script>
export default {
    components: { ClipboardIcon, CheckIcon },
    props: ['value', 'buttonClass', 'text', 'tooltipDirection'],
    computed: {
        tooltipClass() {
            switch (this.tooltipDirection) {
                case 'top': return 'bottom-[100%] right-0 translate-x-3'
                case 'right': return 'left-[100%] top-0'
            }
            return 'right-[100%] top-0'
        }
    },
    data() {
        return {
            hover: false,
            copied: false,
            timer: null
        }
    },
    methods: {
        copy() {
            this.removeTimer()
            this.copied = true
            this.$util.copyToClipboard(this.value)
            this.timer = setTimeout(() => { this.copied = false }, 2500)
        },
        removeTimer() {
            if (this.timer) clearTimeout(this.timer)
        },
    }
}
</script>