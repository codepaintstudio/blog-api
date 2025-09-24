# 安全机制

了解 Blog API 的安全机制对于安全地使用接口非常重要。

## 认证机制

### JWT Token 认证

Blog API 使用 JWT (JSON Web Token) 进行身份认证，主要包含：

- **Access Token**: 有效期 1 小时，用于日常 API 调用
- **Refresh Token**: 有效期 7 天，用于刷新 Access Token

### Token 使用方式

在需要认证的请求中，需要在请求头中添加 Bearer Token：

```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Token 刷新

Access Token 过期后，使用 Refresh Token 获取新的 Access Token：

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your_refresh_token_here"
  }'
```

## 权限控制

### 基于角色的访问控制 (RBAC)

Blog API 实现了基于角色的权限控制：

- **普通用户 (user)**: 基本功能权限
- **管理员 (admin)**: 管理功能权限

### 权限级别

1. **公开接口**: 无需认证
2. **用户接口**: 需要登录用户
3. **管理员接口**: 需要管理员权限

## 数据安全

### 密码安全

- 密码使用 Bcrypt 算法加密存储
- 支持密码强度验证
- 密码历史记录检查，防止重复使用

### 输入验证

所有 API 端点都实施严格的数据验证：

```javascript
// 示例：创建文章的验证规则
const articleValidationRules = {
  title: {
    required: true,
    minLength: 1,
    maxLength: 200,
    type: 'string'
  },
  content: {
    required: true,
    minLength: 1,
    type: 'string'
  },
  category_id: {
    required: false,
    type: 'number',
    min: 1
  }
};
```

### SQL 注入防护

使用 GORM ORM 框架，所有数据库查询都使用预编译语句，有效防止 SQL 注入。

### XSS 防护

- 对用户输入的内容进行 HTML 转义
- 使用 Content Security Policy (CSP) 头

## 通信安全

### HTTPS 传输

生产环境中建议使用 HTTPS 确保数据传输安全。

### 敏感信息处理

- Token 不会明文存储在日志中
- 敏感信息在日志中会被脱敏处理

## 限流与防护

### 接口限流

- **IP 限流**: 每个 IP 每分钟最多 100 次请求
- **用户限流**: 每个用户每小时最多 1000 次请求

### 防护机制

- 请求频率限制
- 异常行为检测
- 恶意请求识别

## 安全最佳实践

### 1. Token 安全存储

客户端应安全地存储 Token：

```javascript
// 在浏览器中使用 httpOnly Cookie (推荐)
// 或在 localStorage/sessionStorage 中存储（需要注意 XSS 风险）
function storeTokens(accessToken, refreshToken) {
  // 存储 access token
  localStorage.setItem('accessToken', accessToken);
  // 存储 refresh token
  localStorage.setItem('refreshToken', refreshToken);
}

function clearTokens() {
  localStorage.removeItem('accessToken');
  localStorage.removeItem('refreshToken');
}
```

### 2. Token 使用安全

```javascript
// 在请求中安全添加 token
const apiClient = {
  async request(url, options = {}) {
    const token = localStorage.getItem('accessToken');
    
    // 验证 token 是否存在
    if (!token) {
      throw new Error('No authentication token');
    }
    
    // 添加认证头
    const config = {
      ...options,
      headers: {
        ...options.headers,
        'Authorization': `Bearer ${token}`
      }
    };
    
    const response = await fetch(url, config);
    
    // 处理 token 过期
    if (response.status === 401) {
      // 尝试刷新 token
      const refreshResult = await this.refreshToken();
      if (refreshResult.success) {
        // 重新发起原请求
        config.headers['Authorization'] = `Bearer ${refreshResult.token}`;
        return fetch(url, config);
      } else {
        // 跳转到登录页
        window.location.href = '/login';
      }
    }
    
    return response;
  }
};
```

### 3. 前端安全

- 验证用户输入（虽然后端也会验证，但前端验证可以提供更好体验）
- 防止 XSS 攻击
- 使用现代前端框架的安全特性

### 4. 错误处理

不要在错误响应中暴露敏感信息：

```javascript
// 正确的错误处理
app.use((err, req, res, next) => {
  // 记录详细错误信息到服务器日志
  logger.error('API Error:', {
    error: err.message,
    stack: err.stack,
    url: req.url,
    method: req.method,
    ip: req.ip
  });

  // 只返回通用错误信息给客户端
  res.status(500).json({
    code: 500,
    message: '服务器内部错误',
    data: null,
    timestamp: new Date().toISOString()
  });
});
```

## 安全监控

### 日志记录

系统记录以下安全相关事件：

- 登录尝试（成功/失败）
- Token 刷新
- 权限相关错误
- 异常请求模式

### 安全审计

- 定期审查权限配置
- 监控异常访问模式
- 定期更新安全策略

## 常见安全威胁与防护

### 1. CSRF 攻击

在 Web 应用中使用适当的 CSRF 防护措施：

```javascript
// 获取 CSRF Token
const csrfToken = document.querySelector('meta[name="csrf-token"]').getAttribute('content');

// 在请求中包含 CSRF Token
fetch('/api/v1/sensitive-action', {
  method: 'POST',
  headers: {
    'X-CSRF-Token': csrfToken,
    'Authorization': 'Bearer ' + accessToken
  }
});
```

### 2. 会话固定攻击

- 每次登录成功后都生成新的 Token
- 实现 Token 黑名单机制

### 3. 暴力破解

- 登录失败次数限制
- 临时账户锁定机制

### 4. 数据泄露

- 实施最小权限原则
- 敏感数据脱敏
- 访问日志记录

## 安全事件响应

### 发现安全问题时

1. 立即断开相关认证 Token
2. 通知系统管理员
3. 记录详细事件信息
4. 实施临时防护措施

### Token 泄露处理

如果怀疑 Token 被泄露：

```javascript
// 立即清除本地存储的 Token
localStorage.removeItem('accessToken');
localStorage.removeItem('refreshToken');

// 调用后端接口使 Token 失效
await fetch('/api/v1/auth/logout', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer ' + leakedToken
  }
});
```

## 合规性

Blog API 遵循以下安全标准：

- 数据最小化原则
- 用户隐私保护
- 适当的访问控制
- 安全的密码存储

通过理解这些安全机制，你可以更好地使用 Blog API 并保护用户数据的安全性。