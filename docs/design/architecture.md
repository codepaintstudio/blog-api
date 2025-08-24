# Blog API 架构设计文档

## 1. 系统架构概述

### 1.1 架构风格

采用经典的三层架构模式：

- **表现层 (Presentation Layer)**：HTTP API 接口
- **业务逻辑层 (Business Logic Layer)**：核心业务处理
- **数据访问层 (Data Access Layer)**：数据库操作

### 1.2 技术栈选型

- **编程语言**：Go (Golang)
- **Web 框架**：Gin
- **数据库**：MySQL 8.0+
- **ORM 框架**：GORM
- **认证**：JWT
- **文件存储**：本地存储/七牛云 OSS
- **缓存**：Redis
- **日志**：Logrus
- **文档**：Swagger

## 2. 项目目录结构

```
blog-api/
├── cmd/                    # 应用程序入口
│   └── server/
│       └── main.go        # 主程序入口
├── internal/              # 内部代码，不对外暴露
│   ├── controllers/       # 控制器层
│   │   ├── auth.go       # 认证相关
│   │   ├── user.go       # 用户管理
│   │   ├── article.go    # 文章管理
│   │   ├── category.go   # 分类管理
│   │   ├── comment.go    # 评论管理
│   │   ├── favorite.go   # 收藏管理
│   │   └── admin.go      # 管理员功能
│   ├── services/         # 业务逻辑层
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── article.go
│   │   ├── category.go
│   │   ├── comment.go
│   │   ├── favorite.go
│   │   └── admin.go
│   ├── repositories/     # 数据访问层
│   │   ├── user.go
│   │   ├── article.go
│   │   ├── category.go
│   │   ├── comment.go
│   │   └── favorite.go
│   ├── models/          # 数据模型
│   │   ├── user.go
│   │   ├── article.go
│   │   ├── category.go
│   │   ├── comment.go
│   │   └── favorite.go
│   ├── middleware/      # 中间件
│   │   ├── auth.go      # 认证中间件
│   │   ├── cors.go      # 跨域处理
│   │   ├── rate_limit.go # 限流中间件
│   │   ├── logging.go   # 日志中间件
│   │   └── recovery.go  # 错误恢复
│   ├── utils/           # 工具函数
│   │   ├── jwt.go       # JWT工具
│   │   ├── password.go  # 密码处理
│   │   ├── upload.go    # 文件上传
│   │   └── validator.go # 数据验证
│   └── config/          # 配置管理
│       └── config.go
├── pkg/                 # 公共代码库
│   ├── response/        # 统一响应格式
│   │   └── response.go
│   ├── errors/          # 错误定义
│   │   └── errors.go
│   └── database/        # 数据库连接
│       └── database.go
├── api/                 # API文档
│   └── swagger/
├── configs/             # 配置文件
│   ├── config.yaml
│   └── config.prod.yaml
├── scripts/             # 脚本文件
│   └── init.sql        # 数据库初始化脚本
├── docs/               # 文档
│   └── design/         # 设计文档
├── uploads/           # 本地文件存储目录
├── go.mod            # Go模块定义
└── README.md         # 项目说明
```

## 3. 数据库设计

### 3.1 数据库选择

- **主数据库**：MySQL 8.0

  - 成熟稳定，支持事务
  - 丰富的生态和工具支持
  - 适合结构化数据存储

- **缓存数据库**：Redis
  - 用于会话存储
  - 接口限流计数
  - 热点数据缓存

### 3.2 表结构设计

#### 3.2.1 用户表 (users)

```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(50),
    avatar VARCHAR(255),
    bio TEXT,
    role ENUM('admin', 'user') DEFAULT 'user',
    status ENUM('active', 'inactive') DEFAULT 'active',
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_username (username),
    INDEX idx_status (status)
);
```

#### 3.2.2 文章表 (articles)

```sql
CREATE TABLE articles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    content LONGTEXT,
    description TEXT,
    cover_image VARCHAR(255),
    category_id BIGINT,
    status ENUM('draft', 'published') DEFAULT 'draft',
    visibility ENUM('public', 'private') DEFAULT 'public',
    view_count INT DEFAULT 0,
    like_count INT DEFAULT 0,
    favorite_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_category_id (category_id),
    INDEX idx_status (status),
    INDEX idx_visibility (visibility),
    INDEX idx_created_at (created_at)
);
```

#### 3.2.3 分类表 (categories)

```sql
CREATE TABLE categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_category (user_id, name),
    INDEX idx_user_id (user_id)
);
```

#### 3.2.4 评论表 (comments)

