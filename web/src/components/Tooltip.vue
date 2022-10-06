<template>
  <div class="relative inline-block" @mouseenter="hover=true" @mouseout="hover=false">
    <slot />
    <div v-if="hover"
      class="absolute w-40 bottom-[100%] mb-2 translate-x-[-50%] opacity-75 z-10 bg-black text-white text-center py-1 px-2 rounded"
      v-html="tip"
    >
    </div>
  </div>
</template>
<script>
export default {
  props: {
    tip: String,
  },
  data() {
    return {
      hover: false,
    }
  }
}
</script>
