<script setup>

import Control from "./Control.vue";
import Log from "./Log.vue";
import {onMounted, ref} from "vue";
import {getUrl} from "../global.js";
import axios from "axios";

const info = ref("未初始化")
const state = ref(0)
const archives = ref([])

function translate(state) {
  switch (state) {
    case 0:
      return "监控中"
    case 1:
      return "同步中"
    case 2:
      return "存档中"
    case 3:
      return "暂停中"
  }
}

function reloadInfo() {
  axios.get(getUrl() + '/state').then(res => {
    info.value = translate(res.data.state)
    state.value = res.data.state
    if (res.data.reloadFlag) {
      reloadArchive()
    }
  }, err => {
    info.value = "无法获取后端信息,请重新刷新页面"
    state.value = -1
    stopTimer()
    console.log(err)
  })
}

function reloadArchive() {
  axios.get(getUrl() + '/archives').then(res => {
    archives.value = res.data.archives
  }, err => {
    console.log(err)
  })
}


let myInterval;

onMounted(() => {
  startTimer()
  reloadArchive()
})

// 启动定时器
function startTimer() {
  myInterval = setInterval(reloadInfo, 1000);
}

// 停止定时器
function stopTimer() {
  clearInterval(myInterval);
}

// 在页面卸载前清除定时器
window.addEventListener('beforeunload', function () {
  stopTimer(); // 在页面卸载前清除定时器
});

</script>

<template>
  {{ info }}
  <div>
    <div>
      <Control
          :state="state"
          :archives="archives"
      />
    </div>
  </div>
<!--  <div>-->
<!--    <Log/>-->
<!--  </div>-->
</template>

<style scoped>

</style>