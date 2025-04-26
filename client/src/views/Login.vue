<script setup>
import { ref } from 'vue'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const handleSubmit = async () => {
  if (loading.value) return

  loading.value = true
  error.value = ''

  try {
    await userStore.login(email.value, password.value)
  } catch (err) {
    error.value = err.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-[calc(100vh-4rem)] flex items-center justify-center">
    <div class="max-w-md w-full bg-[var(--color-canvas-subtle)] shadow-2xl overflow-hidden">
      <div class="p-6 space-y-6 text-[var(--color-accent-emphasis)]">
        <h2 class="text-3xl font-bold text-center">登录</h2>

        <!-- 错误提示 -->
        <div v-if="error" class="p-4 bg-[#ff000010] border border-red-500/20 shadow-xl">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                <path
                  fill-rule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>
            <div class="ml-3">
              <p class="text-sm font-medium text-red-500">{{ error }}</p>
            </div>
          </div>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-6">
          <div>
            <label for="email" class="block text-sm font-medium text-[var(--color-accent-emphasis)]"
              >邮箱</label
            >
            <input
              id="email"
              v-model="email"
              type="email"
              required
              :disabled="loading"
              class="mt-1 block w-full border bg-[var(--color-canvas-default)] border-[var(--color-accent-emphasis)] px-3 py-2 text-[var(--color-accent-emphasis)] focus:border-[#1f80ff] focus:outline-none focus:ring-1 focus:ring-[#1f80ff] disabled:opacity-50 disabled:cursor-not-allowed"
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
              :disabled="loading"
              class="mt-1 block w-full border bg-[var(--color-canvas-default)] border-[var(--color-accent-emphasis)] px-3 py-2 text-[var(--color-accent-emphasis)] focus:border-[#1f80ff] focus:outline-none focus:ring-1 focus:ring-[#1f80ff] disabled:opacity-50 disabled:cursor-not-allowed"
            />
          </div>

          <button type="submit" :disabled="loading" class="btn-primary w-full">
            <svg
              v-if="loading"
              class="animate-spin -ml-1 mr-2 h-4 w-4"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </form>

        <div class="text-center">
          <router-link
            to="/register"
            class="text-sm text-[var(--color-accent-emphasis)] hover:text-[var(--color-accent-emphasis)]/90 transition-colors"
          >
            还没有账号？点击注册
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>
