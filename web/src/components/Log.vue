<script setup>
  import {ref} from "vue";

  const logContext=ref("未初始化\n".repeat(100))

</script>

<template>
 <textarea v-model="logContext" readonly="readonly" >

 </textarea>

</template>

<style scoped>
textarea{
  width: 100%;
  height: 100px;
}
</style>