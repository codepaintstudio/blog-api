<template>
  <div class="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
    <!-- 加载状态 -->
    <div v-if="loading" class="text-center py-12">
      <p class="text-[var(--color-fg-muted)]">加载中...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="text-center py-12">
      <p class="text-red-500">{{ error }}</p>
      <div class="mt-6 flex justify-center space-x-4">
        <button @click="fetchArticle()" class="btn-primary">重试</button>
        <router-link to="/" class="btn-secondary"> 返回首页 </router-link>
      </div>
    </div>

    <!-- 帖子内容 -->
    <div
      v-else-if="article"
      class="bg-[var(--color-canvas-subtle)] border border-[var(--color-border-default)] rounded-lg overflow-hidden"
    >
      <div class="p-6">
        <div class="flex justify-between items-start">
          <h1 class="text-3xl font-bold text-[var(--color-fg-default)]">
            {{ article.title }}
          </h1>
          <div
            class="flex space-x-2"
            v-if="isLoggedIn && userInfo?.id == article?.user_id"
          >
            <router-link
              :to="'/article/edit/' + article.id"
              class="text-[var(--color-accent-fg)] hover:text-[var(--color-accent-emphasis)] transition-colors"
            >
              编辑
            </router-link>
            <button
              @click="handleDelete"
              class="text-red-500 hover:text-red-400 transition-colors"
            >
              删除
            </button>
          </div>
        </div>

        <div
          class="mt-4 flex items-center text-sm text-[var(--color-fg-muted)]"
        >
          <span class="mx-1">·</span>
          <span>{{ new Date(article.created_at).toLocaleDateString() }}</span>
        </div>

        <div class="mt-6">
          <ArticleContent v-if="article" :content="article.content" />
        </div>
      </div>
    </div>

    <!-- 返回按钮 -->
    <div v-if="article" class="mt-6">
      <router-link
        to="/"
        class="inline-flex items-center text-[var(--color-accent-fg)] hover:text-[var(--color-accent-emphasis)] transition-colors"
      >
        <svg class="w-5 h-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
          <path
            fill-rule="evenodd"
            d="M9.707 16.707a1 1 0 01-1.414 0l-6-6a1 1 0 010-1.414l6-6a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l4.293 4.293a1 1 0 010 1.414z"
            clip-rule="evenodd"
          />
        </svg>
        返回帖子列表
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useUserStore } from "../stores/user";
import { getArticleById, deleteArticle } from "../apis/articles";
import ArticleContent from "../components/ArticleContent.vue";

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const article = ref(null);
const loading = ref(false);
const error = ref("");

const isLoggedIn = computed(() => userStore.isLoggedIn);
const userInfo = computed(() => userStore.userInfo);

const fetchArticle = async () => {
  loading.value = true;
  error.value = "";
  try {
    const data = await getArticleById(route.params.id);
    article.value = data;
  } catch (err) {
    error.value = err.message || "获取帖子失败";
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchArticle();
});

const handleDelete = async () => {
  if (!confirm("确定要删除这篇帖子吗？")) {
    return;
  }

  loading.value = true;
  try {
    await deleteArticle(article.value.id);
    router.push("/");
  } catch (err) {
    error.value = err.message || "删除帖子失败";
    loading.value = false;
  }
};
</script>
