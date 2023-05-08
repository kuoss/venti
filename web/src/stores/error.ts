import { defineStore } from 'pinia';
import type ErrorData from '@/types/error';

export const useErrorStore = defineStore('error', {
  state: () => ({
    errorData: {} as ErrorData,
  }),
  actions: {
    set(errorData: ErrorData) {
      this.errorData = errorData;
    },
    clear() {
      this.errorData = {} as ErrorData;
    },
  }
});