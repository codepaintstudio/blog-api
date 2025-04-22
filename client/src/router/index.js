import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('../views/Home.vue'),
      meta: { title: '首页' }
    },
    {
      path: '/login',
      component: () => import('../views/Login.vue'),
      meta: { title: '登录' }
    },
    {
      path: '/register',
      component: () => import('../views/Register.vue'),
      meta: { title: '注册' }
    },
    {
      path: '/article/:id',
      component: () => import('../views/ArticleDetail.vue'),
      meta: { title: '帖子详情' }
    },
    {
      path: '/article/create',
      component: () => import('../views/ArticleCreate.vue'),
      meta: { title: '创建帖子', requiresAuth: true }
    },
    {
      path: '/article/edit/:id',
      component: () => import('../views/ArticleEdit.vue'),
      meta: { title: '编辑帖子', requiresAuth: true }
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - Network Demo` : 'Network Demo'

  // 检查是否需要登录
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
  } else {
    next()
  }
})

export default router
