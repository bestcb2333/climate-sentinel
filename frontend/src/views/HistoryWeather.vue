<script setup lang="ts">
import {request} from '@/axios';
import type {History} from '@/tables';
import {LineChart} from 'echarts/charts';
import {GridComponent, LegendComponent, TitleComponent, TooltipComponent} from 'echarts/components';
import {use} from 'echarts/core';
import {CanvasRenderer} from 'echarts/renderers';
import {computed, reactive, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import VChart from 'vue-echarts';
import useSessionStore from '@/stores/session';
import dayjs from 'dayjs';
import type {StatsItem} from '@/types';
import {useRouteQuery} from '@vueuse/router';

use([
  LineChart,
  TitleComponent,
  LegendComponent,
  TooltipComponent,
  GridComponent,
  CanvasRenderer,
])

const session = useSessionStore()

const futureOnly = useRouteQuery('future', 'false', {transform: Boolean})
const regionId = useRouteQuery('region_id', '0', {transform: Number})
const type = useRouteQuery('type', 'all')
const page = useRouteQuery('page', '1', {transform: Number})
const pageSize = useRouteQuery('page_size', '10', {transform: Number})
const total = ref(0)
const histories = ref<History[]>([])
async function loadTable() {
  const res = await request.get<any, {
    data: History[],
    total: number,
  }>(`/histories?page=${page.value}&page_size=${pageSize.value}&region_id=${regionId.value}&future=${futureOnly.value}&type=${type.value}`)
  total.value = res.total
  histories.value = res.data
}
watch([page, pageSize, regionId, futureOnly, type], loadTable, {immediate: true})

interface Trend {
  name: string,
  series: {
    name: string,
    type: string,
    stack: string,
    data: number[],
  }[],
}

const types = [
  'sunny', 'rainy', 'cloudy', 'foggy', 'snowy', 'windy', 'overcast',
]

const typeColorMap: Record<string, string> = {
  sunny: '#FFD700',
  rainy: '#4A90E2',
  cloudy: '#B0BEC5',
  foggy: '#CFD8DC',
  snowy: '#E0F7FA',
  windy: '#81D4FA',
  overcast: '#78909C',
}

const trends = ref<Trend[]>([])
request.get<any, Trend[]>('/trends').then(res => {
  trends.value = res
}).catch(() => {})

const lineChartOptions = computed(() => trends.value.map((trend) => ({
  title: {
    text: trend.name
  },
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    top: '10%',
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
    tableTitle: '天气状况列表',
    add: '添加记录',
    region: '区域',
    type: '天气类型',
    time: '时间',
    temperature: '气温(最高/平均/最低) ℃',
    windSpeed: '风速',
    rainFall: '降水量',
    visibility: '能见度',
    severity: '严重性',
    source: '数据源',
    futureOnly: '仅显示未来预测',
    allRecords: '显示所有数据',
    allRegions: '所有区域',
    all: '所有天气',
    sunny: '晴天',
    rainy: '雨天',
    cloudy: '多云',
    foggy: '雾天',
    snowy: '下雪',
    windy: '有风',
    overcast: '阴天',
  },
}})

const StatsItems = ref<StatsItem[]>([])
request.get<any, StatsItem[]>('/stats/histories').then(res => {
  StatsItems.value =res
}).catch(() => {})

const isDialogOpen = ref(false)

const addItemForm = reactive({
  time: new Date(),
  regionId: 1,
  type: 'sunny',
  maxTemperature: 30,
  minTemperature: 10,
  avgTemperature: 20,
  windSpeed: 3,
  visibility: 5,
  rainFall: 0,
})

async function addItem() {
  try {
    await request.post('/histories', addItemForm)
    loadTable()
  } catch {}
}
</script>

<template>
  <div class="h-full grid grid-cols-[auto_1fr_1fr] grid-rows-[auto_1fr] gap-2">

    <div class="row-span-2 flex flex-col justify-between">
      <el-button type="primary" round @click="isDialogOpen=true">
        {{t('add')}}
      </el-button>
      <el-segmented :options="['all', ...types].map(type=>({label:t(type),value:type}))" direction="vertical" v-model="type">

      </el-segmented>
    </div>

    <el-card shadow="hover" body-class="h-full flex justify-around items-center">
      <div v-for="item in StatsItems" :key="item.label" class="flex flex-col items-center">
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
        <v-chart :option="option" autoresize />
      </el-card>
    </el-card>

    <el-card shadow="never" class="grow basis-0 flex flex-col"
      header-class="flex"
      body-class="grow overflow-y-auto" footer-class="flex justify-end"
    >

      <template #header>
        <div>{{t('tableTitle')}}</div>
        <el-select class="ms-auto w-32" v-model="regionId">
          <el-option :label="t('allRegions')" :value="0" />
          <el-option
            v-for="region in session.regions" :key="region.id"
            :label="region.name" :value="region.id"
          />
        </el-select>
        <el-switch class="ms-2" :active-text="t('futureOnly')" :inactive-text="t('allRecords')" v-model="futureOnly" />
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
            <el-tag class="text-black" :color="typeColorMap[row.type]">
              {{t(row.type)}}
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
  <el-dialog v-model="isDialogOpen">
    <el-form>
      <el-form-item :label="t('time')">
        <el-date-picker v-model="addItemForm.time" />
      </el-form-item>
      <el-form-item :label="t('region')">
        <el-select v-model="addItemForm.regionId">
          <el-option
            v-for="region in session.regions"
            :key="region.id"
            :label="region.name"
            :value="region.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('type')">
        <el-select v-model="addItemForm.type">
          <el-option
            v-for="type in types"
            :key="type"
            :label="t(type)"
            :value="type"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('temperature')">
        <el-input-number v-model="addItemForm.maxTemperature" />
        <el-input-number class="ms-2" v-model="addItemForm.minTemperature" />
        <el-input-number class="ms-2" v-model="addItemForm.avgTemperature" />
      </el-form-item>
      <el-form-item :label="t('windSpeed')">
        <el-input-number v-model="addItemForm.windSpeed" />
      </el-form-item>
      <el-form-item :label="t('visibility')">
        <el-input-number v-model="addItemForm.visibility" />
      </el-form-item>
      <el-form-item :label="t('rainFall')">
        <el-input-number v-model="addItemForm.rainFall" />
      </el-form-item>
      <el-form-item>
        <el-button class="ms-auto" type="primary" round @click="addItem">
          {{t('confirm')}}
        </el-button>
        <el-button round @click="isDialogOpen=false">
          {{t('cancel')}}
        </el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>
