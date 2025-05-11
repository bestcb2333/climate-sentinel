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

use([
  RadarChart,
  TitleComponent,
  LegendComponent,
  CanvasRenderer,
])

const session = useSessionStore()

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const resources = ref<Resource[]>([])
watch([page, pageSize], async ([page, pageSize]) => {
  const res = await request.get<any, {
    data: Resource[],
    total: number,
  }>('/resources', {params: {
    page, pageSize,
  }})
  resources.value = res.data
  total.value = res.total
}, {immediate: true})

const isDialogOpen = ref(false)

const currentType = ref('towVehicle')

const resourceTypes = [
  ['vehicle', '救援车辆'],
  ['personnel', '救援人员'],
  ['comm', '通讯工具'],
  ['weather', '天气监测设备'],
  ['maintain', '道路维护设备'],
  ['other', '其他资源'],
].map(([value, label]) => ({value, label}))

const addResourceForm = reactive({
  type: '',
  name: '',
  quantity: 0,
  regionId: 1,
  coordinate: [0, 0],
  available: false,
})

const board = computed(() => ([
  ['资源总数', resources.value.length],
  ['可用资源数量', resources.value.filter(resource => resource.available).length],
  ['资源类型数', 6],
].map(([label, value]) => ({label, value}))))

async function addResource() {

}

interface RadarData {
  value: number[],
  name: string,
}

const radarData = ref<RadarData[]>([])

request.get<any, RadarData[]>('/radar').then(res => {
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
    indicator: resourceTypes.map(({label}) => ({
      name: label,
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
    addResource: '添加资源',
    resourceList: '资源列表',
    type: '资源类型',
    name: '资源名称',
    quantity: '资源数量',
    region: '所在区域',
    coordinate: '坐标',
    available: '是否可用',
  },
}})
</script>

<template>
  <div class="h-full grid grid-cols-[auto_auto_2fr_2fr] grid-rows-[2fr_3fr] gap-2">

    <div class="row-span-2 flex flex-col justify-between">
      <el-button type="primary" round @click="isDialogOpen=true">
        {{t('addResource')}}
      </el-button>
      <el-segmented :options="resourceTypes" direction="vertical" v-model="currentType" />
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

      <region-map :markers="resources" />
    </el-card>

    <el-card
      shadow="never"
      class="col-span-2 flex flex-col"
      body-class="grow overflow-y-auto"
      header-class="font-bold text-lg"
      footer-class="flex justify-end"
    >
      <template #header>
        {{t('resourceList')}}
      </template>

      <el-table :data="resources">

        <el-table-column
          :label="$t('resource.updatedAt')"
          prop="updatedAt"
          :formatter="formatDate"
        />

        <el-table-column :label="$t('resource.name')" prop="name" />
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

  <el-dialog v-model="isDialogOpen" :title="t('addResource')">
    <el-form :model="addResourceForm">
      <el-form-item :label="t('type')">
        <el-select v-model="addResourceForm.type">
          <el-option
            v-for="type in resourceTypes"
            :key="type.value"
            :label="type.label"
            :value="type.value"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('name')">
        <el-input v-model="addResourceForm.name" />
      </el-form-item>
      <el-form-item :label="t('quantity')">
        <el-input-number v-model="addResourceForm.quantity" />
      </el-form-item>
      <el-form-item :label="t('region')">
        <el-select v-model="addResourceForm.regionId">
          <el-option
            v-for="region in session.regions"
            :key="region.id"
            :label="region.name"
            :value="region.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('coordinate')">
        <el-input-number v-model="addResourceForm.coordinate[0]" />
        <el-input-number v-model="addResourceForm.coordinate[1]" class="ms-4" />
      </el-form-item>
      <el-form-item :label="t('available')">
        <el-switch v-model="addResourceForm.available" />
      </el-form-item>
      <el-form-item>
        <el-button round class="ms-auto" type="primary" @click="addResource">
          {{$t('confirm')}}
        </el-button>
        <el-button round @click="isDialogOpen=false">
          {{$t('cancel')}}
        </el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>
