<template>
    <div class="orders">
      <h2>我的订单</h2>
      <table>
        <thead>
          <tr>
            <th>订单号</th>
            <th>商品</th>
            <th>价格</th>
            <th>状态</th>
            <th>下单时间</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="order in orders" :key="order.id">
            <td>{{ order.id }}</td>
            <td>{{ order.productName }}</td>
            <td>{{ order.price }}</td>
            <td>{{ order.status }}</td>
            <td>{{ order.time }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </template>
  
  // ... existing code ...
<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
const orders = ref([])

onMounted(async () => {
  // 从后端获取订单列表
  const res = await axios.get('/api/orders')
  orders.value = res.data
})
</script>
// ... existing code ...
  
  <style>
  .orders { padding: 20px; }
  table { width: 100%; border-collapse: collapse; }
  th, td { border: 1px solid #eee; padding: 8px; text-align: center; }
  </style>