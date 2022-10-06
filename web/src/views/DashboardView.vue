<script setup>
import Dropdown from "@/components/Dropdown.vue"
import Panel from '@/components/Panel.vue'
import RunButton from "@/components/RunButton.vue"
import SidePanel from '@/components/SidePanel.vue'
import TimeRangePicker from "@/components/TimeRangePicker.vue"

import { useDashboardStore } from "@/stores/dashboard";
import { useFilterStore } from "@/stores/filter";
import { useSidePanelStore } from '@/stores/sidePanel'
import { useTimeStore } from "@/stores/time";
</script>

<template>
  <nav class="fixed left-0 w-full bg-white border-b shadow p-2 pl-52 pr-4 z-10">
    <div class="flex gap-2">
      <div class="flex-none">
        <Dropdown
          :options="namespaces"
          :currentDropdown="currentDropdown"
          @select="selectNamespace"
          @open="onDropdownOpen"
        ></Dropdown>
      </div>
      <div class="flex-none">
        <Dropdown
          :options="nodes"
          :currentDropdown="currentDropdown"
          @select="selectNode"
          @open="onDropdownOpen"
        ></Dropdown>
      </div>
      <div class="grow font-bold text-center align-middle">
        {{ dashboard.title }}
        <button
          class="ml-2 px-2 py-2 border rounded"
          @click="useSidePanelStore().toggleDashboardInfo()"
        >
          <i class="mdi mdi-information-outline"></i>
        </button>
      </div>
      <div class="flex-none">
        <TimeRangePicker @updateTimeRange="updateTimeRange" />
      </div>
      <div class="flex-none">
        <RunButton
          :disabled="false"
          btnText="Refresh"
          @execute="execute"
          @changeInterval="changeInterval"
        />
      </div>
    </div>
  </nav>
  <div class="flex w-full mt-12 text-xs">
    <div class="flex-1" ref="root">
      <div class="p-4 pt-5">
        <div
          v-for="(row, i) in dashboard.rows"
          class="pb-3 grid gap-3"
          :class="['flex', ['', 'grid-cols-1', 'grid-cols-2', 'grid-cols-3', 'grid-cols-4', 'grid-cols-5', 'grid-cols-6', 'grid-cols-7', 'grid-cols-8', 'grid-cols-9', 'grid-cols-10', 'grid-cols-11', 'grid-cols-12',][row.panels.length]]"
        >
          <div class="flex-1 bg-white border" v-for="(panel, j) in row.panels">
            <Panel
              :position="`${i+1}${j+1}`"
              :count="count"
              :panelConfig="panel"
              :panelWidth="clientWidth / row.panels.length"
              :timeRange="timeRange"
            />
          </div>
        </div>
      </div>
    </div>
    <SidePanel @resize="resize" />
  </div>
</template>

<script>
export default {
  data() {
    return {
      count: 0,
      range: [],
      timeRange: [],
      intervalSeconds: 0,
      currentDropdown: -1,
      dashboards: [],
      dashboard: {},
      namespaces: [],
      nodes: [],
      clientWidth: 100,
      rows: [],
    }
  },
  methods: {
    updateTimeRange(r) {
      this.range = r
    },
    async execute() {
      this.timeRange = await useTimeStore().toTimeRangeForQuery(this.range)
      this.count++
      if (this.intervalSeconds > 0) {
        setTimeout(() => this.timerHandler(), this.intervalSeconds * 1000)
      }
    },
    timerHandler() {
      if (useTimeStore().timerManager != 'DashboardView' || this.intervalSeconds == 0) return
      this.execute()
    },
    changeInterval(i) {
      this.intervalSeconds = i
      this.execute()
    },
    onDropdownOpen(uid) {
      this.currentDropdown = uid
    },
    selectNamespace(ns) {
      useFilterStore().selectedNamespace = ns
      this.execute()
    },
    selectNode(node) {
      useFilterStore().selectedNode = node
      this.execute()
    },
    mousemove(e) {
      this.legendPosition = [e.clientX + 50, e.clientY + 50]
    },
    resize() {
      // wait for browser rendering
      setTimeout(() => { this.clientWidth = this.$refs.root.clientWidth }, 80)
    },
    async init() {
      this.namespaces = await useFilterStore().getNamespaces()
      this.nodes = await useFilterStore().getNodes()
      this.dashboards = await useDashboardStore().getDashboards()
      this.renderDashboard()
      this.execute()
      this.resize()
    },
    renderDashboard() {
      this.dashboards.forEach((d) => {
        if (d.title == this.$route.params.name) {
          this.dashboard = d
          useSidePanelStore().updateDashboardInfo(d)
        }
      })
    },
  },
  mounted() {
    useTimeStore().timerManager = 'DashboardView'
    this.$watch(
      () => this.$route.params, () => { this.renderDashboard() },
    )
    window.addEventListener("resize", this.resize)
    this.init()
  },
  unmounted() {
    window.removeEventListener("resize", this.resize)
  },
}
</script>
