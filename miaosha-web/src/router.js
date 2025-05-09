import { createRouter, createWebHistory } from 'vue-router';
import SeckillList from './components/SeckillList.vue';
import SeckillDetail from './components/SeckillDetail.vue';
import OrderList from './components/OrderList.vue';
import LoginModal from './components/LoginModal.vue';

const routes = [
  { path: '/', component: SeckillList },
  { path: '/detail', component: SeckillDetail, name: 'SeckillDetail' }, 
  { path: '/orders', component: OrderList },
  { path: '/login', component: LoginModal }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

// 全局前置守卫，检查用户是否登录
router.beforeEach((to, from, next) => {
  const isLoggedIn = localStorage.getItem('token') !== null;
  // 使用 to.name 准确判断是否跳转到商品详情页
  if (to.name === 'SeckillDetail' && !isLoggedIn) { 
    alert('请先登录！');
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    });
  } else {
    next();
  }
});

export default router;