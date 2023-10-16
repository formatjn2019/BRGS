<script setup>

import {btRecoverTime, getUrl} from "../global.js";
import axios from "axios";
import {ref} from "vue";

let prop = defineProps(["state", "archives"]);


const flagObj = ref(
    {
      backupFlag: false,
      recoverFlag: false,
      scanningFlag: false,
      pauseFlag: false,
      watchFlag: false,
      syncTempFlag: false,
      reSyncTempFlag: false,
    }
)

function cellServer(serverPath, type) {
  flagObj.value[type] = true
  axios.get(getUrl() + serverPath).then(res => {
    setTimeout(() => {
      flagObj.value[type] = false
    }, btRecoverTime)
  }, err => {
    console.log(err)
    setTimeout(() => {
      flagObj.value[type] = false
    }, btRecoverTime)
  })
}

function recover() {
  let formData = new FormData();
  formData.set("archiveName", selectedArchive.value)
  flagObj.value.recoverFlag = true
  axios.post(getUrl() + "/recover", formData).then(res => {
    setTimeout(() => {
      flagObj.value.recoverFlag = false
    }, btRecoverTime)
  }, err => {
    console.log(err)
    setTimeout(() => {
      flagObj.value.recoverFlag = false
    }, btRecoverTime)
  })
}


const selectedArchive = ref("")

</script>

<template>
  <div class="container">
    <div class="control-column">
      <ul>
        <li>
          <button @click="cellServer('/zipBackup','backupFlag')" :disabled="flagObj.backupFlag || prop.state !== 0">
            zip备份
          </button>
        </li>
        <li>
          <button @click="cellServer('/hardLinkBackup','backupFlag')" :disabled="flagObj.backupFlag || prop.state !== 0">
            硬链接备份
          </button>
        </li>
        <li>
          <button @click="recover" :disabled="flagObj.recoverFlag || prop.state !== 0 || selectedArchive.length === 0">
            还原
          </button>
        </li>
        <li>
          <button @click="cellServer('/scanning','scanningFlag')" :disabled="flagObj.scanningFlag || prop.state !== 0">
            扫描改动
          </button>
        </li>
        <li>
          <button @click="cellServer('/pause','')" :disabled="flagObj.pauseFlag || prop.state !== 0">
            暂停
          </button>
        </li>
        <li>
          <button @click="cellServer('/continue','watchFlag')" :disabled="flagObj.watchFlag || prop.state !== 3">
            继续
          </button>
        </li>
        <li>
          <button @click="cellServer('/sync','syncTempFlag')" :disabled="flagObj.syncTempFlag || prop.state !== 0">
            同步中转
          </button>
        </li>
        <li>
          <button @click="cellServer('/reSync','reSyncTempFlag')"
                  :disabled="flagObj.reSyncTempFlag || prop.state !== 0">
            反向同步
          </button>
        </li>
      </ul>

    </div>
    <div class="backup-column">
      <div v-for="info in archives" class="archive-list">
        <input type="radio" :id="info.name" :value="info.name" v-model="selectedArchive"/>
        <label :for="info.name" v-text="info.type +'\t'+info.name"></label>
      </div>
    </div>
  </div>

</template>

<style scoped>
.container {
  display: flex;
}

.control-column {
  flex: 1;
}

.backup-column {
  flex: 1;
  text-align: left;
  padding-top: 10%;
}


button {
  width: 150px;
}

ul {
  list-style-type: none;
}

li {
  margin-bottom: 10px;
}
</style>