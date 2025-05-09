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
  try {
    const res = await axios.get('/api/good/list')
    console.log('接口返回内容:', res.data)
    products.value = res.data?.response?.list || []
    console.log('商品列表:', products.value[0].goodID)  
  } catch (error) {
    if (error.response && error.response.status === 401) {
      console.error('未授权，请重新登录', error)
    } else {
      console.error('获取商品列表数据失败', error)
    }
  }
})

function viewDetail(id) {
  console.log('查看详情的商品ID:', id)
  try {
    router.push({
      name: 'SeckillDetail',
      query: { goodID: id }
    });
  } catch (error) {
    console.error('viewDetail函数执行出错', error);
  }
}
</script>

<style>
.seckill-list { padding: 20px; }
.list { display: flex; gap: 20px; }
</style>