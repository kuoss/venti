<template>
  <div class="w-80 mx-auto mt-10">
    <h1 class="text-lg text-center py-10">
      Welcome to venti
    </h1>
    <form class="bg-slate-200 shadow-md rounded px-8 pt-6 pb-8 mb-4" @submit.prevent="onSubmit">
      <div class="mb-4">
        <label class="block text-gray-700 font-bold mb-2" for="username">Username</label>
        <input id="username" type="text" placeholder="Username"
          class="shadow appearance-none border border-slate-400 rounded w-full py-2 px-3 bg-white text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          :class="{ 'border-red-500': error.error }" v-model="credentials.username" :disabled="loading">
      </div>
      <div class="mb-6">
        <label class="block text-gray-700 font-bold mb-2" for="password">Password</label>
        <input id="password" type="password" placeholder="******************" autocomplete="current-password"
          class="shadow appearance-none border rounded w-full py-2 px-3 bg-white text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          :class="[error.error ? 'border-red-500' : 'border-slate-400']" v-model="credentials.password"
          :disabled="loading">
        <p v-if="error.error" class="text-red-500 text-xs italic">{{ error.error }}</p>
      </div>
      <div class="flex items-center justify-between">
        <button type="submit" :disabled="loading"
          class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
          <div v-if="loading" class="spinner-border mx-3 spinner-border-sm" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
          <div v-else>Login</div>
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { useErrorStore } from "@/stores/error";

const credentials = ref({});
const loading = ref(false);
const router = useRouter();
const error = useErrorStore();
const onSubmit = () => {
  loading.value = !loading.value;
  useAuthStore()
    .login(credentials.value)
    .then(() => router.push({ name: "home" }))
    .catch(() => (loading.value = !loading.value));
};
onBeforeUnmount(() => error.$reset());
</script>
