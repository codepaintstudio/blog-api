# 认证接口

认证接口用于用户注册、登录和 Token 刷新等功能。

## 用户注册

### 请求信息

- **接口地址**: `/auth/register`
- **请求方式**: `POST`
- **Content-Type**: `application/json`
- **权限要求**: 公开接口

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名，3-20个字符 |
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码，至少6位 |

### 示例请求

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "nickname": "testuser",
      "avatar": "",
      "bio": "",
      "role": "user",
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 用户登录

### 请求信息

- **接口地址**: `/auth/login`
- **请求方式**: `POST`
- **Content-Type**: `application/json`
- **权限要求**: 公开接口

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名或邮箱 |
| password | string | 是 | 密码 |

### 示例请求

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "nickname": "testuser",
      "avatar": "",
      "bio": "",
      "role": "user",
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Token 刷新

### 请求信息

- **接口地址**: `/auth/refresh`
- **请求方式**: `POST`
- **Content-Type**: `application/json`
- **权限要求**: 公开接口

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| refresh_token | string | 是 | 刷新令牌 |

### 示例请求

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your_refresh_token_here"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "Token 刷新成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "new_refresh_token_here",
    "expires_in": 3600
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 响应码说明

| HTTP 状态码 | 说明 |
|-------------|------|
| 200 | 操作成功 |
| 400 | 请求参数错误 |
| 401 | 用户名或密码错误 |
| 409 | 用户名或邮箱已存在（注册时） |
| 429 | 请求过于频繁 |

## 注意事项

1. 密码在传输过程中应使用 HTTPS 加密
2. 接收到的 Token 需要妥善保存，避免泄露
3. Access Token 有效期为 1 小时，过期后需要使用 Refresh Token 刷新
4. Refresh Token 有效期为 7 天，过期后需要重新登录