# API 参考

本节提供了 Blog API 的详细接口文档，包括请求方式、参数说明和响应格式。

## API 基本信息

- **基础 URL**: `http://localhost:8080/api/v1/` (开发环境)
- **生产环境 URL**: `https://blog-api.hub.feashow.cn/api/v1/` (生产环境)
- **数据格式**: JSON
- **认证方式**: JWT Bearer Token
- **版本**: v1

## 通用响应格式

所有 API 响应都遵循以下格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 通用错误码

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或 Token 过期 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 429 | 请求过于频繁（限流） |
| 500 | 服务器内部错误 |

## 接口分类

- [认证接口](./auth)
- [用户接口](./users)
- [文章接口](./articles)
- [分类接口](./categories)
- [评论接口](./comments)
- [点赞接口](./likes)
- [收藏接口](./favorites)
- [文件接口](./files)
- [管理员接口](./admin)

在使用具体接口前，请先阅读 [入门指南](../getting-started) 了解基本概念和认证方式。