package services

import (
	"blog-api/internal/models"
	"blog-api/internal/repositories"
	"fmt"
	"time"
	"gorm.io/gorm"
)

// FavoriteService 收藏服务接口
type FavoriteService interface {
	// 收藏夹管理
	CreateFolder(userID uint, req *CreateFolderRequest) (*FolderResponse, error)
	GetUserFolders(userID uint) ([]*FolderResponse, error)
	UpdateFolder(userID, folderID uint, req *UpdateFolderRequest) (*FolderResponse, error)
	DeleteFolder(userID, folderID uint) error
	GetFolderDetails(userID, folderID uint) (*FolderDetailResponse, error)
	
	// 文章收藏管理
	ToggleFavorite(userID, articleID, folderID uint) (*FavoriteResult, error)
	GetFavoriteStatus(userID, articleID uint) (*FavoriteStatus, error)
	GetUserFavorites(userID uint, page, size int) ([]*UserFavoriteResponse, int64, error)
	GetFolderFavorites(userID, folderID uint, page, size int) ([]*FolderFavoriteResponse, int64, error)
	MoveFavorite(userID, articleID, targetFolderID uint) error
	
	// Controller所需方法
	AddToFavorite(userID uint, req *FavoriteAddRequest) error
	RemoveFromFavorite(userID, articleID uint, folderID *uint) error
	GetFolderArticles(userID, folderID uint, page, size int) ([]*FolderFavoriteResponse, int64, error)
	CheckFavoriteStatus(userID, articleID uint) (*FavoriteStatusResponse, error)
	
	// 初始化用户默认收藏夹
	InitUserDefaultFolder(userID uint) (*models.FavoriteFolder, error)
}

// CreateFolderRequest 创建收藏夹请求
type CreateFolderRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description" validate:"max=200"`
	IsPublic    bool   `json:"is_public"`
}

// UpdateFolderRequest 更新收藏夹请求
type UpdateFolderRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description" validate:"max=200"`
	IsPublic    bool   `json:"is_public"`
}

// FavoriteAddRequest 添加收藏请求
type FavoriteAddRequest struct {
	ArticleID uint `json:"article_id" validate:"required"`
	FolderID  uint `json:"folder_id" validate:"required"`
}

// FavoriteStatusResponse 收藏状态响应
type FavoriteStatusResponse struct {
	IsFavorited bool                 `json:"is_favorited"`
	Folders     []*FolderBasicInfo   `json:"folders,omitempty"`
}

