<script setup lang="ts">

</script>

<template>
  <el-menu
    mode="horizontal" router :default-active="$route.path"
    class="h-10 bg-transparent !border-0 flex gap-2"
  >
    <el-menu-item class="font-bold text-lg">
      {{$t('navbar.logo')}}
    </el-menu-item>
    <template v-for="route in $router.getRoutes()" :key="route.name">
      <el-menu-item v-if="route.meta.label" :index="route.path">
        {{$t(route.meta.label)}}
      </el-menu-item>
    </template>
  </el-menu>
</template>
