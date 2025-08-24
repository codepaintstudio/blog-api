# 📝 Blog API - 前端学习专用接口

> 🎯 **专为前端开发者设计** - 免费的博客 API 接口，用于学习和练习前端网络请求、状态管理、用户认证等技术

## 📚 目录
- [🌟 项目背景](#🌟-项目背景)
- [🚀 快速开始](#🚀-快速开始)
  - [基础信息](#基础信息)
  - [统一响应格式](#统一响应格式)
- [📚 完整 API 文档](#📚-完整-api-文档)
  - [🔐 用户认证](#🔐-用户认证)
    - [1. 用户注册](#1-用户注册)
    - [2. 用户登录](#2-用户登录)
  - [👤 用户信息查询](#👤-用户信息查询)
    - [1. 获取用户基本信息](#1-获取用户基本信息)
    - [2. 获取用户详情（含文章统计）](#2-获取用户详情含文章统计)
  - [📝 文章管理](#📝-文章管理)
    - [1. 获取文章列表（支持高级搜索）](#1-获取文章列表支持高级搜索)
    - [2. 获取文章详情](#2-获取文章详情)
    - [3. 创建文章 🔒 (需要认证)](#3-创建文章-🔒-需要认证)
    - [4. 更新文章 🔒 (需要认证 + 作者权限)](#4-更新文章-🔒-需要认证--作者权限)
    - [5. 删除文章 🔒 (需要认证 + 作者权限)](#5-删除文章-🔒-需要认证--作者权限)
  - [📊 统计信息](#📊-统计信息)
    - [获取系统统计](#获取系统统计)
- [💡 前端开发最佳实践](#💡-前端开发最佳实践)
  - [1. Token 管理](#1-token-管理)
  - [2. 错误处理](#2-错误处理)
  - [3. TypeScript 接口定义](#3-typescript-接口定义)
- [⚠️ 重要注意事项](#⚠️-重要注意事项)
  - [🚦 限流规则](#🚦-限流规则)
  - [🔒 安全提醒](#🔒-安全提醒)
  - [🌐 跨域支持](#🌐-跨域支持)
  - [📱 移动端适配](#📱-移动端适配)
- [🎯 学习建议](#🎯-学习建议)
- [📞 技术支持](#📞-技术支持)


## 🌟 项目背景

这是一个**完全免费**的博客 API 服务，专门为前端开发者提供：

- ✅ **真实的后端接口**：完整的用户认证 + 文章 CRUD 功能
- ✅ **无需搭建后端**：直接调用线上接口，专注前端开发
- ✅ **学习友好**：接口设计规范，错误信息清晰
- ✅ **功能完整**：注册登录、文章管理、搜索排序一应俱全

**适合学习场景：**

- React/Vue/Angular 项目练手
- 网络请求库使用 (axios, fetch)
- 用户认证流程实现 (JWT)
- 状态管理练习 (Redux, Vuex, Pinia)
- TypeScript 接口定义

## 🚀 快速开始

### 基础信息

- **API 地址**：`https://network-demo.hub.feashow.cn`
- **请求格式**：JSON
- **认证方式**：JWT Bearer Token
- **限流规则**：每秒 20 个请求，突发 30 个

### 统一响应格式

所有接口都返回统一的 JSON 格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    /* 具体数据 */
  }
}
```

**状态码说明：**

- `200` - 成功
- `400` - 请求参数错误
- `401` - 未认证或 token 无效
- `403` - 无权限操作
- `404` - 资源不存在
- `429` - 请求频率限制
- `500` - 服务器错误

## 📚 完整 API 文档

### 🔐 用户认证

#### 1. 用户注册

```http
POST /api/user/register
```

**请求示例：**

```javascript
// 使用 fetch
const registerUser = async () => {
  const response = await fetch(
    "https://network-demo.hub.feashow.cn/api/user/register",
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: "测试用户",
        email: "test@example.com",
        password: "123456",
        re_password: "123456",
      }),
    }
  );

  const data = await response.json();
  console.log(data);
};

// 使用 axios
import axios from "axios";

const registerUser = async () => {
  try {
    const response = await axios.post(
      "https://network-demo.hub.feashow.cn/api/user/register",
      {
        username: "测试用户",
        email: "test@example.com",
        password: "123456",
        re_password: "123456",
      }
    );
    console.log(response.data);
  } catch (error) {
    console.error("注册失败:", error.response.data);
  }
};
```

**响应示例：**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "测试用户",
    "email": "test@example.com"
  }
}
```

#### 2. 用户登录

```http
POST /api/user/login
```

**请求示例：**

```javascript
const loginUser = async () => {
  const response = await fetch(
    "https://network-demo.hub.feashow.cn/api/user/login",
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email: "test@example.com",
        password: "123456",
      }),
    }
  );

  const data = await response.json();

  if (data.code === 200) {
    // 保存 token 到 localStorage
    localStorage.setItem("token", data.data.token);
    localStorage.setItem("user", JSON.stringify(data.data.user));
  }
};
```

**响应示例：**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "测试用户",
      "email": "test@example.com"
    }
  }
}
```

### 👤 用户信息查询

#### 1. 获取用户基本信息

```http
GET /api/users/:id
```

**请求示例：**

```javascript
const getUserInfo = async (userId) => {
  const response = await fetch(
    `https://network-demo.hub.feashow.cn/api/users/${userId}`
  );
  const data = await response.json();
  return data;
};
```

#### 2. 获取用户详情（含文章统计）

```http
GET /api/users/:id/detail
```

**响应示例：**

```json
{
  "code": 200,
  "message": "Get user detail successfully",
  "data": {
    "id": 1,
    "username": "测试用户",
    "email": "test@example.com",
    "article_count": 5
  }
}
```

### 📝 文章管理

#### 1. 获取文章列表（支持高级搜索）

```http
GET /api/articles
```

**查询参数：**

- `page` - 页码（默认 1）
- `size` - 每页数量（默认 10）
- `search` - 搜索关键词
- `user_id` - 按用户筛选
- `sort_by` - 排序字段（created_at, updated_at, title）
- `order` - 排序方向（asc, desc）

**请求示例：**

```javascript
// 基础获取文章列表
const getArticles = async (page = 1, size = 10) => {
  const response = await fetch(
    `https://network-demo.hub.feashow.cn/api/articles?page=${page}&size=${size}`
  );
  return response.json();
};

// 搜索文章
const searchArticles = async (keyword) => {
  const response = await fetch(
    `https://network-demo.hub.feashow.cn/api/articles?search=${encodeURIComponent(
      keyword
    )}`
  );
  return response.json();
};

// 获取某用户的文章
const getUserArticles = async (userId) => {
  const response = await fetch(
    `https://network-demo.hub.feashow.cn/api/articles?user_id=${userId}&sort_by=created_at&order=desc`
  );
  return response.json();
};
```

**响应示例：**

```json
{
  "code": 200,
  "message": "Get articles successfully",
  "data": {
    "articles": [
      {
        "id": 1,
        "title": "我的第一篇文章",
        "content": "这是文章内容...",
        "user_id": 1,
        "user": {
          "id": 1,
          "username": "测试用户",
          "email": "test@example.com"
        },
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "size": 10
  }
}
```

#### 2. 获取文章详情

```http
GET /api/articles/:id
```

**请求示例：**

```javascript
const getArticleDetail = async (articleId) => {
  const response = await fetch(
    `https://network-demo.hub.feashow.cn/api/articles/${articleId}`
  );
  return response.json();
};
```

#### 3. 创建文章 🔒 (需要认证)

```http
POST /api/articles
Authorization: Bearer {token}
```

**请求示例：**

```javascript
const createArticle = async (title, content) => {
  const token = localStorage.getItem("token");

  const response = await fetch(
    "https://network-demo.hub.feashow.cn/api/articles",
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        title: title,
        content: content,
      }),
    }
  );

  return response.json();
};
```

#### 4. 更新文章 🔒 (需要认证 + 作者权限)

```http
PUT /api/articles/:id
Authorization: Bearer {token}
```

**请求示例：**

```javascript
const updateArticle = async (articleId, title, content) => {
  const token = localStorage.getItem("token");

  const response = await fetch(
    `https://network-demo.hub.feashow.cn/api/articles/${articleId}`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        title: title,
        content: content,
      }),
    }
  );

  return response.json();
};
```

#### 5. 删除文章 🔒 (需要认证 + 作者权限)

```http
DELETE /api/articles/:id
Authorization: Bearer {token}
```

**请求示例：**

```javascript
const deleteArticle = async (articleId) => {
  const token = localStorage.getItem("token");

  const response = await fetch(
    `https://network-demo.hub.feashow.cn/api/articles/${articleId}`,
    {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  return response.json();
};
```

### 📊 统计信息

#### 获取系统统计

```http
GET /api/articles/stats
```

**响应示例：**

```json
{
  "code": 200,
  "message": "Get stats successfully",
  "data": {
    "total_articles": 156,
    "total_users": 23
  }
}
```

## 💡 前端开发最佳实践

### 1. Token 管理

```javascript
// 封装API请求工具
class ApiClient {
  constructor() {
    this.baseURL = "https://network-demo.hub.feashow.cn";
  }

  // 获取token
  getToken() {
    return localStorage.getItem("token");
  }

  // 通用请求方法
  async request(url, options = {}) {
    const token = this.getToken();

    const config = {
      headers: {
        "Content-Type": "application/json",
        ...(token && { Authorization: `Bearer ${token}` }),
        ...options.headers,
      },
      ...options,
    };

    const response = await fetch(`${this.baseURL}${url}`, config);
    const data = await response.json();

    // 处理认证失败
    if (data.code === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      // 重定向到登录页
      window.location.href = "/login";
    }

    return data;
  }
}

const api = new ApiClient();
```

### 2. 错误处理

```javascript
// React 示例
const useApi = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleApiCall = async (apiFunction) => {
    try {
      setLoading(true);
      setError(null);
      const result = await apiFunction();
      return result;
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { loading, error, handleApiCall };
};
```

### 3. TypeScript 接口定义

```typescript
// 用户相关接口
interface User {
  id: number;
  username: string;
  email: string;
}

interface LoginRequest {
  email: string;
  password: string;
}

interface LoginResponse {
  token: string;
  user: User;
}

// 文章相关接口
interface Article {
  id: number;
  title: string;
  content: string;
  user_id: number;
  user: User;
  created_at: string;
  updated_at: string;
}

interface ArticleListParams {
  page?: number;
  size?: number;
  search?: string;
  user_id?: number;
  sort_by?: "created_at" | "updated_at" | "title";
  order?: "asc" | "desc";
}

// API响应格式
interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}
```

## ⚠️ 重要注意事项

### 🚦 限流规则

- **正常请求**：每秒最多 20 个请求
- **突发请求**：允许短时间内 30 个请求
- **违规处理**：5 次超限后 IP 将被封禁 30 分钟
- **建议**：在循环或批量操作中添加适当延时

```javascript
// 批量操作示例 - 避免触发限流
const batchCreateArticles = async (articles) => {
  for (let i = 0; i < articles.length; i++) {
    await createArticle(articles[i]);

    // 添加延时避免触发限流
    if (i < articles.length - 1) {
      await new Promise((resolve) => setTimeout(resolve, 100));
    }
  }
};
```

### 🔒 安全提醒

- Token 存储在 localStorage 中，生产环境建议考虑更安全的存储方式
- 仅在 HTTPS 环境下传输敏感信息
- 定期检查 token 有效性，及时处理过期情况

### 🌐 跨域支持

API 已配置 CORS，支持跨域请求，无需额外配置。

### 📱 移动端适配

所有接口均支持移动端访问，响应格式统一。

## 🎯 学习建议

1. **入门练习**：从用户注册登录开始，理解 JWT 认证流程
2. **进阶功能**：实现文章的增删改查，掌握 RESTful API 使用
3. **高级特性**：使用搜索和排序功能，学习复杂参数传递
4. **状态管理**：结合 React/Vue 状态管理库，实现完整的前端应用

## 📞 技术支持

如遇到问题或有改进建议：

- 📧 通过项目 Issues 反馈
- 💬 欢迎提出新功能需求

---

**🎉 祝你学习愉快！这个 API 将陪伴你从前端新手成长为高手！**