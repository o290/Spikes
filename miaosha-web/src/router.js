import { createRouter, createWebHistory } from 'vue-router';
import SeckillList from './components/SeckillList.vue';
import SeckillDetail from './components/SeckillDetail.vue';
import OrderList from './components/OrderList.vue';
import LoginModal from './components/LoginModal.vue';
import OrderDetail from './components/OrderDetail.vue'; 
// 这里是你的路由配置

const routes = [
  { path: '/', component: SeckillList },
  { path: '/detail', component: SeckillDetail, name: 'SeckillDetail' }, 
  { path: '/orders', component: OrderList },
  { path: '/login', component: LoginModal },
  {path:'/order/detail',component:OrderDetail,name:'OrderDetail'}, // 添加订单详情路由
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

// 全局前置守卫，检查用户是否登录
router.beforeEach((to, from, next) => {
  const isLoggedIn = localStorage.getItem('token') !== null;
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