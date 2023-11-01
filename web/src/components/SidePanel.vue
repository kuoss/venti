<script setup>
import { useSidePanelStore } from '@/stores/sidePanel';
import { XIcon } from '@heroicons/vue/solid';
import ButtonClipboard from '@/components/ButtonClipboard.vue';
import yaml from 'js-yaml';
import Util from '@/lib/util';
import { computed, ref, onMounted, watch } from 'vue';

const sidePanelStore = useSidePanelStore()

const show = ref(false)
const type = ref('')
const dataTable = ref({})
const panelInfo = ref({})
const dashboardInfo = ref({})
const currentPosition = ref('')

const panelRefs = ref([])

const width = computed(() => {
  switch (type.value) {
    case 'DataTable':
      return '300px';
  }
  return '600px';
})

const title = computed(() => {
  switch (type.value) {
    case 'DataTable':
      return dataTable.value.title;
    case 'DashboardInfo':
      return dashboardInfo.value.dashboardConfig.title;
  }
  return panelInfo.value.panelConfig.title;
})

const emit = defineEmits(['resize'])

function dumpYAML(j, flowLevel) {
  return yaml.dump(j, { noArrayIndent: true, flowLevel: flowLevel }).replaceAll('>-', '|').replaceAll('>', '|');
}

function yamlDashboard(j) {
  return dumpYAML(j, 5);
}

function yamlPanel(j) {
  return indentText(dumpYAML([j], 4), 2);
}

function yamlTarget(j, type) {
  let x = cloneObject(j);
  const expr = x.expr;
  const path = type == 'logs' ? 'logs' : 'metrics';
  x.expr = 'DUMMY_EXPR';
  x = dumpYAML([x], 3);
  x = x.replace(
    'DUMMY_EXPR',
    `<span class="text-cyan-500"><a class="hover:underline" href="/#/${path}?query=${encodeURIComponent(
      expr,
    )}">${yaml
      .dump(expr)
      .replaceAll('>-', '|')
      .replaceAll('>', '|')
      .replaceAll('\n', '\n  ')
      .trimRight()}</a></span>`,
  );
  return indentText(x, 4)
    .replaceAll('$node', '<span class="text-yellow-500 font-bold">$node</span>')
    .replaceAll('$namespace', '<span class="text-green-500 font-bold">$namespace</span>');
}

function yamlChartOptions(j) {
  return '  chartOptions:\n' + indentText(dumpYAML(j), 4);
}

function indentText(t, level) {
  return t
    .split('\n')
    .map(x => ' '.repeat(level) + x)
    .join('\n')
    .trimRight();
}

function cloneObject(o) {
  return JSON.parse(JSON.stringify(o));
}

function goToPanelConfig() {
  const el = panelRefs.value[Number(currentPosition.value)];
  el.scrollIntoView();
  el.classList.add('highlight');
  setTimeout(() => {
    el.classList.remove('highlight');
  }, 5000);
}

sidePanelStore.$subscribe((_, state) => {
  const needResize = show.value != state.show || type.value != state.type;
  const positionChanged = currentPosition.value != state.currentPosition;
  show.value = state.show;
  type.value = state.type;
  dataTable.value = state.dataTable;
  panelInfo.value = state.panelInfo;
  dashboardInfo.value = state.dashboardInfo;
  currentPosition.value = state.currentPosition;
  if (needResize) emit('resize');
  if (positionChanged) goToPanelConfig();
})
</script>

