<script setup lang="ts">
import { ref, computed } from 'vue';
import { mdiClipboardOutline } from '@mdi/js';
import copyToClipboard from '@/lib/clipboard';

import TheIcon from '@/components/TheIcon.vue';

const props = defineProps({
  value: { type: String, default: '' },
  buttonClass: { type: String, default: '' },
  text: { type: String, default: '' },
  tooltipDirection: { type: String, default: 'left' },
})

const hover = ref(false)
const copied = ref(false)

let timerID = 0;

const tooltipClass = computed(() => {
  switch (props.tooltipDirection) {
    case 'top':
      return 'bottom-[100%] right-0 translate-x-3';
    case 'right':
      return 'left-[100%] top-0';
  }
  return 'right-[100%] top-0';
})

function onClick() {
  console.log('onClick')
  copied.value = true;
  copyToClipboard(props.value);

  if (timerID) {
    clearTimeout(timerID);
  }
  timerID = setTimeout(() => {
    copied.value = false;
  }, 2500);
}

</script>

<template>
  <div class="relative">
    <button :class="`p-1 rounded border leading-4 ${buttonClass}`" @mouseenter="hover = true" @mouseleave="hover = false"
      @click="onClick">
      <TheIcon :path="mdiClipboardOutline" :size="15" />
      <span v-if="text" class="px-2">{{ text }}</span>
    </button>
    <div v-if="hover || copied" class="absolute w-32 px-3 py-1 z-50 bg-white rounded text-center m-[2px] mx-[3px]"
      :class="tooltipClass">
      {{ copied ? 'Copied!' : 'Copy to clipboard' }}
    </div>
  </div>
</template>
