# 入门指南

## 项目简介

Blog API 是一个功能完整的博客平台后端系统，专门为前端开发者和客户端开发者提供了一套完整的博客功能接口。本指南将帮助你快速了解和使用这个 API。

## API 基础知识

### 基本格式

Blog API 采用 RESTful 风格设计，所有接口都遵循统一的格式：

- **基础 URL**: `http://localhost:8080/api/v1/` (开发环境)
- **数据格式**: JSON
- **认证方式**: JWT Bearer Token

### 统一响应格式

所有 API 响应都遵循以下格式：

`json
{
  \"code\": 200,
  \"message\": \"success\",
  \"data\": {},
  \"timestamp\": \"2024-01-01T00:00:00Z\"
}
`

- `code`: HTTP 状态码或自定义业务状态码
- `message`: 响应消息
- `data`: 实际响应数据，可能为空对象
- `timestamp`: 响应时间戳

## 快速开始

### 1. 用户注册

首先需要注册一个账户：

`bash
curl -X POST http://localhost:8080/api/v1/auth/register \\
  -H \"Content-Type: application/json\" \\
  -d '{
    \"username\": \"your_username\",
    \"email\": \"your_email@example.com\",
    \"password\": \"your_password\"
  }'
`

成功注册后会返回用户信息和 Token。

### 2. 用户登录

使用注册的账户登录获取访问令牌：

`bash
curl -X POST http://localhost:8080/api/v1/auth/login \\
  -H \"Content-Type: application/json\" \\
  -d '{
    \"username\": \"your_username\",
    \"password\": \"your_password\"
  }'
`

响应会包含访问令牌（access_token）和刷新令牌（refresh_token）。

### 3. 访问受保护的接口

获取 Token 后，需要在请求头中添加认证信息：

`bash
curl -X GET http://localhost:8080/api/v1/users/profile \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"
`

## 身份认证

### Token 机制

Blog API 使用 JWT（JSON Web Token）进行身份认证：

- **访问令牌（Access Token）**: 有效期 1 小时
- **刷新令牌（Refresh Token）**: 有效期 7 天

当访问令牌过期时，可以使用刷新令牌获取新的访问令牌：

`bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \\
  -H \"Content-Type: application/json\" \\
  -d '{
    \"refresh_token\": \"your_refresh_token\"
  }'
`

### 权限说明

- 公开接口：无需登录即可访问（如文章列表、文章详情）
- 需要登录：需要有效的访问令牌（如创建文章、发表评论）
- 管理员接口：需要管理员权限（如用户管理、内容审核）

## 基本操作

### 创建文章

登录后可以创建文章：

`bash
curl -X POST http://localhost:8080/api/v1/articles \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"title\": \"我的第一篇文章\",
    \"content\": \"这是文章内容\",
    \"description\": \"文章描述\",
    \"category_id\": 1
  }'
`

### 获取文章列表

无需登录即可获取公开文章列表：

`bash
curl -X GET \"http://localhost:8080/api/v1/articles?page=1&per_page=10\"
`

### 发表评论

对文章发表评论（需要登录）：

`bash
curl -X POST http://localhost:8080/api/v1/articles/1/comments \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"content\": \"这是一条评论内容\"
  }'
`

## 错误处理

API 可能返回以下错误码：

- `400`: 请求参数错误
- `401`: 未认证或 Token 过期
- `403`: 权限不足
- `404`: 资源不存在
- `429`: 请求过于频繁（限流）
- `500`: 服务器内部错误

## 限流机制

为保护服务器，API 实施了限流机制：

- IP 限流：100 次/分钟
- 用户限流：1000 次/小时

当达到限流阈值时，会返回 `429 Too Many Requests` 错误。

## 下一步

现在你已经了解了基本的使用方法，可以开始：

1. 查看详细的 [API 参考](./api-reference/)
2. 了解各种 [功能特性](./features)
3. 学习 [安全机制](./security)

如果你在使用过程中遇到任何问题，请查看 [错误处理](./error-handling) 页面或在 GitHub 上提交 Issue。
