<script setup lang="ts">
defineProps({
  path: { type: String, required: true },
  size: { type: Number, default: 24 },
})
</script>

<template>
  <svg fill="currentColor" class="inline" :width="size" :height="size" viewBox="0 0 24 24"
    xmlns="http://www.w3.org/2000/svg">
    <path :d="path" />
  </svg>
</template>