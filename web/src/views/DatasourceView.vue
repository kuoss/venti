<script>
import LetterAvatar from '@/components/LetterAvatar.vue'
export default {
  components: {
    LetterAvatar,
  },
  data() {
    return {
      datasources: [],
      isLoading: false,
    }
  },
  methods: {
    refresh() {
      this.fetchData()
    },
    async fetchData() {
      try {
        this.isLoading = true
        const response = await fetch('/api/datasources')
        let datasources = await response.json()

        const response2 = await fetch('/api/datasources/targets')
        let data = await response2.json()
        data
          .map(x => JSON.parse(x))
          .forEach((x, i) => {
            if (x.status == 'error') {
              datasources[i].health = false
              datasources[i].targets = []
            } else {
              const targets = x.data.activeTargets.map(x => {
                x.age = ((new Date() - new Date(x.lastScrape)) / 1000).toFixed()
                x.job = x.discoveredLabels.job
                for (const [k, v] of Object.entries(x.discoveredLabels)) {
                  if (k == '__meta_kubernetes_service_name') {
                    x.icon = '🐕‍🦺'
                    x.name = v
                  } else if (k == '__meta_kubernetes_pod_name') {
                    x.icon = '🐱'
                    x.name = v
                  } else if (k == '__meta_kubernetes_node_name') {
                    x.icon = '🏝'
                    x.name = v
                  } else if (k == '__meta_kubernetes_namespace') {
                    x.icon = '🖼️'
                    x.name = v
                  }
                }
                x.icon ??= '💼'
                x.name ??= x.job
                return x
              })
              datasources[i].health = true
              datasources[i].targets = targets
            }
          })
        this.datasources = datasources
        this.isLoading = false
      } catch (error) {
        console.error(error)
      }
    },
  },
  mounted() {
    this.fetchData()
  },
}
</script>

<template>
  <header class="fixed right-0 w-full bg-white border-b shadow z-30 p-2 pl-52">
    <div class="flex items-center flex-row">
      <div><i class="mdi mdi-18px mdi-database-outline"></i> Datasource</div>
      <div class="flex ml-auto">
        <div class="inline-flex">
          <button
            @click="refresh()"
            class="h-rounded-group py-2 px-4 text-gray-900 bg-white border border-common hover:bg-gray-100 hover:text-blue-500 focus:text-blue-500"
          >
            <i class="mdi mdi-refresh"></i>
          </button>
        </div>
      </div>
    </div>
  </header>

  <div class="mt-12 w-full p-8">
    <h1 class="mt-4">Datasources</h1>
    <table class="w-full bg-white border" v-if="datasources">
      <tr class="border-b bg-slate-50" :class="{ 'is-loading': isLoading }">
        <th>Name</th>
        <th>Type</th>
        <th>URL</th>
        <th>Discovered</th>
        <th>Up</th>
      </tr>
      <tr class="border-b" v-for="d in datasources">
        <td class="px-2">
          <letterAvatar :bgcolor="$util.string2color(d.name)" />
          {{ d.name }}
        </td>
        <td>{{ d.type == 'prometheus' ? '🔥' : '💧' }} {{ d.type }}</td>
        <td>{{ d.url }}</td>
        <td class="text-center">{{ d.is_discovered ? '✔️' : '-' }}</td>
        <td class="text-center">
          <span :class="[d.health ? 'text-green-400' : 'text-red-400']">●</span>
        </td>
      </tr>
    </table>
    <h1 class="mt-4">Targets</h1>
    <table class="w-full bg-white border">
      <tr class="border-b bg-slate-50" :class="{ 'is-loading': isLoading }">
        <th>Datasource</th>
        <th>Job</th>
        <th>Address</th>
        <th>Name</th>
        <th>Last scrape</th>
        <th>Up</th>
      </tr>
      <template v-for="d of datasources">
        <tr class="border-b" v-for="t in d.targets.sort((a, b) => (a.job + a.name < b.job + b.name ? -1 : 1))">
          <td class="px-2" v-if="d.name">
            <LetterAvatar :bgcolor="$util.string2color(d.name)" />
            {{ d.name }}
          </td>
          <td>{{ t.discoveredLabels.job }}</td>
          <td>{{ t.discoveredLabels.__address__ }}</td>
          <td>{{ t.icon }} {{ t.name }}</td>
          <td class="text-right pr-10">{{ t.age }}s</td>
          <td class="text-center">
            <span :class="[t.health == 'up' ? 'text-green-400' : 'text-red-400']">●</span>
          </td>
        </tr>
      </template>
    </table>
  </div>
</template>
