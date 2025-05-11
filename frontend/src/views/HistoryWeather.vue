<script setup lang="ts">
import {request} from '@/axios';
import type {History} from '@/tables';
import {Menu} from '@element-plus/icons-vue';
import {LineChart} from 'echarts/charts';
import {GridComponent, LegendComponent, TitleComponent, TooltipComponent} from 'echarts/components';
import {use} from 'echarts/core';
import {CanvasRenderer} from 'echarts/renderers';
import {computed, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import VChart from 'vue-echarts';
import useSessionStore from '@/stores/session';
import dayjs from 'dayjs';

use([
  LineChart,
  TitleComponent,
  LegendComponent,
  TooltipComponent,
  GridComponent,
  CanvasRenderer,
])

const session = useSessionStore()


const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const histories = ref<History[]>([])
watch([page, pageSize], async ([page, pageSize]) => {
  const res = await request.get<any, {
    data: History[],
    total: number,
  }>('/histories', {params: {
    page, pageSize,
  }})
  total.value = res.total
  histories.value = res.data
}, {immediate: true})

interface Trend {
  name: string,
  series: {
    name: string,
    type: string,
    stack: string,
    data: number[],
  }[],
}

const weatherTypes: Record<string, string[]> = {
  sunny: ['晴天', '#FFD700'],
  rainy: ['雨天', '#4A90E2'],
  cloudy: ['多云', '#B0BEC5'],
  foggy: ['雾天', '#CFD8DC'],
  snowy: ['下雪', '#E0F7FA'],
  windy: ['有风', '#81D4FA'],
  overcast: ['阴天', '#78909C'],
}

const board = computed(() => [
  ['晴天', 'sunny'],
  ['雨天', 'rainy'],
  ['多云', 'cloudy'],
  ['雾天', 'foggy'],
  ['下雪', 'snowy'],
  ['有风', 'windy'],
  ['阴天', 'overcast'],
].map(([label, id]) => ({
  label,
  value: histories.value.filter(history => history.type===id).length,
})))

const trends = ref<Trend[]>([])
request.get<any, Trend[]>('/trends').then(res => {
  trends.value = res
}).catch(() => {})

const lineChartOptions = computed(() => trends.value.map(trend => ({
  title: {
    text: trend.name
  },
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    data: session.regions.map(region => region.name)
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
    containLabel: true
  },
  toolbox: {
    feature: {
      saveAsImage: {}
    }
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  },
  yAxis: {
    type: 'value'
  },
  series: trend.series
})))

const {t} = useI18n({messages: {
  zh: {
    region: '区域',
    type: '天气类型',
    time: '时间',
    temperature: '气温(最高/平均/最低) ℃',
    windSpeed: '风速',
    rainFall: '降水量',
    visibility: '能见度',
    severity: '严重性',
    source: '数据源',
  },
}})
</script>

<template>
  <div class="h-full grid grid-cols-2 grid-rows-[auto_1fr] gap-2">

    <el-card shadow="hover" body-class="h-full flex justify-around items-center">
      <div v-for="item in board" :key="item.label" class="flex flex-col items-center">
        <div class="font-bold">
          {{item.label}}
        </div>
        <div>
          {{item.value}}
        </div>
      </div>
    </el-card>

    <el-card shadow="never"
      class="row-span-2 flex flex-col" body-class="grow space-y-2 overflow-y-auto"
    >
      <el-card v-for="option in lineChartOptions" :key="option.title.text" shadow="hover"
        class="aspect-video" body-class="h-full"
      >
        <v-chart :option="option" />
      </el-card>
    </el-card>

    <el-card shadow="never" class="grow basis-0 flex flex-col" body-class="grow overflow-y-auto" footer-class="flex justify-end">
      <template #header>
        <card-title :title="t('title')" :icon="Menu">
          <el-button circle :icon="Menu" />
          <el-button circle :icon="Menu" />
        </card-title>
      </template>

      <el-table :data="histories">
        <el-table-column
          :label="t('time')"
          prop="time"
          :formatter="(_1, _2, date) => dayjs(date).format('MM月DD日')"
        />
        <el-table-column :label="t('region')" prop="region.name" />
        <el-table-column :label="t('type')" prop="type">
          <template #default="{row}">
            <el-tag class="text-black" :color="weatherTypes[row.type][1]">
              {{weatherTypes[row.type][0]}}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('temperature')" width="120">
          <template #default="{row}">
            {{row.maxTemperature}}/{{row.avgTemperature}}/{{row.minTemperature}}
          </template>
        </el-table-column>
        <el-table-column :label="t('windSpeed')" prop="windSpeed" />
        <el-table-column :label="t('visibility')" prop="visibility" />
        <el-table-column :label="t('rainFall')" prop="rainFall" />
        <el-table-column :label="t('source')" prop="source" />
      </el-table>

      <template #footer>
        <el-pagination layout="sizes, prev, pager, next, total"
          :total="total"
          v-model:current-page="page"
          v-model:page-size="pageSize"
        />
      </template>
    </el-card>

  </div>
</template>