// FolderResponse 收藏夹响应
type FolderResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	IsDefault    bool      `json:"is_default"`
	IsPublic     bool      `json:"is_public"`
	Sort         int       `json:"sort"`
	ArticleCount int64     `json:"article_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// FolderDetailResponse 收藏夹详情响应
type FolderDetailResponse struct {
	*FolderResponse
	Articles []*FolderFavoriteResponse `json:"articles"`
}

// FavoriteResult 收藏操作结果
type FavoriteResult struct {
	IsFavorited   bool   `json:"is_favorited"`
	FavoriteCount int64  `json:"favorite_count"`
	FolderName    string `json:"folder_name"`
	Message       string `json:"message"`
}

// FavoriteStatus 收藏状态
type FavoriteStatus struct {
	IsFavorited   bool   `json:"is_favorited"`
	FavoriteCount int64  `json:"favorite_count"`
	FolderID      uint   `json:"folder_id,omitempty"`
	FolderName    string `json:"folder_name,omitempty"`
}

// UserFavoriteResponse 用户收藏响应
type UserFavoriteResponse struct {
	ID        uint                `json:"id"`
	Article   *ArticleBasicInfo   `json:"article"`
	Folder    *FolderBasicInfo    `json:"folder"`
	CreatedAt time.Time           `json:"created_at"`
}

// FolderFavoriteResponse 收藏夹文章响应
type FolderFavoriteResponse struct {
	ID        uint              `json:"id"`
	Article   *ArticleBasicInfo `json:"article"`
	CreatedAt time.Time         `json:"created_at"`
}

// FolderBasicInfo 收藏夹基本信息
type FolderBasicInfo struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
	IsPublic  bool   `json:"is_public"`
}

// favoriteService 收藏服务实现
type favoriteService struct {
	favoriteRepo repositories.FavoriteRepository
	articleRepo  repositories.ArticleRepository
}

// NewFavoriteService 创建收藏服务实例
func NewFavoriteService() FavoriteService {
	return &favoriteService{
		favoriteRepo: repositories.NewFavoriteRepository(),
		articleRepo:  repositories.NewArticleRepository(),
	}
}

// CreateFolder 创建收藏夹
func (s *favoriteService) CreateFolder(userID uint, req *CreateFolderRequest) (*FolderResponse, error) {
	// 检查收藏夹名称是否已存在
	exists, err := s.favoriteRepo.FolderExists(userID, req.Name)
	if err != nil {
		return nil, fmt.Errorf("检查收藏夹名称失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("收藏夹名称已存在")
	}

	// 创建收藏夹
	folder := &models.FavoriteFolder{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		IsDefault:   false,
		IsPublic:    req.IsPublic,
		Sort:        0,
	}

	if err := s.favoriteRepo.CreateFolder(folder); err != nil {
		return nil, fmt.Errorf("创建收藏夹失败: %w", err)
	}

	return &FolderResponse{
		ID:           folder.ID,
		Name:         folder.Name,
		Description:  folder.Description,
		IsDefault:    folder.IsDefault,
		IsPublic:     folder.IsPublic,
		Sort:         folder.Sort,
		ArticleCount: 0,
		CreatedAt:    folder.CreatedAt,
		UpdatedAt:    folder.UpdatedAt,
	}, nil
}

// GetUserFolders 获取用户收藏夹列表
func (s *favoriteService) GetUserFolders(userID uint) ([]*FolderResponse, error) {
	folders, err := s.favoriteRepo.GetUserFolders(userID)
	if err != nil {
		return nil, fmt.Errorf("获取收藏夹列表失败: %w", err)
	}

	var responses []*FolderResponse
	for _, folder := range folders {
		// 这里可以异步获取文章数量，为了简化先设为0
		responses = append(responses, &FolderResponse{
			ID:           folder.ID,
			Name:         folder.Name,
			Description:  folder.Description,
			IsDefault:    folder.IsDefault,
			IsPublic:     folder.IsPublic,
			Sort:         folder.Sort,
			ArticleCount: 0, // TODO: 获取实际文章数量
			CreatedAt:    folder.CreatedAt,
			UpdatedAt:    folder.UpdatedAt,
		})
	}

	return responses, nil
}

// UpdateFolder 更新收藏夹
func (s *favoriteService) UpdateFolder(userID, folderID uint, req *UpdateFolderRequest) (*FolderResponse, error) {
	// 获取收藏夹
	folder, err := s.favoriteRepo.GetFolderByID(folderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("收藏夹不存在")
		}
		return nil, fmt.Errorf("获取收藏夹失败: %w", err)
	}

	// 检查权限
	if folder.UserID != userID {
		return nil, fmt.Errorf("无权限操作此收藏夹")
	}

	// 检查是否是默认收藏夹（默认收藏夹不能重命名）
	if folder.IsDefault && folder.Name != req.Name {
		return nil, fmt.Errorf("默认收藏夹不能重命名")
	}

	// 检查新名称是否已存在（排除当前收藏夹）
	if folder.Name != req.Name {
		exists, err := s.favoriteRepo.FolderExists(userID, req.Name)
		if err != nil {
			return nil, fmt.Errorf("检查收藏夹名称失败: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("收藏夹名称已存在")
		}
	}

	// 更新收藏夹信息
	folder.Name = req.Name
	folder.Description = req.Description
	folder.IsPublic = req.IsPublic

	if err := s.favoriteRepo.UpdateFolder(folder); err != nil {
		return nil, fmt.Errorf("更新收藏夹失败: %w", err)
	}

	return &FolderResponse{
		ID:           folder.ID,
		Name:         folder.Name,
		Description:  folder.Description,
		IsDefault:    folder.IsDefault,
		IsPublic:     folder.IsPublic,
		Sort:         folder.Sort,
		ArticleCount: 0,
		CreatedAt:    folder.CreatedAt,
		UpdatedAt:    folder.UpdatedAt,
	}, nil
}

// DeleteFolder 删除收藏夹
func (s *favoriteService) DeleteFolder(userID, folderID uint) error {
	// 获取收藏夹
	folder, err := s.favoriteRepo.GetFolderByID(folderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("收藏夹不存在")
		}
		return fmt.Errorf("获取收藏夹失败: %w", err)
	}

	// 检查权限
	if folder.UserID != userID {
		return fmt.Errorf("无权限操作此收藏夹")
	}

	// 不能删除默认收藏夹
	if folder.IsDefault {
		return fmt.Errorf("不能删除默认收藏夹")
	}

	// TODO: 处理收藏夹中的文章（移动到默认收藏夹或删除收藏）
	
	// 删除收藏夹
	if err := s.favoriteRepo.DeleteFolder(folderID); err != nil {
		return fmt.Errorf("删除收藏夹失败: %w", err)
	}

	return nil
}

// ToggleFavorite 切换文章收藏状态
func (s *favoriteService) ToggleFavorite(userID, articleID, folderID uint) (*FavoriteResult, error) {
	// 检查文章是否存在
	article, err := s.articleRepo.GetByID(articleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("文章不存在")
		}
		return nil, fmt.Errorf("查询文章失败: %w", err)
	}

	// 检查是否已收藏
	isFavorited, err := s.favoriteRepo.IsFavorited(userID, articleID)
	if err != nil {
		return nil, fmt.Errorf("检查收藏状态失败: %w", err)
	}

	var folderName string
	var message string

	if isFavorited {
		// 取消收藏
		if err := s.favoriteRepo.RemoveFavorite(userID, articleID); err != nil {
			return nil, fmt.Errorf("取消收藏失败: %w", err)
		}

		// 更新文章收藏数
		if article.FavoriteCount > 0 {
			article.FavoriteCount--
			s.articleRepo.Update(article)
		}

		message = "取消收藏成功"
		isFavorited = false
	} else {
		// 添加收藏
		// 如果没有指定收藏夹，使用默认收藏夹
		if folderID == 0 {
			defaultFolder, err := s.favoriteRepo.GetDefaultFolder(userID)
			if err != nil {
				// 如果没有默认收藏夹，创建一个
				if err == gorm.ErrRecordNotFound {
					defaultFolder, err = s.InitUserDefaultFolder(userID)
					if err != nil {
						return nil, fmt.Errorf("创建默认收藏夹失败: %w", err)
					}
				} else {
					return nil, fmt.Errorf("获取默认收藏夹失败: %w", err)
				}
			}
			folderID = defaultFolder.ID
			folderName = defaultFolder.Name
		} else {
			// 验证收藏夹是否属于当前用户
			folder, err := s.favoriteRepo.GetFolderByID(folderID)
			if err != nil {
				return nil, fmt.Errorf("收藏夹不存在")
			}
			if folder.UserID != userID {
				return nil, fmt.Errorf("无权限操作此收藏夹")
			}
			folderName = folder.Name
		}

		favorite := &models.ArticleFavorite{
			UserID:    userID,
			ArticleID: articleID,
			FolderID:  folderID,
			CreatedAt: time.Now(),
		}

		if err := s.favoriteRepo.AddFavorite(favorite); err != nil {
			return nil, fmt.Errorf("收藏失败: %w", err)
		}

		// 更新文章收藏数
		article.FavoriteCount++
		s.articleRepo.Update(article)

		message = "收藏成功"
		isFavorited = true
	}

	// 获取最新收藏数
	favoriteCount, err := s.favoriteRepo.GetFavoriteCount(articleID)
	if err != nil {
		favoriteCount = int64(article.FavoriteCount) // 使用缓存的计数
	}

	return &FavoriteResult{
		IsFavorited:   isFavorited,
		FavoriteCount: favoriteCount,
		FolderName:    folderName,
		Message:       message,
	}, nil
}

// GetFavoriteStatus 获取文章收藏状态
func (s *favoriteService) GetFavoriteStatus(userID, articleID uint) (*FavoriteStatus, error) {
	isFavorited, err := s.favoriteRepo.IsFavorited(userID, articleID)
	if err != nil {
		return nil, fmt.Errorf("检查收藏状态失败: %w", err)
	}

	favoriteCount, err := s.favoriteRepo.GetFavoriteCount(articleID)
	if err != nil {
		return nil, fmt.Errorf("获取收藏数失败: %w", err)
	}

	status := &FavoriteStatus{
		IsFavorited:   isFavorited,
		FavoriteCount: favoriteCount,
	}

	// 如果已收藏，获取收藏夹信息
	if isFavorited {
		favorite, err := s.favoriteRepo.GetFavorite(userID, articleID)
		if err == nil && favorite.Folder.ID != 0 {
			status.FolderID = favorite.FolderID
			status.FolderName = favorite.Folder.Name
		}
	}

	return status, nil
}

// GetUserFavorites 获取用户收藏列表
func (s *favoriteService) GetUserFavorites(userID uint, page, size int) ([]*UserFavoriteResponse, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	favorites, total, err := s.favoriteRepo.GetUserFavorites(userID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户收藏列表失败: %w", err)
	}

	var responses []*UserFavoriteResponse
	for _, favorite := range favorites {
		response := &UserFavoriteResponse{
			ID:        favorite.ID,
			CreatedAt: favorite.CreatedAt,
		}

		// 文章信息
		if favorite.Article.ID != 0 {
			response.Article = &ArticleBasicInfo{
				ID:          favorite.Article.ID,
				Title:       favorite.Article.Title,
				Description: favorite.Article.Description,
			}

			if favorite.Article.User.ID != 0 {
				response.Article.Author = &UserBasicInfo{
					ID:       favorite.Article.User.ID,
					Username: favorite.Article.User.Username,
					Nickname: favorite.Article.User.Nickname,
					Avatar:   favorite.Article.User.Avatar,
				}
			}
		}

		// 收藏夹信息
		if favorite.Folder.ID != 0 {
			response.Folder = &FolderBasicInfo{
				ID:        favorite.Folder.ID,
				Name:      favorite.Folder.Name,
				IsDefault: favorite.Folder.IsDefault,
				IsPublic:  favorite.Folder.IsPublic,
			}
		}

		responses = append(responses, response)
	}

	return responses, total, nil
}

// GetFolderFavorites 获取收藏夹中的文章
func (s *favoriteService) GetFolderFavorites(userID, folderID uint, page, size int) ([]*FolderFavoriteResponse, int64, error) {
	// 验证收藏夹权限
	folder, err := s.favoriteRepo.GetFolderByID(folderID)
	if err != nil {
		return nil, 0, fmt.Errorf("收藏夹不存在")
	}
	if folder.UserID != userID && !folder.IsPublic {
		return nil, 0, fmt.Errorf("无权限访问此收藏夹")
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	favorites, total, err := s.favoriteRepo.GetFolderFavorites(folderID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取收藏夹文章失败: %w", err)
	}

	var responses []*FolderFavoriteResponse
	for _, favorite := range favorites {
		response := &FolderFavoriteResponse{
			ID:        favorite.ID,
			CreatedAt: favorite.CreatedAt,
		}

		if favorite.Article.ID != 0 {
			response.Article = &ArticleBasicInfo{
				ID:          favorite.Article.ID,
				Title:       favorite.Article.Title,
				Description: favorite.Article.Description,
			}

			if favorite.Article.User.ID != 0 {
				response.Article.Author = &UserBasicInfo{
					ID:       favorite.Article.User.ID,
					Username: favorite.Article.User.Username,
					Nickname: favorite.Article.User.Nickname,
					Avatar:   favorite.Article.User.Avatar,
				}
			}
		}

		responses = append(responses, response)
	}

	return responses, total, nil
}

// MoveFavorite 移动收藏到其他收藏夹
func (s *favoriteService) MoveFavorite(userID, articleID, targetFolderID uint) error {
	// 检查是否已收藏
	isFavorited, err := s.favoriteRepo.IsFavorited(userID, articleID)
	if err != nil {
		return fmt.Errorf("检查收藏状态失败: %w", err)
	}
	if !isFavorited {
		return fmt.Errorf("文章未收藏")
	}

	// 验证目标收藏夹
	targetFolder, err := s.favoriteRepo.GetFolderByID(targetFolderID)
	if err != nil {
		return fmt.Errorf("目标收藏夹不存在")
	}
	if targetFolder.UserID != userID {
		return fmt.Errorf("无权限操作目标收藏夹")
	}

	// 移动收藏
	if err := s.favoriteRepo.MoveFavorite(userID, articleID, targetFolderID); err != nil {
		return fmt.Errorf("移动收藏失败: %w", err)
	}

	return nil
}

// InitUserDefaultFolder 初始化用户默认收藏夹
func (s *favoriteService) InitUserDefaultFolder(userID uint) (*models.FavoriteFolder, error) {
	defaultFolder := &models.FavoriteFolder{
		UserID:      userID,
		Name:        "默认收藏夹",
		Description: "系统默认创建的收藏夹",
		IsDefault:   true,
		IsPublic:    false,
		Sort:        0,
	}

	if err := s.favoriteRepo.CreateFolder(defaultFolder); err != nil {
		return nil, err
	}

	return defaultFolder, nil
}

// GetFolderDetails 获取收藏夹详情
func (s *favoriteService) GetFolderDetails(userID, folderID uint) (*FolderDetailResponse, error) {
	// 验证收藏夹权限
	folder, err := s.favoriteRepo.GetFolderByID(folderID)
	if err != nil {
		return nil, fmt.Errorf("收藏夹不存在")
	}
	if folder.UserID != userID && !folder.IsPublic {
		return nil, fmt.Errorf("无权限访问此收藏夹")
	}

	// 获取收藏夹文章（前10篇）
	articles, total, err := s.GetFolderFavorites(userID, folderID, 1, 10)
	if err != nil {
		return nil, err
	}

	return &FolderDetailResponse{
		FolderResponse: &FolderResponse{
			ID:           folder.ID,
			Name:         folder.Name,
			Description:  folder.Description,
			IsDefault:    folder.IsDefault,
			IsPublic:     folder.IsPublic,
			Sort:         folder.Sort,
			ArticleCount: total,
			CreatedAt:    folder.CreatedAt,
			UpdatedAt:    folder.UpdatedAt,
		},
		Articles: articles,
	}, nil
}

// AddToFavorite 添加文章到收藏夹
func (s *favoriteService) AddToFavorite(userID uint, req *FavoriteAddRequest) error {
	// 验证文章是否存在
	_, err := s.articleRepo.GetByID(req.ArticleID)
	if err != nil {
		return fmt.Errorf("文章不存在")
	}

	// 验证收藏夹是否存在且属于用户
	folder, err := s.favoriteRepo.GetFolderByID(req.FolderID)
	if err != nil {
		return fmt.Errorf("收藏夹不存在")
	}
	if folder.UserID != userID {
		return fmt.Errorf("无权限操作此收藏夹")
	}

	// 检查是否已收藏
	isFavorited, err := s.favoriteRepo.IsFavorited(userID, req.ArticleID)
	if err != nil {
		return fmt.Errorf("检查收藏状态失败: %w", err)
	}
	if isFavorited {
		return fmt.Errorf("文章已收藏")
	}

	// 添加收藏
	favorite := &models.ArticleFavorite{
		UserID:    userID,
		ArticleID: req.ArticleID,
		FolderID:  req.FolderID,
	}

	if err := s.favoriteRepo.AddFavorite(favorite); err != nil {
		return fmt.Errorf("添加收藏失败: %w", err)
	}

	// 更新文章收藏数
	if err := s.articleRepo.IncrementFavoriteCount(req.ArticleID); err != nil {
		// 记录错误但不影响主流程
	}

	return nil
}

// RemoveFromFavorite 从收藏夹移除文章
func (s *favoriteService) RemoveFromFavorite(userID, articleID uint, folderID *uint) error {
	// 检查是否已收藏
	isFavorited, err := s.favoriteRepo.IsFavorited(userID, articleID)
	if err != nil {
		return fmt.Errorf("检查收藏状态失败: %w", err)
	}
	if !isFavorited {
		return fmt.Errorf("文章未收藏")
	}

	// 移除收藏
	if err := s.favoriteRepo.RemoveFavorite(userID, articleID); err != nil {
		return fmt.Errorf("移除收藏失败: %w", err)
	}

	// 更新文章收藏数
	if err := s.articleRepo.DecrementFavoriteCount(articleID); err != nil {
		// 记录错误但不影响主流程
	}

	return nil
}

// GetFolderArticles 获取收藏夹中的文章列表
func (s *favoriteService) GetFolderArticles(userID, folderID uint, page, size int) ([]*FolderFavoriteResponse, int64, error) {
	return s.GetFolderFavorites(userID, folderID, page, size)
}

// CheckFavoriteStatus 检查文章收藏状态
func (s *favoriteService) CheckFavoriteStatus(userID, articleID uint) (*FavoriteStatusResponse, error) {
	// 检查是否已收藏
	isFavorited, err := s.favoriteRepo.IsFavorited(userID, articleID)
	if err != nil {
		return nil, fmt.Errorf("检查收藏状态失败: %w", err)
	}

	status := &FavoriteStatusResponse{
		IsFavorited: isFavorited,
	}

	// 如果已收藏，获取收藏夹信息
	if isFavorited {
		favorite, err := s.favoriteRepo.GetFavorite(userID, articleID)
		if err == nil && favorite.Folder.ID != 0 {
			status.Folders = append(status.Folders, &FolderBasicInfo{
				ID:        favorite.Folder.ID,
				Name:      favorite.Folder.Name,
				IsDefault: favorite.Folder.IsDefault,
				IsPublic:  favorite.Folder.IsPublic,
			})
		}
	}

	return status, nil
}