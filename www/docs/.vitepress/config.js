import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Blog API Documentation",
  description: "Comprehensive API documentation for the Blog API project",
  base: '/docs/', // This ensures the docs are served from /docs path
  cleanUrls: true,
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: '首页', link: '/' },
      { text: 'API 参考', link: '/api-reference/' },
      { text: '入门指南', link: '/getting-started' }
    ],

    sidebar: [
      {
        text: '介绍',
        items: [
          { text: '概述', link: '/' },
          { text: '功能', link: '/features' },
          { text: '架构', link: '/architecture' },
          { text: '入门指南', link: '/getting-started' }
        ]
      },
      {
        text: 'API 参考',
        items: [
          { text: '认证接口', link: '/api-reference/auth' },
          { text: '用户接口', link: '/api-reference/users' },
          { text: '文章接口', link: '/api-reference/articles' },
          { text: '分类接口', link: '/api-reference/categories' },
          { text: '评论接口', link: '/api-reference/comments' },
          { text: '点赞接口', link: '/api-reference/likes' },
          { text: '收藏接口', link: '/api-reference/favorites' },
          { text: '文件接口', link: '/api-reference/files' },
          { text: '管理员接口', link: '/api-reference/admin' }
        ]
      },
      {
        text: '指南',
        items: [
          { text: '错误处理', link: '/error-handling' },
          { text: '限流机制', link: '/rate-limiting' },
          { text: '安全机制', link: '/security' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/codepaintstudio/blog-api' }
    ]
  }
})