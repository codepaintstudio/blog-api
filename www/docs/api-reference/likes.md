# 点赞接口

点赞接口用于管理文章和评论的点赞功能。

## 获取用户点赞的文章列表

### 请求信息

- **接口地址**: `/likes/articles`
- **请求方式**: `GET`
- **权限要求**: 需要登录
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认为1 |
| per_page | int | 否 | 每页数量，默认为10，最大为50 |

### 示例请求

```bash
curl -X GET "http://localhost:8080/api/v1/likes/articles?page=1&per_page=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取点赞文章列表成功",
  "data": {
    "articles": [
      {
        "id": 1,
        "title": "Go语言学习指南",
        "description": "Go语言的基础知识介绍",
        "cover_image": "http://localhost:8080/uploads/covers/go-guide.jpg",
        "user": {
          "id": 2,
          "username": "author",
          "nickname": "作者昵称"
        },
        "like_count": 15,
        "liked_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "per_page": 10,
      "total": 5,
      "total_pages": 1
    }
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 获取用户点赞的评论列表

### 请求信息

- **接口地址**: `/likes/comments`
- **请求方式**: `GET`
- **权限要求**: 需要登录
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认为1 |
| per_page | int | 否 | 每页数量，默认为10，最大为50 |

### 示例请求

```bash
curl -X GET "http://localhost:8080/api/v1/likes/comments?page=1&per_page=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取点赞评论列表成功",
  "data": {
    "comments": [
      {
        "id": 1,
        "content": "这是一条被点赞的评论",
        "like_count": 10,
        "user": {
          "id": 2,
          "username": "commenter",
          "nickname": "评论者"
        },
        "article": {
          "id": 1,
          "title": "文章标题"
        },
        "liked_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "per_page": 10,
      "total": 3,
      "total_pages": 1
    }
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 点赞/取消点赞文章

### 请求信息

- **接口地址**: `/articles/{id}/like`
- **请求方式**: `POST`
- **权限要求**: 需要登录
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 文章ID |

### 示例请求

```bash
curl -X POST http://localhost:8080/api/v1/articles/1/like \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应（点赞）

```json
{
  "code": 200,
  "message": "文章点赞成功",
  "data": {
    "article_id": 1,
    "is_liked": true,
    "like_count": 16
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 成功响应（取消点赞）

```json
{
  "code": 200,
  "message": "文章取消点赞成功",
  "data": {
    "article_id": 1,
    "is_liked": false,
    "like_count": 15
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 点赞/取消点赞评论

### 请求信息

- **接口地址**: `/comments/{id}/like`
- **请求方式**: `POST`
- **权限要求**: 需要登录
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 评论ID |

### 示例请求

```bash
curl -X POST http://localhost:8080/api/v1/comments/1/like \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应（点赞）

```json
{
  "code": 200,
  "message": "评论点赞成功",
  "data": {
    "comment_id": 1,
    "is_liked": true,
    "like_count": 6
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 成功响应（取消点赞）

```json
{
  "code": 200,
  "message": "评论取消点赞成功",
  "data": {
    "comment_id": 1,
    "is_liked": false,
    "like_count": 5
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 获取文章点赞用户列表

### 请求信息

- **接口地址**: `/articles/{id}/likers`
- **请求方式**: `GET`
- **权限要求**: 公开接口
- **认证方式**: 无需认证

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 文章ID |

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认为1 |
| per_page | int | 否 | 每页数量，默认为10，最大为50 |

### 示例请求

```bash
curl -X GET "http://localhost:8080/api/v1/articles/1/likers?page=1&per_page=10"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取文章点赞用户列表成功",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "user1",
        "nickname": "用户1",
        "avatar": "http://localhost:8080/uploads/avatars/user1.jpg",
        "liked_at": "2024-01-01T00:00:00Z"
      },
      {
        "id": 2,
        "username": "user2",
        "nickname": "用户2",
        "avatar": "http://localhost:8080/uploads/avatars/user2.jpg",
        "liked_at": "2024-01-01T00:00:10Z"
      }
    ],
    "pagination": {
      "page": 1,
      "per_page": 10,
      "total": 15,
      "total_pages": 2
    }
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 响应码说明

| HTTP 状态码 | 说明 |
|-------------|------|
| 200 | 操作成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或 Token 过期 |
| 403 | 权限不足 |
| 404 | 文章或评论不存在 |
| 429 | 请求过于频繁 |

## 注意事项

1. 点赞和取消点赞使用同一个接口，系统会自动判断当前状态并执行相应操作
2. 不能给自己发布的内容点赞
3. 点赞状态会实时更新到文章/评论的统计数据中
4. 用户可以查看自己点赞过的文章和评论列表
5. 文章作者可以查看给其文章点赞的用户列表