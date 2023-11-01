import { defineStore } from 'pinia'
import type { Datasource, Target } from '@/types/datasource'
import { ref } from 'vue'

export const useDatasourceStore = defineStore('datasource', () => {

  const datasources = ref([] as Datasource[])
  const status = ref({ loaded: false, loading: false })

  // unexported
  async function fetchDatasources() {
    status.value.loading = true;
    try {
      const response = await fetch('/api/v1/datasources')
      datasources.value = await response.json()
      status.value.loaded = true
    } catch (error) {
      console.error(error)
      status.value.loading = false
    }
  }

  async function waitForLoaded() {
    let tries = 0;
    if (!status.value.loading) fetchDatasources()
    while (!status.value.loaded) {
      await new Promise(resolve => setTimeout(resolve, 100 * ++tries))
    }
  }

  // exported
  async function getDatasources() {
    await waitForLoaded()
    return datasources.value
  }

  async function getDatasourcesByType(typ: string) {
    const dss = await getDatasources()
    return dss.filter(x => x.type == typ)
  }

  function getDatasourceByName(name: string) {
    for (const ds of datasources.value as Datasource[]) {
      if (ds.name == name) {
        return ds
      }
    }
    return null
  }

  async function getDatasourceHealthy(datasource: Datasource): Promise<Number> {
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
  }

  async function getTargets(datasource: Datasource) {
    try {
      const resp = await fetch(`/api/v1/datasources/targets/${datasource.name}`)
      const jsonString = await resp.json()
      const obj = JSON.parse(jsonString)
      const targets = obj.data.activeTargets as Target[]
      return targets
    } catch (error) {
      console.error(error)
    }
  }
  return { getDatasources, getDatasourcesByType, getDatasourceByName, getDatasourceHealthy, getTargets }
})
