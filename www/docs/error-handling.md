# 错误处理

了解 API 的错误处理机制对于正确使用接口非常重要。

## 通用响应格式

所有 API 响应都遵循以下格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": "2024-01-01T00:00:00Z"
}
```

- `code`: HTTP 状态码或自定义业务状态码
- `message`: 响应消息，描述操作结果
- `data`: 实际响应数据，可能为空对象或 null
- `timestamp`: 响应时间戳

## HTTP 状态码

### 2xx 成功

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 资源创建成功 |
| 204 | 请求成功，无响应内容 |

### 4xx 客户端错误

| 状态码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未认证或 Token 过期 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 405 | 请求方法不被允许 |
| 409 | 资源冲突 |
| 413 | 请求实体过大 |
| 422 | 请求参数验证失败 |
| 429 | 请求过于频繁（限流） |

### 5xx 服务器错误

| 状态码 | 说明 |
|--------|------|
| 500 | 服务器内部错误 |
| 502 | 网关错误 |
| 503 | 服务不可用 |
| 504 | 网关超时 |

## 业务错误码

除了 HTTP 状态码，API 还定义了以下业务错误码：

| 业务码 | 说明 |
|--------|------|
| 10000 | 通用错误 |
| 10001 | 参数验证错误 |
| 10002 | Token 过期 |
| 10003 | Token 无效 |
| 10004 | 权限不足 |
| 10005 | 资源不存在 |
| 10006 | 资源已存在 |
| 10007 | 请求过于频繁 |
| 10008 | 文件上传错误 |
| 10009 | 文件大小超限 |
| 10010 | 文件类型不支持 |

## 常见错误示例

### 参数验证错误

```json
{
  "code": 400,
  "message": "请求参数错误",
  "data": {
    "errors": {
      "username": [
        "用户名为必填项",
        "用户名长度必须在3-20个字符之间"
      ],
      "email": [
        "邮箱格式不正确"
      ]
    }
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 未认证错误

```json
{
  "code": 401,
  "message": "未认证或 Token 过期",
  "data": null,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 权限不足错误

```json
{
  "code": 403,
  "message": "权限不足，无法执行此操作",
  "data": null,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 资源不存在错误

```json
{
  "code": 404,
  "message": "请求的资源不存在",
  "data": null,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 错误处理最佳实践

### 1. 请求前验证

在发送请求前，先在客户端验证参数的有效性，减少不必要的请求。

```javascript
// 错误示例：没有验证就发送请求
function createUser(userData) {
  return fetch('/api/v1/users', {
    method: 'POST',
    body: JSON.stringify(userData)
  });
}

// 正确示例：先验证再发送请求
function createUser(userData) {
  // 验证参数
  if (!userData.email || !isValidEmail(userData.email)) {
    throw new Error('邮箱格式不正确');
  }
  
  return fetch('/api/v1/users', {
    method: 'POST',
    body: JSON.stringify(userData)
  });
}
```

### 2. 错误响应处理

正确处理错误响应，给用户友好的提示。

```javascript
async function handleApiResponse(response) {
  const data = await response.json();
  
  if (!response.ok) {
    // 根据错误码处理不同类型的错误
    switch (response.status) {
      case 400:
        showMessage('请求参数错误：' + data.message, 'error');
        break;
      case 401:
        // Token 过期，重定向到登录页
        redirectToLogin();
        break;
      case 403:
        showMessage('权限不足', 'error');
        break;
      case 429:
        showMessage('请求太频繁，请稍后重试', 'warning');
        break;
      case 500:
        showMessage('服务器错误，请稍后重试', 'error');
        break;
      default:
        showMessage('请求失败：' + data.message, 'error');
    }
    return null;
  }
  
  return data;
}
```

### 3. Token 过期处理

实现 Token 自动刷新机制。

```javascript
let isRefreshing = false;
let failedQueue = [];

function processQueue(error, token = null) {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  
  failedQueue = [];
}

// 请求拦截器中处理 Token 过期
axios.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config;
    
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // 如果正在刷新 Token，将请求加入队列
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then(token => {
          originalRequest.headers['Authorization'] = 'Bearer ' + token;
          return axios(originalRequest);
        });
      }
      
      originalRequest._retry = true;
      isRefreshing = true;
      
      try {
        const refreshToken = localStorage.getItem('refreshToken');
        const response = await refreshTokenApi(refreshToken);
        const newToken = response.data.access_token;
        
        localStorage.setItem('accessToken', newToken);
        axios.defaults.headers.common['Authorization'] = 'Bearer ' + newToken;
        originalRequest.headers['Authorization'] = 'Bearer ' + newToken;
        
        processQueue(null, newToken);
        return axios(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        // 重定向到登录页
        redirectToLogin();
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }
    
    return Promise.reject(error);
  }
);
```

### 4. 限流处理

实现请求重试机制。

```javascript
async function makeRequestWithRetry(url, options, maxRetries = 3) {
  let attempts = 0;
  
  while (attempts < maxRetries) {
    try {
      const response = await fetch(url, options);
      
      if (response.status === 429) {
        // 限流错误，等待后重试
        const retryAfter = response.headers.get('Retry-After') || 1;
        await new Promise(resolve => setTimeout(resolve, retryAfter * 1000));
      } else {
        return response;
      }
    } catch (error) {
      attempts++;
      if (attempts >= maxRetries) {
        throw error;
      }
      // 指数退避
      await new Promise(resolve => setTimeout(resolve, Math.pow(2, attempts) * 1000));
    }
  }
}
```

## 调试技巧

### 1. 使用开发工具

利用浏览器开发者工具查看请求和响应的详细信息。

### 2. 记录错误日志

在客户端记录错误信息，便于调试和分析。

```javascript
function logError(error, context = '') {
  console.error(`[API Error] ${context}`, {
    message: error.message,
    stack: error.stack,
    timestamp: new Date().toISOString()
  });
  
  // 可选：发送错误到日志服务
  // sendErrorToLogService(error);
}
```

### 3. 详细的错误信息

在开发环境提供详细的错误信息，生产环境则提供用户友好的错误信息。

通过理解这些错误处理机制，你可以更好地使用 Blog API 并构建更可靠的应用程序。