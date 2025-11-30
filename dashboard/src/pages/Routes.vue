<template>
  <el-card>
    <div style="margin-bottom: 12px;">
      <el-button type="primary" @click="load" :loading="loading">刷新</el-button>
      <el-button @click="save" type="success">保存</el-button>
    </div>
    <el-alert
      v-if="clicked && !connected"
      title="网关未连接或接口不可用"
      type="warning"
      description="已加载本地或默认配置。请确保网关运行在 http://localhost:8080 后点击刷新。"
      show-icon
      style="margin-bottom:12px"/>
    <el-alert
      v-if="connected"
      title="网关已连接"
      type="success"
      description="可以编辑并保存配置到网关。"
      show-icon
      style="margin-bottom:12px"/>
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
        <span style="margin-left:8px;color:#909399">开启后按权重加权随机；关闭为轮询</span>
      </el-form-item>

      <el-form-item label="v1 权重 A/B/C">
        <div class="weights">
          <div class="witem"><span class="wlabel">A</span><el-input-number v-model="v1[0].weight" :min="0" :max="100" controls-position="right" /></div>
          <div class="witem"><span class="wlabel">B</span><el-input-number v-model="v1[1].weight" :min="0" :max="100" controls-position="right" /></div>
          <div class="witem"><span class="wlabel">C</span><el-input-number v-model="v1[2].weight" :min="0" :max="100" controls-position="right" /></div>
        </div>
      </el-form-item>

      <el-form-item label="v2 权重 A/B/C">
        <div class="weights">
          <div class="witem"><span class="wlabel">A</span><el-input-number v-model="v2[0].weight" :min="0" :max="100" controls-position="right" /></div>
          <div class="witem"><span class="wlabel">B</span><el-input-number v-model="v2[1].weight" :min="0" :max="100" controls-position="right" /></div>
          <div class="witem"><span class="wlabel">C</span><el-input-number v-model="v2[2].weight" :min="0" :max="100" controls-position="right" /></div>
        </div>
      </el-form-item>

      <el-form-item label="API Key">
        <el-switch v-model="form.plugins.auth.enabled" />
        <el-input v-model="form.plugins.auth.key" placeholder="在请求头中使用 X-API-Key" style="margin-left:8px" />
      </el-form-item>

      <el-form-item label="IP 白名单">
        <el-switch v-model="form.plugins.ipWhitelist.enabled" />
        <el-input v-model="ips" placeholder="逗号分隔，如 127.0.0.1, ::1" style="margin-left:8px" />
      </el-form-item>

      <el-form-item label="限流">
        <el-switch v-model="form.plugins.rateLimit.enabled" />
        <el-input-number v-model="form.plugins.rateLimit.windowSec" :min="1" :max="600" controls-position="right" />
        <span style="margin:0 8px;color:#909399">秒内最多</span>
        <el-input-number v-model="form.plugins.rateLimit.max" :min="1" :max="100000" controls-position="right" />
        <span style="margin-left:8px;color:#909399">次请求</span>
      </el-form-item>

      <el-form-item label="CORS 设置">
        <el-switch v-model="form.plugins.cors.enabled" />
        <span style="margin-left:8px;color:#909399">启用跨域</span>
        <div style="margin-top:8px">
          <el-radio-group v-model="corsMode" :disabled="!form.plugins.cors.enabled">
            <el-radio label="all">允许所有来源（*）</el-radio>
            <el-radio label="custom">自定义来源</el-radio>
          </el-radio-group>
          <div v-if="corsMode==='custom'" style="margin-top:8px">
            <el-input v-model="corsOrigins" placeholder="逗号分隔 Origin，如 https://example.com, http://localhost:3000" />
          </div>
        </div>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { api } from '../api/admin'
import { ElMessage } from 'element-plus'

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

const v1 = ref<Upstream[]>([
  { name: 'A', url: 'http://localhost:9001', weight: 1 },
  { name: 'B', url: 'http://localhost:9002', weight: 1 },
  { name: 'C', url: 'http://localhost:9003', weight: 1 }
])
const v2 = ref<Upstream[]>([
  { name: 'A', url: 'http://localhost:9011', weight: 1 },
  { name: 'B', url: 'http://localhost:9012', weight: 1 },
  { name: 'C', url: 'http://localhost:9013', weight: 1 }
])
const ips = ref('')
const loading = ref(false)
const connected = ref(false)
const clicked = ref(false)
const corsMode = ref<'all'|'custom'>('all')
const corsOrigins = ref('')

function merge() {
  form.v1 = v1.value
  form.v2 = v2.value
  form.plugins.ipWhitelist.ips = ips.value.split(',').map(s => s.trim()).filter(Boolean)
  form.plugins.cors.enabled = form.plugins.cors.enabled
  if (corsMode.value === 'all') {
    form.plugins.cors.allowAll = true
    form.plugins.cors.origins = []
  } else {
    form.plugins.cors.allowAll = false
    form.plugins.cors.origins = corsOrigins.value.split(',').map(s => s.trim()).filter(Boolean)
  }
}

function ensureMinThree(arr: Upstream[], baseUrls: string[]): Upstream[] {
  const names = ['A', 'B', 'C']
  const out = arr.slice()
  for (let i = out.length; i < 3; i++) {
    out.push({ name: names[i], url: baseUrls[i], weight: 1 })
  }
  return out
}

function prefillFromLocal() {
  try {
    const raw = localStorage.getItem('lumagate.route.users')
    if (!raw) return
    const rt = JSON.parse(raw)
    Object.assign(form, rt)
    v1.value = ensureMinThree(rt.v1 || [], ['http://localhost:9001','http://localhost:9002','http://localhost:9003'])
    v2.value = ensureMinThree(rt.v2 || [], ['http://localhost:9011','http://localhost:9012','http://localhost:9013'])
    ips.value = (rt.plugins?.ipWhitelist?.ips || []).join(',')
    corsMode.value = rt.plugins?.cors?.allowAll ? 'all' : 'custom'
    corsOrigins.value = (rt.plugins?.cors?.origins || []).join(',')
  } catch {}
}

async function load() {
  loading.value = true
  clicked.value = true
  try {
    const { data } = await api.getRoutes()
    connected.value = true
    const rt = (data as any[]).find(r => r.id === 'users') || (data as any[])[0]
    if (rt) {
      Object.assign(form, rt)
      v1.value = ensureMinThree(rt.v1 || [], ['http://localhost:9001','http://localhost:9002','http://localhost:9003'])
      v2.value = ensureMinThree(rt.v2 || [], ['http://localhost:9011','http://localhost:9012','http://localhost:9013'])
      ips.value = (rt.plugins?.ipWhitelist?.ips || []).join(',')
      corsMode.value = rt.plugins?.cors?.allowAll ? 'all' : 'custom'
      corsOrigins.value = (rt.plugins?.cors?.origins || []).join(',')
      try { localStorage.setItem('lumagate.route.users', JSON.stringify(rt)) } catch {}
    }
  } catch (e) {
    connected.value = false
  } finally {
    loading.value = false
  }
}

async function save() {
  try {
    merge()
    await api.updateRoute(form)
    try { localStorage.setItem('lumagate.route.users', JSON.stringify(form)) } catch {}
    ElMessage.success('保存成功')
  } catch (e) {
    ElMessage.error('保存失败')
  }
}

onMounted(() => {
  prefillFromLocal()
})
load()
</script>

<style scoped>
.weights { display: flex; gap: 8px }
</style>
