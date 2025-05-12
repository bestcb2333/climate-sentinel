<script setup lang="ts">
import {request} from '@/axios';
import type {Region} from '@/tables';
import {ref} from 'vue';
import RegionMap from '@/components/RegionMap.vue';

const series = ref<{
  name: string
  value: number
}[]>([])

request.get<any, {data: Region[]}>('/regions?history=true').then(res => {
  series.value = res.data.map(region => ({
  name: region.name,
  value: getAverage(region.histories.map(history => Math.max(
    riskIndex(history.maxTemperature, 24, 10),
    riskIndex(history.minTemperature, 24, 10),
    riskIndex(history.windSpeed, 0, 5),
    riskIndex(history.rainfall??0, 0, 5),
    riskIndex(history.visibility, 0, 5)
  ))),
}))
}).catch(() => {})

function getAverage(numbers: number[]): number {
  if (numbers.length === 0) return 0; // 避免除以 0
  const sum = numbers.reduce((acc, val) => acc + val, 0);
  return sum / numbers.length;
}
function riskIndex(value: number, optimal: number, k: number): number {
    const deviation = Math.abs(value - optimal);
    return deviation / (deviation + k);
}
</script>

<template>
  <div class="h-full grid grid-cols-[3fr_3fr_2fr] grid-rows-[1fr_3fr] gap-2">
    <el-card class="row-span-2" shadow="never">
    </el-card>
    <el-card shadow="never">
    </el-card>
    <el-card class="row-span-2" shadow="never" body-class="h-full">
      <region-map :series="series" />
    </el-card>
    <el-card shadow="never">
    </el-card>
  </div>
</template>
