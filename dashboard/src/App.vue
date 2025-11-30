<template>
  <el-container style="height:100vh">
    <el-header>
      <el-menu mode="horizontal" :default-active="active">
        <el-menu-item index="routes" @click="active='routes'">路由</el-menu-item>
        <el-menu-item index="logs" @click="active='logs'">日志</el-menu-item>
      </el-menu>
    </el-header>
    <el-main>
      <Routes v-if="active==='routes'" />
      <Logs v-else />
    </el-main>
  </el-container>
  
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Routes from './pages/Routes.vue'
import Logs from './pages/Logs.vue'
const active = ref('routes')
</script>

