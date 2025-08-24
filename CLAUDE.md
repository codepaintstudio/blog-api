# 博客 API 项目 - Claude 助手指南

## 项目概述
这是一个基于 Go 语言的博客 API 系统，提供用户认证和文章 CRUD 操作。项目使用 Gin Web 框架、GORM ORM 和 MySQL 数据库。

## 技术栈
- **开发语言**: Go 1.23+
- **Web 框架**: Gin (HTTP Web 框架)
- **数据库**: MySQL + GORM ORM
- **身份认证**: JWT 令牌
- **密码安全**: bcrypt 哈希加密

## 项目结构
```
blog-api/
├── main.go                 # 应用程序入口
├── go.mod                  # Go 模块依赖
├── internal/
│   ├── controllers/        # HTTP 请求处理器
│   │   ├── article.go     # 文章 CRUD 操作
│   │   └── user.go        # 用户注册/登录
│   ├── models/            # 数据库模型
│   │   ├── article.go     # 文章模型
│   │   └── user.go        # 用户模型  
│   ├── routes/            # 路由定义
│   │   └── routes.go      # API 路由设置
│   └── services/          # 业务逻辑层
│       ├── article.go     # 文章服务层
│       └── user.go        # 用户服务层
└── pkg/
    ├── middleware/        # HTTP 中间件
    │   ├── auth.go        # JWT 身份验证
    │   └── cors.go        # 跨域处理
    ├── response/          # 标准 API 响应
    │   └── response.go    # 响应格式化
    └── utils/             # 工具函数
        ├── jwt.go         # JWT 令牌操作
        └── password.go    # 密码哈希处理
```

## 核心功能
1. **用户管理**: 用户注册和登录，使用 JWT 身份验证
2. **文章管理**: 博客文章的 CRUD 操作，包含作者权限验证
3. **安全特性**: 密码哈希、JWT 令牌、授权中间件
4. **API 规范**: 统一的 JSON 响应格式

## 开发环境配置

### 环境要求
- Go 1.23+
- MySQL 数据库
- 环境变量配置

### 环境变量设置
创建 `.env` 文件：
```
DB_HOST=localhost
DB_PORT=3306
DB_USER=数据库用户名
DB_PASSWORD=数据库密码
DB_NAME=数据库名
PORT=7777
JWT_SECRET=JWT密钥
```

### 运行应用程序
```bash
# 安装依赖
go mod tidy

# 启动服务器
go run main.go
```

服务器将在 7777 端口启动（或 .env 中指定的 PORT）。

## API 接口

### 用户相关接口
- `POST /api/user/register` - 用户注册
- `POST /api/user/login` - 用户登录

### 用户信息查询接口
- `GET /api/users/:id` - 获取用户基本信息
- `GET /api/users/:id/detail` - 获取用户详情（包含文章统计）

### 文章相关接口  
- `GET /api/articles` - 获取文章列表（支持搜索、排序、过滤、分页）
- `GET /api/articles/:id` - 获取指定文章
- `GET /api/articles/stats` - 获取文章统计信息
- `POST /api/articles` - 创建文章（需要身份验证）
- `PUT /api/articles/:id` - 更新文章（需要身份验证 + 所有权验证）
- `DELETE /api/articles/:id` - 删除文章（需要身份验证 + 所有权验证）

## 高级查询功能

### 文章列表查询参数
`GET /api/articles` 支持以下查询参数：

- `page` - 页码（默认: 1）
- `size` - 每页数量（默认: 10）
- `search` - 搜索关键词（搜索标题和内容）
- `user_id` - 按用户ID过滤文章
- `sort_by` - 排序字段（created_at, updated_at, title）
- `order` - 排序方向（asc, desc，默认: desc）

**示例：**
```bash
# 搜索包含"技术"的文章
GET /api/articles?search=技术

# 获取用户ID为1的所有文章
GET /api/articles?user_id=1

# 按标题升序排列
GET /api/articles?sort_by=title&order=asc

# 组合查询：搜索用户1的技术文章，按时间倒序
GET /api/articles?search=技术&user_id=1&sort_by=created_at&order=desc
```

## 常用开发命令

### 运行测试
```bash
go test ./...
```

### 数据库迁移
应用程序使用 GORM 的 AutoMigrate 功能在启动时自动创建/更新数据库表。

### 添加新功能
1. 在 `internal/models/` 中创建模型
2. 在 `internal/services/` 中添加服务逻辑
3. 在 `internal/controllers/` 中创建控制器
4. 在 `internal/routes/routes.go` 中注册路由

## 安全说明
- 所有密码使用 bcrypt 进行哈希加密
- 使用 JWT 令牌进行身份验证
- 文章修改需要所有权验证
- 实现了 CORS 中间件处理跨域请求

## 安全防护功能

### IP 限流保护
系统集成了 `golang.org/x/time/rate` 官方限流库，提供强大的反爬虫和防刷保护：

**限流策略：**
- **正常限制**：每秒最多 10 个请求
- **突发处理**：允许突发 20 个请求
- **违规检测**：5 次违规触发 IP 封禁
- **封禁时长**：30 分钟自动解封

**防护特性：**
- 🛡️ IP级别限流，防止单个IP恶意请求
- 🚫 自动封禁异常IP，有效防止攻击
- 🧹 自动清理过期记录，节省内存
- 📊 实时监控统计，便于运维管理

### 监控接口
可通过以下接口查看限流状态：
```bash
GET /api/system/rate-limit-stats
```

响应示例：
```json
{
  "code": 200,
  "message": "Rate limit stats",
  "data": {
    "active_ips": 15,     // 活跃IP数量
    "banned_count": 2,    // 被封禁IP数量  
    "rate_limit": 10.0,   // 每秒请求限制
    "burst_size": 20      // 突发请求大小
  }
}
```

## 部署信息
当前应用部署地址: https://network-demo.hub.feashow.cn/

本地开发默认端口: 7777