```sql
CREATE TABLE comments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    article_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    parent_id BIGINT NULL,
    content TEXT NOT NULL,
    like_count INT DEFAULT 0,
    status ENUM('published', 'hidden', 'deleted') DEFAULT 'published',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE,
    INDEX idx_article_id (article_id),
    INDEX idx_user_id (user_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_created_at (created_at)
);
```

#### 3.2.5 评论点赞表 (comment_likes)

```sql
CREATE TABLE comment_likes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    comment_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_comment (user_id, comment_id),
    INDEX idx_comment_id (comment_id),
    INDEX idx_user_id (user_id)
);
```

#### 3.2.6 文章点赞表 (article_likes)

```sql
CREATE TABLE article_likes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    article_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_article (user_id, article_id),
    INDEX idx_article_id (article_id),
    INDEX idx_user_id (user_id)
);
```

#### 3.2.7 收藏夹表 (favorite_folders)

```sql
CREATE TABLE favorite_folders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_folder (user_id, name),
    INDEX idx_user_id (user_id)
);
```

#### 3.2.8 文章收藏表 (article_favorites)

```sql
CREATE TABLE article_favorites (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    article_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    folder_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (folder_id) REFERENCES favorite_folders(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_article_folder (user_id, article_id, folder_id),
    INDEX idx_article_id (article_id),
    INDEX idx_user_id (user_id),
    INDEX idx_folder_id (folder_id),
    INDEX idx_created_at (created_at)
);
```

#### 3.2.9 文件管理表 (files)
```sql
CREATE TABLE files (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    filename VARCHAR(255) NOT NULL COMMENT '原始文件名',
    hash VARCHAR(64) UNIQUE NOT NULL COMMENT '文件哈希值(SHA256)',
    path VARCHAR(500) NOT NULL COMMENT '文件存储路径',
    url VARCHAR(500) NOT NULL COMMENT '访问 URL',
    mime_type VARCHAR(100) NOT NULL COMMENT '文件 MIME 类型',
    size BIGINT NOT NULL COMMENT '文件大小(字节)',
    upload_user_id BIGINT NOT NULL COMMENT '上传用户ID',
    reference_count INT DEFAULT 1 COMMENT '引用计数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (upload_user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE INDEX uk_hash (hash),
    INDEX idx_upload_user_id (upload_user_id),
    INDEX idx_mime_type (mime_type),
    INDEX idx_created_at (created_at)
);
```

## 4. API 设计

### 4.1 API 版本控制

- 使用 URL 路径版本控制：`/api/v1/`
- 向后兼容原则
- 版本废弃通知机制

### 4.2 统一响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 4.3 API 路由设计

#### 4.3.1 认证相关

```
POST   /api/v1/auth/register     # 用户注册
POST   /api/v1/auth/login        # 用户登录
POST   /api/v1/auth/logout       # 用户登出
POST   /api/v1/auth/refresh      # 刷新Token
```

#### 4.3.2 用户管理

```
GET    /api/v1/user/profile      # 获取用户信息
PUT    /api/v1/user/profile      # 更新用户信息
POST   /api/v1/user/avatar       # 上传头像
PUT    /api/v1/user/password     # 修改密码
```

#### 4.3.3 文章管理

```
GET    /api/v1/articles          # 获取公开文章列表（支持分页、分类筛选、排序）
POST   /api/v1/articles          # 创建文章
GET    /api/v1/articles/:id      # 获取文章详情（包含浏览量统计）
PUT    /api/v1/articles/:id      # 更新文章
DELETE /api/v1/articles/:id      # 删除文章
GET    /api/v1/articles/my       # 获取我的文章
GET    /api/v1/articles/search   # 文章搜索
```

#### 4.3.4 分类管理

```
GET    /api/v1/categories        # 获取分类列表
POST   /api/v1/categories        # 创建分类
PUT    /api/v1/categories/:id    # 更新分类
DELETE /api/v1/categories/:id    # 删除分类
GET    /api/v1/categories/:id/articles # 获取分类下的文章列表
```

#### 4.3.5 评论管理

```
GET    /api/v1/articles/:id/comments     # 获取文章评论列表
POST   /api/v1/articles/:id/comments     # 发表评论
POST   /api/v1/comments/:id/reply        # 回复评论
DELETE /api/v1/comments/:id             # 删除评论
POST   /api/v1/comments/:id/like         # 点赞评论
DELETE /api/v1/comments/:id/like        # 取消点赞
GET    /api/v1/comments/:id/replies      # 获取评论回复列表
```

#### 4.3.6 文章互动功能

```
POST   /api/v1/articles/:id/like         # 点赞文章
DELETE /api/v1/articles/:id/like        # 取消点赞
GET    /api/v1/articles/:id/likes        # 获取文章点赞列表
POST   /api/v1/articles/:id/favorite     # 收藏文章
DELETE /api/v1/articles/:id/favorite    # 取消收藏
```

