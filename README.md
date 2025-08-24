# Blog API

一个基于 Go 语言构建的博客后端 API 系统，专为前端开发学习者提供免费的接口服务。

## 快速开始

### 1. 环境准备
```bash
# Go 1.19+
# MySQL 8.0+  
# Redis 6.0+
```

### 2. 项目初始化
```bash
# 克隆项目
git clone <repository-url>
cd blog-api

# 安装依赖
go mod download

# 配置文件设置
cp configs/config.example.yaml configs/config.yaml
# 编辑 configs/config.yaml，修改数据库连接等配置
```

### 3. 数据库初始化
```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE blog_api;"

# 导入表结构（开发完成后提供）
mysql -u root -p blog_api < scripts/init.sql
```

### 4. 运行项目
```bash
# 开发环境运行
go run cmd/server/main.go

# 编译运行
go build -o blog-api cmd/server/main.go
./blog-api
```

## 项目文档

- **[项目概览](CLAUDE.md)** - 快速了解项目架构和功能
- **[开发计划](TODO.md)** - 完整的开发任务清单
- **[产品设计](docs/design/requirements.md)** - 详细的功能需求
- **[架构设计](docs/design/architecture.md)** - 技术架构和实现方案

## 功能特性

- ✅ 用户认证：注册、登录、JWT Token
- ✅ 文章管理：CRUD、分类、搜索、排序
- ✅ 互动功能：点赞、收藏、评论系统  
- ✅ 文件上传：图片上传、去重机制
- ✅ 安全防护：限流、参数验证、XSS防护
- ✅ 系统管理：用户管理、内容审核

## 技术栈

- **后端框架**: Go + Gin
- **数据库**: MySQL 8.0 + Redis  
- **认证方式**: JWT Token
- **文件存储**: 本地文件系统
- **架构模式**: 三层架构（MVC）

## API 文档

开发完成后将提供完整的 API 文档，包含所有接口的请求参数、响应格式和使用示例。

## 开发状态

当前处于设计阶段，已完成：
- [x] 产品需求分析
- [x] 架构设计方案  
- [x] 数据库表结构设计
- [x] API 接口设计
- [ ] 代码实现（待开始）

查看详细进度请参考 [开发计划](TODO.md)。

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建功能分支
3. 提交代码变更
4. 创建 Pull Request

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件。

## 联系方式

如有问题或建议，请通过 Issue 联系我们。