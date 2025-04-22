import request from '../utils/request'

// 获取帖子列表
export const getArticles = (page = 1, size = 10) => {
  return request.get('/articles', { params: { page, size } })
}

// 获取帖子详情
export const getArticleById = (id) => {
  return request.get(`/articles/${id}`)
}

// 创建帖子
export const createArticle = (title, content) => {
  return request.post('/articles', { title, content })
}

// 更新帖子
export const updateArticle = (id, title, content) => {
  return request.put(`/articles/${id}`, { title, content })
}

// 删除帖子
export const deleteArticle = (id) => {
  return request.delete(`/articles/${id}`)
}
