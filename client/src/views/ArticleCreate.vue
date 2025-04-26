<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { createArticle } from '../apis/articles'
import Editor from '../components/Editor.vue'

const router = useRouter()
const title = ref('')
const content = ref('')
const error = ref('')
const submitting = ref(false)

const handleSubmit = async () => {
  if (submitting.value) return

  submitting.value = true
  error.value = ''

  try {
    await createArticle(title.value, content.value)
    router.push('/')
  } catch (err) {
    error.value = err.message || '发布帖子失败'
    submitting.value = false
  }
}
</script>

<template>
  <div class="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
    <div
      class="bg-[var(--color-canvas-subtle)] shadow-2xl overflow-hidden"
    >
      <div class="p-6">
        <div class="flex justify-between items-center mb-8">
          <h1 class="text-3xl font-bold text-[var(--color-accent-emphasis)]">创建帖子</h1>
          <router-link
            to="/"
            class="inline-flex items-center text-[var(--color-accent-emphasis)] hover:text-[var(--color-accent-emphasis)]/90 transition-colors"
          >
            <svg class="w-5 h-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
              <path
                fill-rule="evenodd"
                d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                clip-rule="evenodd"
              />
            </svg>
            取消
          </router-link>
        </div>

        <!-- 错误提示 -->
        <div v-if="error" class="mb-6 p-4 bg-[#ff000010] border border-red-500/20 shadow-xl">
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
            <label for="title" class="block text-sm font-medium text-[var(--color-accent-emphasis)]"
              >标题</label
            >
            <input
              id="title"
              v-model="title"
              type="text"
              required
              :disabled="submitting"
              class="mt-1 block w-full border bg-[var(--color-canvas-default)] border-[var(--color-accent-emphasis)] px-3 py-2 text-[var(--color-accent-emphasis)] focus:border-[#1f80ff] focus:outline-none focus:ring-1 focus:ring-[#1f80ff] disabled:opacity-50 disabled:cursor-not-allowed"
            />
          </div>

          <div>
            <label
              for="content"
              class="block text-sm font-medium text-[var(--color-accent-emphasis)] mb-2"
              >内容</label
            >
            <Editor v-model="content" :editable="!submitting" />
          </div>

          <div class="flex justify-end space-x-4">
            <router-link to="/" class="btn-secondary" :disabled="submitting"> 取消 </router-link>
            <button type="submit" :disabled="submitting" class="btn-primary">
              <svg
                v-if="submitting"
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
              {{ submitting ? '发布中...' : '发布' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
