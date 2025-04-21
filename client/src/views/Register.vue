<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="max-w-md w-full p-6 bg-white rounded-lg shadow-md">
      <h2 class="text-3xl font-bold text-center text-gray-900 mb-8">注册</h2>
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div>
          <label for="username" class="block text-sm font-medium text-gray-700">用户名</label>
          <input
            id="username"
            v-model="username"
            type="text"
            required
            class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none"
          />
        </div>
        
        <div>
          <label for="email" class="block text-sm font-medium text-gray-700">邮箱</label>
          <input
            id="email"
            v-model="email"
            type="email"
            required
            class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none"
          />
        </div>
        
        <div>
          <label for="password" class="block text-sm font-medium text-gray-700">密码</label>
          <input
            id="password"
            v-model="password"
            type="password"
            required
            class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none"
          />
        </div>
        <div>
          <label for="re_password" class="block text-sm font-medium text-gray-700">确认密码</label>
          <input
            id="re_password"
            v-model="re_password"
            type="password"
            required
            class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none"
          />
        </div>
        
        <div>
          <button
            type="submit"
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            注册
          </button>
        </div>
      </form>
      
      <div class="mt-4 text-center">
        <router-link to="/login" class="text-sm text-indigo-600 hover:text-indigo-500">
          已有账号？点击登录
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '../apis/user'

const router = useRouter()
const username = ref('')
const email = ref('')
const password = ref('')
const re_password = ref('')

const handleSubmit = async () => {
  if (password.value !== re_password.value) {
    alert('两次输入的密码不一致')
    return
  }
  try {
    const data = await register(username.value, email.value, password.value, re_password.value)
    localStorage.setItem('token', data.token)
    localStorage.setItem('userInfo', JSON.stringify(data.user))
    router.push('/')
  } catch (error) {
    alert(error.message)
  }
}
</script>
