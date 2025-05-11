<script setup lang="ts">
import RegionMap from '@/components/RegionMap.vue';
import {useI18n} from 'vue-i18n';
import type {Route} from '@/tables';
import {ref} from 'vue';
import {watch} from 'vue';
import {request} from '@/axios';
import {computed} from 'vue';

const {t} = useI18n({messages: {
  zh: {
    type: '路线类型',
    name: '路线名称',
    available: '是否可用',
    rate: '限速',
    region: '所在区域',
    description: '描述',
    pleaseSelect: '请选中一项路线查看信息',
    routesMap: '路线地图与描述',
    routesCount: '总路线数量',
    availableRoutesCount: '可用路线数量',
  },
}})

const currentRow = ref<Route|undefined>(undefined)
const routes = ref<Route[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
watch(([page, pageSize]), async ([page, pageSize]) => {
  try {
    const res = await request.get<any, {
      data: Route[],
      total: number,
    }>('/routes', {params: {
      page, pageSize,
    }})
    routes.value = res.data
    total.value = res.total
  } catch {}
}, {immediate: true})

const board = computed(() => [
  ['总路线数量', routes.value.length],
  ['高速路线数量', routes.value.filter(route => route.rate>50).length],
  ['可用路线数量', routes.value.filter(route => route.available).length],
].map(([label, value]) => ({label, value})))
</script>

<template>
  <div class="h-full grid grid-rows-[auto_1fr] grid-cols-[auto_3fr_2fr] gap-2">

    <div class="row-span-2">

    </div>

    <el-card shadow="hover" body-class="h-full flex justify-evenly items-center">
      <div v-for="item in board" :key="item.label" class="flex flex-col items-center">
        <div class="font-bold">
          {{item.label}}
        </div>
        <div>
          {{item.value}}
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="row-span-2 flex flex-col"
      body-class="grow"
    >

      <template #header>
        {{t('routesMap')}}
      </template>

      <region-map :routes="routes" />

      <template #footer>
        {{currentRow?currentRow.description:t('pleaseSelect')}}
      </template>

    </el-card>

    <el-card shadow="hover" class="flex flex-col"
      body-class="grow"
    >

      <template #header>
        救援路线列表
      </template>

      <el-table :data="routes" highlight-current-row @current-change="val=>currentRow=val">
        <el-table-column :label="t('type')" prop="type" />
        <el-table-column :label="t('name')" prop="name" />
        <el-table-column :label="t('available')" prop="available">
          <template #default="{row}">
            <el-tag :type="row.available?'success':'danger'">
              {{row.available?$t('yes'):$t('no')}}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('rate')" prop="rate">
          <template #default="{row}">
            {{row.rate}} km/h
          </template>
        </el-table-column>
        <el-table-column :label="t('region')" prop="region.name" />
        <el-table-column :label="t('description')" prop="description" show-overflow-tooltip />
      </el-table>

      <template #footer>
        <el-pagination layout="sizes, prev, pager, next, total" :total="total"
          v-model:current-page="page" v-model:page-size="pageSize"
        />
      </template>

    </el-card>

  </div>
</template>
