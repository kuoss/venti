import { defineStore } from 'pinia'
import type { Datasource, Target } from '@/types/datasource'

export const useDatasourceStore = defineStore('datasource', {
  state: () => ({
    datasources: [] as Datasource[],
    status: { loaded: false, loading: false },
  }),
  actions: {
    async getDatasources() {
      await this.waitForLoaded()
      return this.datasources
    },
    async waitForLoaded() {
      let tries = 0;
      if (!this.status.loading) this.fetchDatasources()
      while (!this.status.loaded) {
        await new Promise(resolve => setTimeout(resolve, 100 * ++tries))
      }
    },
    async fetchDatasources() {
      this.status.loading = true;
      try {
        const response = await fetch('/api/v1/datasources')
        this.datasources = await response.json()
        this.status.loaded = true
      } catch (error) {
        console.error(error)
        this.status.loading = false
      }
    },
    getDatasourceByName(name: string) {
      for (const ds of this.datasources as Datasource[]) {
        if (ds.name == name) {
          return ds
        }
      }
      return null
    },
    async getDatasourceHealthy(datasource: Datasource): Promise<Number> {
      try {
        const resp = await fetch(`/api/v1/remote/healthy?dsName=${datasource.name}`)
        const text = await resp.text()
        if (text.includes(' is Healthy.')) {
          return 1
        }
        return 0
      } catch (error) {
        console.error(error)
      }
      return 2
    },
    async getTargets(datasource: Datasource) {
      try {
        const resp = await fetch(`/api/v1/datasources/targets/${datasource.name}`)
        const jsonString = await resp.json()
        const obj = JSON.parse(jsonString)
        const targets = obj.data.activeTargets as Target[]
        console.log('targets', targets)
        return targets
      } catch (error) {
        console.error(error)
      }
    }
  },
})
