# 文件接口

文件接口用于文件上传和管理功能。

## 文件上传

### 请求信息

- **接口地址**: `/files/upload`
- **请求方式**: `POST`
- **Content-Type**: `multipart/form-data`
- **权限要求**: 需要登录
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 表单参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| file | file | 是 | 要上传的文件 |
| type | string | 否 | 文件类型，如 image, document 等，默认根据文件扩展名自动识别 |

### 示例请求

```bash
curl -X POST http://localhost:8080/api/v1/files/upload \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -F "file=@/path/to/your/file.jpg" \
  -F "type=image"
```

### 成功响应

```json
{
  "code": 200,
  "message": "文件上传成功",
  "data": {
    "id": 1,
    "filename": "file.jpg",
    "original_name": "original_filename.jpg",
    "hash": "a1b2c3d4e5f6...",
    "path": "./uploads/images/2024/01/01/a1b2c3d4e5f6.jpg",
    "url": "http://localhost:8080/uploads/images/2024/01/01/a1b2c3d4e5f6.jpg",
    "mime_type": "image/jpeg",
    "size": 102400,
    "upload_user_id": 1,
    "reference_count": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 获取用户文件列表

### 请求信息

- **接口地址**: `/files`
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
| type | string | 否 | 文件类型筛选，如 image, document 等 |

### 示例请求

```bash
curl -X GET "http://localhost:8080/api/v1/files?page=1&per_page=10&type=image" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取文件列表成功",
  "data": {
    "files": [
      {
        "id": 1,
        "filename": "file1.jpg",
        "original_name": "original1.jpg",
        "hash": "a1b2c3d4e5f6...",
        "path": "./uploads/images/2024/01/01/a1b2c3d4e5f6.jpg",
        "url": "http://localhost:8080/uploads/images/2024/01/01/a1b2c3d4e5f6.jpg",
        "mime_type": "image/jpeg",
        "size": 102400,
        "reference_count": 3,
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

## 获取文件信息

### 请求信息

- **接口地址**: `/files/{id}`
- **请求方式**: `GET`
- **权限要求**: 需要登录且为文件上传者
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 文件ID |

### 示例请求

```bash
curl -X GET http://localhost:8080/api/v1/files/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "获取文件信息成功",
  "data": {
    "id": 1,
    "filename": "file.jpg",
    "original_name": "original_filename.jpg",
    "hash": "a1b2c3d4e5f6...",
    "path": "./uploads/images/2024/01/01/a1b2c3d4e5f6.jpg",
    "url": "http://localhost:8080/uploads/images/2024/01/01/a1b2c3d4e5f6.jpg",
    "mime_type": "image/jpeg",
    "size": 102400,
    "upload_user": {
      "id": 1,
      "username": "uploader",
      "nickname": "上传者"
    },
    "reference_count": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 删除文件

### 请求信息

- **接口地址**: `/files/{id}`
- **请求方式**: `DELETE`
- **权限要求**: 需要登录且为文件上传者
- **认证方式**: Bearer Token

### 请求头

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Authorization | string | 是 | Bearer Token，格式："Bearer {token}" |

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 文件ID |

### 示例请求

```bash
curl -X DELETE http://localhost:8080/api/v1/files/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 成功响应

```json
{
  "code": 200,
  "message": "文件删除成功",
  "data": null,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 响应码说明

| HTTP 状态码 | 说明 |
|-------------|------|
| 200 | 操作成功 |
| 400 | 请求参数错误或文件类型不支持 |
| 401 | 未认证或 Token 过期 |
| 403 | 权限不足（非文件上传者） |
| 404 | 文件不存在 |
| 413 | 文件大小超过限制（10MB） |
| 429 | 请求过于频繁 |

## 支持的文件类型

### 图片类型
- JPG/JPEG: image/jpeg
- PNG: image/png
- GIF: image/gif
- WebP: image/webp
- SVG: image/svg+xml

### 文档类型
- PDF: application/pdf
- DOC/DOCX: application/msword, application/vnd.openxmlformats-officedocument.wordprocessingml.document
- XLS/XLSX: application/vnd.ms-excel, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
- PPT/PPTX: application/vnd.ms-powerpoint, application/vnd.openxmlformats-officedocument.presentationml.presentation

### 其他类型
- TXT: text/plain
- JSON: application/json

## 文件限制

1. **文件大小**: 单个文件最大 10MB
2. **文件类型**: 仅支持上述列出的安全文件类型
3. **文件去重**: 系统会根据文件哈希值检查重复文件，如果文件已存在，将增加引用计数而不是重新上传
4. **引用计数**: 每次文件被使用时引用计数加1，删除时减1，当引用计数为0时文件将被物理删除

## 注意事项

1. 上传的文件会根据类型和日期自动存储到相应的子目录
2. 系统会自动对上传的图片进行安全检查
3. 用户只能管理和删除自己上传的文件
4. 删除文件时，如果还有其他地方在使用该文件（引用计数大于1），则只减少引用计数而不物理删除文件