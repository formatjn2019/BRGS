
<script setup>
import QRCodeVue3 from "qrcode-vue3";
import { computed} from "vue";

let locationUrl = location.href
const isDesktop = computed(() => {
  return !Boolean(navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i))
})


</script>

<template>
  <div v-show="isDesktop">
    <QRCodeVue3
        :value="locationUrl"
    />
    {{ locationUrl }}
  </div>
</template>
