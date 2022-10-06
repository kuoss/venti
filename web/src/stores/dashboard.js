import { defineStore } from 'pinia'
import axios from 'axios'

export const useDashboardStore = defineStore('dashboard', {
    state: () => ({
        dashboards: [],
        status: { loaded: false, loading: false },
    }),
    actions: {
        async getDashboards() {
            await this.waitForLoaded()
            return this.dashboards
        },
        async waitForLoaded() {
            let tries = 0
            if (!this.status.loading) this.loadData()
            while (!this.status.loaded) await new Promise((resolve) => setTimeout(resolve, 100 * (++tries)))
        },
        async loadData() {
            this.status.loading = true
            try {
                const response = await axios.get('api/config/dashboards')
                this.dashboards = response.data
                this.status.loaded = true
            } catch (error) {
                console.error(error)
                this.status.loading = false
            }
        },
    },
})