<template>
  <div v-if="show" :style="{ width: width }">
    <div style="height: calc(100vh - 56px)" :style="{ width: width }" class="fixed right-0 bottom-0 bg-slate-300">
      <div class="h-[44px]">
        <button class="float-right px-2 py-1 cursor-pointer hover:bg-slate-400" @click="sidePanelStore.close()">
          <XIcon class="w-5 h-4" />
        </button>
        <div class="flex-1 h-full">
          <div v-if="type == 'DataTable'">
            <div class="font-bold text-center p-1">
              {{ dataTable.title }}
            </div>
            <div class="float-right px-2">
              {{ Util.dateTimeAsLocal(dataTable.time) }}
            </div>
            <div class="px-2">{{ dataTable.rows.length }} rows</div>
          </div>
          <div v-else class="flex justify-center py-2">
            <span class="font-bold text-center p-1">{{ dashboardInfo.dashboardConfig.title }}</span>
            <span class="ml-2">
              <ButtonClipboard text="Dashboard" tooltip-direction="right"
                :value="yamlDashboard(dashboardInfo.dashboardConfig)"
                button-class="inline border-slate-400 hover:bg-slate-400" />
            </span>
          </div>
        </div>
      </div>
      <div
        class="overflow-y-auto border-l border-2 w-full bg-slate-300 border-b scrollbar-thin scrollbar-track-transparnt scrollbar-thumb-slate-400 dark:scrollbar-thumb-slate-500"
        style="height: calc(100vh - 100px)">
        <div class="bg-slate-200 pb-8">
          <div v-if="type == 'DataTable'">
            <table v-if="dataTable.time" class="w-full">
              <tr v-for="row in dataTable.rows">
                <td class="pl-2" :style="{ color: Util.string2color(row[0]) }">■</td>
                <td class="break-all">
                  {{ row[0] }}
                </td>
                <td class="pr-3 text-right">
                  {{ Number.parseFloat(row[1]).toFixed(1) }}
                </td>
              </tr>
            </table>
          </div>
          <div v-else-if="type == 'PanelInfo'">
            <pre>{{ panelConfigYAML }}</pre>
          </div>
          <div v-else-if="type == 'DashboardInfo'">
            <table class="w-full font-mono">
              <tr class="align-top">
                <td class="w-6 text-center" />
                <td class="bg-slate-100 whitespace-pre-wrap">title:
                  {{ dashboardInfo.dashboardConfig.title + '\nrows:' }}
                </td>
              </tr>
              <template v-for="(row, i) in dashboardInfo.dashboardConfig.rows">
                <tr class="border-b align-top">
                  <td />
                  <td class="bg-slate-100 whitespace-pre-wrap">- panels:</td>
                </tr>
                <tr v-for="(panel, j) in row.panels" class="border-b align-top">
                  <td class="text-center">
                    <div class="p-1 bg-cyan-100" style="user-select: none">{{ i + 1 }}{{ j + 1 }}</div>
                  </td>
                  <td class="bg-slate-100 highlight-base transition-colors duration-[5000ms]"
                    :ref="(el) => { panelRefs[10 * (i + 1) + (j + 1)] = el }">
                    <div class="float-right my-1 mr-3">
                      <ButtonClipboard :value="yamlPanel(panel)" button-class="hover:bg-slate-300" />
                    </div>
                    <pre class="whitespace-pre-wrap">  - title: {{ panel.title }}</pre>
                    <pre class="whitespace-pre-wrap">    type: {{ panel.type }}</pre>
                    <template v-if="panel.targets">
                      <pre class="whitespace-pre-wrap">    targets:</pre>
                      <template v-for="target in panel.targets">
                        <pre class="whitespace-pre-wrap" v-html="yamlTarget(target, panel.type)" />
                      </template>
                    </template>
                    <pre v-if="panel.chartOptions" class="whitespace-pre-wrap"
                      v-html="yamlChartOptions(panel.chartOptions)" />
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

<style scoped>
.highlight-base {
  @apply relative z-[1];
}

.highlight-base::before {
  content: '';
  @apply bg-gradient-to-r from-slate-100 to-yellow-100 absolute top-0 left-0 w-full h-full opacity-0 z-[-1] transition-opacity duration-500;
}

.dark .highlight-base::before {
  @apply bg-gradient-to-r from-slate-800 to-teal-800;
}

.highlight::before {
  @apply opacity-100;
}
</style>