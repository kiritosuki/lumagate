<template>
  <el-card>
    <div style="margin-bottom: 12px;">
      <el-input-number v-model="tail" :min="10" />
      <el-button type="primary" @click="load" style="margin-left:8px">刷新</el-button>
    </div>
    <el-table :data="rows" height="70vh">
      <el-table-column prop="line" label="日志" />
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { api } from '../api/admin'

const tail = ref(100)
const rows = ref<{ line: string }[]>([])

async function load() {
  const { data } = await api.tailLogs(tail.value)
  rows.value = (data as string[]).map(s => ({ line: s }))
}

load()
</script>

