<template>
  <div class="item">
    <img :src="product.img" alt="商品图片" />
    <h3>{{ product.goodName }}</h3>
    <p>原价: <del>{{ product.originPrice }}</del> 元</p>
    <p>秒杀价: <b>{{ product.price }}</b> 元</p>
    <p>库存: {{ product.stock }}</p>
    <p>秒杀状态：{{ countdownText }}</p>
    <p>秒杀开始时间：{{ formatDateTime(product.startTime) }}</p>
    <p>秒杀结束时间：{{ formatDateTime(product.endTime) }}</p>
    <button @click="$emit('view-detail', product.goodID)">查看详情</button>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
const props = defineProps({ product: Object })
const countdownText = ref('')

function formatDateTime(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const pad = n => n.toString().padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} `
       + `${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

let timer
function updateCountdown() {
  const end = new Date(props.product.endTime).getTime()
  const now = Date.now()
  const left = end - now
  if (left > 0) {
    const s = Math.floor(left / 1000)
    countdownText.value = `${Math.floor(s / 60)}分${s % 60}秒`
  } else {
    countdownText.value = '已结束'
    clearInterval(timer)
  }
}
onMounted(() => {
  updateCountdown()
  timer = setInterval(updateCountdown, 1000)
})
onUnmounted(() => clearInterval(timer))
</script>