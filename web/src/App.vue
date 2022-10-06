<template>
  <transition name="layout" mode="out-in">
    <component :is="$route.meta.layout ?? 'nosidebar'">
      <router-view v-slot="{ Component }">
        <transition name="route" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </component>
  </transition>
</template>

<style>
@import "@/assets/base.css";
@import "@mdi/font/css/materialdesignicons.css";
</style>