<template>
  <div id="app">
    <header>
      <h1>秒杀系统</h1>
      <nav>
        <router-link to="/">秒杀活动</router-link>
        <router-link to="/orders">我的订单</router-link>
        <button v-if="!isLogin" @click="showLogin = true">登录/注册</button>
        <span v-else>欢迎，{{ nickName }}</span>
      </nav>
    </header>
    <router-view @login="onLogin" />
    <LoginModal v-if="showLogin" @close="showLogin = false" @login="onLogin" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import LoginModal from './components/LoginModal.vue'

const showLogin = ref(false)
const isLogin = ref(false)
const nickName = ref('')

function onLogin(name) {
  isLogin.value = true
  nickName.value = name
  showLogin.value = false
}
</script>

<style>
#app { font-family: Arial, sans-serif; }
header { display: flex; justify-content: space-between; align-items: center; padding: 10px 20px; background: #f5f5f5; }
nav { display: flex; gap: 20px; align-items: center; }
nav a { text-decoration: none; color: #333; }
button { padding: 5px 10px; }
</style>