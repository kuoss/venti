<template>
  <transition name="auth" mode="out-in">
    <div v-if="$auth.loggedIn" class="flex h-full w-full">
      <div class="m-auto">
        <span>Welcome, {{ $auth.username }}</span>
      </div>
    </div>
  </transition>
</template>
