# 分类接口

分类接口用于管理文章分类。

## 获取分类列表

### 请求信息

- **接口地址**: `/categories`
- **请求方式**: `GET`
- **权限要求**: 公开接口
- **认证方式**: 无需认证

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认为 1 |
| per_page | int | 否 | 每页数量，默认为 10，最大为 100 |
| status | string | 否 | 分类状态，active 或 inactive，默认为 active |

### 示例请求

```bash
curl -X GET "http://localhost:8080/api/v1/categories?page=1&per_page=10"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取分类列表成功",
  "data": {
    "categories": [
      {
        "id": 1,
        "name": "技术分享",
        "description": "技术相关的文章分享",
        "article_count": 50,
        "status": "active",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
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

## 获取活跃分类列表

### 请求信息

- **接口地址**: `/categories/active`
- **请求方式**: `GET`
- **权限要求**: 公开接口
- **认证方式**: 无需认证

### 示例请求

```bash
curl -X GET http://localhost:8080/api/v1/categories/active
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取活跃分类列表成功",
  "data": [
    {
      "id": 1,
      "name": "技术分享",
      "description": "技术相关的文章分享",
      "article_count": 50
    },
    {
      "id": 2,
      "name": "生活随笔",
      "description": "生活中的感悟和思考",
      "article_count": 30
    }
  ],
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 获取分类详情

### 请求信息

- **接口地址**: `/categories/{id}`
- **请求方式**: `GET`
- **权限要求**: 公开接口
- **认证方式**: 无需认证

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 分类 ID |

### 示例请求

```bash
curl -X GET http://localhost:8080/api/v1/categories/1
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取分类详情成功",
  "data": {
    "id": 1,
    "name": "技术分享",
    "description": "技术相关的文章分享",
    "article_count": 50,
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 创建分类

### 请求信息

- **接口地址**: `/categories`
- **请求方式**: `POST`
- **Content-Type**: `application/json`
- **权限要求**: 管理员
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 分类名称，最多 50 个字符 |
| description | string | 否 | 分类描述，最多 200 个字符 |
| status | string | 否 | 分类状态，active 或 inactive，默认为 active |

### 示例请求

```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "前端开发",
    "description": "前端技术相关的文章",
    "status": "active"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "分类创建成功",
  "data": {
    "id": 3,
    "name": "前端开发",
    "description": "前端技术相关的文章",
    "article_count": 0,
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 更新分类

### 请求信息

- **接口地址**: `/categories/{id}`
- **请求方式**: `PUT`
- **Content-Type**: `application/json`
- **权限要求**: 管理员
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 分类 ID |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 否 | 分类名称，最多 50 个字符 |
| description | string | 否 | 分类描述，最多 200 个字符 |
| status | string | 否 | 分类状态，active 或 inactive |

### 示例请求

```bash
curl -X PUT http://localhost:8080/api/v1/categories/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "更新的分类名称",
    "description": "更新的分类描述"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "分类更新成功",
  "data": {
    "id": 1,
    "name": "更新的分类名称",
    "description": "更新的分类描述",
    "article_count": 50,
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  },
  "timestamp": "2024-01-02T00:00:00Z"
}
```

## 删除分类

### 请求信息

- **接口地址**: `/categories/{id}`
- **请求方式**: `DELETE`
- **权限要求**: 管理员
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 分类 ID |

### 示例请求

```bash
curl -X DELETE http://localhost:8080/api/v1/categories/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "分类删除成功",
  "data": null,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 获取分类下的文章

### 请求信息

- **接口地址**: `/categories/{id}/articles`
- **请求方式**: `GET`
- **权限要求**: 公开接口
- **认证方式**: 无需认证

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 分类 ID |

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认为 1 |
| per_page | int | 否 | 每页数量，默认为 10，最大为 100 |

### 示例请求

```bash
curl -X GET "http://localhost:8080/api/v1/categories/1/articles?page=1&per_page=10"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取分类下文章成功",
  "data": {
    "articles": [
      {
        "id": 1,
        "title": "Go语言入门指南",
        "description": "Go语言的基础知识介绍",
        "cover_image": "http://localhost:8080/uploads/covers/go-guide.jpg",
        "view_count": 200,
        "like_count": 15,
        "comment_count": 8,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "per_page": 10,
      "total": 50,
      "total_pages": 5
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
| 403 | 权限不足（非管理员） |
| 404 | 分类不存在 |
| 409 | 分类名称已存在 |
| 429 | 请求过于频繁 |

## 注意事项

1. 只有管理员可以创建、更新和删除分类
2. 分类名称在系统中必须唯一
3. 删除分类时，该分类下的文章不会被删除，但分类 ID 会被设为 null
4. 通过 `/categories/active` 接口可以快速获取所有活跃分类，适用于前端下拉选择等场景