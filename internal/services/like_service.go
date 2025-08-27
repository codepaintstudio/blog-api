package services

import (
	"blog-api/internal/models"
	"blog-api/internal/repositories"
	"fmt"
	"time"
	"gorm.io/gorm"
)

// LikeService 点赞服务接口
type LikeService interface {
	// 文章点赞相关
	ToggleArticleLike(userID, articleID uint) (*ArticleLikeResult, error)
	GetArticleLikeStatus(userID, articleID uint) (*ArticleLikeStatus, error)
	GetUserArticleLikes(userID uint, page, size int) ([]*UserLikeResponse, int64, error)
	GetArticleLikers(articleID uint, page, size int) ([]*UserResponse, int64, error)
	
	// 评论点赞相关
	ToggleCommentLike(userID, commentID uint) (*CommentLikeResult, error)
	GetCommentLikeStatus(userID, commentID uint) (*CommentLikeStatus, error)
	GetUserCommentLikes(userID uint, page, size int) ([]*UserCommentLikeResponse, int64, error)
}

// ArticleLikeResult 文章点赞操作结果
type ArticleLikeResult struct {
	IsLiked   bool  `json:"is_liked"`
	LikeCount int64 `json:"like_count"`
	Message   string `json:"message"`
}

// CommentLikeResult 评论点赞操作结果
type CommentLikeResult struct {
	IsLiked   bool  `json:"is_liked"`
	LikeCount int64 `json:"like_count"`
	Message   string `json:"message"`
}

// ArticleLikeStatus 文章点赞状态
type ArticleLikeStatus struct {
	IsLiked   bool  `json:"is_liked"`
	LikeCount int64 `json:"like_count"`
}

// CommentLikeStatus 评论点赞状态
type CommentLikeStatus struct {
	IsLiked   bool  `json:"is_liked"`
	LikeCount int64 `json:"like_count"`
}

// UserLikeResponse 用户点赞文章响应
type UserLikeResponse struct {
	ID        uint      `json:"id"`
	ArticleID uint      `json:"article_id"`
	Article   *ArticleBasicInfo `json:"article"`
	CreatedAt time.Time `json:"created_at"`
}

