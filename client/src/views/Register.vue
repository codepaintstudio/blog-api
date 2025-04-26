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

<template>
  <div class="min-h-[calc(100vh-4rem)] flex items-center justify-center">
    <div class="max-w-md w-full p-6 bg-[var(--color-canvas-subtle)] shadow-2xl">
      <h2 class="text-3xl font-bold text-center text-[var(--color-accent-emphasis)] mb-8">注册</h2>
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div>
          <label
            for="username"
            class="block text-sm font-medium text-[var(--color-accent-emphasis)]"
            >用户名</label
          >
          <input
            id="username"
            v-model="username"
            type="text"
            required
            class="mt-1 block w-full border bg-[var(--color-canvas-default)] border-[var(--color-accent-emphasis)] px-3 py-2 text-[var(--color-accent-emphasis)] focus:border-[#1f80ff] focus:outline-none focus:ring-1 focus:ring-[#1f80ff] shadow-sm"
          />
        </div>

        <div>
          <label for="email" class="block text-sm font-medium text-[var(--color-accent-emphasis)]"
            >邮箱</label
          >
          <input
            id="email"
            v-model="email"
            type="email"
            required
            class="mt-1 block w-full border bg-[var(--color-canvas-default)] border-[var(--color-accent-emphasis)] px-3 py-2 text-[var(--color-accent-emphasis)] focus:border-[#1f80ff] focus:outline-none focus:ring-1 focus:ring-[#1f80ff] shadow-sm"
          />
        </div>

        <div>
          <label
            for="password"
            class="block text-sm font-medium text-[var(--color-accent-emphasis)]"
            >密码</label
          >
          <input
            id="password"
            v-model="password"
            type="password"
            required
            class="mt-1 block w-full border bg-[var(--color-canvas-default)] border-[var(--color-accent-emphasis)] px-3 py-2 text-[var(--color-accent-emphasis)] focus:border-[#1f80ff] focus:outline-none focus:ring-1 focus:ring-[#1f80ff] shadow-sm"
          />
        </div>
        <div>
          <label
            for="re_password"
            class="block text-sm font-medium text-[var(--color-accent-emphasis)]"
            >确认密码</label
          >
          <input
            id="re_password"
            v-model="re_password"
            type="password"
            required
            class="mt-1 block w-full border bg-[var(--color-canvas-default)] border-[var(--color-accent-emphasis)] px-3 py-2 text-[var(--color-accent-emphasis)] focus:border-[#1f80ff] focus:outline-none focus:ring-1 focus:ring-[#1f80ff] shadow-sm"
          />
        </div>

        <div>
          <button type="submit" class="btn-primary w-full">注册</button>
        </div>
      </form>

      <div class="mt-4 text-center">
        <router-link
          to="/login"
          class="text-sm text-[var(--color-accent-emphasis)] hover:text-[var(--color-accent-emphasis)]/90 transition-colors"
        >
          已有账号？点击登录
        </router-link>
      </div>
    </div>
  </div>
</template>
