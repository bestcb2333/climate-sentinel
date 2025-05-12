<script setup lang="ts">
import {request} from '@/axios';
import {Plus, Minus} from '@element-plus/icons-vue';
import {reactive, ref, watch} from 'vue';
import { formatDate } from '@/utils';
import {useI18n} from 'vue-i18n';
import RegionMap from '@/components/RegionMap.vue';
import type {Event} from '@/tables';
import dayjs from 'dayjs';
import useSessionStore from '@/stores/session';
import type {StatsItem} from '@/types';
import {computed} from 'vue';
import {useRouteQuery} from '@vueuse/router';

const session = useSessionStore()
const isOnlyCurrent = ref(false)

const disasterTypes = [
  'blizzard', 'typhoon', 'hail', 'fog', 'thunder', 'others',
]

const severities = [
  'safe', 'low', 'medium', 'high',
]

const severityColorMap = {
  safe: 'green',
  low: 'cyan',
  medium: 'orange',
  high: 'red',
}

const currentRow = ref<Event|null>(null)
const page = useRouteQuery('page', '1', {transform: Number})
const pageSize = useRouteQuery('pag_size', '10', {transform: Number})
const type = useRouteQuery('type', 'all')
const regionId = useRouteQuery('region_id', '0', {transform: Number})
const total = ref(0)
const events = ref<Event[]>([])
async function loadTable() {
  const res = await request.get<any, {
    data: Event[],
    total: number,
  }>(`/events?page=${page.value}&page_size=${pageSize.value}&region_id=${regionId.value}&type=${type.value}`)
  events.value = res.data
  total.value = res.total
}
watch([page, pageSize, type, regionId], loadTable, {immediate: true})

const isDialogOpen = ref(false)

const addItemForm = reactive({
  type: '',
  startTime: new Date(),
  endTime: new Date(),
  severity: 0,
  regionId: 1,
  coordinate: [0, 0],
  description: '',
})

const markers = computed(() => events.value.map(({
  name, coordinate, severity,
}) => ({
  name, coordinate, color: severityColorMap[severity],
})))

const statsItems = ref<StatsItem[]>([])
async function getStatsItems() {
  try {
    const res = await request.get<any, StatsItem[]>('/events/stats')
    statsItems.value = res
  } catch {}
}
getStatsItems()

async function addEvent() {
  try {
    await request.post('/events', addItemForm)
    isDialogOpen.value = false
    loadTable()
    getStatsItems()
  } catch {}
}

const selected = ref<number[]>([])
async function deleteItems() {
  try {
    await request.delete('/events', {params: {id: selected.value}})
    selected.value = []
    loadTable()
    getStatsItems()
  } catch {}
}

const { t } = useI18n({messages: {
  zh: {
    reportEvent: '报告事件',
    onlyCurrent: '仅当前事件',
    disasterType: '事件类型',
    startTime: '开始时间',
    endTime: '结束时间',
    type: '灾害类型',
    severity: '严重性',
    coordinate: '坐标',
    description: '描述',
    notEnd: '未结束',
    region: '区域名',
    eventList: '事件列表',
    all: '所有灾害',
    allRegion: '所有区域',
    blizzard: '暴雪',
    typhoon: '台风',
    hail: '冰雹',
    fog: '大雾',
    thunder: '雷雨',
    others: '其他',
    safe: '无风险',
    low: '低风险',
    medium: '中风险',
    high: '高风险',
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
      <el-segmented class="mt-auto" direction="vertical" v-model="type"
        :options="['all', ...disasterTypes].map(type=>({label:t(type),value:type}))"
      >
        <template #default="scope">
          {{(scope.item as Record<string, string>).label}}
        </template>
      </el-segmented>
    </div>

    <el-card shadow="hover" body-class="h-full flex justify-around items-center">
      <div v-for="item in statsItems" :key="item.label" class="flex flex-col items-center">
        <div class="font-bold">
          {{t(item.label)}}
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

      <region-map :markers="markers" />

      <template #footer>
        {{currentRow ? currentRow.description : '请选择事件查看信息'}}
      </template>

    </el-card>

      <el-card
        shadow="never"
        class="flex flex-col"
        body-class="grow overflow-y-auto"
        header-class="flex"
        footer-class="flex justify-end"
      >
        <template #header>
          <div>
            {{t('eventList')}}
          </div>
          <el-select class="w-32 ms-auto" v-model="regionId">
            <el-option :label="t('allRegion')" :value="0" />
            <el-option
              v-for="region in session.regions"
              :key="region.id"
              :label="region.name"
              :value="region.id"
            />
          </el-select>
          <el-button class="ms-2" type="danger" round @click="deleteItems">
            删除所选事件
          </el-button>
        </template>

        <el-table :data="events"
          highlight-current-row @current-change="val=>currentRow=val"
          @selection-change="(val:Event[])=>selected=val.map(item=>item.id)"
        >
          <el-table-column type="selection" />
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
          <el-table-column :label="t('type')" prop="type" :formatter="(_1,_2,val)=>t(val)" />
          <el-table-column :label="$t('event.severity')" prop="severity">
            <template #default="{row}">
              <el-tag :color="severityColorMap[row.severity as 'low']" class="text-black">
                {{t(row.severity)}}
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
        <el-select v-model="addItemForm.type">
          <el-option
            v-for="type in disasterTypes"
            :key="type"
            :label="t(type)"
            :value="type"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('startTime')">
        <el-date-picker v-model="addItemForm.startTime" type="datetime" />
      </el-form-item>
      <el-form-item :label="t('endTime')">
        <el-date-picker v-model="addItemForm.endTime" type="datetime" />
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
      <el-form-item :label="t('severity')">
        <el-select v-model="addItemForm.severity">
          <el-option
            v-for="severity, index in severities"
            :key="severity"
            :label="t(severity)"
            :value="index"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('coordinate')">
        <el-input-number v-model="addItemForm.coordinate[0]" />
        <el-input-number v-model="addItemForm.coordinate[1]" class="ms-4" />
      </el-form-item>
      <el-form-item :label="t('description')">
        <el-input type="textarea" v-model="addItemForm.description" />
      </el-form-item>
      <el-form-item>
        <el-button @click="addEvent" round class="ms-auto" type="primary">
          {{t('confirm')}}
        </el-button>
        <el-button @click="isDialogOpen=false" round>
          {{t('cancel')}}
        </el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>
