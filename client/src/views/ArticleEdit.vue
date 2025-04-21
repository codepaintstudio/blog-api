<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
      <div class="bg-white shadow rounded-lg p-6">
        <h1 class="text-3xl font-bold text-gray-900 mb-8">编辑文章</h1>
        
        <form v-if="article" @submit.prevent="handleSubmit" class="space-y-6">
          <div>
            <label for="title" class="block text-sm font-medium text-gray-700">标题</label>
            <input
              id="title"
              v-model="title"
              type="text"
              required
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none"
            />
          </div>
          
          <div>
            <label for="content" class="block text-sm font-medium text-gray-700">内容</label>
            <textarea
              id="content"
              v-model="content"
              rows="10"
              required
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none"
            ></textarea>
          </div>
          
          <div class="flex justify-end space-x-4">
            <router-link
              :to="'/article/' + route.params.id"
              class="inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              取消
            </router-link>
            <button
              type="submit"
              class="inline-flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              保存
            </button>
          </div>
        </form>
        
        <div v-else class="text-center py-12">
          <p class="text-gray-500">加载中...</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticleById, updateArticle } from '../apis/articles'

const route = useRoute()
const router = useRouter()
const article = ref(null)
const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))
const title = ref('')
const content = ref('')

onMounted(async () => {
  try {
    const data = await getArticleById(route.params.id)
    article.value = data
    
    // 检查文章所有权
    if (userInfo.value.id !== article.value.user_id) {
      alert('您没有权限编辑这篇文章')
      router.push('/')
      return
    }
    
    title.value = article.value.title
    content.value = article.value.content
  } catch (error) {
    alert(error.message)
    router.push('/')
  }
})

const handleSubmit = async () => {
  try {
    await updateArticle(route.params.id, title.value, content.value)
    router.push('/article/' + route.params.id)
  } catch (error) {
    alert(error.message)
  }
}</script>
