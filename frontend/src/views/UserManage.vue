<script setup lang="ts">
import {request} from '@/axios';
import type {User} from '@/tables';
import type {StatsItem} from '@/types';
import {formatDate} from '@/utils';
import {ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const regionId = ref(0)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const users = ref<User[]>([])
const statsItems = ref<StatsItem[]>([])

watch([page, pageSize, regionId], loadTable, {immediate: true})
loadStats()

async function loadTable() {
  try {
    const res = await request.get<any, {
      total: number
      data: User[]
    }>(`/users?page=${page.value}&page_size=${pageSize.value}&region_id=${regionId.value}`)
    total.value = res.total
    users.value = res.data
  } catch {}
}

async function loadStats() {
  try {
    statsItems.value = await request.get<any, StatsItem[]>('/users/stats')
  } catch {}
}

const {t} = useI18n({messages: {
  zh: {
    createdAt: '创建时间',
    name: '用户名',
    email: '邮箱',
    admin: '是否管理员',
    region: '志愿服务的区域',
    total: '总用户数',
    volunteers: '志愿者数量',
    admins: '管理员数量',
    userList: '用户列表',
    loginSignup: '登录/注册',
  },
}})
</script>

<template>
  <div class="h-full grid grid-cols-2 grid-rows-[auto_1fr] gap-2">

    <el-card shadow="hover" body-class="flex gap-2">
      <el-card v-for="item in statsItems" :key="item.label" shadow="hover">
        <div>{{t(item.label)}}</div>
        <div>{{item.value}}</div>
      </el-card>
    </el-card>

    <el-card class="row-span-2" shadow="hover">
    </el-card>

    <el-card shadow="hover" class="flex flex-col col-span-2"
      body-class="grow min-h-0 overflow-y-auto"
    >

      <template #header>
        {{t('userList')}}
      </template>

      <el-table :data="users">
        <el-table-column :label="t('createdAt')" prop="createdAt" :formatter="formatDate" />
        <el-table-column :label="t('name')" prop="name" />
        <el-table-column :label="t('email')" prop="email" width="200" />
        <el-table-column :label="t('admin')">
          <template #default="{row}">
            <el-tag :type="row.admin?'success':'warning'">
              {{t(row.admin?'yes':'no')}}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('region')" prop="region.name" />
      </el-table>

      <template #footer>
        <el-pagination layout="sizes,prev,pager,next,total" :total="total"
          v-model:current-page="page" v-model:page-size="pageSize"
        />
      </template>

    </el-card>

  </div>
</template>
