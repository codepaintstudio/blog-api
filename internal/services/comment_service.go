package services

import (
	"errors"
	"strings"

	"blog-api/internal/models"
	"blog-api/internal/repositories"
	"blog-api/pkg/logger"

	"gorm.io/gorm"
)

// CommentCreateRequest 创建评论请求
type CommentCreateRequest struct {
	ArticleID uint   `json:"article_id" validate:"required"`
	ParentID  *uint  `json:"parent_id"`  // 父评论ID，可选
	Content   string `json:"content" validate:"required,max=1000"`
}

// CommentUpdateRequest 更新评论请求  
type CommentUpdateRequest struct {
	Content string `json:"content" validate:"required,max=1000"`
}

// CommentListResponse 评论列表响应
type CommentListResponse struct {
	ID        uint                   `json:"id"`
	UserID    uint                   `json:"user_id"`
	ArticleID uint                   `json:"article_id"`
	ParentID  *uint                  `json:"parent_id"`
	Content   string                 `json:"content"`
	Status    string                 `json:"status"`
	LikeCount int                    `json:"like_count"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
	User      *CommentUserInfo       `json:"user"`
	Children  []*CommentListResponse `json:"children,omitempty"`
	IsLiked   bool                   `json:"is_liked"` // 当前用户是否已点赞
}

// CommentUserInfo 评论用户信息
type CommentUserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// CommentService 评论服务接口
type CommentService interface {
	CreateComment(userID uint, req *CommentCreateRequest, clientIP, userAgent string) (*models.Comment, error)
	UpdateComment(userID, commentID uint, req *CommentUpdateRequest) error
	DeleteComment(userID, commentID uint) error
	GetCommentByID(commentID uint, userID *uint) (*CommentListResponse, error)
	GetArticleComments(articleID uint, page, size int, userID *uint) ([]*CommentListResponse, int64, error)
	GetCommentReplies(parentID uint, page, size int, userID *uint) ([]*CommentListResponse, int64, error)
	GetUserComments(userID uint, page, size int) ([]*CommentListResponse, int64, error)
	ApproveComment(commentID uint) error
	RejectComment(commentID uint) error
}

// commentService 评论服务实现
type commentService struct {
	commentRepo       repositories.CommentRepository
	articleRepo       repositories.ArticleRepository
	commentLikeRepo   repositories.CommentLikeRepository
}

// NewCommentService 创建评论服务实例
func NewCommentService(
	commentRepo repositories.CommentRepository,
	articleRepo repositories.ArticleRepository,
	commentLikeRepo repositories.CommentLikeRepository,
) CommentService {
	return &commentService{
		commentRepo:     commentRepo,
		articleRepo:     articleRepo,
		commentLikeRepo: commentLikeRepo,
	}
}

// CreateComment 创建评论
func (s *commentService) CreateComment(userID uint, req *CommentCreateRequest, clientIP, userAgent string) (*models.Comment, error) {
	// 验证文章是否存在
	article, err := s.articleRepo.GetByID(req.ArticleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		logger.Error("获取文章失败", logger.Err("error", err))
		return nil, errors.New("获取文章失败")
	}

	// 检查文章是否允许评论
	if !article.AllowComment {
		return nil, errors.New("该文章不允许评论")
	}

	// 验证父评论是否存在（如果是回复）
	if req.ParentID != nil {
		parent, err := s.commentRepo.GetByID(*req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("父评论不存在")
			}
			logger.Error("获取父评论失败", logger.Err("error", err))
			return nil, errors.New("获取父评论失败")
		}

		// 确保父评论属于同一篇文章
		if parent.ArticleID != req.ArticleID {
			return nil, errors.New("父评论与文章不匹配")
		}
	}

	// 内容安全检查和处理
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errors.New("评论内容不能为空")
	}

	// 创建评论对象
	comment := &models.Comment{
		UserID:    userID,
		ArticleID: req.ArticleID,
		ParentID:  req.ParentID,
		Content:   content,
		Status:    "approved", // 默认审核通过，可根据需要修改
		IP:        clientIP,
		UserAgent: userAgent,
	}

	// 保存评论
	if err := s.commentRepo.Create(comment); err != nil {
		logger.Error("创建评论失败", logger.Err("error", err))
		return nil, errors.New("创建评论失败")
	}

	// 更新文章评论数
	if err := s.articleRepo.IncrementCommentCount(req.ArticleID); err != nil {
		logger.Error("更新文章评论数失败", logger.Err("error", err))
	}

	logger.Info("创建评论成功", 
		logger.Uint("user_id", userID),
		logger.Uint("article_id", req.ArticleID),
		logger.Uint("comment_id", comment.ID),
	)

	return comment, nil
}

// UpdateComment 更新评论
func (s *commentService) UpdateComment(userID, commentID uint, req *CommentUpdateRequest) error {
	// 获取评论
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		logger.Error("获取评论失败", logger.Err("error", err))
		return errors.New("获取评论失败")
	}

	// 检查权限（只能修改自己的评论）
	if comment.UserID != userID {
		return errors.New("无权限修改此评论")
	}

	// 内容处理
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return errors.New("评论内容不能为空")
	}

	// 更新评论
	comment.Content = content
	if err := s.commentRepo.Update(comment); err != nil {
		logger.Error("更新评论失败", logger.Err("error", err))
		return errors.New("更新评论失败")
	}

	logger.Info("更新评论成功", 
		logger.Uint("user_id", userID),
		logger.Uint("comment_id", commentID),
	)

	return nil
}

// DeleteComment 删除评论
func (s *commentService) DeleteComment(userID, commentID uint) error {
	// 获取评论
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		logger.Error("获取评论失败", logger.Err("error", err))
		return errors.New("获取评论失败")
	}

	// 检查权限（只能删除自己的评论）
	if comment.UserID != userID {
		return errors.New("无权限删除此评论")
	}

	// 删除评论
	if err := s.commentRepo.Delete(commentID); err != nil {
		logger.Error("删除评论失败", logger.Err("error", err))
		return errors.New("删除评论失败")
	}

	// 更新文章评论数
	if err := s.articleRepo.DecrementCommentCount(comment.ArticleID); err != nil {
		logger.Error("更新文章评论数失败", logger.Err("error", err))
	}

	logger.Info("删除评论成功", 
		logger.Uint("user_id", userID),
		logger.Uint("comment_id", commentID),
	)

	return nil
}

// GetCommentByID 获取评论详情
func (s *commentService) GetCommentByID(commentID uint, userID *uint) (*CommentListResponse, error) {
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评论不存在")
		}
		logger.Error("获取评论失败", logger.Err("error", err))
		return nil, errors.New("获取评论失败")
	}

	response := s.convertToCommentResponse(comment, userID)
	return response, nil
}

// GetArticleComments 获取文章评论列表
func (s *commentService) GetArticleComments(articleID uint, page, size int, userID *uint) ([]*CommentListResponse, int64, error) {
	comments, total, err := s.commentRepo.GetByArticleID(articleID, page, size)
	if err != nil {
		logger.Error("获取文章评论列表失败", logger.Err("error", err))
		return nil, 0, errors.New("获取评论列表失败")
	}

	responses := make([]*CommentListResponse, len(comments))
	for i, comment := range comments {
		responses[i] = s.convertToCommentResponse(comment, userID)
		
		// 转换子评论
		if len(comment.Children) > 0 {
			children := make([]*CommentListResponse, len(comment.Children))
			for j, child := range comment.Children {
				children[j] = s.convertToCommentResponse(&child, userID)
			}
			responses[i].Children = children
		}
	}

	return responses, total, nil
}

// GetCommentReplies 获取评论回复列表
func (s *commentService) GetCommentReplies(parentID uint, page, size int, userID *uint) ([]*CommentListResponse, int64, error) {
	replies, total, err := s.commentRepo.GetReplies(parentID, page, size)
	if err != nil {
		logger.Error("获取评论回复列表失败", logger.Err("error", err))
		return nil, 0, errors.New("获取回复列表失败")
	}

	responses := make([]*CommentListResponse, len(replies))
	for i, reply := range replies {
		responses[i] = s.convertToCommentResponse(reply, userID)
	}

	return responses, total, nil
}

// GetUserComments 获取用户评论列表
func (s *commentService) GetUserComments(userID uint, page, size int) ([]*CommentListResponse, int64, error) {
	comments, total, err := s.commentRepo.GetUserComments(userID, page, size)
	if err != nil {
		logger.Error("获取用户评论列表失败", logger.Err("error", err))
		return nil, 0, errors.New("获取评论列表失败")
	}

	responses := make([]*CommentListResponse, len(comments))
	for i, comment := range comments {
		responses[i] = s.convertToCommentResponse(comment, &userID)
	}

	return responses, total, nil
}

// ApproveComment 审核通过评论
func (s *commentService) ApproveComment(commentID uint) error {
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		return err
	}

	comment.Status = "approved"
	return s.commentRepo.Update(comment)
}

// RejectComment 拒绝评论
func (s *commentService) RejectComment(commentID uint) error {
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		return err
	}

	comment.Status = "rejected"
	return s.commentRepo.Update(comment)
}

// convertToCommentResponse 转换评论模型为响应格式
func (s *commentService) convertToCommentResponse(comment *models.Comment, userID *uint) *CommentListResponse {
	response := &CommentListResponse{
		ID:        comment.ID,
		UserID:    comment.UserID,
		ArticleID: comment.ArticleID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		Status:    comment.Status,
		LikeCount: comment.LikeCount,
		CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: comment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// 转换用户信息
	if comment.User.ID > 0 {
		response.User = &CommentUserInfo{
			ID:       comment.User.ID,
			Username: comment.User.Username,
			Nickname: comment.User.Nickname,
			Avatar:   comment.User.Avatar,
		}
	}

	// 检查当前用户是否已点赞
	if userID != nil {
		isLiked, err := s.commentLikeRepo.IsLikedByUser(*userID, comment.ID)
		if err == nil {
			response.IsLiked = isLiked
		}
	}

	return response
}