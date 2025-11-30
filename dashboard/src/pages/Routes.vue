<template>
  <el-card>
    <div style="margin-bottom: 12px;">
      <el-button type="primary" @click="load">刷新</el-button>
      <el-button @click="save">保存</el-button>
    </div>
    <el-form :model="form" label-width="140px">
      <el-form-item label="前缀">
        <el-input v-model="form.prefix" />
      </el-form-item>

      <el-form-item label="灰度 v1 百分比">
        <el-slider v-model="form.plugins.trafficSplit.v1Percent" :min="0" :max="100" />
        <el-switch v-model="form.plugins.trafficSplit.enabled" style="margin-left:8px" />
      </el-form-item>

      <el-form-item label="负载均衡权重启用">
        <el-switch v-model="form.lbEnabled" />
      </el-form-item>

      <el-form-item label="v1 权重 A/B/C">
        <div class="weights">
          <el-input-number v-model="v1[0].weight" :min="0" />
          <el-input-number v-model="v1[1].weight" :min="0" />
          <el-input-number v-model="v1[2].weight" :min="0" />
        </div>
      </el-form-item>

      <el-form-item label="v2 权重 A/B/C">
        <div class="weights">
          <el-input-number v-model="v2[0].weight" :min="0" />
          <el-input-number v-model="v2[1].weight" :min="0" />
          <el-input-number v-model="v2[2].weight" :min="0" />
        </div>
      </el-form-item>

      <el-form-item label="API Key">
        <el-switch v-model="form.plugins.auth.enabled" />
        <el-input v-model="form.plugins.auth.key" placeholder="X-API-Key" style="margin-left:8px" />
      </el-form-item>

      <el-form-item label="IP 白名单">
        <el-switch v-model="form.plugins.ipWhitelist.enabled" />
        <el-input v-model="ips" placeholder="逗号分隔" style="margin-left:8px" />
      </el-form-item>

      <el-form-item label="限流">
        <el-switch v-model="form.plugins.rateLimit.enabled" />
        <el-input-number v-model="form.plugins.rateLimit.windowSec" :min="1" />
        <el-input-number v-model="form.plugins.rateLimit.max" :min="1" />
      </el-form-item>

      <el-form-item label="CORS 允许全部">
        <el-switch v-model="form.plugins.cors.enabled" />
        <el-switch v-model="form.plugins.cors.allowAll" style="margin-left:8px" />
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { api } from '../api/admin'

type Upstream = { name: string; url: string; weight: number }

const form = reactive<any>({
  id: 'users',
  prefix: '/api/users',
  v1: [],
  v2: [],
  lbEnabled: false,
  plugins: {
    auth: { enabled: false, key: '' },
    ipWhitelist: { enabled: false, ips: [] },
    rateLimit: { enabled: false, windowSec: 1, max: 5 },
    cors: { enabled: false, allowAll: true, origins: [] },
    trafficSplit: { enabled: false, v1Percent: 100 }
  }
})

const v1 = ref<Upstream[]>([])
const v2 = ref<Upstream[]>([])
const ips = ref('')

function merge() {
  form.v1 = v1.value
  form.v2 = v2.value
  form.plugins.ipWhitelist.ips = ips.value.split(',').map(s => s.trim()).filter(Boolean)
}

async function load() {
  const { data } = await api.getRoutes()
  const rt = (data as any[]).find(r => r.id === 'users') || (data as any[])[0]
  if (rt) {
    Object.assign(form, rt)
    v1.value = rt.v1 || []
    v2.value = rt.v2 || []
    ips.value = (rt.plugins?.ipWhitelist?.ips || []).join(',')
  } else {
    v1.value = [
      { name: 'A', url: 'http://localhost:9001', weight: 1 },
      { name: 'B', url: 'http://localhost:9002', weight: 1 },
      { name: 'C', url: 'http://localhost:9003', weight: 1 }
    ]
    v2.value = [
      { name: 'A', url: 'http://localhost:9011', weight: 1 },
      { name: 'B', url: 'http://localhost:9012', weight: 1 },
      { name: 'C', url: 'http://localhost:9013', weight: 1 }
    ]
  }
}

async function save() {
  merge()
  await api.updateRoute(form)
}

load()
</script>

<style scoped>
.weights { display: flex; gap: 8px }
</style>

