<template>
  <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
    <div class="px-4 py-6 sm:px-0">
      <!-- 加载状态 -->
      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-500">加载中...</p>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="text-center py-12">
        <p class="text-red-500">{{ error }}</p>
        <button
          @click="fetchArticles()"
          class="mt-4 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
        >
          重试
        </button>
      </div>

      <!-- 空状态 -->
      <div v-else-if="articles.length === 0" class="text-center py-12">
        <p class="text-gray-500">暂无文章</p>
        <router-link
          v-if="isLoggedIn"
          to="/article/create"
          class="mt-4 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
        >
          发布第一篇文章
        </router-link>
      </div>

      <!-- 文章列表 -->
      <div v-else class="space-y-6">
        <div
          v-for="article in articles"
          :key="article.id"
          class="bg-white shadow rounded-lg p-6"
        >
          <div class="flex justify-between items-start">
            <div class="flex-1">
              <h2 class="text-xl font-semibold text-gray-900">
                <router-link
                  :to="'/article/' + article.id"
                  class="hover:text-indigo-600"
                >
                  {{ article.title }}
                </router-link>
              </h2>
              <p class="mt-2 text-gray-600 line-clamp-3">
                {{ article.content }}
              </p>
              <div class="mt-4 flex items-center text-sm text-gray-500">
                <span>{{ article.author?.name || "未知用户" }}</span>
                <span class="mx-1">·</span>
                <span>{{
                  new Date(article.created_at).toLocaleDateString()
                }}</span>
              </div>
            </div>
            <div class="flex space-x-2 ml-4" v-if="isLoggedIn">
              <router-link
                :to="'/article/edit/' + article.id"
                class="text-indigo-600 hover:text-indigo-900"
              >
                编辑
              </router-link>
              <button
                @click="handleDelete(article.id)"
                class="text-red-600 hover:text-red-900"
              >
                删除
              </button>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div
          v-if="total > pageSize"
          class="mt-6 flex justify-between items-center"
        >
          <div class="text-sm text-gray-700">共 {{ total }} 篇文章</div>
          <nav
            class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px"
          >
            <button
              @click="handlePageChange(currentPage - 1)"
              :disabled="currentPage === 1"
              class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              上一页
            </button>
            <div
              class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700"
            >
              {{ currentPage }} / {{ Math.ceil(total / pageSize) }}
            </div>
            <button
              @click="handlePageChange(currentPage + 1)"
              :disabled="currentPage * pageSize >= total"
              class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              下一页
            </button>
          </nav>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import { useUserStore } from "../stores/user";
import { getArticles, deleteArticle } from "../apis/articles";

const router = useRouter();
const userStore = useUserStore();

const articles = ref([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);
const loading = ref(false);
const error = ref("");

const isLoggedIn = computed(() => userStore.isLoggedIn);
const currentUser = computed(() => userStore.userInfo);

const fetchArticles = async (page = 1) => {
  loading.value = true;
  error.value = "";
  try {
    const {
      articles: list,
      total: count,
      page: current,
    } = await getArticles(page, pageSize.value);
    articles.value = list;
    total.value = count;
    currentPage.value = current;
  } catch (err) {
    error.value = err.message || "获取文章列表失败";
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchArticles();
});

const handlePageChange = (page) => {
  if (page >= 1 && page * pageSize.value <= total.value) {
    fetchArticles(page);
  }
};

const handleDelete = async (id) => {
  if (confirm("确定要删除这篇文章吗？")) {
    try {
      await deleteArticle(id);
      await fetchArticles(currentPage.value);
    } catch (error) {
      alert(error.message);
    }
  }
};
</script>
