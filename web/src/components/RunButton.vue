<script setup lang="ts">
import { ref } from 'vue'
import { age2seconds } from '@/lib/util2'

defineProps({
  disabled: { type: Boolean, default: false },
  btnText: { type: String, required: true },
})

const emit = defineEmits<{
  (e: 'execute'): void
  (e: 'changeInterval', intervalSeconds: number): void
}>()

const text = ref('')
const isOpen = ref(false)
const intervalSeconds = ref(0)
const intervals = ref(['Off', '5s', '10s', '30s', '1m', '5m', '15m', '30m', '1h'])

function close() {
  isOpen.value = false
}

function onClick() {
  emit('execute')
}

function onSelectInterval(i: string) {
  text.value = i
  intervalSeconds.value = age2seconds(i)
  emit('changeInterval', intervalSeconds.value)
  close()
}
</script>

<template>
  <div v-click-outside="close" class="relative">
    <button class="border-x-group p-2 px-4 bg-blue-200 dark:bg-blue-900 hover:bg-blue-300" :disabled="disabled"
      :class="{ 'cursor-not-allowed': disabled }" @click="onClick">
      <i class="mdi mdi-refresh" :class="{ 'mdi-spin': intervalSeconds > 0 }" />
      {{ btnText }}
    </button>
    <button class="border-x-group p-2 bg-blue-200 dark:bg-blue-900 hover:bg-blue-300" @click="isOpen = !isOpen">
      <span v-if="intervalSeconds > 0" class="px-1">{{ text }}</span>
      <i class="mdi" :class="[isOpen ? 'mdi-chevron-up' : 'mdi-chevron-down']" />
    </button>
    <div v-if="isOpen" class="absolute top-[100%] right-0 bg-white dark:bg-stone-900 border">
      <div v-for="i in intervals" class="p-1 px-2 cursor-pointer hover:bg-gray-200" @click="onSelectInterval(i)">
        {{ i }}
      </div>
    </div>
  </div>
</template>
