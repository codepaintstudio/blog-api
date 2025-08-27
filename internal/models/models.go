package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含公共字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User 用户模型
type User struct {
	BaseModel
	Username      string `gorm:"uniqueIndex;size:50;not null" json:"username" validate:"required,min=3,max=50"`
	Email         string `gorm:"uniqueIndex;size:100;not null" json:"email" validate:"required,email"`
	Password      string `gorm:"size:255;not null" json:"-" validate:"required,min=6"`
	Nickname      string `gorm:"size:50" json:"nickname" validate:"max=50"`
	Avatar        string `gorm:"size:255" json:"avatar"`
	Bio           string `gorm:"type:text" json:"bio" validate:"max=500"`
	Role          string `gorm:"type:enum('admin','user');default:'user'" json:"role"`
	Status        string `gorm:"type:enum('active','inactive');default:'active'" json:"status"`
	EmailVerified bool   `gorm:"default:false" json:"email_verified"`
	LoginAt       *time.Time `json:"login_at"`
	
	// 关联关系
	Articles       []Article       `gorm:"foreignKey:UserID" json:"-"`
	Comments       []Comment       `gorm:"foreignKey:UserID" json:"-"`
	ArticleLikes   []ArticleLike   `gorm:"foreignKey:UserID" json:"-"`
	CommentLikes   []CommentLike   `gorm:"foreignKey:UserID" json:"-"`
	Favorites      []ArticleFavorite `gorm:"foreignKey:UserID" json:"-"`
	FavoriteFolders []FavoriteFolder `gorm:"foreignKey:UserID" json:"-"`
}

// Category 文章分类模型
type Category struct {
	BaseModel
	Name        string `gorm:"size:50;not null" json:"name" validate:"required,max=50"`
	Description string `gorm:"type:text" json:"description"`
	Color       string `gorm:"size:7;default:'#409EFF'" json:"color"` // 分类颜色
	Icon        string `gorm:"size:50" json:"icon"`                   // 分类图标
	Sort        int    `gorm:"default:0" json:"sort"`                 // 排序权重
	IsActive    bool   `gorm:"default:true" json:"is_active"`         // 是否启用
	
	// 关联关系
	Articles []Article `gorm:"foreignKey:CategoryID" json:"-"`
}

// Article 文章模型
type Article struct {
	BaseModel
	UserID        uint   `gorm:"not null;index" json:"user_id"`
	CategoryID    *uint  `gorm:"index" json:"category_id"`
	Title         string `gorm:"size:255;not null" json:"title" validate:"required,max=255"`
	Content       string `gorm:"type:longtext" json:"content" validate:"required"`
	Description   string `gorm:"type:text" json:"description" validate:"max=500"`
	CoverImage    string `gorm:"size:255" json:"cover_image"`
	Status        string `gorm:"type:enum('draft','published');default:'draft'" json:"status"`
	Visibility    string `gorm:"type:enum('public','private');default:'public'" json:"visibility"`
	IsTop         bool   `gorm:"default:false" json:"is_top"`        // 是否置顶
	AllowComment  bool   `gorm:"default:true" json:"allow_comment"`  // 是否允许评论
	ViewCount     int    `gorm:"default:0" json:"view_count"`
	LikeCount     int    `gorm:"default:0" json:"like_count"`
	CommentCount  int    `gorm:"default:0" json:"comment_count"`
	FavoriteCount int    `gorm:"default:0" json:"favorite_count"`
	PublishedAt   *time.Time `json:"published_at"`
	
	// 关联关系
	User         User            `gorm:"foreignKey:UserID" json:"user"`
	Category     *Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Comments     []Comment       `gorm:"foreignKey:ArticleID" json:"-"`
	Likes        []ArticleLike   `gorm:"foreignKey:ArticleID" json:"-"`
	Favorites    []ArticleFavorite `gorm:"foreignKey:ArticleID" json:"-"`
	Tags         []Tag           `gorm:"many2many:article_tags;" json:"tags"`
}

// Tag 标签模型
type Tag struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;size:50;not null" json:"name" validate:"required,max=50"`
	Color       string `gorm:"size:7;default:'#909399'" json:"color"`
	Description string `gorm:"type:text" json:"description"`
	UseCount    int    `gorm:"default:0" json:"use_count"` // 使用次数
	
	// 关联关系
	Articles []Article `gorm:"many2many:article_tags;" json:"-"`
}

