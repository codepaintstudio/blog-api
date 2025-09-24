# 管理员接口

管理员接口用于系统管理和内容审核功能。

## 获取用户列表

### 请求信息

- **接口地址**: `/admin/users`
- **请求方式**: `GET`
- **权限要求**: 管理员
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
| status | string | 否 | 用户状态筛选，如 active, inactive |
| role | string | 否 | 用户角色筛选，如 user, admin |
| search | string | 否 | 搜索关键词，匹配用户名或邮箱 |

### 示例请求

```bash
curl -X GET "http://localhost:8080/api/v1/admin/users?page=1&per_page=10&status=active" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取用户列表成功",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "admin_user",
        "email": "admin@example.com",
        "nickname": "管理员",
        "avatar": "http://localhost:8080/uploads/avatars/admin.jpg",
        "bio": "系统管理员",
        "role": "admin",
        "status": "active",
        "email_verified": true,
        "article_count": 50,
        "comment_count": 100,
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

## 获取用户详情

### 请求信息

- **接口地址**: `/admin/users/{id}`
- **请求方式**: `GET`
- **权限要求**: 管理员
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 用户ID |

### 示例请求

```bash
curl -X GET http://localhost:8080/api/v1/admin/users/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取用户详情成功",
  "data": {
    "id": 1,
    "username": "user",
    "email": "user@example.com",
    "nickname": "普通用户",
    "avatar": "http://localhost:8080/uploads/avatars/user.jpg",
    "bio": "普通用户",
    "role": "user",
    "status": "active",
    "email_verified": true,
    "last_login_at": "2024-01-01T00:00:00Z",
    "last_login_ip": "192.168.1.1",
    "article_count": 10,
    "comment_count": 25,
    "like_count": 50,
    "favorite_count": 15,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 更新用户状态

### 请求信息

- **接口地址**: `/admin/users/{id}/status`
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
| id | int | 是 | 用户ID |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| status | string | 是 | 用户状态，active 或 inactive |

### 示例请求

```bash
curl -X PUT http://localhost:8080/api/v1/admin/users/1/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "status": "inactive"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "用户状态更新成功",
  "data": {
    "id": 1,
    "username": "user",
    "status": "inactive"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 更新用户角色

### 请求信息

- **接口地址**: `/admin/users/{id}/role`
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
| id | int | 是 | 用户ID |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| role | string | 是 | 用户角色，user 或 admin |

### 示例请求

```bash
curl -X PUT http://localhost:8080/api/v1/admin/users/1/role \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "role": "admin"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "用户角色更新成功",
  "data": {
    "id": 1,
    "username": "user",
    "role": "admin"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 删除用户

### 请求信息

- **接口地址**: `/admin/users/{id}`
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
| id | int | 是 | 用户ID |

### 示例请求

```bash
curl -X DELETE http://localhost:8080/api/v1/admin/users/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "用户删除成功",
  "data": null,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 获取系统统计

### 请求信息

- **接口地址**: `/admin/stats/system`
- **请求方式**: `GET`
- **权限要求**: 管理员
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 示例请求

```bash
curl -X GET http://localhost:8080/api/v1/admin/stats/system \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取系统统计成功",
  "data": {
    "system": {
      "status": "running",
      "version": "1.0.0",
      "uptime": "24h 10m 30s",
      "go_version": "go1.24.5",
      "golang_version": "1.24.5",
      "goroutines": 15,
      "memory_usage": "123.4 MB"
    },
    "total_users": 1000,
    "total_articles": 5000,
    "total_comments": 10000,
    "total_likes": 50000,
    "total_favorites": 20000,
    "today_registrations": 10,
    "today_articles": 50,
    "today_comments": 200
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 获取用户统计

### 请求信息

- **接口地址**: `/admin/stats/users`
- **请求方式**: `GET`
- **权限要求**: 管理员
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| start_date | string | 否 | 开始日期，格式：YYYY-MM-DD |
| end_date | string | 否 | 结束日期，格式：YYYY-MM-DD |

### 示例请求

```bash
curl -X GET "http://localhost:8080/api/v1/admin/stats/users?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取用户统计成功",
  "data": {
    "registration_trend": [
      {
        "date": "2024-01-01",
        "count": 5
      },
      {
        "date": "2024-01-02",
        "count": 8
      }
    ],
    "user_growth": {
      "total_users": 1000,
      "new_users_today": 10,
      "new_users_week": 50,
      "new_users_month": 200
    },
    "active_users": {
      "today": 150,
      "week": 400,
      "month": 800
    }
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 获取内容统计

### 请求信息

- **接口地址**: `/admin/stats/content`
- **请求方式**: `GET`
- **权限要求**: 管理员
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| start_date | string | 否 | 开始日期，格式：YYYY-MM-DD |
| end_date | string | 否 | 结束日期，格式：YYYY-MM-DD |

### 示例请求

```bash
curl -X GET "http://localhost:8080/api/v1/admin/stats/content?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取内容统计成功",
  "data": {
    "content_trend": [
      {
        "date": "2024-01-01",
        "articles": 10,
        "comments": 50
      }
    ],
    "content_stats": {
      "total_articles": 5000,
      "total_comments": 10000,
      "new_articles_today": 20,
      "new_comments_today": 100
    },
    "popular_content": {
      "top_articles": [
        {
          "id": 1,
          "title": "热门文章标题",
          "view_count": 1000,
          "like_count": 100,
          "comment_count": 50
        }
      ],
      "top_commented_articles": [
        {
          "id": 2,
          "title": "评论最多的文章",
          "comment_count": 200
        }
      ]
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
| 404 | 用户不存在 |
| 429 | 请求过于频繁 |

## 注意事项

1. 只有管理员角色的用户才能访问管理员接口
2. 在禁用用户或删除用户时要谨慎操作，这将影响用户的正常使用
3. 更新用户角色是一个敏感操作，需要确认操作的必要性
4. 系统统计数据提供了重要的运营指标，需要定期关注
5. 删除用户会删除其创建的所有内容，这是一个不可逆的操作