<template>
    <div class="flex h-full">
        <div class="m-auto">
            <span>Please wait...</span>
        </div>
    </div>
</template>

<script setup>
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { onMounted } from "@vue/runtime-core";
const auth = useAuthStore();
const router = useRouter();
onMounted(() => auth.logout().then(() => router.push({ name: "login" })));
</script>
