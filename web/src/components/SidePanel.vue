<script setup>
import { useSidePanelStore } from '@/stores/sidePanel'
import { XIcon } from '@heroicons/vue/solid'
import ButtonClipboard from '@/components/ButtonClipboard.vue'
import yaml from 'js-yaml'
   </script>

<template>
  <div :style="{ width: width }" v-if="show">
    <div
      style="height: calc(100vh - 56px)"
      :style="{ width: width }"
      class="fixed right-0 bottom-0 bg-slate-300"
    >
      <div class="h-[44px]">
        <button
          class="float-right px-2 py-1 cursor-pointer hover:bg-slate-400"
          @click="useSidePanelStore().close()"
        >
          <XIcon class="w-5 h-4" />
        </button>
        <div class="flex-1 h-full">
          <div v-if="type == 'DataTable'">
            <div class="font-bold text-center p-1">{{ dataTable.title }}</div>
            <div class="float-right px-2">{{ $util.dateTimeAsLocal(dataTable.time) }}</div>
            <div class="px-2">{{ dataTable.rows.length }} rows</div>
          </div>
          <div class="flex justify-center py-2" v-else>
            <span class="font-bold text-center p-1">{{ dashboardInfo.dashboardConfig.title }}</span>
            <span class="ml-2">
              <ButtonClipboard
                text="Dashboard"
                tooltipDirection="right"
                :value="yamlDashboard(dashboardInfo.dashboardConfig)"
                buttonClass="inline border-slate-400 hover:bg-slate-400"
              />
            </span>
          </div>
        </div>
      </div>
      <div
        class="overflow-y-auto border-l border-2 w-full bg-slate-300 border-b scrollbar-thin scrollbar-track-transparnt scrollbar-thumb-slate-400 dark:scrollbar-thumb-slate-500"
        style="height: calc(100vh - 100px)"
      >
        <div class="bg-slate-200 pb-8">
          <div v-if="type == 'DataTable'">
            <table class="w-full" v-if="dataTable.time">
              <tr v-for="row in dataTable.rows">
                <td class="pl-2" :style="{ 'color': $util.string2color(row[0]) }">■</td>
                <td class="break-all">{{ row[0] }}</td>
                <td class="pr-3 text-right">{{ Number.parseFloat(row[1]).toFixed(1) }}</td>
              </tr>
            </table>
          </div>
          <div v-else-if="type == 'PanelInfo'">
            <pre>{{ panelConfigYAML }}</pre>
          </div>
          <div v-else-if="type == 'DashboardInfo'">
            <table class="w-full font-mono">
              <tr class="align-top">
                <td class="w-6 text-center"></td>
                <td class="bg-slate-100 whitespace-pre-wrap">title: {{ dashboardInfo.dashboardConfig.title + '\nrows:' }}</td>
              </tr>
              <template v-for="(row, i) in dashboardInfo.dashboardConfig.rows">
                <tr class="border-b align-top">
                  <td></td>
                  <td class="bg-slate-100 whitespace-pre-wrap">- panels:</td>
                </tr>
                <tr class="border-b align-top" v-for="(panel, j) in row.panels">
                  <td class="text-center">
                    <div class="p-1 bg-cyan-100" style="user-select: none">{{ i + 1 }}{{ j + 1 }}</div>
                  </td>
                  <td
                    :ref="`ref${i + 1}${j + 1}`"
                    class="bg-slate-100 highlight-base transition-colors duration-[5000ms]"
                  >
                    <div class="float-right my-1 mr-3">
                      <ButtonClipboard :value="yamlPanel(panel)" buttonClass="hover:bg-slate-300" />
                    </div>
                    <pre class="whitespace-pre-wrap">  - title: {{ panel.title }}</pre>
                    <pre class="whitespace-pre-wrap">    type: {{ panel.type }}</pre>
                    <template v-if="panel.targets">
                      <pre class="whitespace-pre-wrap">    targets:</pre>
                      <template v-for="target in panel.targets">
                        <pre class="whitespace-pre-wrap" v-html="yamlTarget(target, panel.type)"></pre>
                      </template>
                    </template>
                    <pre
                      class="whitespace-pre-wrap"
                      v-if="panel.chartOptions"
                      v-html="yamlChartOptions(panel.chartOptions)"
                    ></pre>
                  </td>
                </tr>
              </template>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      show: false,
      type: '',
      dataTable: {},
      panelInfo: {},
      dashboardInfo: {},
      currentPosition: null,
    }
  },
  computed: {
    width() {
      switch (this.type) {
        case 'DataTable': return '300px'
      }
      return '600px'
    },
    title() {
      switch (this.type) {
        case 'DataTable': return this.dataTable.title
        case 'DashboardInfo': return this.dashboardInfo.dashboardConfig.title
        case 'PanelInfo': return this.panelInfo.panelConfig.title
      }
    },
  },
  methods: {
    dumpYAML(j, flowLevel) { return yaml.dump(j, { noArrayIndent: true, flowLevel: flowLevel }).replaceAll('>-', '|').replaceAll('>', '|') },
    yamlDashboard(j) { return this.dumpYAML(j, 5) },
    yamlPanel(j) { return this.indentText(this.dumpYAML([j], 4), 2) },
    yamlTarget(j, type) {
      let x = this.cloneObject(j)
      const expr = x.expr
      const path = (type == 'logs') ? 'logs' : 'metrics'
      x.expr = 'DuMmYeXpR'
      x = this.dumpYAML([x], 3)
      x = x.replace('DuMmYeXpR', `<span class="text-cyan-500"><a class="hover:underline" href="/#/${path}?query=${encodeURIComponent(expr)}">${yaml.dump(expr).replaceAll('>-', '|').replaceAll('>', '|').replaceAll("\n", "\n  ").trimRight()}</a></span>`)
      return this.indentText(x, 4).replaceAll('$node', '<span class="text-yellow-500 font-bold">$node</span>').replaceAll('$namespace', '<span class="text-green-500 font-bold">$namespace</span>')
    },
    yamlChartOptions(j) { return "  chartOptions:\n" + this.indentText(this.dumpYAML(j), 4) },
    indentText(t, level) {
      return t.split("\n").map(x => ' '.repeat(level) + x).join("\n").trimRight()
    },
    cloneObject(o) {
      return JSON.parse(JSON.stringify(o))
    },
    goToPanelConfig(position, retries = 0) {
      if (this.type != 'DashboardInfo') {
        console.log('this.type=', this.type)
        return
      }
      const el = this.$refs[`ref${position}`]
      if (!el || !el[0] || !el[0]?.classList) return
      el[0].scrollIntoView()
      el[0].classList.add('highlight')
      setTimeout(() => {
        if (!el || !el[0] || !el[0]?.classList) return
        el[0].classList.remove('highlight')
      }, 5000)
    },
  },
  mounted() {
    useSidePanelStore().$subscribe((mutation, state) => {
      const needResize = (this.show != state.show) || (this.type != state.type)
      this.show = state.show
      this.type = state.type
      this.dataTable = state.dataTable
      this.panelInfo = state.panelInfo
      this.dashboardInfo = state.dashboardInfo
      if (state.currentPosition && this.currentPosition != state.currentPosition) {
        // console.log('SidePanel.vue: mounted: $subscribe: state=', state)
        this.goToPanelConfig(state.currentPosition)
      }
      this.currentPosition = state.currentPosition
      if (needResize) this.$emit('resize')
    })
  }
}
</script>
