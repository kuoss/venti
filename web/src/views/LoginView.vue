<script lang="ts">
import type ErrorData from '@/types/error';
import { useErrorStore } from '@/stores/error';

export default {
  data() {
    return {
      errorData: {} as ErrorData,
      loading: false,
      username: '',
      password: '',
    };
  },
  methods: {
    async login() {
      this.loading = !this.loading;
      await this.$auth.login(this.username, this.password);
      this.loading = !this.loading;
    },
  },
  mounted() {
    this.$watch(
      () => this.$auth.loggedIn,
      () => {
        if (this.$auth.loggedIn) {
          // @ts-ignore
          this.$router.push({ name: 'home' });
        }
      },
    );
    this.$watch(
      () => useErrorStore().errorData,
      () => {
        if (useErrorStore().errorData) {
          this.errorData = useErrorStore().errorData;
        }
      },
    );
  },
};
</script>

<template>
  <div class="w-80 mx-auto mt-10">
    <h1 class="text-lg text-center py-10">Welcome to venti</h1>
    <form class="bg-slate-200 shadow-md rounded px-8 pt-6 pb-8 mb-4" @submit.prevent="login">
      <div class="mb-4">
        <label class="block text-gray-700 font-bold mb-2" for="username">Username</label>
        <input
          id="username"
          type="text"
          placeholder="Username"
          autocomplete="username"
          class="shadow appearance-none border border-slate-400 rounded w-full py-2 px-3 bg-white text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          :class="{ 'border-red-500': errorData.error }"
          v-model="username"
          :disabled="loading"
        />
      </div>
      <div class="mb-6">
        <label class="block text-gray-700 font-bold mb-2" for="password">Password</label>
        <input
          id="password"
          type="password"
          placeholder="******************"
          autocomplete="current-password"
          class="shadow appearance-none border rounded w-full py-2 px-3 bg-white text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          :class="[errorData.error ? 'border-red-500' : 'border-slate-400']"
          v-model="password"
          :disabled="loading"
        />
        <p v-if="errorData.error" class="text-red-500 text-xs italic">
          {{ errorData.error }}
        </p>
      </div>
      <div class="flex items-center justify-between">
        <button
          type="submit"
          :disabled="loading"
          class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
        >
          <div v-if="loading" class="spinner-border mx-3 spinner-border-sm" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
          <div v-else>Login</div>
        </button>
      </div>
    </form>
  </div>
</template>
