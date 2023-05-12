<template>
  <header class="fixed right-0 w-full bg-white border-b border-common shadow z-30 p-2 pl-52">
    <div class="flex items-center flex-row">
      <div><i class="mdi mdi-18px mdi-file-outline" /> Config</div>
    </div>
  </header>
  <main class="mt-12 w-full p-4">
    <div class="m-10 bg-white">
      <div class="border border-common">
        <div class="p-3">
          <p>dashboards.yaml</p>
        </div>
      </div>
      <div class="border border-common">
        <div class="p-3">
          <pre class="text-sm">{{ yaml }}</pre>
        </div>
      </div>
    </div>
  </main>
</template>

<script>
export default {
  data() {
    return {
      yaml: '',
    };
  },
  mounted() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      this.axios.get('/api/config/dashboards/yaml').then(response => {
        this.yaml = response.data.yaml;
      });
    },
  },
};
</script>
