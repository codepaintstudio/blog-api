# Network Demo 教学项目

这是一个用于教学目的的全栈网络应用示例项目，展示了现代 Web 开发的最佳实践和常见功能实现。

## 项目概述

本项目是一个简单的文章管理系统，实现了用户认证和文章的基本 CRUD 操作。项目采用前后端分离架构，提供了良好的用户体验和直观的界面设计。

## 技术栈

### 前端

- Vue 3：使用组合式 API 构建用户界面
- Vue Router：实现前端路由管理
- Pinia：状态管理
- TailwindCSS：原子化 CSS 框架，实现响应式设计

### 后端（已部署）

- RESTful API 架构
- 用户认证和授权
- 数据持久化

## 功能特性

1. 用户管理

   - 用户注册
   - 用户登录

2. 文章管理
   - 创建文章 (仅注册用户)
   - 查看文章
   - 编辑文章（仅作者）
   - 删除文章（仅作者）

## API 接口文档

### 用户相关接口

#### 用户注册

- 请求方式：`POST`
- 路径：`/api/register`
- 请求体：
  ```json
  {
    "email": "user@example.com",
    "password": "password"
  }
  ```

#### 用户登录

- 请求方式：`POST`
- 路径：`/api/login`
- 请求体：
  ```json
  {
    "email": "user@example.com",
    "password": "password"
  }
  ```

### 文章相关接口

#### 获取文章列表

- 请求方式：`GET`
- 路径：`/api/articles`
- 响应：
  ```json
  {
    "articles": [
      {
        "id": "article_id",
        "title": "文章标题",
        "content": "文章内容",
        "user_id": "作者ID",
        "created_at": "创建时间"
      }
    ]
  }
  ```

#### 获取单篇文章

- 请求方式：`GET`
- 路径：`/api/articles/:id`
- 响应：
  ```json
  {
    "id": "article_id",
    "title": "文章标题",
    "content": "文章内容",
    "user_id": "作者ID",
    "created_at": "创建时间"
  }
  ```

#### 创建文章

- 请求方式：`POST`
- 路径：`/api/articles`
- 请求头：`Authorization: Bearer {token}`
- 请求体：
  ```json
  {
    "title": "文章标题",
    "content": "文章内容"
  }
  ```
- 响应：
  ```json
  {
    "id": "article_id",
    "title": "文章标题",
    "content": "文章内容",
    "user_id": "作者ID",
    "created_at": "创建时间"
  }
  ```

#### 更新文章

- 请求方式：`PUT`
- 路径：`/api/articles/:id`
- 请求头：`Authorization: Bearer {token}`
- 请求体：
  ```json
  {
    "title": "更新的标题",
    "content": "更新的内容"
  }
  ```
- 响应：
  ```json
  {
    "id": "article_id",
    "title": "更新的标题",
    "content": "更新的内容",
    "user_id": "作者ID",
    "updated_at": "更新时间"
  }
  ```

#### 删除文章

- 请求方式：`DELETE`
- 路径：`/api/articles/:id`
- 请求头：`Authorization: Bearer {token}`
- 响应状态码：204

## 安全性

- 使用 JWT 进行用户认证
- 文章编辑和删除操作需要验证用户身份和所有权
- 密码加密存储

## 部署信息

后端服务已经完成部署，前端项目可以直接连接使用。

- 部署地址：https://network-demo.hub.feashow.cn
- API 基础 URL：https://doc.apipost.net/docs/44a33117dce0000?locale=zh-cn

## 其他：

如要进行二次开发，请 fork 本仓库并提交 Pull Request。
