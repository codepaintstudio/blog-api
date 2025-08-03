# Netwok 教学项目

这是一个用于教学目的的全栈网络应用示例项目，展示了现代 Web 开发的最佳实践和常见功能实现。

## 项目概述

本项目是一个简单的帖子管理系统，实现了用户认证和帖子的基本 CRUD 操作。项目采用前后端分离架构，提供了良好的用户体验和直观的界面设计。

<p align="center">
  <img src="./client/public/logo.png" alt="Logo" width="100">
</p>

## 功能特性

1. 用户管理

   - 用户注册
   - 用户登录

2. 帖子管理
   
   - 创建帖子 (仅注册用户)
   - 查看帖子
   - 编辑帖子（仅作者）
   - 删除帖子（仅作者）

## API 接口文档

### 用户相关接口

#### 用户注册

- 请求方式：`POST`
- 路径：`/api/register`
- 请求体：
  ```json
  {
    "username": "x",
    "email": "xxx@xxxx.com",
    "password": "123456",
    "re_password": "123456"
  }
  ```

#### 用户登录

- 请求方式：`POST`
- 路径：`/api/login`
- 请求体：
  ```json
  {
    "email": "xxx@xxxx.com",
    "password": "123456"
  }
  ```

### 帖子相关接口

#### 获取帖子列表

- 请求方式：`GET`
- 路径：`/api/articles`
- 响应：
  ```json
  {
    "articles": [
      {
        "id": "article_id",
        "title": "帖子标题",
        "content": "帖子内容",
        "user_id": "作者ID",
        "created_at": "创建时间"
      },
      {
        "id": "article_id",
        "title": "帖子标题",
        "content": "帖子内容",
        "user_id": "作者ID",
        "created_at": "创建时间"
      },
      ...
    ]
  }
  ```

#### 获取单篇帖子

- 请求方式：`GET`
- 路径：`/api/articles/:id`
- 响应：
  ```json
  {
    "id": "article_id",
    "title": "帖子标题",
    "content": "帖子内容",
    "user_id": "作者ID",
    "created_at": "创建时间"
  }
  ```

#### 创建帖子

- 请求方式：`POST`
- 路径：`/api/articles`
- 请求头：`Authorization: Bearer {token}`
- 请求体：
  ```json
  {
    "title": "帖子标题",
    "content": "帖子内容"
  }
  ```
- 响应：
  ```json
  {
    "id": "article_id",
    "title": "帖子标题",
    "content": "帖子内容",
    "user_id": "作者ID",
    "created_at": "创建时间"
  }
  ```

#### 更新帖子

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

#### 删除帖子

- 请求方式：`DELETE`
- 路径：`/api/articles/:id`
- 请求头：`Authorization: Bearer {token}`
- 响应状态码：204

## 安全性

- 使用 JWT 进行用户认证
- 帖子编辑和删除操作需要验证用户身份和所有权
- 密码加密存储

## 部署信息

后端服务已经完成部署，前端项目可以直接连接使用。

- 部署地址：https://network-client.hub.feashow.cn
- API 基础 URL：https://doc.apipost.net/docs/44a33117dce0000?locale=zh-cn

## 其他：

如要进行二次开发，请 fork 本仓库并提交 Pull Request。
