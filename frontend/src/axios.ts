// @ts-check
import axios from 'axios'
import { ElMessage } from 'element-plus'
import usePersistedStore from '@/stores/persisted'
import qs from 'qs'

export const request = axios.create({
  paramsSerializer: params => qs.stringify(params, {arrayFormat: 'repeat'})
})

request.interceptors.request.use((config) => {
  const persisted = usePersistedStore()
  config.baseURL = persisted.setting.apiAddr
  if (persisted.token) {
    config.headers.set('Authorization', `Bearer ${persisted.token}`)
  }
  return config
})

request.interceptors.response.use(
  (res) => {
    if (res.data.message) {
      ElMessage({'type': 'success', 'message': res.data.message})
    }
    return res.data.data
  },
  (err) => {
    if (err.response) {
      const data = err.response.data
      if (err.response.data) {
        ElMessage({type: 'error', message: data.message})
        return Promise.reject(err)
      } else {
        return Promise.reject(err)
      }
    } else {
      return Promise.reject(err)
    }
  },
)
