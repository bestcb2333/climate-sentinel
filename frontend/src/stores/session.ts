import { defineStore } from 'pinia'
import type { Region, User } from '@/tables'
import { ref } from 'vue'
import { request } from '@/axios'

export const useSessionStore = defineStore('session', () => {

  const user = ref<User | null>(null)
  async function loadUser() {
    user.value = await request.get<any, User>('/myinfo')
  }

  const regions = ref<Region[]>([])
  const total = ref(0)
  async function loadMap() {
    const res = await request.get<any, {
      data: Region[],
      total: number,
    }>('/regions?page=1&page_size=100')
    regions.value = res.data
    total.value = res.total
  }

  return { user, loadUser, regions, loadMap }
})

export default useSessionStore
