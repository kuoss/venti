import { defineStore } from 'pinia'
import axios from 'axios'

export const useFilterStore = defineStore('filter', {
    state: () => ({
        namespaces: [],
        nodes: [],
        selectedNamespace: 'All namespaces',
        selectedNode: 'All nodes',
        status: { loaded: false, loading: false },
    }),
    actions: {
        async getNamespaces() {
            await this.waitForLoaded()
            return this.namespaces
        },
        async getNodes() {
            await this.waitForLoaded()
            return this.nodes
        },
        async waitForLoaded() {
            let tries = 0
            if (!this.status.loading) this.loadData()
            while (!this.status.loaded) await new Promise((resolve) => setTimeout(resolve, 100 * (++tries)))
        },
        async loadData() {
            this.status.loading = true
            try {
                const response = await axios.get('/api/prometheus/namespaces')
                const namespaces = response.data.data.result.map(x => x.metric.namespace)
                this.namespaces = ['All namespaces', ...namespaces]

                const response2 = await axios.get('/api/prometheus/nodes')
                const nodes = response2.data.data.result.map(x => x.metric.node)
                this.nodes = ['All nodes', ...nodes]

                this.status.loaded = true
            } catch (error) {
                console.error(error)
                this.status.loading = false
            }
        },
        renderExpr(expr) {
            return expr.replaceAll(/\$namespace/g, this.selectedNamespace == 'All namespaces' ? '.*' : this.selectedNamespace).replaceAll(/\$node/g, this.selectedNode == 'All nodes' ? '.*' : this.selectedNode)
        },

    }
})
