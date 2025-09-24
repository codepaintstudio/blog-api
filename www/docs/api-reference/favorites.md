# 收藏接口

收藏接口用于管理文章收藏和收藏夹功能。

## 创建收藏夹

### 请求信息

- **接口地址**: `/favorites/folders`
- **请求方式**: `POST`
- **Content-Type**: `application/json`
- **权限要求**: 需要登录
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 收藏夹名称，1-50个字符 |
| description | string | 否 | 收藏夹描述，最多200个字符 |

### 示例请求

```bash
curl -X POST http://localhost:8080/api/v1/favorites/folders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "技术文章",
    "description": "收藏的技术相关文章"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "收藏夹创建成功",
  "data": {
    "id": 1,
    "name": "技术文章",
    "description": "收藏的技术相关文章",
    "article_count": 0,
    "is_default": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 获取收藏夹列表

### 请求信息

- **接口地址**: `/favorites/folders`
- **请求方式**: `GET`
- **权限要求**: 需要登录
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 示例请求

```bash
curl -X GET http://localhost:8080/api/v1/favorites/folders \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取收藏夹列表成功",
  "data": [
    {
      "id": 1,
      "name": "技术文章",
      "description": "收藏的技术相关文章",
      "article_count": 5,
      "is_default": false,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "name": "生活随笔",
      "description": "生活类文章收藏",
      "article_count": 3,
      "is_default": false,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 更新收藏夹

### 请求信息

- **接口地址**: `/favorites/folders/{id}`
- **请求方式**: `PUT`
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
| id | int | 是 | 收藏夹ID |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 否 | 收藏夹名称，1-50个字符 |
| description | string | 否 | 收藏夹描述，最多200个字符 |

### 示例请求

```bash
curl -X PUT http://localhost:8080/api/v1/favorites/folders/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "更新的技术文章",
    "description": "更新的收藏描述"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "收藏夹更新成功",
  "data": {
    "id": 1,
    "name": "更新的技术文章",
    "description": "更新的收藏描述",
    "article_count": 5,
    "is_default": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  },
  "timestamp": "2024-01-02T00:00:00Z"
}
```

## 删除收藏夹

### 请求信息

- **接口地址**: `/favorites/folders/{id}`
- **请求方式**: `DELETE`
- **权限要求**: 需要登录
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 收藏夹ID |

### 示例请求

```bash
curl -X DELETE http://localhost:8080/api/v1/favorites/folders/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "收藏夹删除成功",
  "data": null,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 添加文章到收藏

### 请求信息

- **接口地址**: `/favorites`
- **请求方式**: `POST`
- **Content-Type**: `application/json`
- **权限要求**: 需要登录
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| article_id | int | 是 | 文章ID |
| folder_id | int | 否 | 收藏夹ID，默认为用户的默认收藏夹 |

### 示例请求

```bash
curl -X POST http://localhost:8080/api/v1/favorites \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "article_id": 1,
    "folder_id": 1
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "文章收藏成功",
  "data": {
    "article_id": 1,
    "folder_id": 1,
    "is_favorited": true
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 从收藏中移除文章

### 请求信息

- **接口地址**: `/favorites/{article_id}`
- **请求方式**: `DELETE`
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

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| folder_id | int | 否 | 收藏夹ID，从指定收藏夹中移除，默认为所有收藏夹 |

### 示例请求

```bash
curl -X DELETE "http://localhost:8080/api/v1/favorites/1?folder_id=1" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "文章取消收藏成功",
  "data": {
    "article_id": 1,
    "is_favorited": false
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 获取收藏夹文章列表

### 请求信息

- **接口地址**: `/favorites/folders/{id}/articles`
- **请求方式**: `GET`
- **权限要求**: 需要登录
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 收藏夹ID |

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认为1 |
| per_page | int | 否 | 每页数量，默认为10，最大为50 |

### 示例请求

```bash
curl -X GET "http://localhost:8080/api/v1/favorites/folders/1/articles?page=1&per_page=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取收藏夹文章列表成功",
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
        "favored_at": "2024-01-01T00:00:00Z"
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

## 获取用户收藏列表

### 请求信息

- **接口地址**: `/favorites/articles`
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
curl -X GET "http://localhost:8080/api/v1/favorites/articles?page=1&per_page=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取收藏文章列表成功",
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
        "folder": {
          "id": 1,
          "name": "技术文章"
        },
        "favored_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "per_page": 10,
      "total": 8,
      "total_pages": 1
    }
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 检查文章收藏状态

### 请求信息

- **接口地址**: `/favorites/{article_id}/status`
- **请求方式**: `GET`
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

### 示例请求

```bash
curl -X GET http://localhost:8080/api/v1/favorites/1/status \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取收藏状态成功",
  "data": {
    "article_id": 1,
    "is_favorited": true,
    "folder_id": 1,
    "favored_at": "2024-01-01T00:00:00Z"
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
| 404 | 文章、收藏夹不存在 |
| 409 | 文章已在收藏夹中 |
| 429 | 请求过于频繁 |

## 注意事项

1. 每个用户有且仅有一个默认收藏夹
2. 文章可以被收藏到多个不同的收藏夹中
3. 删除收藏夹时，其中的文章不会被删除，只是从该收藏夹中移除
4. 收藏夹名称在用户的收藏夹中必须唯一
5. 不能收藏自己的文章到收藏夹