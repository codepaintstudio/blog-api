# Blog API 产品规格

## 项目概述

博客后端 API，为前端开发者提供完整的练习接口。已实现所有核心功能，包含 42 个 API 接口。

**目标用户**：前端开发学习者、学生、初级开发者  
**项目状态**：✅ 已完成并可投入使用

## 功能模块

### 用户系统

- 邮箱注册登录，JWT 认证
- 用户资料管理，密码修改
- 系统管理员自动初始化
- 用户状态管理（活跃/禁用）

### 内容管理

- 文章 CRUD，支持草稿/发布状态
- 文章分类系统，用户自定义分类
- 文章搜索，支持标题和内容搜索
- 文件上传，支持图片去重机制

### 互动功能

- 文章点赞系统，防重复点赞
- 收藏夹管理，支持多个收藏夹
- 多级评论系统，支持 3 级嵌套回复
- 评论点赞功能

### 系统管理

- 用户数据管理和统计
- 内容审核和管理
- 系统数据统计
- 管理员权限控制

## 技术特性

### 安全防护

- 接口限流（IP 级别和用户级别）
- JWT Token 认证
- 密码 bcrypt 加密
- 参数验证和数据校验
- CORS 跨域配置

### 性能优化

- Redis 缓存支持
- 数据库连接池
- 查询优化和索引
- 文件去重机制

### 系统特性

- 配置文件管理（YAML）
- 结构化日志系统
- 优雅关闭机制


## 数据模型

### 核心表结构

**用户表 (users)**

```
id, username, email, password, nickname, avatar, bio
role(admin/user), status(active/inactive)
created_at, updated_at
```

**文章表 (articles)**

```
id, user_id, title, content, description, cover_image
category_id, status(draft/published), visibility(public/private)
view_count, like_count, favorite_count
created_at, updated_at
```

**分类表 (categories)**

```
id, user_id, name, description
created_at, updated_at
```

**评论表 (comments)**

```
id, article_id, user_id, parent_id, content
like_count, status(published/hidden/deleted)
created_at, updated_at
```

**互动表**

- article_likes: 文章点赞记录
- comment_likes: 评论点赞记录
- favorite_folders: 收藏夹
- article_favorites: 文章收藏记录
- files: 文件管理

## API 接口

### 接口概览

- **认证接口**：3 个（注册、登录、刷新）
- **用户接口**：3 个（资料查看、更新、密码）
- **文章接口**：6 个（CRUD、搜索、列表）
- **分类接口**：6 个（CRUD、列表、统计）
- **文件接口**：4 个（上传、列表、详情、删除）
- **评论接口**：8 个（CRUD、回复、审核、管理）
- **点赞接口**：4 个（文章点赞、评论点赞、历史）
- **收藏接口**：8 个（收藏夹、收藏管理、状态）

### 业务规则

- 用户只能管理自己的内容
- 文章作者可以管理文章评论
- 管理员拥有全局管理权限
- 每用户每文章只能点赞一次
- 收藏支持分文件夹组织

## 技术要求

### 环境依赖

- Go 1.19+
- MySQL 8.0+
- Redis 6.0+

### 性能指标

- API 响应时间 < 500ms
- 支持并发用户：1000+
- 文件上传限制：10MB
- 系统可用性：99%+

### 安全要求

- JWT Token 过期机制
- 接口防刷保护
- SQL 注入防护
- XSS 攻击防护
- 数据加密存储

## 部署说明

### 快速启动

```bash
# 安装依赖
go mod download

# 配置文件
cp configs/config.example.yaml configs/config.yaml

# 启动服务
go run cmd/server/main.go
```

### 访问地址

- API 服务：http://localhost:8080

- 健康检查：http://localhost:8080/api/v1/health

---

**项目已完成所有核心功能，可直接用于前端开发练习。**

- password: 密码（加密）
- nickname: 昵称
- avatar: 头像 URL
- bio: 个人简介
- role: 角色（admin/user）
- status: 状态（active/inactive）
- created_at: 创建时间
- updated_at: 更新时间

```

### 5.2 文章表 (articles)

```

