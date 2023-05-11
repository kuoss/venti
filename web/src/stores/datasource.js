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
      while (!this.status.loaded) await new Promise(resolve => setTimeout(resolve, 100 * ++tries));
    },
    async loadData() {
      this.status.loading = true;
      try {
        const resp1 = await fetch('/api/v1/datasources');
        let datasources = await resp1.json();

        const resp2 = await fetch('/api/v1/datasources/targets');
        const temp = await resp2.json();
        const targetInfos = temp.map(x => JSON.parse(x));

        for (let i = 0; i < targetInfos.length; i++) {
          const targetInfo = targetInfos[i];
          if (targetInfo.status != 'success') {
            datasources[i].health = false;
            datasources[i].target = [];
            continue;
          }
          datasources[i].health = true;
          datasources[i].targets = this.getTargets(targetInfo.data.activeTargets);
        }
        this.datasources = datasources;
        this.status.loaded = true;
      } catch (error) {
        console.error(error);
        this.status.loading = false;
      }
    },
    getTargets(activeTargets) {
      let targets = [];
      for (let i = 0; i < activeTargets.length; i++) {
        targets.push(this.getTarget(activeTargets[i]));
      }
      return targets;
    },
    getTarget(x) {
      x.age = ((new Date() - new Date(x.lastScrape)) / 1000).toFixed();
      x.job = x.discoveredLabels.job;
      for (const [k, v] of Object.entries(x.discoveredLabels)) {
        if (k == '__meta_kubernetes_service_name') {
          x.icon = '🐕‍🦺';
          x.name = v;
        } else if (k == '__meta_kubernetes_pod_name') {
          x.icon = '🐱';
          x.name = v;
        } else if (k == '__meta_kubernetes_node_name') {
          x.icon = '🏝';
          x.name = v;
        } else if (k == '__meta_kubernetes_namespace') {
          x.icon = '🖼️';
          x.name = v;
        }
      }
      x.icon = x.icon || '💼';
      x.name = x.name || x.job;
      return x;
    },
  },
});
