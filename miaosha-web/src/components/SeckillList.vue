<template>
    <div class="seckill-list">
      <h2>秒杀商品列表</h2>
      <div class="list">
        <SeckillItem
          v-for="item in products"
          :key="item.id"
          :product="item"
          @view-detail="viewDetail"
        />
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import SeckillItem from './SeckillItem.vue'
  import { useRouter } from 'vue-router'
  import axios from 'axios'
  
  const products = ref([])
  const router = useRouter()
  
  onMounted(async () => {
  const res = await axios.get('/api/good/list')
  console.log('接口返回内容:', res.data)
  products.value = res.data?.response?.list || []
})
  
  function viewDetail(id) {
    router.push(`/detail/${id}`)
  }
  </script>
  
  <style>
  .seckill-list { padding: 20px; }
  .list { display: flex; gap: 20px; }
  </style>