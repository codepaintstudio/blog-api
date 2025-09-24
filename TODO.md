# Blog API 开发计划

## 项目状态

当前项目已完成所有核心功能开发，包括 42 个 API 接口，涵盖 10 个功能模块。系统已可正常运行，API 文档已生成。

## 已完成功能 ✅

### 基础设施

- [x] 项目初始化和配置管理
- [x] 数据库连接（MySQL + Redis）
- [x] 三层架构设计（Controller-Service-Repository）
- [x] 依赖注入容器
- [x] 统一响应格式和错误处理
- [x] 结构化日志系统

### 中间件系统

- [x] JWT 认证中间件
- [x] 权限控制中间件
- [x] 接口限流中间件（令牌桶算法）
- [x] CORS 跨域处理
- [x] 请求日志中间件
- [x] 错误恢复中间件

### 用户认证模块

- [x] 用户注册接口
- [x] 用户登录接口
- [x] JWT Token 刷新
- [x] 密码加密（bcrypt）
- [x] 用户资料管理
- [x] 密码修改功能

### 内容管理模块

- [x] 文章 CRUD 操作
- [x] 文章列表和分页
- [x] 文章搜索功能
- [x] 分类管理系统
- [x] 文件上传和管理
- [x] 文件去重机制（SHA256）

### 互动功能模块

- [x] 文章点赞系统
- [x] 文章收藏功能
- [x] 收藏夹管理
- [x] 多级评论系统（3 级嵌套）
- [x] 评论点赞功能

### 系统管理模块

- [x] 管理员自动初始化
- [x] 用户管理接口
- [x] 系统统计功能
- [x] 内容审核管理

### 文档和部署


- [x] 完整的接口文档（42 个接口）
- [x] 项目部署和运行

## 技术特性

### 安全防护

- [x] 接口限流保护（IP 级别和用户级别）
- [x] 参数验证和数据校验
- [x] JWT Token 认证
- [x] 密码安全存储
- [x] CORS 安全配置

### 性能优化

- [x] Redis 缓存支持
- [x] 数据库连接池
- [x] 查询优化和索引
- [x] 文件上传优化

### 系统特性

- [x] 配置文件管理（YAML）
- [x] 环境变量支持
- [x] 日志轮转和分级
- [x] 优雅关闭机制

## 接口列表

### 认证接口（3 个）

- POST /api/v1/auth/register - 用户注册
- POST /api/v1/auth/login - 用户登录
- POST /api/v1/auth/refresh - Token 刷新

### 用户接口（3 个）

- GET /api/v1/users/profile - 获取用户资料
- PUT /api/v1/users/profile - 更新用户资料
- PUT /api/v1/users/password - 修改密码

### 文章接口（6 个）

- GET /api/v1/articles - 文章列表
- GET /api/v1/articles/search - 文章搜索
- GET /api/v1/articles/:id - 文章详情
- POST /api/v1/articles - 创建文章
- PUT /api/v1/articles/:id - 更新文章
- DELETE /api/v1/articles/:id - 删除文章

### 分类接口（6 个）

- GET /api/v1/categories - 分类列表
- GET /api/v1/categories/active - 活跃分类
- GET /api/v1/categories/:id - 分类详情
- POST /api/v1/categories - 创建分类
- PUT /api/v1/categories/:id - 更新分类
- DELETE /api/v1/categories/:id - 删除分类

### 文件接口（4 个）

- POST /api/v1/files/upload - 文件上传
- GET /api/v1/files - 用户文件列表
- GET /api/v1/files/:id - 文件详情
- DELETE /api/v1/files/:id - 删除文件

### 评论接口（8 个）

- GET /api/v1/comments/:id - 评论详情
- GET /api/v1/comments/:id/replies - 评论回复
- POST /api/v1/comments - 发表评论
- PUT /api/v1/comments/:id - 更新评论
- DELETE /api/v1/comments/:id - 删除评论
- GET /api/v1/comments/mine - 我的评论
- POST /api/v1/comments/:id/approve - 审核通过
- POST /api/v1/comments/:id/reject - 审核拒绝

### 点赞接口（4 个）

- GET /api/v1/likes/articles - 我点赞的文章
- GET /api/v1/likes/comments - 我点赞的评论
- POST /api/v1/articles/:id/like - 文章点赞
- POST /api/v1/comments/:id/like - 评论点赞

### 收藏接口（8 个）

- POST /api/v1/favorites/folders - 创建收藏夹
- GET /api/v1/favorites/folders - 收藏夹列表
- PUT /api/v1/favorites/folders/:id - 更新收藏夹
- DELETE /api/v1/favorites/folders/:id - 删除收藏夹
- GET /api/v1/favorites/folders/:id/articles - 收藏夹文章
- POST /api/v1/favorites - 添加收藏
- DELETE /api/v1/favorites/:article_id - 取消收藏
- GET /api/v1/favorites/articles - 我的收藏