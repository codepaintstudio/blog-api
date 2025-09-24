# 项目架构

## 系统架构

Blog API 采用经典的三层架构模式，各层职责明确，便于维护和扩展。

```
+------------------+
|   控制器层       |  HTTP 请求处理和响应
+------------------+
|   服务层         |  业务逻辑处理
+------------------+
|   数据层         |  数据库操作
+------------------+
```

- **控制器层 (Controller)**：负责处理 HTTP 请求和响应，是外部世界与业务逻辑的桥梁
- **服务层 (Service)**：封装核心业务逻辑，处理复杂的业务规则和流程
- **数据层 (Repository)**：封装数据库操作，提供数据访问接口

## 技术选型

### 核心技术

- **语言**: Go 1.24.5+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis 6.0
- **认证**: JWT

### 关键组件

- **配置管理**: YAML 格式
- **日志系统**: 结构化日志
- **依赖注入**: 自定义容器
- **中间件**: 认证、限流、CORS、日志

## 项目结构

```
blog-api/
├── cmd/
│   └── server/main.go       # 程序入口
├── internal/                # 内部代码
│   ├── controllers/         # 控制器层
│   ├── services/            # 业务逻辑层
│   ├── repositories/        # 数据访问层
│   ├── models/              # 数据模型
│   ├── middlewares/         # 中间件
│   ├── routes/              # 路由配置
│   ├── utils/               # 工具函数
│   └── container/           # 依赖注入
├── pkg/                     # 公共代码
│   ├── config/              # 配置管理
│   ├── database/            # 数据库连接
│   ├── logger/              # 日志系统
│   └── response/            # 响应格式
├── configs/                 # 配置文件
├── docs/                    # 文档
├── uploads/                 # 文件存储
└── logs/                    # 日志文件
```

## 数据库设计

### 核心表结构

**users 用户表**

- 基本信息：id, username, email, password, nickname
- 扩展信息：avatar, bio, role, status, email_verified
- 时间戳：created_at, updated_at
- 索引：email, username

**articles 文章表**

- 内容信息：id, title, content, description, cover_image
- 关联信息：user_id, category_id
- 状态信息：status(draft/published), visibility(public/private)
- 统计信息：view_count, like_count, favorite_count

**categories 分类表**

- 基本信息：id, user_id, name, description
- 约束：用户内分类名唯一

**comments 评论表**

- 内容信息：id, content, like_count
- 关联信息：article_id, user_id, parent_id(支持多级回复)
- 状态信息：status(published/hidden/deleted)

### 互动功能表

- **article_likes**: 文章点赞记录
- **comment_likes**: 评论点赞记录
- **favorite_folders**: 收藏夹
- **article_favorites**: 文章收藏记录
- **files**: 文件管理

## 接口设计

### API 风格

- RESTful 风格
- JSON 数据格式
- 统一响应结构
- HTTP 状态码规范

### 认证机制

- JWT Bearer Token
- Token 过期时间：1 小时
- Refresh Token：7 天
- 头部格式：`Authorization: Bearer <token>`

### 响应格式

```json
{
  "code": 200,
  "message": "成功",
  "data": {},
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 安全设计

### 认证与授权

- JWT Token 认证
- 基于角色的权限控制
- 管理员权限中间件
- 资源所有权校验

### 防护机制

- 接口限流（令牌桶算法）
- 参数验证和数据校验
- SQL 注入防护（GORM 预处理）
- XSS 防护（参数转义）
- CORS 跨域配置

### 数据安全

- 密码 bcrypt 加密
- 敏感信息过滤
- 文件上传类型限制
- 文件大小限制（10MB）

## 性能优化

### 缓存策略

- Redis 缓存热点数据
- 数据库连接池
- 限流计数器缓存

### 数据库优化

- 合理的索引设计
- 分页查询优化
- 预编译语句（GORM）
- 数据库连接复用

### 并发处理

- Gin 框架原生 Goroutine 支持
- 数据库连接池配置
- Redis 连接池管理

这套架构设计确保了系统的可扩展性、可维护性和高性能，为前端开发者提供了稳定可靠的 API 服务。
