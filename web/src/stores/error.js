import { defineStore } from 'pinia';

export const useErrorStore = defineStore('error', {
  state: () => ({
    status: null,
    errorType: null,
    error: null,
  }),
});