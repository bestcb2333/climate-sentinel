<script setup lang="ts">
import NavBar from './components/NavBar.vue';
import usePersistedStore from './stores/persisted';
import useSessionStore from './stores/session';

const session = useSessionStore()
session.loadMap()
const persisted = usePersistedStore()
if (persisted.token) {
  session.loadUser()
}
</script>

<template>
  <div class="h-screen flex flex-col bg-slate-100 p-3 gap-3">
    <nav-bar />
    <div class="grow overflow-y-auto rounded-2xl">
      <router-view />
    </div>
  </div>
</template>
