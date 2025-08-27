package services

import (
	"errors"
	"fmt"
	"time"

	"blog-api/internal/models"
	"blog-api/internal/repositories"
	"blog-api/pkg/cache"
	"blog-api/pkg/logger"
	
	"gorm.io/gorm"
)

// ArticleService 文章服务接口
type ArticleService interface {
	Create(userID uint, req *CreateArticleRequest) (*models.Article, error)
	GetByID(id uint, userID *uint) (*ArticleDetailResponse, error)
	Update(id, userID uint, req *UpdateArticleRequest) (*models.Article, error)
	Delete(id, userID uint) error
	List(req *ListArticleRequest) (*ListArticleResponse, error)
	GetUserArticles(userID, currentUserID uint, req *ListArticleRequest) (*ListArticleResponse, error)
	Publish(id, userID uint) error
	SetTop(id, userID uint, isTop bool) error
	IncrementViewCount(id uint) error
	Search(keyword string, page, pageSize int) (*ListArticleResponse, error)
}

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title       string `json:"title" validate:"required,max=255"`
	Content     string `json:"content" validate:"required"`
	Description string `json:"description" validate:"max=500"`
	CoverImage  string `json:"cover_image"`
	CategoryID  *uint  `json:"category_id"`
	Status      string `json:"status" validate:"omitempty,oneof=draft published"`
	Visibility  string `json:"visibility" validate:"omitempty,oneof=public private"`
	AllowComment bool  `json:"allow_comment"`
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title       string `json:"title" validate:"max=255"`
	Content     string `json:"content"`
	Description string `json:"description" validate:"max=500"`
	CoverImage  string `json:"cover_image"`
	CategoryID  *uint  `json:"category_id"`
	Status      string `json:"status" validate:"omitempty,oneof=draft published"`
	Visibility  string `json:"visibility" validate:"omitempty,oneof=public private"`
	AllowComment *bool `json:"allow_comment"`
}

// ListArticleRequest 文章列表请求
type ListArticleRequest struct {
	Page       int    `json:"page" validate:"min=1"`
	PageSize   int    `json:"page_size" validate:"min=1,max=100"`
	Status     string `json:"status" validate:"omitempty,oneof=draft published"`
	Visibility string `json:"visibility" validate:"omitempty,oneof=public private"`
	CategoryID *uint  `json:"category_id"`
	SortBy     string `json:"sort_by" validate:"omitempty,oneof=created_at updated_at view_count like_count"`
	SortOrder  string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// ArticleDetailResponse 文章详情响应
type ArticleDetailResponse struct {
	*ArticleResponse
	Content     string `json:"content"`
	AllowComment bool  `json:"allow_comment"`
}

// ArticleResponse 文章响应
type ArticleResponse struct {
	ID            uint                 `json:"id"`
	Title         string               `json:"title"`
	Description   string               `json:"description"`
	CoverImage    string               `json:"cover_image"`
	Status        string               `json:"status"`
	Visibility    string               `json:"visibility"`
	IsTop         bool                 `json:"is_top"`
	ViewCount     int                  `json:"view_count"`
	LikeCount     int                  `json:"like_count"`
	CommentCount  int                  `json:"comment_count"`
	FavoriteCount int                  `json:"favorite_count"`
	PublishedAt   *time.Time           `json:"published_at"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
	Author        *UserResponse        `json:"author"`
	Category      *CategoryResponse    `json:"category"`
}

// ListArticleResponse 文章列表响应
type ListArticleResponse struct {
	Articles []*ArticleResponse `json:"articles"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
	Pages    int                `json:"pages"`
}

// articleService 文章服务实现
type articleService struct {
	articleRepo  repositories.ArticleRepository
	categoryRepo repositories.CategoryRepository
	cacheService cache.CacheService
	cacheKey     *cache.CacheKey
}

// NewArticleService 创建新的文章服务
func NewArticleService() ArticleService {
	return &articleService{
		cacheService: cache.NewCacheService(),
		cacheKey:     cache.NewCacheKey(),
		articleRepo:  repositories.NewArticleRepository(),
		categoryRepo: repositories.NewCategoryRepository(),
	}
}

// Create 创建文章
func (s *articleService) Create(userID uint, req *CreateArticleRequest) (*models.Article, error) {
	// 验证分类是否存在
	if req.CategoryID != nil {
		_, err := s.categoryRepo.GetByID(*req.CategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("分类不存在")
			}
			return nil, err
		}
	}
	
	// 设置默认值
	status := "draft"
	if req.Status != "" {
		status = req.Status
	}
	
	visibility := "public"
	if req.Visibility != "" {
		visibility = req.Visibility
	}
	
	// 创建文章
	article := &models.Article{
		UserID:       userID,
		CategoryID:   req.CategoryID,
		Title:        req.Title,
		Content:      req.Content,
		Description:  req.Description,
		CoverImage:   req.CoverImage,
		Status:       status,
		Visibility:   visibility,
		AllowComment: req.AllowComment,
	}
	
	// 如果是发布状态，设置发布时间
	if status == "published" {
		now := time.Now()
		article.PublishedAt = &now
	}
	
	if err := s.articleRepo.Create(article); err != nil {
		logger.Error("创建文章失败", logger.Err("error", err))
		return nil, err
	}
	
	logger.Info("文章创建成功",
		logger.Int("user_id", int(userID)),
		logger.Int("article_id", int(article.ID)),
		logger.String("title", article.Title),
	)
	
	return article, nil
}