#### 4.3.7 收藏管理

```
GET    /api/v1/favorites/folders         # 获取收藏夹列表
POST   /api/v1/favorites/folders         # 创建收藏夹
PUT    /api/v1/favorites/folders/:id     # 更新收藏夹
DELETE /api/v1/favorites/folders/:id    # 删除收藏夹
GET    /api/v1/favorites/articles        # 获取我的收藏文章
GET    /api/v1/favorites/folders/:id/articles # 获取收藏夹下的文章
```

#### 4.3.8 个人中心

```
GET    /api/v1/user/likes               # 获取我的点赞记录
GET    /api/v1/user/favorites           # 获取我的收藏记录
GET    /api/v1/user/stats               # 获取个人统计数据
```

#### 4.3.9 文件管理
```
POST   /api/v1/files/upload             # 上传文件（支持去重）
GET    /api/v1/files/:id                # 获取文件信息
DELETE /api/v1/files/:id             # 删除文件（减少引用计数）
```

#### 4.3.10 管理员功能

```
GET    /api/v1/admin/users       # 获取用户列表
PUT    /api/v1/admin/users/:id   # 管理用户状态
GET    /api/v1/admin/articles    # 获取所有文章
DELETE /api/v1/admin/articles/:id # 删除文章
GET    /api/v1/admin/comments    # 获取所有评论
PUT    /api/v1/admin/comments/:id # 管理评论状态
GET    /api/v1/admin/stats       # 系统统计
```

## 5. 安全与防护机制

### 5.1 认证与授权

- **JWT Token**：

  - Access Token (1 小时有效期)
  - Refresh Token (7 天有效期)
  - Token 黑名单机制

- **权限控制**：
  - RBAC 角色权限控制
  - 资源所有权验证
  - 管理员特权管理

### 5.2 接口防护

- **限流机制**：

  - 使用令牌桶算法
  - IP 级别限流：100 次/分钟
  - 用户级别限流：1000 次/小时
  - Redis 存储限流计数

- **防刷机制**：
  - 滑动窗口算法
  - 验证码机制（连续失败后）
  - 临时封禁机制

### 5.3 数据安全

- **密码安全**：

  - bcrypt 加密存储
  - 密码强度要求
  - 密码历史检查

- **数据验证**：
  - 输入参数验证
  - SQL 注入防护
  - XSS 攻击防护

## 6. 中间件设计

### 6.1 认证中间件

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Token验证逻辑
        // 用户权限检查
        // 设置用户上下文
    }
}
```

### 6.2 限流中间件

```go
func RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // IP限流检查
        // 用户限流检查
        // 计数器更新
    }
}
```

### 6.3 CORS 中间件

```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 跨域策略设置
        // 预检请求处理
    }
}
```

## 7. 配置管理

### 7.1 配置文件结构

```yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 3306
  username: root
  password: password
  dbname: blog_api

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-secret-key
  expire: 3600

upload:
  path: ./uploads
  max_size: 10485760 # 10MB

rate_limit:
  ip_limit: 100
  user_limit: 1000
```

## 8. 部署架构

### 8.1 开发环境

- Docker Compose 部署
- 本地数据库和 Redis
- 热重载开发

### 8.2 生产环境

- 容器化部署
- 数据库集群
- Redis 集群
- 负载均衡
- 日志收集

## 9. 监控与日志

### 9.1 日志策略

- **结构化日志**：JSON 格式
- **日志级别**：DEBUG/INFO/WARN/ERROR
- **日志轮转**：按大小和时间
- **敏感信息**：脱敏处理

### 9.2 监控指标

- API 响应时间
- 错误率统计
- 数据库连接池
- 内存和 CPU 使用率

## 10. 开源组件选型

### 10.1 核心依赖

- **gin-gonic/gin**：Web 框架
- **gorm.io/gorm**：ORM 框架
- **go-redis/redis**：Redis 客户端
- **golang-jwt/jwt**：JWT 处理
- **bcrypt**：密码加密

### 10.2 防护组件

- **gin-rate-limit**：限流中间件
- **validator/v10**：参数验证
- **cors**：跨域处理
- **secure**：安全头设置

## 11. 性能优化策略

### 11.1 数据库优化

- 合理的索引设计
- 查询优化
- 连接池配置
- 读写分离

### 11.2 缓存策略

- Redis 缓存热点数据
- 应用层缓存
- HTTP 缓存头设置

### 11.3 并发处理

- Goroutine 池
- 数据库连接池
- 合理的超时设置
