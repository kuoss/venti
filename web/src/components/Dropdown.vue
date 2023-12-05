<script setup lang="ts">
import { ref, watch } from 'vue';
import { vOnClickOutside } from '@vueuse/components'

const props = defineProps({
  options: { type: Array<String>, required: true },
  index: { type: Number, default: 0 },
})

const emit = defineEmits<{
  (e: 'change', value: String): void
}>()

const myOptions = ref(props.options)
const myIndex = ref(props.index)

const dropdown = ref(false)

function onClickOutside() {
  dropdown.value = false
}

function onClick(idx: number) {
  dropdown.value = false
  myIndex.value = idx
  emit('change', myOptions.value[myIndex.value])
}

watch(() => props.options, (newValue) => {
  myOptions.value = newValue
})
</script>

<template>
  <div class="inline-block relative z-30">
    <button class="border text-gray-700 dark:text-gray-300 py-2 px-4 rounded inline-flex items-center" @click.stop="dropdown = !dropdown">
      <span>{{ myOptions[myIndex] }}</span>
      <i class="mdi mdi-chevron-down" />
    </button>
    <ul v-if="dropdown" v-on-click-outside="onClickOutside"
      class="absolute text-gray-700 dark:text-gray-300 w-max border border-gray-300 dark:border-gray-700 bg-white dark:bg-black max-h-[80vh] overflow-y-auto">
      <li v-for="(option, idx) in options" class="hover:bg-gray-400 dark:hover:bg-gray-600 px-2 cursor-pointer py-[1px]" @click="onClick(idx)">
        {{ option }}
      </li>
    </ul>
  </div>
</template>