// GetByID 根据ID获取文章
func (s *articleService) GetByID(id uint, userID *uint) (*ArticleDetailResponse, error) {
	// 尝试从缓存获取文章
	cacheKey := s.cacheKey.Article(id)
	var article *models.Article
	
	err := s.cacheService.Get(cacheKey, &article)
	if err != nil {
		// 缓存未命中，从数据库获取
		article, err = s.articleRepo.GetByIDWithAuthor(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("文章不存在")
			}
			return nil, err
		}
		
		// 将文章缓存5分钟
		s.cacheService.Set(cacheKey, article, 5*time.Minute)
	}
	
	// 检查访问权限
	if article.Visibility == "private" && (userID == nil || *userID != article.UserID) {
		return nil, errors.New("文章不存在")
	}
	
	if article.Status == "draft" && (userID == nil || *userID != article.UserID) {
		return nil, errors.New("文章不存在")
	}
	
	// 如果不是作者访问，增加浏览量
	if userID == nil || *userID != article.UserID {
		s.articleRepo.IncrementViewCount(id)
		// 异步更新缓存中的浏览量
		go func() {
			article.ViewCount++
			s.cacheService.Set(cacheKey, article, 5*time.Minute)
		}()
	}
	
	return s.toArticleDetailResponse(article), nil
}

// Update 更新文章
func (s *articleService) Update(id, userID uint, req *UpdateArticleRequest) (*models.Article, error) {
	// 获取文章
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}
	
	// 检查权限
	if article.UserID != userID {
		return nil, errors.New("无权限操作")
	}
	
	// 验证分类是否存在
	if req.CategoryID != nil {
		_, err := s.categoryRepo.GetByID(*req.CategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("分类不存在")
			}
			return nil, err
		}
	}
	
	// 更新字段
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.Description != "" {
		article.Description = req.Description
	}
	if req.CoverImage != "" {
		article.CoverImage = req.CoverImage
	}
	if req.CategoryID != nil {
		article.CategoryID = req.CategoryID
	}
	if req.Status != "" {
		oldStatus := article.Status
		article.Status = req.Status
		
		// 如果从草稿变为发布，设置发布时间
		if oldStatus == "draft" && req.Status == "published" {
			now := time.Now()
			article.PublishedAt = &now
		}
	}
	if req.Visibility != "" {
		article.Visibility = req.Visibility
	}
	if req.AllowComment != nil {
		article.AllowComment = *req.AllowComment
	}
	
	if err := s.articleRepo.Update(article); err != nil {
		logger.Error("更新文章失败", logger.Err("error", err))
		return nil, err
	}
	
	// 清除相关缓存
	s.cacheService.Delete(s.cacheKey.Article(id))
	s.cacheService.DeletePattern("articles:list:*")
	s.cacheService.DeletePattern("articles:popular")
	
	logger.Info("文章更新成功",
		logger.Int("user_id", int(userID)),
		logger.Int("article_id", int(article.ID)),
		logger.String("title", article.Title),
	)
	
	return article, nil
}

// Delete 删除文章
func (s *articleService) Delete(id, userID uint) error {
	// 获取文章
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}
	
	// 检查权限
	if article.UserID != userID {
		return errors.New("无权限操作")
	}
	
	if err := s.articleRepo.Delete(id); err != nil {
		logger.Error("删除文章失败", logger.Err("error", err))
		return err
	}
	
	// 清除相关缓存
	s.cacheService.Delete(s.cacheKey.Article(id))
	s.cacheService.DeletePattern("articles:list:*")
	s.cacheService.DeletePattern("articles:popular")
	
	logger.Info("文章删除成功",
		logger.Int("user_id", int(userID)),
		logger.Int("article_id", int(id)),
		logger.String("title", article.Title),
	)
	
	return nil
}

// List 获取文章列表
func (s *articleService) List(req *ListArticleRequest) (*ListArticleResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	
	offset := (req.Page - 1) * req.PageSize
	
	// 构建过滤器
	filters := &repositories.ArticleFilters{
		Status:     req.Status,
		Visibility: req.Visibility,
		CategoryID: req.CategoryID,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
	}
	
	// 公开接口只显示已发布的公开文章
	if filters.Status == "" {
		filters.Status = "published"
	}
	if filters.Visibility == "" {
		filters.Visibility = "public"
	}
	
	articles, total, err := s.articleRepo.List(offset, req.PageSize, filters)
	if err != nil {
		return nil, err
	}
	
	return s.toListArticleResponse(articles, total, req.Page, req.PageSize), nil
}

