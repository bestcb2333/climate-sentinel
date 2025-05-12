<script setup lang="ts">
import {request} from '@/axios';
import usePersistedStore from '@/stores/persisted';
import useSessionStore from '@/stores/session';
import {reactive, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const session = useSessionStore()
const persisted = usePersistedStore()

const isDialogOpen = ref(false)
const currentTab = ref('login')

const loginForm = reactive({
  username: '',
  password: '',
})

async function login() {
  try {
    persisted.token = await request.post<any, string>('/login', loginForm)
    await session.loadUser()
    isDialogOpen.value = false
  } catch {}
}

const signupForm = reactive({
  email: '',
  authcode: '',
  username: '',
  password: '',
})

async function signup() {
  try {
    persisted.token = await request.post<any, string>('/signup', signupForm)
    await session.loadUser()
    isDialogOpen.value = false
  } catch {}
}

const retrieveForm = reactive({
  email: '',
  authcode: '',
  password: '',
})

async function retrieve() {
  try {
    persisted.token = await request.post<any, string>('/retrieve', retrieveForm)
    await session.loadUser()
    isDialogOpen.value = false
  } catch {}
}

async function logout() {
  persisted.token = null
  session.user = null
}


const {t} = useI18n({messages: {
  zh: {
    logo: 'Climate Sentinel',
    dashboard: '仪表盘',
    resource: '救援资源管理',
    weather: '天气状况',
    routes: '救援路线管理',
    events: '灾害事件',
    staff: '用户管理',
    loginSignup: '登录/注册',
    logout: '退出登录',
  },
}})
</script>

<template>
  <el-menu
    mode="horizontal" router :default-active="$route.path"
    class="h-10 bg-transparent !border-0 flex gap-2" :ellipsis="false"
  >
    <el-menu-item class="font-bold text-lg">
      {{t('logo')}}
    </el-menu-item>
    <template v-for="route in $router.getRoutes()" :key="route.name">
      <el-menu-item v-if="route.meta.label" :index="route.path">
        {{t(route.meta.label as string)}}
      </el-menu-item>
    </template>
    <el-menu-item class="!ms-auto">
      <el-button type="primary" @click="logout" v-if="session.user">
        {{t('logout')}}
      </el-button>
      <el-button type="primary" @click="isDialogOpen=true" v-else>
        {{t('loginSignup')}}
      </el-button>
    </el-menu-item>
  </el-menu>
  <el-dialog v-model="isDialogOpen" title="欢迎使用Climate Sentinel">
    <el-tabs v-model="currentTab">
      <el-tab-pane label="登录" name="login">
        <el-form :model="loginForm" label-width="auto">
          <el-form-item label="用户名">
            <el-input v-model="loginForm.username" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="loginForm.password" />
          </el-form-item>
          <el-form-item>
            <el-button class="ms-2" type="primary" @click="login">
              登录
            </el-button>
            <el-button @click="isDialogOpen=false">
              返回
            </el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      <el-tab-pane label="注册" name="signup">
        <el-form :model="signupForm" label-width="auto">
          <el-form-item label="邮箱">
            <el-input v-model="signupForm.email" />
          </el-form-item>
          <el-form-item label="邮箱验证码">
            <el-input v-model="signupForm.authcode" />
          </el-form-item>
          <el-form-item label="用户名">
            <el-input v-model="signupForm.username" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="signupForm.password" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" class="ms-auto" @click="signup">
              注册
            </el-button>
            <el-button @click="isDialogOpen=false">
              返回
            </el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      <el-tab-pane label="找回密码" name="retrieve">
        <el-form :model="retrieveForm" label-width="auto">
          <el-form-item label="邮箱">
            <el-input v-model="retrieveForm.email" />
          </el-form-item>
          <el-form-item label="验证码">
            <el-input v-model="retrieveForm.authcode" />
          </el-form-item>
          <el-form-item label="新密码">
            <el-input v-model="retrieveForm.password" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="retrieve" class="ms-auto">
              找回密码
            </el-button>
            <el-button @click="isDialogOpen=false">
              返回
            </el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </el-dialog>
</template>
