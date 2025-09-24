# 限流机制

为了保护服务器资源和确保服务质量，Blog API 实施了限流机制。

## 限流策略

Blog API 使用令牌桶算法实施限流，主要包含以下两个层面的限制：

### 1. IP 限流

- **限制**: 每个 IP 地址每分钟最多 100 次请求
- **用途**: 防止恶意用户或爬虫程序滥用 API
- **计数器重置**: 每分钟重置

### 2. 用户限流

- **限制**: 每个已认证用户每小时最多 1000 次请求
- **用途**: 防止单个用户过度使用服务
- **计数器重置**: 每小时重置

## 限流响应

当达到限流阈值时，API 会返回 `429 Too Many Requests` 状态码：

```json
{
  "code": 429,
  "message": "请求过于频繁，请稍后重试",
  "data": {
    "retry_after": 60,
    "limit": 100,
    "remaining": 0,
    "reset_time": "2024-01-01T00:01:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

响应字段说明：

- `retry_after`: 需要等待的秒数
- `limit`: 限制次数
- `remaining`: 剩余可请求次数
- `reset_time`: 计数器重置时间

## 限流头信息

某些请求可能包含以下限流相关信息的响应头：

- `X-RateLimit-Limit`: 限制次数
- `X-RateLimit-Remaining`: 剩余次数
- `X-RateLimit-Reset`: 重置时间戳

## 避免限流的建议

### 1. 合理控制请求频率

不要在循环中发送过多请求，特别是在批量操作时：

```javascript
// 错误示例：快速发送大量请求
for (let i = 0; i < 100; i++) {
  await fetch(`/api/v1/articles/${i}`);
}

// 正确示例：添加延迟或批量处理
async function fetchArticles(articleIds) {
  const results = [];
  for (let i = 0; i < articleIds.length; i++) {
    results.push(await fetch(`/api/v1/articles/${articleIds[i]}`));
    // 每次请求后等待 100ms
    await new Promise(resolve => setTimeout(resolve, 100));
  }
  return results;
}
```

### 2. 使用批量接口

如果需要获取多个资源，优先使用支持批量操作的接口：

```javascript
// 不推荐：多次单个请求
const article1 = await fetch('/api/v1/articles/1');
const article2 = await fetch('/api/v1/articles/2');
const article3 = await fetch('/api/v1/articles/3');

// 推荐：如果有批量接口
const articles = await fetch('/api/v1/articles?ids=1,2,3');
```

### 3. 缓存常用数据

对于不经常变化的数据，可以使用客户端缓存：

```javascript
class ApiClient {
  constructor() {
    this.cache = new Map();
    this.cacheTimeout = 5 * 60 * 1000; // 5分钟
  }
  
  async getCached(url) {
    const cached = this.cache.get(url);
    if (cached && Date.now() - cached.timestamp < this.cacheTimeout) {
      return cached.data;
    }
    
    const response = await fetch(url);
    const data = await response.json();
    
    this.cache.set(url, {
      data,
      timestamp: Date.now()
    });
    
    return data;
  }
}
```

### 4. 实施退避策略

在遇到限流时，实施指数退避策略：

```javascript
async function makeRequestWithBackoff(url, options, maxRetries = 3) {
  let attempts = 0;
  
  while (attempts <= maxRetries) {
    try {
      const response = await fetch(url, options);
      
      if (response.status === 429) {
        if (attempts === maxRetries) {
          throw new Error('Request limit exceeded after retries');
        }
        
        // 获取重试时间，如果没有则使用指数退避
        const retryAfter = response.headers.get('Retry-After') || Math.pow(2, attempts) * 1000;
        console.log(`Rate limited, waiting ${retryAfter} seconds...`);
        
        await new Promise(resolve => setTimeout(resolve, retryAfter * 1000));
        attempts++;
      } else {
        return response;
      }
    } catch (error) {
      throw error;
    }
  }
}
```

## 特殊场景处理

### 1. 文件上传

文件上传接口有特殊的限制：
- 单个文件大小不超过 10MB
- 频率限制相对宽松，但不建议频繁上传

### 2. 搜索接口

搜索接口可能会消耗较多资源，建议：
- 合理设置搜索频率
- 使用防抖技术减少不必要的搜索请求

```javascript
function createDebouncedSearch(searchFunction, delay = 300) {
  let timeoutId;
  return function(...args) {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => searchFunction.apply(this, args), delay);
  };
}

const debouncedSearch = createDebouncedSearch(async (keyword) => {
  const results = await fetch(`/api/v1/articles/search?keyword=${keyword}`);
  // 处理搜索结果
}, 500);
```

### 3. 实时功能

对于需要实时更新的场景，考虑使用 WebSocket 或长轮询，而不是频繁的 REST 请求。

## 监控和调试

### 检查请求频率

在开发过程中，可以监控自己的请求频率：

```javascript
class RequestMonitor {
  constructor() {
    this.requests = [];
    this.window = 60 * 1000; // 1分钟窗口
  }
  
  addRequest() {
    const now = Date.now();
    this.requests.push(now);
    
    // 清理窗口外的请求
    this.requests = this.requests.filter(timestamp => now - timestamp < this.window);
    
    const count = this.requests.length;
    console.log(`当前窗口内的请求数: ${count}`);
    
    if (count > 80) { // 接近限制
      console.warn('请求频率接近限制，请考虑减缓请求速度');
    }
  }
  
  getRequestsPerMinute() {
    return this.requests.length;
  }
}

const monitor = new RequestMonitor();
```

## 业务影响

限流机制会影响以下业务场景：

1. **批量数据操作**: 大量文章导入、用户批量注册等
2. **高频交互功能**: 实时聊天、频繁点赞等
3. **数据同步**: 大量数据的同步场景

对于这些场景，需要合理规划 API 调用策略，确保用户体验的同时不违反限流规则。

通过合理使用 API 和遵循最佳实践，可以有效避免限流限制，确保应用程序的稳定运行。