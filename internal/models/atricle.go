package models

import "time"

// UserInfo 用户信息（不含密码）
type UserInfo struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Article struct {
	Id        int       `gorm:"primarykey;column:id" json:"id"`
	Title     string    `gorm:"column:title" json:"title"`
	Content   string    `gorm:"column:content" json:"content"`
	UserId    int       `gorm:"column:user_id" json:"user_id"`
	User      User      `gorm:"foreignKey:UserId" json:"-"`
	UserInfo  UserInfo  `gorm:"-" json:"user"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// 创建帖子request
type CreateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// 修改帖子request
type UpdateArticleRequest struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// 帖子列表request
type ArticleListRequest struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`   // 搜索关键词（标题或内容）
	UserId   int    `form:"user_id"`  // 按用户ID过滤
	SortBy   string `form:"sort_by"`  // 排序字段: created_at, updated_at, title
	Order    string `form:"order"`    // 排序方向: asc, desc
}

// 帖子列表response
type ArticleListResponse struct {
	Articles []Article `json:"articles"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	Size     int       `json:"size"`
}
