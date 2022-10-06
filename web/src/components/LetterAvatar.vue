<template>
  <div
    class="relative inline-block align-middle uppercase rounded-full h-5 w-5 text-white font-mono"
    :style="{ 'background-color': bgcolor }"
  >
    <div
      class="letters font-bold absolute inset-y-auto">{{ letters }}</div>
  </div>
</template>

<script>
export default {
  props: ['letters', 'bgcolor'],
}
</script>

<style>
.letters {
  font-size: 12px;
  left: 50%;
  transform: translate(-50%, 0);
}
</style>