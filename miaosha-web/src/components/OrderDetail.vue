<template>
  <div v-if="order" class="detail">
    <img :src="order.img" />
    <h2>{{ order.goodName }}</h2>
    <p>订单号:{{ order.orderNumber }}</p>
    <p v-if="msg">订单状态: {{ msg }}</p>
    <p>订单创建时间：{{ order.createdAt }}</p>
    <p>订单更新时间：{{ order.updatedAt }}</p>
    <p>商品数量：{{ order.number }}个</p>
    <p>应付款: <b>{{ order.goodPrice }}</b> 元</p>
    <p>实际付款: <b>{{ order.actualPay }}</b> 元</p>
    <p>请在15分钟内完成支付</p>
    <button :disabled="order.status!=1" @click="pay">立即支付</button>
    
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const order = ref(null)
const msg = ref('')
onMounted(async () => {
  // 从路由的 query 参数中获取商品 ID
  const orderID = route.query.orderID
  console.log('获取到的订单 ID:', orderID)
  if (orderID) {
    try {
      const res = await axios.get(`/api/order/detail`,{
        headers: {
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        },
        params: { orderID }
      }
    )
      order.value = res.data.data
      console.log('获取到的订单详情:', order.value)
      if (order.value.status === 0) {
        msg.value = '订单已取消'
      } else if (order.value.status === 1) {
        msg.value = '订单未支付'
      } else if (order.value.status === 2) {
        msg.value = '订单已完成'
      }
    } catch (error) {
      console.error('获取订单详情失败', error)
      msg.value = '获取订单详情失败，请稍后重试'
    }
  } else {
    console.error('未获取到订单ID')
    msg.value = '未获取到订单ID'
  }
})


async function pay() {
  console.log("支付流程进行中")
}
</script>

<style>
.detail { padding: 20px; }
img { width: 300px; height: 200px; object-fit: cover; }
button { margin: 10px 0; }
</style>