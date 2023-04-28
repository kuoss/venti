import { defineStore } from 'pinia';

export const useDatasourceStore = defineStore('datasource', {
  state: () => ({
    datasources: [],
    status: { loaded: false, loading: false },
  }),
  actions: {
    async getDatasources() {
      await this.waitForLoaded();
      return this.datasources;
    },
    async waitForLoaded() {
      let tries = 0;
      if (!this.status.loading) this.loadData();
      while (!this.status.loaded)
        await new Promise(resolve => setTimeout(resolve, 100 * ++tries));
    },
    async loadData() {
      this.status.loading = true;
      try {
        const response = await fetch('/api/datasources');
        this.datasources = await response.json();
        this.status.loaded = true;
      } catch (error) {
        console.error(error);
        this.status.loading = false;
      }
    },
  },
});
