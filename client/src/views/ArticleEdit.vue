<template>
  <div class="min-h-[calc(100vh-4rem)]">
    <div class="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
      <div
        class="bg-[var(--color-canvas-subtle)] border border-[var(--color-border-default)] rounded-lg p-6"
      >
        <h1 class="text-3xl font-bold text-[var(--color-fg-default)] mb-8">
          编辑帖子
        </h1>

        <form v-if="article" @submit.prevent="handleSubmit" class="space-y-6">
          <div>
            <label
              for="title"
              class="block text-sm font-medium text-[var(--color-fg-default)]"
              >标题</label
            >
            <input
              id="title"
              v-model="title"
              type="text"
              required
              class="mt-1 block w-full rounded-md border bg-[var(--color-canvas-default)] border-[var(--color-border-default)] px-3 py-2 text-[var(--color-fg-default)] focus:border-[var(--color-accent-fg)] focus:outline-none focus:ring-1 focus:ring-[var(--color-accent-fg)]"
            />
          </div>

          <div>
            <label
              for="content"
              class="block text-sm font-medium text-[var(--color-fg-default)] mb-2"
              >内容</label
            >
            <Editor v-model="content" />
          </div>

          <div class="flex justify-end space-x-4">
            <router-link
              :to="'/article/' + route.params.id"
              class="btn-secondary"
            >
              取消
            </router-link>
            <button type="submit" class="btn-primary">保存</button>
          </div>
        </form>

        <div v-else class="text-center py-12">
          <p class="text-[var(--color-fg-muted)]">加载中...</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getArticleById, updateArticle } from "../apis/articles";
import Editor from "../components/Editor.vue";

const route = useRoute();
const router = useRouter();
const article = ref(null);
const userInfo = ref(JSON.parse(localStorage.getItem("userInfo") || "{}"));
const title = ref("");
const content = ref("");

onMounted(async () => {
  try {
    const data = await getArticleById(route.params.id);
    article.value = data;

    // 检查帖子所有权
    if (userInfo.value.id !== article.value.user_id) {
      alert("您没有权限编辑这篇帖子");
      router.push("/");
      return;
    }

    title.value = article.value.title;
    content.value = article.value.content;
  } catch (error) {
    alert(error.message);
    router.push("/");
  }
});

const handleSubmit = async () => {
  try {
    await updateArticle(route.params.id, title.value, content.value);
    router.push("/article/" + route.params.id);
  } catch (error) {
    alert(error.message);
  }
};
</script>