// Comment 评论模型
type Comment struct {
	BaseModel
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	ArticleID uint   `gorm:"not null;index" json:"article_id"`
	ParentID  *uint  `gorm:"index" json:"parent_id"`        // 父评论ID，用于回复
	Content   string `gorm:"type:text;not null" json:"content" validate:"required,max=1000"`
	Status    string `gorm:"type:enum('pending','approved','rejected');default:'approved'" json:"status"`
	LikeCount int    `gorm:"default:0" json:"like_count"`
	IP        string `gorm:"size:45" json:"ip"`             // 评论者IP
	UserAgent string `gorm:"size:255" json:"user_agent"`    // 用户代理
	
	// 关联关系
	User     User          `gorm:"foreignKey:UserID" json:"user"`
	Article  Article       `gorm:"foreignKey:ArticleID" json:"article"`
	Parent   *Comment      `gorm:"foreignKey:ParentID" json:"parent"`
	Children []Comment     `gorm:"foreignKey:ParentID" json:"children"`
	Likes    []CommentLike `gorm:"foreignKey:CommentID" json:"-"`
}

// ArticleLike 文章点赞模型
type ArticleLike struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_article,unique" json:"user_id"`
	ArticleID uint      `gorm:"not null;index:idx_user_article,unique" json:"article_id"`
	CreatedAt time.Time `json:"created_at"`
	
	// 关联关系
	User    User    `gorm:"foreignKey:UserID" json:"user"`
	Article Article `gorm:"foreignKey:ArticleID" json:"article"`
}

// CommentLike 评论点赞模型
type CommentLike struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_comment,unique" json:"user_id"`
	CommentID uint      `gorm:"not null;index:idx_user_comment,unique" json:"comment_id"`
	CreatedAt time.Time `json:"created_at"`
	
	// 关联关系
	User    User    `gorm:"foreignKey:UserID" json:"user"`
	Comment Comment `gorm:"foreignKey:CommentID" json:"comment"`
}

// FavoriteFolder 收藏夹模型
type FavoriteFolder struct {
	BaseModel
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	Name        string `gorm:"size:50;not null" json:"name" validate:"required,max=50"`
	Description string `gorm:"type:text" json:"description"`
	IsDefault   bool   `gorm:"default:false" json:"is_default"`     // 是否为默认收藏夹
	IsPublic    bool   `gorm:"default:false" json:"is_public"`      // 是否公开
	Sort        int    `gorm:"default:0" json:"sort"`               // 排序权重
	
	// 关联关系
	User      User              `gorm:"foreignKey:UserID" json:"user"`
	Favorites []ArticleFavorite `gorm:"foreignKey:FolderID" json:"-"`
}

// ArticleFavorite 文章收藏模型
type ArticleFavorite struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	ArticleID uint      `gorm:"not null;index" json:"article_id"`
	FolderID  uint      `gorm:"not null;index" json:"folder_id"`
	CreatedAt time.Time `json:"created_at"`
	
	// 关联关系
	User    User           `gorm:"foreignKey:UserID" json:"user"`
	Article Article        `gorm:"foreignKey:ArticleID" json:"article"`
	Folder  FavoriteFolder `gorm:"foreignKey:FolderID" json:"folder"`
}

// File 文件模型
type File struct {
	BaseModel
	UserID       uint   `gorm:"index" json:"user_id"`             // 上传用户ID
	Filename     string `gorm:"size:255;not null" json:"filename"`
	OriginalName string `gorm:"size:255;not null" json:"original_name"`
	Size         int64  `gorm:"not null" json:"size"`             // 文件大小（字节）
	MimeType     string `gorm:"size:100;not null" json:"mime_type"`
	Path         string `gorm:"size:500;not null" json:"path"`    // 文件存储路径
	URL          string `gorm:"size:500;not null" json:"url"`     // 文件访问URL
	Hash         string `gorm:"size:64;index" json:"hash"`        // 文件哈希值，用于去重
	RefCount     int    `gorm:"default:0" json:"ref_count"`       // 引用计数
	
	// 关联关系
	User *User `gorm:"foreignKey:UserID" json:"user"`
}

// GetAllModels 获取所有模型，用于自动迁移
func GetAllModels() []interface{} {
	return []interface{}{
		&User{},
		&Category{},
		&Article{},
		&Tag{},
		&Comment{},
		&ArticleLike{},
		&CommentLike{},
		&FavoriteFolder{},
		&ArticleFavorite{},
		&File{},
	}
}