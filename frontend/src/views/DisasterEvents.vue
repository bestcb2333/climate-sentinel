<script setup lang="ts">
import {request} from '@/axios';
import {Plus, Minus} from '@element-plus/icons-vue';
import {computed, reactive, ref, watch} from 'vue';
import { formatDate } from '@/utils';
import {useI18n} from 'vue-i18n';
import RegionMap from '@/components/RegionMap.vue';
import type {Event} from '@/tables';
import {map} from 'echarts/types/src/export/api/util.js';
import dayjs from 'dayjs';

const isOnlyCurrent = ref(false)
const disasterType = ref('all')
const disasterTypes = [
  ['all', '所有灾害'],
  ['blizzard', '暴雪'],
  ['typhoon', '台风'],
  ['hail', '冰雹'],
  ['fog', '大雾'],
  ['thunder', '雷雨'],
  ['others', '其他'],
].map(([value, label]) => ({value, label}))

const currentRow = ref<Event|null>(null)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const events = ref<Event[]>([])
watch([page, pageSize], async ([page, pageSize]) => {
  const res = await request.get<any, {
    data: Event[],
    total: number,
  }>('/events', {params: {
    page, pageSize,
  }})
  events.value = res.data
  total.value = res.total
}, {immediate: true})

const severities: {
  type: 'success'|'primary'|'warning'|'danger',
  content: string,
}[] = [
  {type: 'success', content: '无风险'},
  {type: 'primary', content: '低风险'},
  {type: 'warning', content: '中风险'},
  {type: 'danger', content: '高风险'},
]

const isDialogOpen = ref(false)

const addEventForm = reactive({
  type: '',
  startTime: new Date(),
  endTime: new Date(),
  severity: 0,
  regionId: 1,
  coordinate: [0, 0],
  description: '',
})

const board = computed(() => [
  ['事件总数', events.value.length],
  ['高风险事件数', events.value.filter(event => event.severity===3).length],
  ['中风险事件数', events.value.filter(event => event.severity===2).length],
  ['低风险事件数', events.value.filter(event => event.severity===1).length],
  ['无风险事件数', events.value.filter(event => event.severity===0).length],
].map(([label, value]) => ({label, value})))

async function addEvent() {
}

const { t } = useI18n({messages: {
  zh: {
    reportEvent: '报告事件',
    onlyCurrent: '仅当前事件',
    disasterType: '事件类型',
    startTime: '开始时间',
    endTime: '结束时间',
    severity: '严重性',
    coordinate: '坐标',
    description: '描述',
    notEnd: '未结束',
    region: '区域名',
  },
}})
</script>

<template>
  <div class="h-full grid grid-cols-[auto_3fr_2fr] grid-rows-[auto_1fr] gap-2">

    <div class="row-span-2 flex flex-col items-center gap-2">
      <el-button round :type="isOnlyCurrent?'success':''" @click="isOnlyCurrent=!isOnlyCurrent">
        {{t('onlyCurrent')}}
      </el-button>
      <el-button round type="primary" @click="isDialogOpen=true">
        {{t('reportEvent')}}
      </el-button>
      <el-segmented class="mt-auto" direction="vertical" v-model="disasterType" :options="disasterTypes">
        <template #default="scope">
          {{(scope.item as Record<string, string>).label}}
        </template>
      </el-segmented>
    </div>

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

    <el-card class="row-span-2 flex flex-col" shadow="hover"
      header-class="flex justify-between"
      body-class="grow"
    >
      <template #header>
        <div class="font-bold">
          {{('event.mapTitle')}}
        </div>
        <div>
          <el-button circle :icon="Plus" />
          <el-button circle :icon="Minus" />
        </div>
      </template>

      <region-map :markers="events" />

      <template #footer>
        {{currentRow ? currentRow.description : '请选择事件查看信息'}}
      </template>

    </el-card>

      <el-card
        shadow="never"
        class="flex flex-col"
        body-class="grow overflow-y-auto"
        header-class="font-bold text-lg"
        footer-class="flex justify-end"
      >
        <template #header>
          {{$t('event.eventList')}}
        </template>

        <el-table :data="events" highlight-current-row @current-change="val => currentRow = val">
          <el-table-column
            :label="$t('event.startTime')"
            prop="startTime"
            :formatter="formatDate"
          />
          <el-table-column :label="$t('event.endTime')">
            <template #default="{row}">
              <template v-if="row.endTime">
                {{dayjs(row.endTime).format('DD月MM日 HH:mm')}}
              </template>
              <div v-else class="text-red-500">
                {{t('notEnd')}}
              </div>
            </template>
          </el-table-column>
          <el-table-column :label="$t('event.severity')" prop="severity">
            <template #default="scope">
              <el-tag :type="severities[scope.row.severity].type">
                {{severities[scope.row.severity].content}}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('event.user')" prop="user.name" />
          <el-table-column :label="t('region')" prop="region.name" />
          <el-table-column :label="$t('event.description')" prop="description" show-overflow-tooltip />
        </el-table>

        <template #footer>
          <el-pagination layout="sizes, prev, pager, next, total"
            :total="total"
            v-model:current-page="page"
            v-model:page-size="pageSize"
            :page-sizes="[5, 10, 20, 30]"
          />
        </template>
      </el-card>

  </div>

  <el-dialog v-model="isDialogOpen" :title="$t('event.report')">
    <el-form>
      <el-form-item :label="t('disasterType')">
        <el-select v-model="addEventForm.type">
          <el-option
            v-for="type in disasterTypes"
            :key="type.value"
            :label="type.label"
            :value="type.value"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('startTime')">
        <el-date-picker v-model="addEventForm.startTime" type="datetime" />
      </el-form-item>
      <el-form-item :label="t('endTime')">
        <el-date-picker v-model="addEventForm.endTime" type="datetime" />
      </el-form-item>
      <el-form-item :label="t('severity')">
        <el-select v-model="addEventForm.severity">
          <el-option
            v-for="severity, index in severities"
            :key="severity.type"
            :label="severity.content"
            :value="index"
          />
        </el-select>
      </el-form-item>
      <el-form-item>

      </el-form-item>
      <el-form-item :label="t('coordinate')">
        <el-input-number v-model="addEventForm.coordinate[0]" />
        <el-input-number v-model="addEventForm.coordinate[1]" class="ms-4" />
      </el-form-item>
      <el-form-item :label="t('description')">
        <el-input type="textarea" v-model="addEventForm.description" />
      </el-form-item>
      <el-form-item>
        <el-button @click="addEvent" round class="ms-auto" type="primary">
          {{$t('global.confirm')}}
        </el-button>
        <el-button @click="isDialogOpen=false" round>
          {{$t('global.cancel')}}
        </el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>
