<script setup lang="ts">
import {request} from '@/axios';
import type {Notice, Region} from '@/tables';
import {computed, ref, watch} from 'vue';
import RegionMap from '@/components/RegionMap.vue';
import {useI18n} from 'vue-i18n';
import {useRouteQuery} from '@vueuse/router';
import {use} from 'echarts/core';
import {GridComponent, TitleComponent, TooltipComponent} from 'echarts/components';
import {CanvasRenderer} from 'echarts/renderers';
import {BarChart} from 'echarts/charts';
import VChart from 'vue-echarts';
import type {EChartsOption} from 'echarts';

use([
  TitleComponent,
  TooltipComponent,
  BarChart,
  GridComponent,
  CanvasRenderer,
])

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

const {t} = useI18n({messages: {
  zh: {
    noticeTitle: '公告列表',
  },
}})

const notices = ref<Notice[]>([])
const total = ref(0)
const page = useRouteQuery('page', 1, {transform: Number})
const pageSize = useRouteQuery('page_size', 10, {transform: Number})
watch([page, pageSize], loadTable, {immediate: true})
async function loadTable() {
  try {
    const res = await request.get<any, {
      total: number,
      data: Notice[],
    }>('/notices', {params: {
      page: page.value,
      page_size: pageSize.value,
    }})
    notices.value = res.data
    total.value = res.total
  } catch {}
}

const chartOption = computed<EChartsOption>(() => ({
  title: {
    text: '各地区风险指数排行',
  },
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'shadow',
    },
  },
  xAxis: {
    type: 'category',
    data: series.value.map(serie => serie.name),
  },
  yAxis: {
    type: 'value',
  },
  series: [{
    data: series.value.map(serie => serie.value),
    type: 'bar',
  }],
}))
</script>

<template>
  <div class="h-full grid grid-cols-[3fr_3fr_2fr] grid-rows-[1fr_3fr] gap-2">
    <el-card class="row-span-2 flex flex-col" shadow="never"
      header-class="font-bold"
      body-class="grow space-y-2 overflow-y-auto min-h-0"
    >
      <template #header>
        {{t('noticeTitle')}}
      </template>
      <el-card shadow="hover" v-for="notice in notices" :key="notice.id"
        header-class="flex justify-between"
      >
        <template #header>
          <div class="font-bold">{{notice.title}}</div>
          <div>创建者:{{notice.user?.name}}</div>
        </template>
        {{notice.content}}
      </el-card>
      <template #footer>
        <el-pagination layout="sizes, prev, pager, next, total" :total="total"
          v-model:current-page="page" v-model:page-size="pageSize"
        />
      </template>
    </el-card>
    <el-card shadow="never" body-class="h-full flex flex-col justify-center items-center">
      <div class="font-bold text-lg">
        极端天气城市道路救援管理信息系统
      </div>
      <div class="font-bold text-xl">
        Climate Sentinel
      </div>
    </el-card>
    <el-card class="row-span-2" shadow="never" body-class="h-full" header-class="font-bold">
      <template #header>
        武汉市地图
      </template>
      <region-map />
    </el-card>
    <el-card shadow="never" body-class="h-full">
      <v-chart :option="chartOption" autoresize />
    </el-card>
  </div>
</template>
