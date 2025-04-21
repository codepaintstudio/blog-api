import request from '../utils/request'

// 获取文章列表
export const getArticles = (page = 1, size = 10) => {
  return request.get('/articles', { params: { page, size } })
}

// 获取文章详情
export const getArticleById = (id) => {
  return request.get(`/articles/${id}`)
}

// 创建文章
export const createArticle = (title, content) => {
  return request.post('/articles', { title, content })
}

// 更新文章
export const updateArticle = (id, title, content) => {
  return request.put(`/articles/${id}`, { title, content })
}

// 删除文章
export const deleteArticle = (id) => {
  return request.delete(`/articles/${id}`)
}
