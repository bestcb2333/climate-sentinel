<script setup lang="ts">
import {request} from '@/axios';
import useSessionStore from '@/stores/session';
import type {Resource} from '@/tables';
import {formatDate} from '@/utils';
import {reactive, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import RegionMap from '@/components/RegionMap.vue';
import {computed} from 'vue';
import {use} from 'echarts/core';
import {RadarChart} from 'echarts/charts';
import {LegendComponent, TitleComponent} from 'echarts/components';
import {CanvasRenderer} from 'echarts/renderers';
import type {EChartsOption} from 'echarts';
import VChart from 'vue-echarts'
import type {StatsItem} from '@/types';
import {useRouteQuery} from '@vueuse/router';

use([
  RadarChart,
  TitleComponent,
  LegendComponent,
  CanvasRenderer,
])

const session = useSessionStore()

const page = useRouteQuery('page', '1', {transform: Number})
const pageSize = useRouteQuery('page_size', '10', {transform: Number})
const total = ref(0)
const currentType = useRouteQuery('type', 'all')
const regionId = useRouteQuery('region_id', '0', {transform: Number})
const availableOnly = useRouteQuery('available', 'false', {transform: Boolean})
const resources = ref<Resource[]>([])
async function loadTable() {
  try {
    const res = await request.get<any, {
      data: Resource[],
      total: number,
    }>(`/resources?page=${page.value}&page_size=${pageSize.value}&type=${currentType.value}&region_id=${regionId.value}&available=${availableOnly.value}`)
    resources.value = res.data
    total.value = res.total
  } catch {}
}
watch([page, pageSize, currentType, regionId, availableOnly], loadTable, {immediate: true})

const isDialogOpen = ref(false)

const resourceTypes = ['vehicle', 'personnel', 'comm', 'weather', 'maintain', 'other']

const addForm = reactive({
  type: '',
  name: '',
  quantity: 0,
  regionId: 1,
  coordinate: [0, 0],
  available: false,
})

const board = ref<StatsItem[]>([])

request.get<any, StatsItem[]>('/resources/stats').then(res => {
  board.value = res
}).catch(() => {})

async function addItem() {
  try {
    await request.post<any, null>('/resources', addForm)
    isDialogOpen.value = false
    loadTable()
  } catch {}
}

interface RadarData {
  value: number[],
  name: string,
}

const radarData = ref<RadarData[]>([])

request.get<any, RadarData[]>('/resources/radar').then(res => {
  radarData.value = res
}).catch(() => {})

const radarOption = computed((): EChartsOption => ({
  title: {
    text: '资源雷达图'
  },
  legend: {
    right: 0,
    top: 'center',
    orient: 'vertical',
    data: session.regions.map(({name}) => name),
  },
  radar: {
    indicator: resourceTypes.map(type => ({
      name: t(type),
      max: 60,
    })),
  },
  series: [{
    name: '资源雷达图',
    type: 'radar',
    data: radarData.value,
  }],
}))

const {t} = useI18n({messages: {
  zh: {
    add: '添加资源',
    resourceList: '资源列表',
    type: '资源类型',
    name: '资源名称',
    quantity: '资源数量',
    region: '所在区域',
    coordinate: '坐标',
    available: '是否可用',
    all: '所有资源',
    allRegion: '所有区域',
    vehicle: '救援车辆',
    personnel: '救援人员',
    comm: '通讯工具',
    weather: '天气监测设备',
    maintain: '道路维护设备',
    other: '其他资源',
    availableOnly: '仅显示可用',
    showAll: '显示所有',
  },
}})

const markers = computed(() => resources.value.map(({name,coordinate,available})=>({
  name,coordinate,color:available?'green':'red',
})))

const selected = ref<number[]>([])
async function deleteItems() {
  try {
    await request.delete('/resources', {params: {id: selected.value}})
    isDialogOpen.value = false
    loadTable()
  } catch {}
}
</script>

<template>
  <div class="h-full grid grid-cols-[auto_auto_2fr_2fr] grid-rows-[2fr_3fr] gap-2">

    <div class="row-span-2 flex flex-col justify-between">
      <el-button type="primary" round @click="isDialogOpen=true">
        {{t('add')}}
      </el-button>
      <el-segmented
        :options="['all', ...resourceTypes].map(type=>({label:t(type),value:type}))"
        direction="vertical" v-model="currentType"
      />
    </div>

    <el-card shadow="hover" body-class="h-full flex flex-col justify-around items-center">
      <div v-for="item in board" :key="item.label" class="flex flex-col items-center">
        <div class="font-bold">
          {{item.label}}
        </div>
        <div>
          {{item.value}}
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" body-class="h-full">
      <v-chart :option="radarOption" />
    </el-card>

    <el-card shadow="hover" class="row-span-2 flex flex-col" body-class="grow">
      <template #header>
        资源分布图
      </template>

      <region-map
        :markers="markers" />
    </el-card>

    <el-card
      shadow="never"
      class="col-span-2 flex flex-col"
      body-class="grow overflow-y-auto"
      header-class="flex"
      footer-class="flex justify-end"
    >
      <template #header>
        <div>
          {{t('resourceList')}}
        </div>
        <el-switch v-model="availableOnly" class="ms-auto"
          :active-text="t('availableOnly')" :inactive-text="t('showAll')"
        />
        <el-select v-model="regionId" class="w-32 ms-2">
          <el-option :label="t('allRegion')" :value="0" />
          <el-option
            v-for="region in session.regions"
            :key="region.id"
            :label="region.name"
            :value="region.id"
          />
        </el-select>
        <el-button type="danger" @click="deleteItems" class="ms-2">
          删除所选项
        </el-button>
      </template>

      <el-table :data="resources" @selection-change="(val:Resource[])=>selected=val.map(item=>item.id)">

        <el-table-column type="selection" />
        <el-table-column
          :label="$t('resource.updatedAt')"
          prop="updatedAt"
          :formatter="formatDate"
        />
        <el-table-column :label="$t('resource.name')" prop="name" />
        <el-table-column :label="t('type')" prop="type">
          <template #default="{row}">
            {{t(row.type)}}
          </template>
        </el-table-column>
        <el-table-column :label="$t('resource.quantity')" prop="quantity" />
        <el-table-column :label="$t('resource.region')" prop="region.name" />
        <el-table-column :label="$t('resource.coordinate')" prop="coordinate" />
        <el-table-column :label="$t('resource.available')" prop="available">
          <template #default="scope">
            <el-tag :type="scope.row.available?'success':'danger'">
              {{$t(scope.row.available?'yes':'no')}}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <el-pagination layout="sizes, prev, pager, next, total"
          :total="total" v-model:current-page="page" v-model:page-size="pageSize"
        />
      </template>
    </el-card>

  </div>

  <el-dialog v-model="isDialogOpen" :title="t('add')">
    <el-form :model="addForm">
      <el-form-item :label="t('type')">
        <el-select v-model="addForm.type">
          <el-option
            v-for="type in resourceTypes"
            :key="type"
            :label="t(type)"
            :value="type"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('name')">
        <el-input v-model="addForm.name" />
      </el-form-item>
      <el-form-item :label="t('quantity')">
        <el-input-number v-model="addForm.quantity" />
      </el-form-item>
      <el-form-item :label="t('region')">
        <el-select v-model="addForm.regionId">
          <el-option
            v-for="region in session.regions"
            :key="region.id"
            :label="region.name"
            :value="region.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('coordinate')">
        <el-input-number v-model="addForm.coordinate[0]" />
        <el-input-number v-model="addForm.coordinate[1]" class="ms-4" />
      </el-form-item>
      <el-form-item :label="t('available')">
        <el-switch v-model="addForm.available" />
      </el-form-item>
      <el-form-item>
        <el-button round class="ms-auto" type="primary" @click="addItem">
          {{$t('confirm')}}
        </el-button>
        <el-button round @click="isDialogOpen=false">
          {{$t('cancel')}}
        </el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>
