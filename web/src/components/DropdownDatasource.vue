<script setup lang="ts">
import { useDatasourceStore } from '@/stores/datasource';
import Dropdown from './Dropdown.vue';
import { ref } from 'vue';
import { type Datasource } from '@/types/datasource'

const props = defineProps({
  dsType: { type: String, required: true }
})

const emit = defineEmits<{
  (e: 'change', value: String): void
}>()

const dsStore = useDatasourceStore();
const options = ref([]);

async function fetchData() {
  const datasources = await dsStore.getDatasourcesByType(props.dsType)
  // @ts-ignore
  options.value = datasources.map((x: Datasource) => x.name)
  onChange(options.value[0])
}

function onChange(value: String): any {
  emit('change', value)
}

fetchData()
</script>

<template>
  <Dropdown :options="options" @change="onChange" />
</template>