- id: 文章 ID
- user_id: 作者 ID
- title: 标题
- content: 内容
- description: 描述
- cover_image: 封面图片 URL
- category_id: 分类 ID
- status: 状态（draft/published）
- visibility: 可见性（public/private）
- view_count: 阅读量
- created_at: 创建时间
- updated_at: 更新时间

```

### 5.3 分类表 (categories)

```

- id: 分类 ID
- user_id: 用户 ID
- name: 分类名称
- description: 分类描述
- created_at: 创建时间
- updated_at: 更新时间

```

### 5.4 评论表 (comments)

```

- id: 评论 ID
- article_id: 文章 ID
- user_id: 评论用户 ID
- parent_id: 父评论 ID（用于回复功能，顶级评论为 NULL）
- content: 评论内容
- like_count: 点赞数
- status: 状态（published/hidden/deleted）
- created_at: 创建时间
- updated_at: 更新时间

```

### 5.5 评论点赞表 (comment_likes)

```

- id: 点赞 ID
- comment_id: 评论 ID
- user_id: 点赞用户 ID
- created_at: 创建时间

```

### 5.6 文章点赞表 (article_likes)

```

- id: 点赞 ID
- article_id: 文章 ID
- user_id: 点赞用户 ID
- created_at: 创建时间

```

### 5.7 收藏夹表 (favorite_folders)

```

- id: 收藏夹 ID
- user_id: 用户 ID
- name: 收藏夹名称
- description: 收藏夹描述
- is_default: 是否默认收藏夹
- created_at: 创建时间
- updated_at: 更新时间

```

### 5.8 文章收藏表 (article_favorites)

```

- id: 收藏 ID
- article_id: 文章 ID
- user_id: 用户 ID
- folder_id: 收藏夹 ID
- created_at: 收藏旷间

```

## 6. 业务规则

### 6.1 用户规则

- 每个用户只能管理自己的内容
- 用户名和邮箱全局唯一
- 管理员可以查看所有用户数据

### 6.2 文章规则

- 文章只有作者可以编辑
- 私有文章仅作者可见
- 草稿状态文章不对外展示
- 管理员可以管理所有文章

### 6.3 分类规则

- 分类由用户自行创建和管理
- 分类删除时需处理关联文章
- 每个用户的分类名称不可重复

### 6.4 文章浏览规则

- 公开文章所有人可见
- 私有文章仅作者可见
- 文章浏览量实时更新
- 支持游客浏览（无需登录）

### 6.5 评论规则

- 评论需要登录后才能发布
- 用户只能删除自己的评论
- 文章作者可以删除文章下的任何评论
- 管理员可以管理所有评论
- 评论支持多级回复（最多 3 级）
- 每个用户对每条评论只能点赞一次

### 6.6 文章互动规则

- **点赞规则**：

  - 点赞需要登录后才能操作
  - 每个用户对每篇文章只能点赞一次
  - 用户可以取消自己的点赞
  - 文章作者不能对自己的文章点赞

- **收藏规则**：
  - 收藏需要登录后才能操作
  - 用户可以收藏任何公开文章（包括自己的）
  - 每个用户默认有一个“默认收藏夹”
  - 用户可以创建多个收藏夹进行分类管理
  - 收藏夹名称在用户下不可重复

## 7. API 设计原则

### 7.1 RESTful 风格

- 使用标准 HTTP 方法
- 资源导向的 URL 设计
- 统一的响应格式

### 7.2 安全原则

- 所有接口需要适当的权限验证
- 敏感操作需要额外验证
- 实现接口限流和防刷

### 7.3 可扩展性

- 预留扩展字段
- 版本化 API 设计
- 模块化架构设计

## 8. 项目里程碑

### 8.1 第一阶段：基础功能

- 用户注册登录
- 基本文章 CRUD
- 基础分类管理
- 文章公开浏览功能
- 文章点赞功能

### 8.2 第二阶段：完善功能

- 文章可见性控制
- 头像上传
- 文章收藏功能
- 评论系统
- 系统管理功能

### 8.3 第三阶段：优化增强

- 评论点赞和多级回复
- 收藏夹分类管理
- 文章搜索和排序
- 个人中心完善功能
- 性能优化
- 安全加固
- API 文档完善
```
