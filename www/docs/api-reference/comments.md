# 评论接口

评论接口用于管理文章评论和回复。

## 获取文章评论列表

### 请求信息

- **接口地址**: `/articles/{article_id}/comments`
- **请求方式**: `GET`
- **权限要求**: 公开接口
- **认证方式**: 无需认证

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| article_id | int | 是 | 文章ID |

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认为1 |
| per_page | int | 否 | 每页数量，默认为10，最大为50 |
| sort | string | 否 | 排序方式，支持：created_at, like_count，默认为created_at |

### 示例请求

```bash
curl -X GET "http://localhost:8080/api/v1/articles/1/comments?page=1&per_page=10&sort=created_at"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取文章评论成功",
  "data": {
    "comments": [
      {
        "id": 1,
        "content": "这是一条评论内容",
        "like_count": 5,
        "reply_count": 2,
        "status": "published",
        "is_liked": false,
        "user": {
          "id": 1,
          "username": "commenter",
          "nickname": "评论者",
          "avatar": "http://localhost:8080/uploads/avatars/commenter.jpg"
        },
        "replies": [
          {
            "id": 10,
            "content": "这是对评论的回复",
            "parent_id": 1,
            "like_count": 2,
            "user": {
              "id": 2,
              "username": "replier",
              "nickname": "回复者",
              "avatar": "http://localhost:8080/uploads/avatars/replier.jpg"
            },
            "created_at": "2024-01-01T00:00:10Z"
          }
        ],
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "per_page": 10,
      "total": 25,
      "total_pages": 3
    }
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 发表评论

### 请求信息

- **接口地址**: `/articles/{article_id}/comments`
- **请求方式**: `POST`
- **Content-Type**: `application/json`
- **权限要求**: 需要登录
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| article_id | int | 是 | 文章ID |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| content | string | 是 | 评论内容，1-500个字符 |
| parent_id | int | 否 | 父评论ID，用于回复评论 |

### 示例请求

```bash
curl -X POST http://localhost:8080/api/v1/articles/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "content": "这是一条新的评论"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "评论发表成功",
  "data": {
    "id": 1,
    "content": "这是一条新的评论",
    "like_count": 0,
    "reply_count": 0,
    "status": "published",
    "is_liked": false,
    "user": {
      "id": 1,
      "username": "current_user",
      "nickname": "当前用户",
      "avatar": "http://localhost:8080/uploads/avatars/current.jpg"
    },
    "article_id": 1,
    "parent_id": null,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 获取评论详情

### 请求信息

- **接口地址**: `/comments/{id}`
- **请求方式**: `GET`
- **权限要求**: 公开接口
- **认证方式**: 无需认证

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 评论ID |

### 示例请求

```bash
curl -X GET http://localhost:8080/api/v1/comments/1
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取评论详情成功",
  "data": {
    "id": 1,
    "content": "这是一条评论内容",
    "like_count": 5,
    "reply_count": 2,
    "status": "published",
    "is_liked": false,
    "user": {
      "id": 1,
      "username": "commenter",
      "nickname": "评论者",
      "avatar": "http://localhost:8080/uploads/avatars/commenter.jpg"
    },
    "article": {
      "id": 1,
      "title": "文章标题"
    },
    "parent_id": null,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 获取评论回复列表

### 请求信息

- **接口地址**: `/comments/{id}/replies`
- **请求方式**: `GET`
- **权限要求**: 公开接口
- **认证方式**: 无需认证

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 评论ID |

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认为1 |
| per_page | int | 否 | 每页数量，默认为10，最大为50 |

### 示例请求

```bash
curl -X GET "http://localhost:8080/api/v1/comments/1/replies?page=1&per_page=10"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取评论回复成功",
  "data": {
    "replies": [
      {
        "id": 10,
        "content": "这是对评论的回复",
        "like_count": 2,
        "user": {
          "id": 2,
          "username": "replier",
          "nickname": "回复者",
          "avatar": "http://localhost:8080/uploads/avatars/replier.jpg"
        },
        "created_at": "2024-01-01T00:00:10Z"
      }
    ],
    "pagination": {
      "page": 1,
      "per_page": 10,
      "total": 2,
      "total_pages": 1
    }
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 更新评论

### 请求信息

- **接口地址**: `/comments/{id}`
- **请求方式**: `PUT`
- **Content-Type**: `application/json`
- **权限要求**: 需要登录且为评论作者或管理员
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 评论ID |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| content | string | 是 | 评论内容，1-500个字符 |

### 示例请求

```bash
curl -X PUT http://localhost:8080/api/v1/comments/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "content": "更新后的评论内容"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "评论更新成功",
  "data": {
    "id": 1,
    "content": "更新后的评论内容",
    "like_count": 5,
    "reply_count": 2,
    "status": "published",
    "is_liked": false,
    "user": {
      "id": 1,
      "username": "commenter",
      "nickname": "评论者",
      "avatar": "http://localhost:8080/uploads/avatars/commenter.jpg"
    },
    "article_id": 1,
    "parent_id": null,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  },
  "timestamp": "2024-01-02T00:00:00Z"
}
```

## 删除评论

### 请求信息

- **接口地址**: `/comments/{id}`
- **请求方式**: `DELETE`
- **权限要求**: 需要登录且为评论作者或管理员
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
curl -X DELETE http://localhost:8080/api/v1/comments/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "评论删除成功",
  "data": null,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 获取用户评论列表

### 请求信息

- **接口地址**: `/comments/mine`
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
curl -X GET "http://localhost:8080/api/v1/comments/mine?page=1&per_page=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取我的评论成功",
  "data": {
    "comments": [
      {
        "id": 1,
        "content": "我发表的评论",
        "like_count": 5,
        "reply_count": 0,
        "status": "published",
        "article": {
          "id": 1,
          "title": "文章标题"
        },
        "created_at": "2024-01-01T00:00:00Z"
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

## 审核评论（管理员）

### 请求信息

- **接口地址**: `/comments/{id}/approve` 或 `/comments/{id}/reject`
- **请求方式**: `POST`
- **权限要求**: 管理员
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 评论ID |

### 示例请求（审核通过）

```bash
curl -X POST http://localhost:8080/api/v1/comments/1/approve \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 示例请求（审核拒绝）

```bash
curl -X POST http://localhost:8080/api/v1/comments/1/reject \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "评论审核操作成功",
  "data": null,
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
| 404 | 评论不存在 |
| 429 | 请求过于频繁 |

## 注意事项

1. 评论支持多级回复，通过 `parent_id` 字段关联
2. 评论发布后可能需要审核（根据系统配置）
3. 用户只能修改或删除自己的评论，管理员可以管理所有评论
4. 评论内容长度限制为1-500个字符
5. 评论和回复是两种不同的实体，都存储在 comments 表中，通过 parent_id 区分