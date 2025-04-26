<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '../stores/user'
import { getArticles, deleteArticle } from '../apis/articles'

const userStore = useUserStore()

const articles = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const error = ref('')

const isLoggedIn = computed(() => userStore.isLoggedIn)
const userInfo = computed(() => userStore.userInfo)

const fetchArticles = async (page = 1) => {
  loading.value = true
  error.value = ''
  try {
    const { articles: list, total: count, page: current } = await getArticles(page, pageSize.value)
    articles.value = list
    total.value = count
    currentPage.value = current
  } catch (err) {
    error.value = err.message || '获取帖子列表失败'
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  if (page >= 1 && page <= Math.ceil(total.value / pageSize.value)) {
    fetchArticles(page)
  }
}

const handleDelete = async (id) => {
  if (confirm('确定要删除这篇帖子吗？')) {
    try {
      await deleteArticle(id)
      await fetchArticles(currentPage.value)
    } catch (error) {
      alert(error.message)
    }
  }
}

onMounted(() => {
  fetchArticles()
})
</script>

<template>
  <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
    <div class="px-4 py-6 sm:px-0">
      <!-- 加载状态 -->
      <div v-if="loading" class="text-center py-12">
        <p class="text-[var(--color-fg-muted)]">加载中...</p>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="text-center py-12">
        <p class="text-red-500">{{ error }}</p>
        <button @click="fetchArticles()" class="btn-primary mt-4">重试</button>
      </div>

      <!-- 空状态 -->
      <div v-else-if="articles.length === 0" class="text-center py-12">
        <p class="text-[var(--color-fg-muted)]">暂无帖子</p>
        <router-link v-if="isLoggedIn" to="/article/create" class="btn-primary mt-4">
          发布第一篇帖子
        </router-link>
      </div>

      <!-- 帖子列表 -->
      <div v-else class="space-y-6">
        <div
          v-for="article in articles"
          :key="`${article.id}-${article.created_at}`"
          class="duration-500 bg-[var(--color-canvas-subtle)] p-6 shadow-xl"
        >
          <div class="flex justify-between items-start">
            <div class="flex-1">
              <h2 class="text-xl font-semibold text-[var(--color-accent-emphasis)]">
                <router-link
                  :to="'/article/' + article.id"
                  class="hover:text-[#1f80ff] transition-colors"
                >
                  {{ article.title }}
                </router-link>
              </h2>
              <p class="mt-2 line-clamp-3 text-black dark:text-white">
                {{ article.content }}
              </p>
              <div class="mt-4 flex items-center text-sm text-[var(--color-fg-muted)]">
                <span>{{ new Date(article.created_at).toLocaleString() }}</span>
              </div>
            </div>
            <div class="flex space-x-4 ml-4" v-if="isLoggedIn && userInfo?.id === article.user_id">
              <router-link
                :to="'/article/edit/' + article.id"
                class="text-[#1f80ff] hover:text-[#1f80ff]/90 transition-colors"
              >
                编辑
              </router-link>
              <button
                @click="handleDelete(article.id)"
                class="text-red-500 hover:text-red-400 transition-colors"
              >
                删除
              </button>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div class="mt-6 flex justify-between items-center">
          <div class="text-sm text-[var(--color-fg-muted)]">共 {{ total }} 篇帖子</div>
          <nav class="flex rounded-md -space-x-px">
            <button
              @click="handlePageChange(currentPage - 1)"
              :disabled="currentPage === 1"
              class="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
            >
              上一页
            </button>
            <div
              class="relative inline-flex items-center px-4 py-2 border border-[var(--color-accent-emphasis)] bg-[var(--color-canvas-subtle)] text-sm font-medium text-[var(--color-fg-default)]"
            >
              {{ currentPage }} / {{ Math.ceil(total / pageSize) }}
            </div>
            <button
              @click="handlePageChange(currentPage + 1)"
              :disabled="currentPage * pageSize >= total"
              class="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
            >
              下一页
            </button>
          </nav>
        </div>
      </div>
    </div>
  </div>
</template>
