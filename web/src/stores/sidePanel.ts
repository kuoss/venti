import { defineStore } from 'pinia';
import type DataTable from '@/types/datatable'

export const useSidePanelStore = defineStore('sidePanel', {
  state: () => ({
    show: false,
    title: '',
    type: '',
    dataTable: {},
    dashboardInfo: {},
    currentPosition: '',
  }),
  actions: {
    close() {
      this.show = false;
    },
    updatetDataTable(t: DataTable) {
      this.dataTable.title = t.title;
      this.dataTable.time = t.time;
      this.dataTable.rows = t.rows.sort((a, b) => (a[1] > b[1] ? -1 : 1));
    },
    updateDashboardInfo(dashboardConfig: any) {
      this.dashboardInfo = { dashboardConfig: dashboardConfig };
    },
    toggleShow(type: string) {
      if (this.type != type) {
        this.type = type;
        this.show = true;
        return;
      }
      this.show = !this.show;
    },
    toggleDashboardInfo() {
      if (this.type != 'DashboardInfo') this.show = true;
      else this.show = !this.show;
      if (this.show) this.type = 'DashboardInfo';
    },
    goToPanelConfig(position: string) {
      this.type = 'DashboardInfo';
      this.show = true;
      setTimeout(() => (this.currentPosition = position), 300);
    },
  },
});
