package container

import (
	"blog-api/internal/controllers"
	"blog-api/internal/repositories"
	"blog-api/internal/services"
	"blog-api/internal/utils"
	"blog-api/pkg/config"
)

// Container 依赖注入容器
type Container struct {
	config *config.Config

	// Repositories
	userRepo        repositories.UserRepository
	articleRepo     repositories.ArticleRepository
	categoryRepo    repositories.CategoryRepository
	fileRepo        repositories.FileRepository
	commentRepo     repositories.CommentRepository
	articleLikeRepo repositories.ArticleLikeRepository
	commentLikeRepo repositories.CommentLikeRepository
	favoriteRepo    repositories.FavoriteRepository

	// Services
	userService      services.UserService
	articleService   services.ArticleService
	categoryService  services.CategoryService
	fileService      services.FileService
	commentService   services.CommentService
	likeService      services.LikeService
	favoriteService  services.FavoriteService
	adminService     services.AdminService
	adminInitService *services.AdminInitService

	// Utils
	fileUtils *utils.FileUtils

	// Controllers
	authController     *controllers.AuthController
	userController     *controllers.UserController
	articleController  *controllers.ArticleController
	categoryController *controllers.CategoryController
	fileController     *controllers.FileController
	commentController  *controllers.CommentController
	likeController     *controllers.LikeController
	favoriteController *controllers.FavoriteController
	adminController    *controllers.AdminController
}

// NewContainer 创建依赖注入容器
func NewContainer() *Container {
	return &Container{}
}

// InitRepositories 初始化Repository层
func (c *Container) InitRepositories() {
	c.userRepo = repositories.NewUserRepository()
	c.articleRepo = repositories.NewArticleRepository()
	c.categoryRepo = repositories.NewCategoryRepository()
	c.fileRepo = repositories.NewFileRepository()
	c.commentRepo = repositories.NewCommentRepository()
	c.articleLikeRepo = repositories.NewArticleLikeRepository()
	c.commentLikeRepo = repositories.NewCommentLikeRepository()
	c.favoriteRepo = repositories.NewFavoriteRepository()
}

// InitServices 初始化Service层
func (c *Container) InitServices() {
	cfg := config.GetConfig()
	c.config = cfg

	// 初始化文件工具
	c.fileUtils = utils.NewFileUtils(
		cfg.Upload.Path,
		cfg.Upload.BaseURL,
		cfg.Upload.MaxSize,
		cfg.Upload.AllowedTypes,
	)

	// 初始化Services
	c.userService = services.NewUserService()
	c.articleService = services.NewArticleService()
	c.categoryService = services.NewCategoryService()
	c.fileService = services.NewFileServiceDefault()
	c.commentService = services.NewCommentService(c.commentRepo, c.articleRepo, c.commentLikeRepo)
	c.likeService = services.NewLikeService()
	c.favoriteService = services.NewFavoriteService()
	c.adminService = services.NewAdminService()
	c.adminInitService = services.NewAdminInitService(c.userRepo)

}

// InitControllers 初始化Controller层
func (c *Container) InitControllers() {
	c.authController = controllers.NewAuthController(c.userService)
	c.userController = controllers.NewUserController(c.userService)
	c.articleController = controllers.NewArticleController(c.articleService)
	c.categoryController = controllers.NewCategoryController(c.categoryService)
	c.fileController = controllers.NewFileController(c.fileService)
	c.commentController = controllers.NewCommentController(c.commentService)
	c.likeController = controllers.NewLikeController(c.likeService)
	c.favoriteController = controllers.NewFavoriteController(c.favoriteService)
	c.adminController = controllers.NewAdminController(c.adminService)
}

// GetControllers 获取所有Controller
func (c *Container) GetControllers() (
	*controllers.AuthController,
	*controllers.UserController,
	*controllers.ArticleController,
	*controllers.CategoryController,
	*controllers.FileController,
	*controllers.CommentController,
	*controllers.LikeController,
	*controllers.FavoriteController,
	*controllers.AdminController,
) {
	return c.authController, c.userController, c.articleController, c.categoryController, c.fileController, c.commentController, c.likeController, c.favoriteController, c.adminController
}

// GetAdminInitService 获取管理员初始化服务
func (c *Container) GetAdminInitService() *services.AdminInitService {
	return c.adminInitService
}
