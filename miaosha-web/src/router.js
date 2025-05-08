// import { createRouter, createWebHistory } from 'vue-router'
// import SeckillList from './components/SeckillList.vue'
// import SeckillDetail from './components/SeckillDetail.vue'
// import OrderList from './components/OrderList.vue'

// const routes = [
//   { path: '/', component: SeckillList },
//   { path: '/detail/:id', component: SeckillDetail },
//   { path: '/orders', component: OrderList }
// ]

// const router = createRouter({
//   history: createWebHistory(),
//   routes
// })

// export default router
import { createRouter, createWebHistory } from 'vue-router'
import SeckillList from './components/SeckillList.vue'
import SeckillDetail from './components/SeckillDetail.vue'
import OrderList from './components/OrderList.vue'
import LoginModal from './components/LoginModal.vue'

const routes = [
  { path: '/', component: SeckillList },
  { path: '/detail/:id', component: SeckillDetail },
  { path: '/orders', component: OrderList },
  { path: '/login', component: LoginModal }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局前置守卫，检查用户是否登录
router.beforeEach((to, from, next) => {
  const isLoggedIn = localStorage.getItem('token') !== null
  if (to.path.includes('/detail/') && !isLoggedIn) {
    alert('请先登录！')
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
  } else {
    next()
  }
})

export default router