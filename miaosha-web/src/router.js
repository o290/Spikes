import { createRouter, createWebHistory } from 'vue-router'
import SeckillList from './components/SeckillList.vue'
import SeckillDetail from './components/SeckillDetail.vue'
import OrderList from './components/OrderList.vue'

const routes = [
  { path: '/', component: SeckillList },
  { path: '/detail/:id', component: SeckillDetail },
  { path: '/orders', component: OrderList }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router