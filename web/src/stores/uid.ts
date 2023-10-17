import { defineStore } from 'pinia';

export const useUIDStore = defineStore('uid', {
  state: () => {
    return {
      lastUID: 1,
      currentUID: 0,
    }
  },
});
