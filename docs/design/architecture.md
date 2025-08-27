# Blog API 架构设计

## 系统架构

采用经典的三层架构模式：

- **控制器层 (Controller)**：HTTP 请求处理和响应
- **服务层 (Service)**：业务逻辑处理
- **数据层 (Repository)**：数据库操作

## 技术选型

### 核心技术

- **语言**: Go 1.19+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis 6.0
- **认证**: JWT
- **文档**: Swagger

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
├── internal/              # 内部代码
│   ├── controllers/       # 控制器层
│   ├── services/          # 业务逻辑层
│   ├── repositories/      # 数据访问层
│   ├── models/           # 数据模型
│   ├── middlewares/      # 中间件
│   ├── routes/           # 路由配置
│   ├── utils/            # 工具函数
│   └── container/        # 依赖注入
├── pkg/                  # 公共代码
│   ├── config/           # 配置管理
│   ├── database/         # 数据库连接
│   ├── logger/           # 日志系统
│   └── response/         # 响应格式
├── configs/              # 配置文件
├── docs/                 # 文档和Swagger
├── uploads/              # 文件存储
└── logs/                 # 日志文件
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

---

**项目已完成所有核心架构设计并投入使用。**

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
