# Blog API 项目概览

## 项目简介

博客后端 API，为前端开发者提供完整的练习接口。基于 Go 语言开发，采用现代化技术栈。

**项目状态**：✅ 已完成所有核心功能，可直接使用

## 快速开始

```bash
# 克隆和运行
git clone <repo-url>
cd blog-api
go mod download
cp configs/config.example.yaml configs/config.yaml
# 编辑配置文件
go run cmd/server/main.go
```

访问 http://localhost:8080/swagger/index.html 查看 API 文档。

## 核心功能

**用户系统**

- 注册登录，JWT 认证
- 个人资料管理
- 系统管理员自动初始化

**内容管理**

- 文章 CRUD，支持分类
- 文章搜索和排序
- 文件上传和管理

**互动功能**

- 点赞和收藏系统
- 多级评论系统
- 收藏夹管理

## 技术架构

**核心技术**

- 语言：Go 1.19+
- 框架：Gin + GORM
- 数据库：MySQL 8.0 + Redis 6.0
- 认证：JWT Token

**架构模式**

- 三层架构（Controller-Service-Repository）
- 依赖注入容器
- 中间件链式处理

## 项目结构

```
blog-api/
├── cmd/server/main.go      # 程序入口
├── internal/              # 业务代码
│   ├── controllers/       # 控制器层
│   ├── services/          # 业务逻辑层
│   ├── repositories/      # 数据访问层
│   ├── models/           # 数据模型
│   ├── middlewares/      # 中间件
│   └── routes/           # 路由配置
├── pkg/                  # 公共代码
│   ├── config/           # 配置管理
│   ├── database/         # 数据库连接
│   └── logger/           # 日志系统
├── configs/              # 配置文件
├── docs/                 # API文档
└── uploads/              # 文件存储
```

## API 接口

**已实现 42 个接口，涵盖 10 个功能模块**

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/articles` - 文章列表
- `POST /api/v1/articles` - 创建文章
- `GET /api/v1/categories` - 分类列表
- `POST /api/v1/comments` - 发表评论
- `POST /api/v1/files/upload` - 文件上传
- `POST /api/v1/favorites` - 添加收藏
- `POST /api/v1/articles/:id/like` - 文章点赞
- ...

完整接口文档：http://localhost:8080/swagger/index.html

## 数据库设计

**核心表结构**

- users: 用户基本信息
- articles: 文章内容和状态
- categories: 用户自定义分类
- comments: 多级评论系统
- files: 文件管理和去重

**互动功能表**

- article_likes: 文章点赞记录
- favorite_folders: 收藏夹管理
- article_favorites: 文章收藏关系

## 安全特性

- JWT Token 认证和授权
- 接口限流（令牌桶算法）
- 密码 bcrypt 加密
- 参数验证和数据校验
- SQL 注入防护
- CORS 跨域配置

## 性能优化

- Redis 缓存支持
- 数据库连接池
- 查询优化和索引
- 文件去重机制
- 结构化日志系统

## 配置管理

YAML 格式配置文件，支持环境变量。

```yaml
# 数据库配置
database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  dbname: blog_api

# Redis配置
redis:
  host: localhost
  port: 6379

# JWT配置
jwt:
  secret: your-secret-key
  access_expire: 3600s
```

## 部署说明

### 开发环境

```bash
go run cmd/server/main.go
```

### 生产环境

```bash
go build -o blog-api cmd/server/main.go
./blog-api
```

### 系统服务

可使用 systemd、Docker 或直接部署。

## 相关文档

- [README.md](README.md) - 项目介绍和快速开始
- [TODO.md](TODO.md) - 开发进度和任务列表
- [docs/design/requirements.md](docs/design/requirements.md) - 产品需求规格
- [docs/design/architecture.md](docs/design/architecture.md) - 技术架构设计

---

**项目已完成并可直接使用，适合前端开发者练习和学习。**
