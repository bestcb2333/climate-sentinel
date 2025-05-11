// @ts-check
import { defineStore } from 'pinia'
import {ref} from 'vue'

const usePersistedStore = defineStore('persisted', () => {

  const apiAddr = ref('http://axtl.cn:8700')
  const darkMode = ref(false)

  const token = ref(null)

  return {apiAddr, darkMode, token}
}, {persist: true})

export default usePersistedStore
