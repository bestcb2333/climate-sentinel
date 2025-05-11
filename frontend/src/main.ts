import 'element-plus/dist/index.css'
import '@/styles.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

import App from './App.vue'
import router from './router'
import {createI18n} from 'vue-i18n'

import zh from './locale/zh'

const app = createApp(App)

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)
app.use(pinia)

const i18n = createI18n({
  locale: 'zh',
  fallbackLocale: 'zh',
  messages: {zh},
})
app.use(i18n)

app.use(ElementPlus)
app.use(router)

app.mount('#app')
