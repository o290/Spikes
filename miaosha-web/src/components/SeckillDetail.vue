<template>
  <div v-if="product" class="detail">
    <img :src="product.img" />
    <h2>{{ product.name }}</h2>
    <p>{{ product.desc || '商品描述...' }}</p>
    <p>原价: <del>{{ product.originPrice }}</del> 元</p>
    <p>秒杀价: <b>{{ product.seckillPrice }}</b> 元</p>
    <p>库存: {{ product.stock }}</p>
    <p>倒计时: {{ countdownText }}</p>
    <button :disabled="!canSeckill" @click="seckill">立即抢购</button>
    <p v-if="msg">{{ msg }}</p>
    <h4>秒杀规则：</h4>
    <ul>
      <li>每人限购1件</li>
      <li>先到先得，售完即止</li>
    </ul>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const product = ref(null)
const countdownText = ref('')
const canSeckill = ref(true)
const msg = ref('')

let timer

onMounted(async () => {
  // 从路由的 query 参数中获取商品 ID
  const goodID = route.query.goodID
  console.log('获取到的商品 ID:', goodID)
  if (goodID) {
    try {
      const res = await axios.get(`/api/good/detail?goodID=${goodID}`,{
        headers:{
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
      })
      product.value = res.data
      updateCountdown()
      timer = setInterval(updateCountdown, 1000)
    } catch (error) {
      console.error('获取商品详情失败', error)
      msg.value = '获取商品详情失败，请稍后重试'
    }
  } else {
    console.error('未获取到商品 ID')
    msg.value = '未获取到商品 ID，请检查链接'
  }
})

function updateCountdown() {
  if (!product.value) return
  const left = new Date(product.value.endTime).getTime() - Date.now()
  if (left > 0) {
    const s = Math.floor(left / 1000)
    countdownText.value = `${Math.floor(s / 60)}分${s % 60}秒`
    canSeckill.value = true
  } else {
    countdownText.value = '已结束'
    canSeckill.value = false
    clearInterval(timer)
  }
}
onUnmounted(() => clearInterval(timer))

async function seckill() {
  try {
    const res = await axios.post('/api/order/spikes', 
      { goodID: product.value.id },
      {
        headers: {
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
      }
    )
    msg.value = '秒杀成功！请到订单页支付'
    setTimeout(() => router.push('/orders'), 1000)
  } catch (e) {
    msg.value = e.response?.data?.message || '秒杀失败'
  }
}
</script>

<style>
.detail { padding: 20px; }
img { width: 300px; height: 200px; object-fit: cover; }
button { margin: 10px 0; }
</style>