// GetUserArticles 获取用户文章列表
func (s *articleService) GetUserArticles(userID, currentUserID uint, req *ListArticleRequest) (*ListArticleResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	
	offset := (req.Page - 1) * req.PageSize
	
	// 构建过滤器
	filters := &repositories.ArticleFilters{
		UserID:     &userID,
		CategoryID: req.CategoryID,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
	}
	
	// 如果是查看自己的文章，可以看到所有状态
	if userID == currentUserID {
		filters.Status = req.Status
		filters.Visibility = req.Visibility
	} else {
		// 查看他人文章只能看到已发布的公开文章
		filters.Status = "published"
		filters.Visibility = "public"
	}
	
	articles, total, err := s.articleRepo.List(offset, req.PageSize, filters)
	if err != nil {
		return nil, err
	}
	
	return s.toListArticleResponse(articles, total, req.Page, req.PageSize), nil
}

// Publish 发布文章
func (s *articleService) Publish(id, userID uint) error {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}
	
	// 检查权限
	if article.UserID != userID {
		return errors.New("无权限操作")
	}
	
	// 更新状态和发布时间
	article.Status = "published"
	if article.PublishedAt == nil {
		now := time.Now()
		article.PublishedAt = &now
	}
	
	if err := s.articleRepo.Update(article); err != nil {
		logger.Error("发布文章失败", logger.Err("error", err))
		return err
	}
	
	logger.Info("文章发布成功",
		logger.Int("user_id", int(userID)),
		logger.Int("article_id", int(id)),
	)
	
	return nil
}

// SetTop 设置文章置顶
func (s *articleService) SetTop(id, userID uint, isTop bool) error {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}
	
	// 检查权限（只有作者可以置顶自己的文章）
	if article.UserID != userID {
		return errors.New("无权限操作")
	}
	
	article.IsTop = isTop
	
	if err := s.articleRepo.Update(article); err != nil {
		logger.Error("设置文章置顶失败", logger.Err("error", err))
		return err
	}
	
	logger.Info("文章置顶设置成功",
		logger.Int("user_id", int(userID)),
		logger.Int("article_id", int(id)),
		logger.String("is_top", fmt.Sprintf("%t", isTop)),
	)
	
	return nil
}

// IncrementViewCount 增加浏览量
func (s *articleService) IncrementViewCount(id uint) error {
	return s.articleRepo.IncrementViewCount(id)
}

// Search 搜索文章
func (s *articleService) Search(keyword string, page, pageSize int) (*ListArticleResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	
	offset := (page - 1) * pageSize
	
	articles, total, err := s.articleRepo.Search(keyword, offset, pageSize)
	if err != nil {
		return nil, err
	}
	
	return s.toListArticleResponse(articles, total, page, pageSize), nil
}

// toArticleDetailResponse 转换为文章详情响应
func (s *articleService) toArticleDetailResponse(article *models.Article) *ArticleDetailResponse {
	resp := &ArticleDetailResponse{
		ArticleResponse: s.toArticleResponse(article),
		Content:         article.Content,
		AllowComment:    article.AllowComment,
	}
	return resp
}

// toArticleResponse 转换为文章响应
func (s *articleService) toArticleResponse(article *models.Article) *ArticleResponse {
	resp := &ArticleResponse{
		ID:            article.ID,
		Title:         article.Title,
		Description:   article.Description,
		CoverImage:    article.CoverImage,
		Status:        article.Status,
		Visibility:    article.Visibility,
		IsTop:         article.IsTop,
		ViewCount:     article.ViewCount,
		LikeCount:     article.LikeCount,
		CommentCount:  article.CommentCount,
		FavoriteCount: article.FavoriteCount,
		PublishedAt:   article.PublishedAt,
		CreatedAt:     article.CreatedAt,
		UpdatedAt:     article.UpdatedAt,
	}
	
	// 添加作者信息
	if article.User.ID != 0 {
		resp.Author = &UserResponse{
			ID:       article.User.ID,
			Username: article.User.Username,
			Nickname: article.User.Nickname,
			Avatar:   article.User.Avatar,
		}
	}
	
	// 添加分类信息
	if article.Category != nil {
		resp.Category = &CategoryResponse{
			ID:          article.Category.ID,
			Name:        article.Category.Name,
			Description: article.Category.Description,
			Color:       article.Category.Color,
			Icon:        article.Category.Icon,
		}
	}
	
	return resp
}

// toListArticleResponse 转换为文章列表响应
func (s *articleService) toListArticleResponse(articles []*models.Article, total int64, page, pageSize int) *ListArticleResponse {
	articleResponses := make([]*ArticleResponse, len(articles))
	for i, article := range articles {
		articleResponses[i] = s.toArticleResponse(article)
	}
	
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}
	
	return &ListArticleResponse{
		Articles: articleResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}
}