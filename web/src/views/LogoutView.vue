<script lang="ts">
export default {
  mounted() {
    this.$auth.logout().then(() => {
      // @ts-ignore
      this.$router.push({ name: 'login' });
    });
  },
};
</script>

<template>
  <div class="flex h-full">
    <div class="m-auto">
      <span>Please wait...</span>
    </div>
  </div>
</template>
