# 文章接口

文章接口用于创建、获取、更新和删除文章。

## 获取文章列表

### 请求信息

- **接口地址**: `/articles`
- **请求方式**: `GET`
- **权限要求**: 公开接口
- **认证方式**: 无需认证

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认为1 |
| per_page | int | 否 | 每页数量，默认为10，最大为100 |
| category_id | int | 否 | 分类ID，用于筛选特定分类的文章 |
| search | string | 否 | 搜索关键词，对标题和内容进行模糊匹配 |
| sort | string | 否 | 排序方式，支持：created_at, updated_at, view_count, like_count，默认为created_at |

### 示例请求

```bash
curl -X GET "http://localhost:8080/api/v1/articles?page=1&per_page=10&category_id=1"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取文章列表成功",
  "data": {
    "articles": [
      {
        "id": 1,
        "title": "文章标题",
        "description": "文章描述",
        "content": "文章内容摘要",
        "cover_image": "http://localhost:8080/uploads/covers/article1.jpg",
        "view_count": 100,
        "like_count": 10,
        "comment_count": 5,
        "favorite_count": 3,
        "status": "published",
        "visibility": "public",
        "user": {
          "id": 1,
          "username": "author_name",
          "nickname": "作者昵称",
          "avatar": "http://localhost:8080/uploads/avatars/author.jpg"
        },
        "category": {
          "id": 1,
          "name": "技术分享"
        },
        "tags": ["go", "gin", "api"],
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "per_page": 10,
      "total": 100,
      "total_pages": 10
    }
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 获取文章详情

### 请求信息

- **接口地址**: `/articles/{id}`
- **请求方式**: `GET`
- **权限要求**: 公开接口
- **认证方式**: 无需认证

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 文章ID |

### 示例请求

```bash
curl -X GET http://localhost:8080/api/v1/articles/1
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取文章详情成功",
  "data": {
    "id": 1,
    "title": "文章标题",
    "content": "完整文章内容",
    "description": "文章描述",
    "cover_image": "http://localhost:8080/uploads/covers/article1.jpg",
    "view_count": 101,
    "like_count": 10,
    "comment_count": 5,
    "favorite_count": 3,
    "status": "published",
    "visibility": "public",
    "user": {
      "id": 1,
      "username": "author_name",
      "nickname": "作者昵称",
      "avatar": "http://localhost:8080/uploads/avatars/author.jpg"
    },
    "category": {
      "id": 1,
      "name": "技术分享"
    },
    "tags": ["go", "gin", "api"],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 创建文章

### 请求信息

- **接口地址**: `/articles`
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
| title | string | 是 | 文章标题，最多200个字符 |
| content | string | 是 | 文章内容 |
| description | string | 否 | 文章描述，最多500个字符 |
| category_id | int | 否 | 分类ID |
| cover_image | string | 否 | 封面图URL |
| tags | array | 否 | 标签数组，最多5个标签 |
| visibility | string | 否 | 可见性，public或private，默认为public |
| status | string | 否 | 状态，published或draft，默认为published |

### 示例请求

```bash
curl -X POST http://localhost:8080/api/v1/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "title": "新文章标题",
    "content": "文章的详细内容",
    "description": "文章摘要",
    "category_id": 1,
    "tags": ["go", "api"],
    "visibility": "public",
    "status": "published"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "文章创建成功",
  "data": {
    "id": 1,
    "title": "新文章标题",
    "content": "文章的详细内容",
    "description": "文章摘要",
    "cover_image": "http://localhost:8080/uploads/covers/article1.jpg",
    "view_count": 0,
    "like_count": 0,
    "comment_count": 0,
    "favorite_count": 0,
    "status": "published",
    "visibility": "public",
    "user_id": 1,
    "category_id": 1,
    "tags": ["go", "api"],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 更新文章

### 请求信息

- **接口地址**: `/articles/{id}`
- **请求方式**: `PUT`
- **Content-Type**: `application/json`
- **权限要求**: 需要登录且为文章作者
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 文章ID |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| title | string | 否 | 文章标题，最多200个字符 |
| content | string | 否 | 文章内容 |
| description | string | 否 | 文章描述，最多500个字符 |
| category_id | int | 否 | 分类ID |
| cover_image | string | 否 | 封面图URL |
| tags | array | 否 | 标签数组，最多5个标签 |
| visibility | string | 否 | 可见性，public或private |
| status | string | 否 | 状态，published或draft |

### 示例请求

```bash
curl -X PUT http://localhost:8080/api/v1/articles/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "title": "更新后的文章标题",
    "content": "更新后的文章内容",
    "description": "更新后的摘要"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "文章更新成功",
  "data": {
    "id": 1,
    "title": "更新后的文章标题",
    "content": "更新后的文章内容",
    "description": "更新后的摘要",
    "cover_image": "http://localhost:8080/uploads/covers/article1.jpg",
    "view_count": 101,
    "like_count": 10,
    "comment_count": 5,
    "favorite_count": 3,
    "status": "published",
    "visibility": "public",
    "user_id": 1,
    "category_id": 1,
    "tags": ["go", "api"],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  },
  "timestamp": "2024-01-02T00:00:00Z"
}
```

## 删除文章

### 请求信息

- **接口地址**: `/articles/{id}`
- **请求方式**: `DELETE`
- **权限要求**: 需要登录且为文章作者或管理员
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
curl -X DELETE http://localhost:8080/api/v1/articles/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "文章删除成功",
  "data": null,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 搜索文章

### 请求信息

- **接口地址**: `/articles/search`
- **请求方式**: `GET`
- **权限要求**: 公开接口
- **认证方式**: 无需认证

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| keyword | string | 是 | 搜索关键词 |
| page | int | 否 | 页码，默认为1 |
| per_page | int | 否 | 每页数量，默认为10，最大为100 |

### 示例请求

```bash
curl -X GET "http://localhost:8080/api/v1/articles/search?keyword=Go&page=1&per_page=10"
```

### 成功响应

```json
{
  "code": 200,
  "message": "搜索文章成功",
  "data": {
    "articles": [
      {
        "id": 1,
        "title": "Go语言学习指南",
        "description": "关于Go语言的详细学习指南",
        "cover_image": "http://localhost:8080/uploads/covers/go-tutorial.jpg",
        "view_count": 150,
        "like_count": 20,
        "comment_count": 8,
        "user": {
          "id": 1,
          "username": "go_teacher",
          "nickname": "Go语言导师"
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

## 响应码说明

| HTTP 状态码 | 说明 |
|-------------|------|
| 200 | 操作成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或 Token 过期 |
| 403 | 权限不足（非文章作者或管理员） |
| 404 | 文章不存在 |
| 429 | 请求过于频繁 |

## 注意事项

1. 文章内容较长时，列表接口会返回内容摘要而非完整内容
2. 创建和更新文章时，所有字段都是可选的，只提供需要更新的字段
3. 文章的浏览量会在获取详情时自动增加
4. 只有文章作者和管理员可以修改或删除文章