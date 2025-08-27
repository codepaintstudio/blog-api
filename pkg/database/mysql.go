package database

import (
	"fmt"
	"time"

	"blog-api/pkg/config"
	"blog-api/pkg/logger"
	"blog-api/internal/models"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitMySQL 初始化MySQL连接
func InitMySQL() error {
	cfg := config.GetConfig()
	if cfg == nil {
		return fmt.Errorf("配置未加载")
	}

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.Charset,
	)

	// 设置GORM日志级别
	var logLevel gormLogger.LogLevel
	switch cfg.Log.Level {
	case "debug":
		logLevel = gormLogger.Info
	case "info":
		logLevel = gormLogger.Info
	case "warn":
		logLevel = gormLogger.Warn
	case "error":
		logLevel = gormLogger.Error
	default:
		logLevel = gormLogger.Warn
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("连接MySQL失败: %w", err)
	}

	// 获取底层sql.DB实例
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层sql.DB失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.Database.ConnMaxIdleTime) * time.Second)

	DB = db
	logger.Info("MySQL连接成功", 
		logger.String("host", cfg.Database.Host),
		logger.Int("port", cfg.Database.Port),
		logger.String("database", cfg.Database.DBName),
	)
	
	// 执行自动迁移
	if err := AutoMigrate(); err != nil {
		return fmt.Errorf("数据库自动迁移失败: %w", err)
	}
	
	return nil
}

// AutoMigrate 执行数据库自动迁移
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}
	
	logger.Info("开始执行数据库自动迁移...")
	
	// 获取所有模型
	allModels := models.GetAllModels()
	
	// 执行迁移
	if err := DB.AutoMigrate(allModels...); err != nil {
		return fmt.Errorf("自动迁移失败: %w", err)
	}
	
	logger.Info("数据库自动迁移完成")
	
	// 创建默认数据
	if err := createDefaultData(); err != nil {
		logger.Warn("创建默认数据失败", logger.Err("error", err))
	}
	
	return nil
}

// createDefaultData 创建默认数据
func createDefaultData() error {
	// 创建默认分类
	var categoryCount int64
	DB.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount == 0 {
		defaultCategories := []models.Category{
			{
				Name:        "技术",
				Description: "技术相关文章",
				Color:       "#409EFF",
				Icon:        "tech",
				Sort:        1,
				IsActive:    true,
			},
			{
				Name:        "生活",
				Description: "生活随笔",
				Color:       "#67C23A",
				Icon:        "life",
				Sort:        2,
				IsActive:    true,
			},
			{
				Name:        "随笔",
				Description: "个人随笔",
				Color:       "#E6A23C",
				Icon:        "note",
				Sort:        3,
				IsActive:    true,
			},
		}
		
		for _, category := range defaultCategories {
			if err := DB.Create(&category).Error; err != nil {
				return fmt.Errorf("创建默认分类失败: %w", err)
			}
		}
		
		logger.Info("创建默认分类成功")
	}
	
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseMySQL 关闭MySQL连接
func CloseMySQL() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}