// UserCommentLikeResponse 用户评论点赞响应
type UserCommentLikeResponse struct {
	ID        uint      `json:"id"`
	CommentID uint      `json:"comment_id"`
	Comment   *CommentBasicInfo `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// ArticleBasicInfo 文章基本信息
type ArticleBasicInfo struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      *UserBasicInfo `json:"author"`
}

// CommentBasicInfo 评论基本信息
type CommentBasicInfo struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
	Author  *UserBasicInfo `json:"author"`
	Article *ArticleBasicInfo `json:"article"`
}

// UserBasicInfo 用户基本信息
type UserBasicInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// likeService 点赞服务实现
type likeService struct {
	articleLikeRepo repositories.ArticleLikeRepository
	commentLikeRepo repositories.CommentLikeRepository
	articleRepo     repositories.ArticleRepository
	commentRepo     repositories.CommentRepository
}

// NewLikeService 创建点赞服务实例
func NewLikeService() LikeService {
	return &likeService{
		articleLikeRepo: repositories.NewArticleLikeRepository(),
		commentLikeRepo: repositories.NewCommentLikeRepository(),
		articleRepo:     repositories.NewArticleRepository(),
		commentRepo:     repositories.NewCommentRepository(),
	}
}

// ToggleArticleLike 切换文章点赞状态
func (s *likeService) ToggleArticleLike(userID, articleID uint) (*ArticleLikeResult, error) {
	// 检查文章是否存在
	article, err := s.articleRepo.GetByID(articleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("文章不存在")
		}
		return nil, fmt.Errorf("查询文章失败: %w", err)
	}

	// 检查是否已点赞
	isLiked, err := s.articleLikeRepo.IsLikedByUser(userID, articleID)
	if err != nil {
		return nil, fmt.Errorf("检查点赞状态失败: %w", err)
	}

	var message string
	if isLiked {
		// 取消点赞
		if err := s.articleLikeRepo.Delete(userID, articleID); err != nil {
			return nil, fmt.Errorf("取消点赞失败: %w", err)
		}
		
		// 更新文章点赞数
		if article.LikeCount > 0 {
			article.LikeCount--
			s.articleRepo.Update(article)
		}
		
		message = "取消点赞成功"
		isLiked = false
	} else {
		// 添加点赞
		like := &models.ArticleLike{
			UserID:    userID,
			ArticleID: articleID,
			CreatedAt: time.Now(),
		}
		
		if err := s.articleLikeRepo.Create(like); err != nil {
			return nil, fmt.Errorf("点赞失败: %w", err)
		}
		
		// 更新文章点赞数
		article.LikeCount++
		s.articleRepo.Update(article)
		
		message = "点赞成功"
		isLiked = true
	}

	// 获取最新点赞数
	likeCount, err := s.articleLikeRepo.GetLikeCount(articleID)
	if err != nil {
		likeCount = int64(article.LikeCount) // 使用缓存的计数
	}

	return &ArticleLikeResult{
		IsLiked:   isLiked,
		LikeCount: likeCount,
		Message:   message,
	}, nil
}

// GetArticleLikeStatus 获取文章点赞状态
func (s *likeService) GetArticleLikeStatus(userID, articleID uint) (*ArticleLikeStatus, error) {
	isLiked, err := s.articleLikeRepo.IsLikedByUser(userID, articleID)
	if err != nil {
		return nil, fmt.Errorf("检查点赞状态失败: %w", err)
	}

	likeCount, err := s.articleLikeRepo.GetLikeCount(articleID)
	if err != nil {
		return nil, fmt.Errorf("获取点赞数失败: %w", err)
	}

	return &ArticleLikeStatus{
		IsLiked:   isLiked,
		LikeCount: likeCount,
	}, nil
}

// GetUserArticleLikes 获取用户点赞的文章列表
func (s *likeService) GetUserArticleLikes(userID uint, page, size int) ([]*UserLikeResponse, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	likes, total, err := s.articleLikeRepo.GetUserLikes(userID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户点赞列表失败: %w", err)
	}

	var responses []*UserLikeResponse
	for _, like := range likes {
		response := &UserLikeResponse{
			ID:        like.ID,
			ArticleID: like.ArticleID,
			CreatedAt: like.CreatedAt,
		}

		if like.Article.ID != 0 {
			response.Article = &ArticleBasicInfo{
				ID:          like.Article.ID,
				Title:       like.Article.Title,
				Description: like.Article.Description,
			}

			if like.Article.User.ID != 0 {
				response.Article.Author = &UserBasicInfo{
					ID:       like.Article.User.ID,
					Username: like.Article.User.Username,
					Nickname: like.Article.User.Nickname,
					Avatar:   like.Article.User.Avatar,
				}
			}
		}

		responses = append(responses, response)
	}

	return responses, total, nil
}

// GetArticleLikers 获取文章点赞用户列表
func (s *likeService) GetArticleLikers(articleID uint, page, size int) ([]*UserResponse, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	users, total, err := s.articleLikeRepo.GetArticleLikers(articleID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取点赞用户列表失败: %w", err)
	}

	var responses []*UserResponse
	for _, user := range users {
		responses = append(responses, &UserResponse{
			ID:            user.ID,
			Username:      user.Username,
			Nickname:      user.Nickname,
			Avatar:        user.Avatar,
			Bio:           user.Bio,
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
		})
	}

	return responses, total, nil
}

// ToggleCommentLike 切换评论点赞状态
func (s *likeService) ToggleCommentLike(userID, commentID uint) (*CommentLikeResult, error) {
	// 检查评论是否存在
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("评论不存在")
		}
		return nil, fmt.Errorf("查询评论失败: %w", err)
	}

	// 检查是否已点赞
	isLiked, err := s.commentLikeRepo.IsLikedByUser(userID, commentID)
	if err != nil {
		return nil, fmt.Errorf("检查点赞状态失败: %w", err)
	}

	var message string
	if isLiked {
		// 取消点赞
		if err := s.commentLikeRepo.Delete(userID, commentID); err != nil {
			return nil, fmt.Errorf("取消点赞失败: %w", err)
		}
		
		// 更新评论点赞数
		if comment.LikeCount > 0 {
			comment.LikeCount--
			s.commentRepo.Update(comment)
		}
		
		message = "取消点赞成功"
		isLiked = false
	} else {
		// 添加点赞
		like := &models.CommentLike{
			UserID:    userID,
			CommentID: commentID,
			CreatedAt: time.Now(),
		}
		
		if err := s.commentLikeRepo.Create(like); err != nil {
			return nil, fmt.Errorf("点赞失败: %w", err)
		}
		
		// 更新评论点赞数
		comment.LikeCount++
		s.commentRepo.Update(comment)
		
		message = "点赞成功"
		isLiked = true
	}

	// 获取最新点赞数
	likeCount, err := s.commentLikeRepo.GetLikeCount(commentID)
	if err != nil {
		likeCount = int64(comment.LikeCount) // 使用缓存的计数
	}

	return &CommentLikeResult{
		IsLiked:   isLiked,
		LikeCount: likeCount,
		Message:   message,
	}, nil
}

// GetCommentLikeStatus 获取评论点赞状态
func (s *likeService) GetCommentLikeStatus(userID, commentID uint) (*CommentLikeStatus, error) {
	isLiked, err := s.commentLikeRepo.IsLikedByUser(userID, commentID)
	if err != nil {
		return nil, fmt.Errorf("检查点赞状态失败: %w", err)
	}

	likeCount, err := s.commentLikeRepo.GetLikeCount(commentID)
	if err != nil {
		return nil, fmt.Errorf("获取点赞数失败: %w", err)
	}

	return &CommentLikeStatus{
		IsLiked:   isLiked,
		LikeCount: likeCount,
	}, nil
}

// GetUserCommentLikes 获取用户点赞的评论列表
func (s *likeService) GetUserCommentLikes(userID uint, page, size int) ([]*UserCommentLikeResponse, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	likes, total, err := s.commentLikeRepo.GetUserCommentLikes(userID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户点赞评论列表失败: %w", err)
	}

	var responses []*UserCommentLikeResponse
	for _, like := range likes {
		response := &UserCommentLikeResponse{
			ID:        like.ID,
			CommentID: like.CommentID,
			CreatedAt: like.CreatedAt,
		}

		if like.Comment.ID != 0 {
			response.Comment = &CommentBasicInfo{
				ID:      like.Comment.ID,
				Content: like.Comment.Content,
			}

			if like.Comment.User.ID != 0 {
				response.Comment.Author = &UserBasicInfo{
					ID:       like.Comment.User.ID,
					Username: like.Comment.User.Username,
					Nickname: like.Comment.User.Nickname,
					Avatar:   like.Comment.User.Avatar,
				}
			}

			if like.Comment.Article.ID != 0 {
				response.Comment.Article = &ArticleBasicInfo{
					ID:          like.Comment.Article.ID,
					Title:       like.Comment.Article.Title,
					Description: like.Comment.Article.Description,
				}
			}
		}

		responses = append(responses, response)
	}

	return responses, total, nil
}