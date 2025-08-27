# Blog API

一个博客后端 API，为前端开发者提供练习接口。

## 快速开始

### 环境要求

- Go 1.19+
- MySQL 8.0+
- Redis 6.0+

### 安装运行

```bash
# 克隆项目
git clone https://github.com/codepaintstudio/blog-api.git
cd blog-api

# 安装依赖
go mod download

# 复制配置文件
cp configs/config.example.yaml configs/config.yaml

# 编辑配置文件，修改数据库连接信息
vim configs/config.yaml

# 运行项目
go run cmd/server/main.go
```

服务启动后，访问 http://localhost:8080/swagger/index.html 查看 API 文档。

## 功能特性

**用户系统**

- 注册登录，JWT 认证
- 个人资料管理
- 系统管理员自动初始化

**内容管理**

- 文章 CRUD，支持分类
- 文章搜索和排序
- 文件上传和管理

**互动功能**

- 文章点赞和收藏
- 收藏夹管理
- 多级评论系统

**系统特性**

- 接口限流和安全防护
- Redis 缓存优化
- 完整的 API 文档

## 技术架构

- **语言框架**: Go + Gin
- **数据存储**: MySQL + Redis
- **认证方式**: JWT Token
- **文档工具**: Swagger
- **架构模式**: 三层架构（Controller-Service-Repository）

## API 文档

启动项目后访问：http://localhost:8080/swagger/index.html

包含 42 个接口，涵盖 10 个功能模块的完整 API 文档。

## 项目状态

✅ **已完成**

- 用户认证系统
- 文章管理功能
- 分类管理
- 评论系统
- 点赞收藏功能
- 文件上传
- 系统管理
- API 文档生成

详细开发进度见 [TODO.md](TODO.md)。
