<template>
    <div class="orders">
      <h2>我的订单</h2>
      <table>
        <thead>
          <tr>
            <th>订单号</th>
            <th>商品名</th>
            <th>价格</th>
            <th>购买数量</th>
            <th>状态</th>
            <th>下单时间</th>
            <th>更新时间</th>
            <th>订单详情</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="order in orders" :key="order.id">
            <td>{{ order.orderNumber }}</td>
            <td>{{ order.goodName }}</td>
            <td>{{ order.goodPrice }}</td>
            <td>{{ order.number }}</td>
            <td>{{ order.status === 0 ? '已取消' : order.status === 1 ? '未支付' : '已完成' }}</td>
            <td>{{ order.createdAt }}</td>
            <td>{{ order.updatedAt }}</td>
             <router-link :to="{ name: 'OrderDetail',query:{orderID:order.ID} }"class="detail-link" >
              查看详情
            </router-link>
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
  const res = await axios.get(`/api/order/list`,{
        headers:{
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
      })
  orders.value = res.data.response.list
  console.log(res.data.response.list)
})

</script>
  
  <style>
  .orders { padding: 20px; }
  table { width: 100%; border-collapse: collapse; }
  th, td { border: 1px solid #eee; padding: 8px; text-align: center; }
  </style>