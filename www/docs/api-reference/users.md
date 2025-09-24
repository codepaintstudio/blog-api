# 用户接口

用户接口用于获取和更新用户信息。

## 获取用户信息

### 请求信息

- **接口地址**: `/users/profile`
- **请求方式**: `GET`
- **权限要求**: 需要登录
- **认证方式**: Bearer Token

### 请求头

| 参数名        | 类型   | 必填 | 说明                                 |
| ------------- | ------ | ---- | ------------------------------------ |
| Authorization | string | 是   | Bearer Token，格式："Bearer {token}" |

### 示例请求

```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取用户信息成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "avatar": "http://localhost:8080/uploads/avatars/test.jpg",
    "bio": "这是一个个人简介",
    "role": "user",
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  },
  "timestamp": "2024-01-02T00:00:00Z"
}
```

## 更新用户信息

### 请求信息

- **接口地址**: `/users/profile`
- **请求方式**: `PUT`
- **Content-Type**: `application/json`
- **权限要求**: 需要登录
- **认证方式**: Bearer Token

### 请求头

| 参数名        | 类型   | 必填 | 说明                                 |
| ------------- | ------ | ---- | ------------------------------------ |
| Authorization | string | 是   | Bearer Token，格式："Bearer {token}" |

### 请求参数

| 参数名   | 类型   | 必填 | 说明                      |
| -------- | ------ | ---- | ------------------------- |
| nickname | string | 否   | 昵称，最多 50 个字符      |
| bio      | string | 否   | 个人简介，最多 200 个字符 |

### 示例请求

```bash
curl -X PUT http://localhost:8080/api/v1/users/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "nickname": "新昵称",
    "bio": "更新后的个人简介"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "用户信息更新成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "新昵称",
    "avatar": "http://localhost:8080/uploads/avatars/test.jpg",
    "bio": "更新后的个人简介",
    "role": "user",
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  },
  "timestamp": "2024-01-02T00:00:00Z"
}
```

## 修改密码

### 请求信息

- **接口地址**: `/users/password`
- **请求方式**: `PUT`
- **Content-Type**: `application/json`
- **权限要求**: 需要登录
- **认证方式**: Bearer Token

### 请求头

| 参数名        | 类型   | 必填 | 说明                                 |
| ------------- | ------ | ---- | ------------------------------------ |
| Authorization | string | 是   | Bearer Token，格式："Bearer {token}" |

### 请求参数

| 参数名       | 类型   | 必填 | 说明              |
| ------------ | ------ | ---- | ----------------- |
| old_password | string | 是   | 旧密码            |
| new_password | string | 是   | 新密码，至少 6 位 |

### 示例请求

```bash
curl -X PUT http://localhost:8080/api/v1/users/password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "old_password": "old_password123",
    "new_password": "new_password123"
  }'
```

### 成功响应

```json
{
  "code": 200,
  "message": "密码修改成功",
  "data": null,
  "timestamp": "2024-01-02T00:00:00Z"
}
```

## 响应码说明

| HTTP 状态码 | 说明                       |
| ----------- | -------------------------- |
| 200         | 操作成功                   |
| 400         | 请求参数错误               |
| 401         | 未认证或 Token 过期        |
| 403         | 权限不足                   |
| 422         | 密码验证失败（修改密码时） |
| 429         | 请求过于频繁               |

## 注意事项

1. 更新用户信息时，所有字段都是可选的，只提供需要更新的字段
2. 修改密码时需要验证旧密码的正确性
3. 用户名和邮箱一旦注册后无法修改
4. 头像上传需要通过文件上传接口完成
