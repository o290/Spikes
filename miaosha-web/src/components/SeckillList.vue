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

// 配置axios请求拦截器，确保每个请求都携带token
axios.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

onMounted(async () => {
  try {
    const res = await axios.get('/api/good/list')
    console.log('接口返回内容:', res.data)
    products.value = res.data?.response?.list || []
  } catch (error) {
    if (error.response && error.response.status === 401) {
      // 401状态码表示未授权，通常是token无效或过期，这里可以引导用户重新登录
      console.error('未授权，请重新登录', error)
      // 可以添加跳转到登录页的逻辑，例如：
      // router.push('/login')
    } else {
      console.error('获取商品列表数据失败', error)
    }
  }
})

function viewDetail(id) {
  router.push(`/detail/${id}`)
}
</script>

<style>
.seckill-list { padding: 20px; }
.list { display: flex; gap: 20px; }
</style>