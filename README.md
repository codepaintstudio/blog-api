# Netwok 教学项目

这是一个用于教学目的的全栈网络应用示例项目，展示了现代 Web 开发的最佳实践和常见功能实现。

## 项目概述

本项目是一个简单的帖子管理系统，实现了用户认证和帖子的基本 CRUD 操作。项目采用前后后端分离架构，提供了良好的用户体验和直观的界面设计。

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

BaseURL: https://network-demo.hub.feashow.cn/

所有 API 响应都遵循统一的格式：
```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

### 用户相关接口

#### 用户注册

```bash
curl -X POST https://network-demo.hub.feashow.cn/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username": "example_user", "email": "user@example.com", "password": "123456", "re_password": "123456"}'
```

成功响应示例：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "example_user",
    "email": "user@example.com"
  }
}
```

#### 用户登录

```bash
curl -X POST https://network-demo.hub.feashow.cn/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "123456"}'
```

成功响应示例：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "example_user",
      "email": "user@example.com"
    }
  }
}
```

### 帖子相关接口

#### 获取帖子列表

```bash
curl -X GET "https://network-demo.hub.feashow.cn/api/articles?page=1&size=10"
```

成功响应示例：
```json
{
  "code": 200,
  "message": "Get articles successfully",
  "data": {
    "articles": [
      {
        "id": 1,
        "title": "帖子标题",
        "content": "帖子内容",
        "user_id": 1,
        "user": {
          "id": 1,
          "username": "example_user",
          "email": "user@example.com"
        },
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "size": 10
  }
}
```

#### 获取单篇帖子

```bash
curl -X GET https://network-demo.hub.feashow.cn/api/articles/1
```

成功响应示例：
```json
{
  "code": 200,
  "message": "Get article successfully",
  "data": {
    "id": 1,
    "title": "帖子标题",
    "content": "帖子内容",
    "user_id": 1,
    "user": {
      "id": 1,
      "username": "example_user",
      "email": "user@example.com"
    },
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

#### 创建帖子

```bash
curl -X POST https://network-demo.hub.feashow.cn/api/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title": "新帖子标题", "content": "新帖子内容"}'
```

成功响应示例：
```json
{
  "code": 200,
  "message": "Create article successfully",
  "data": {
    "id": 1,
    "title": "新帖子标题",
    "content": "新帖子内容",
    "user_id": 1,
    "user": {
      "id": 1,
      "username": "example_user",
      "email": "user@example.com"
    },
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

#### 更新帖子

```bash
curl -X PUT https://network-demo.hub.feashow.cn/api/articles/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title": "更新的标题", "content": "更新的内容"}'
```

成功响应示例：
```json
{
  "code": 200,
  "message": "Update article successfully",
  "data": {
    "id": 1,
    "title": "更新的标题",
    "content": "更新的内容",
    "user_id": 1,
    "user": {
      "id": 1,
      "username": "example_user",
      "email": "user@example.com"
    },
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-02T00:00:00Z"
  }
}
```

#### 删除帖子

```bash
curl -X DELETE https://network-demo.hub.feashow.cn/api/articles/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

成功响应示例：
```json
{
  "code": 200,
  "message": "Delete article successfully",
  "data": null
}
```

## 安全性

- 使用 JWT 进行用户认证
- 帖子编辑和删除操作需要验证用户身份和所有权
- 密码加密存储
- 错误处理和日志记录

## 其他：

如要进行二次开发，请 fork 本仓库并提交 Pull Request。