# Blog API 项目概览

## 项目简介

这是一个基于 Go + MySQL + Redis 构建的博客 API 系统，主要为前端开发学习者提供免费的后端接口服务。

## 核心功能

- 用户认证：注册、登录、JWT Token 认证
- 文章管理：CRUD、分类、公开浏览、搜索排序
- 互动功能：文章点赞、收藏、评论系统
- 个人中心：我的文章、收藏夹、点赞记录
- 文件管理：图片上传、去重机制、本地存储
- 系统管理：用户管理、内容审核、数据统计

## 技术架构

- **语言框架**: Go + Gin
- **数据库**: MySQL 8.0 + Redis
- **认证**: JWT Token
- **存储**: 本地文件存储
- **架构**: 三层架构（Controller + Service + Model）

## 目录结构

```
blog-api/
├── cmd/server/main.go      # 程序入口
├── internal/
│   ├── controllers/        # 控制器层
│   ├── services/          # 业务逻辑层
│   └── models/            # 数据模型层（含数据库操作）
├── pkg/                   # 公共代码库
├── configs/               # 配置文件
├── uploads/               # 文件存储目录
└── docs/design/           # 设计文档
```

## 数据库表

- users: 用户信息
- articles: 文章内容
- categories: 文章分类
- comments: 评论系统
- article_likes: 文章点赞
- favorite_folders: 收藏夹
- article_favorites: 文章收藏
- comment_likes: 评论点赞
- files: 文件管理

## 核心 API 路由

- `/api/v1/auth/*` - 用户认证
- `/api/v1/user/*` - 用户管理
- `/api/v1/articles/*` - 文章管理
- `/api/v1/categories/*` - 分类管理
- `/api/v1/comments/*` - 评论管理
- `/api/v1/favorites/*` - 收藏管理
- `/api/v1/files/*` - 文件管理
- `/api/v1/admin/*` - 管理员功能

## 安全防护

- JWT 认证授权
- 接口限流防刷（Redis + 令牌桶算法）
- 数据验证和 SQL 注入防护
- XSS 攻击防护
- 密码 bcrypt 加密

## 部署方式

- 开发环境：本地直接运行
- 生产环境：二进制部署 + Nginx + systemd

## 配置管理

项目使用 YAML 格式配置文件：
- `configs/config.example.yaml` - 配置文件模板（会提交到仓库）
- `configs/config.yaml` - 实际配置文件（包含敏感信息，不会提交）

首次开发需要：
```bash
# 复制配置文件模板
cp configs/config.example.yaml configs/config.yaml
# 然后修改 config.yaml 中的数据库密码、JWT密钥等敏感信息
```

配置包含：数据库连接、Redis缓存、JWT认证、文件上传、限流策略、日志配置等。

## 开发命令

```bash
# 安装依赖
go mod download

# 开发运行
go run cmd/server/main.go

# 编译
go build -o blog-api cmd/server/main.go

# 数据库初始化
mysql < scripts/init.sql
```

## 文档位置

- 产品设计：`docs/design/requirements.md`
- 架构设计：`docs/design/architecture.md`
- 开发计划：`TODO.md`
