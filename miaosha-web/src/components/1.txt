<template>
  <div class="modal">
    <div class="modal-content">
      <h3>登录/注册</h3>
      <input v-model="idInput" placeholder="用户ID" />
      <input v-model="password" type="password" placeholder="密码" />
      <!-- 添加确认密码输入框 -->
      <input v-model="confirmPassword" type="password" placeholder="确认密码" v-if="isRegister" />
      <button @click="login">登录</button>
      <!-- 添加注册按钮 -->
      <button @click="register" v-if="isRegister">注册</button>
      <!-- 添加切换登录/注册模式的按钮 -->
      <button @click="toggleMode">{{ isRegister ? '切换到登录' : '切换到注册' }}</button>
      <button @click="$emit('close')">关闭</button>
      <p v-if="msg">{{ msg }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'

const idInput = ref('')
const password = ref('')
const confirmPassword = ref('')
const msg = ref('')
const isRegister = ref(false)

const emit = defineEmits(['close', 'login'])

// 辅助函数：将输入值转换为uint类型
function convertToUint(value) {
  const num = parseInt(value, 10)
  if (!isNaN(num) && num >= 0) {
    return num
  }
  return null
}

async function login() {
  const id = convertToUint(idInput.value)
  if (id!== null && password.value) {
    try {
      const res = await axios.post('/api/user/login', {
        id: id,
        password: password.value
      })
      if (res.data.message === '登录成功') {
        msg.value = '登录成功'
        const nickName = res.data.nickname
        setTimeout(() => {
          msg.value = ''
          emit('login', nickName)
        }, 500)
      } else {
        msg.value = res.data.message
      }
    } catch (e) {
      msg.value = e.response?.data?.message || '登录失败'
    }
  } else {
    msg.value = '请输入正确的用户ID和密码'
  }
}

async function register() {
  const id = convertToUint(idInput.value)
  if (id!== null && password.value && confirmPassword.value) {
    if (password.value!== confirmPassword.value) {
      msg.value = '两次输入的密码不一致'
      return
    }
    try {
      const res = await axios.post('/api/user/register', {
        id: id,
        password: password.value
      })
      msg.value = '注册成功，请登录'
      isRegister.value = false
      setTimeout(() => {
        msg.value = ''
      }, 1500)
    } catch (e) {
      msg.value = e.response?.data?.message || '注册失败'
    }
  } else {
    msg.value = '请输入正确的用户ID、密码和确认密码'
  }
}

function toggleMode() {
  isRegister.value =!isRegister.value
  msg.value = ''
}
</script>

<style>
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
}
.modal-content {
  background: #fff;
  padding: 30px;
  border-radius: 8px;
  min-width: 300px;
}
input {
  display: block;
  margin: 10px 0;
  width: 100%;
  padding: 5px;
}
button {
  margin-right: 10px;
}
</style>