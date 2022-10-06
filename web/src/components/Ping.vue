<template>
    <template v-if="loading">
        <span class="absolute right-2 top-2 h-1 w-1 rounded-full bg-sky-500 opacity-50"></span>
        <span class="animate-ping absolute right-2 top-2 h-1 w-1 rounded-full bg-sky-500"></span>
    </template>
</template>

<script>
export default {
    props: {
        loading: Boolean,
    },
}
</script>